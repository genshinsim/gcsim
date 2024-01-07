package simulator

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"regexp"
	"runtime/debug"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/genshinsim/gcsim/pkg/agg"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/gcs/ast"
	"github.com/genshinsim/gcsim/pkg/model"
	"github.com/genshinsim/gcsim/pkg/result"
	"github.com/genshinsim/gcsim/pkg/stats"
	"github.com/genshinsim/gcsim/pkg/worker"
)

// Options sets out the settings to run the sim by (such as debug mode, etc..)
type Options struct {
	ResultSaveToPath string // file name (excluding ext) to save the result file; if "" then nothing is saved to file
	GZIPResult       bool   // should the result file be gzipped; only if ResultSaveToPath is not ""
	ConfigPath       string // path to the config file to read
}

var (
	sha1ver   string
	buildTime string
	modified  bool
)

func init() {
	info, _ := debug.ReadBuildInfo()
	for _, bs := range info.Settings {
		if bs.Key == "vcs.revision" {
			sha1ver = bs.Value
		}
		if bs.Key == "vcs.time" {
			buildTime = bs.Value
		}
		if bs.Key == "vcs.modified" {
			bv, _ := strconv.ParseBool(bs.Value)
			modified = bv
		}
	}
}

func Version() string {
	return sha1ver
}

func Parse(cfg string) (*info.ActionList, ast.Node, error) {
	parser := ast.New(cfg)
	simcfg, gcsl, err := parser.Parse()
	if err != nil {
		return &info.ActionList{}, nil, err
	}

	// check other errors as well
	if len(simcfg.Errors) != 0 {
		fmt.Println("The config has the following errors: ")
		errMsgs := ""
		for _, v := range simcfg.Errors {
			errMsg := fmt.Sprintf("\t%v\n", v)
			fmt.Println(errMsg)
			errMsgs += errMsg
		}
		return &info.ActionList{}, nil, errors.New(errMsgs)
	}

	return simcfg, gcsl, nil
}

// Run will run the simulation given number of times
func Run(ctx context.Context, opts Options) (*model.SimulationResult, error) {
	start := time.Now()

	cfg, err := ReadConfig(opts.ConfigPath)
	if err != nil {
		return &model.SimulationResult{}, err
	}

	simcfg, gcsl, err := Parse(cfg)
	if err != nil {
		return &model.SimulationResult{}, err
	}

	return RunWithConfig(ctx, cfg, simcfg, gcsl, opts, start)
}

// Runs the simulation with a given parsed config
// TODO: cfg string should be in the action list instead
// TODO: need to add a context here to avoid infinite looping
func RunWithConfig(ctx context.Context, cfg string, simcfg *info.ActionList, gcsl ast.Node, opts Options, start time.Time) (*model.SimulationResult, error) {
	// initialize aggregators
	var aggregators []agg.Aggregator
	for _, aggregator := range agg.Aggregators() {
		enabled := simcfg.Settings.CollectStats
		if len(enabled) > 0 && !slices.Contains(enabled, aggregator.Name) {
			continue
		}
		a, err := aggregator.New(simcfg)
		if err != nil {
			return &model.SimulationResult{}, err
		}
		aggregators = append(aggregators, a)
	}

	// set up a pool
	respCh := make(chan stats.Result)
	errCh := make(chan error)
	pool := worker.New(simcfg.Settings.NumberOfWorkers, respCh, errCh)
	pool.StopCh = make(chan bool)

	// spin off a go func that will queue jobs for as long as the total queued < iter
	// this should block as queue gets full
	go func() {
		// make all the seeds
		wip := 0
		for wip < simcfg.Settings.Iterations {
			pool.QueueCh <- worker.Job{
				Cfg:     simcfg.Copy(),
				Actions: gcsl.Copy(),
				Seed:    CryptoRandSeed(),
			}
			wip++
		}
	}()

	defer close(pool.StopCh)

	// start reading respCh, queueing a new job until wip == number of iterations
	count := 0
	for count < simcfg.Settings.Iterations {
		select {
		case result := <-respCh:
			for _, a := range aggregators {
				a.Add(result)
			}
			count += 1
		case err := <-errCh:
			// error encountered
			return &model.SimulationResult{}, err
		case <-ctx.Done():
			return &model.SimulationResult{}, ctx.Err()
		}
	}

	result, err := GenerateResult(cfg, simcfg, opts)
	if err != nil {
		return result, err
	}

	// generate final agg results
	stats := &model.SimulationStatistics{}
	for _, a := range aggregators {
		a.Flush(stats)
	}
	result.Statistics = stats

	return result, nil
}

// Runs the simulation with a given parsed config and custom stat collector and aggregator
// TODO: cfg string should be in the action list instead
// TODO: need to add a context here to avoid infinite looping
func RunWithConfigCustomStats[T any](ctx context.Context, cfg string, simcfg *info.ActionList, gcsl ast.Node, opts Options, start time.Time, cstat stats.NewStatsFuncCustomStats[T], cagg func(T)) (*model.SimulationResult, error) {
	// initialize aggregators
	var aggregators []agg.Aggregator
	for _, aggregator := range agg.Aggregators() {
		enabled := simcfg.Settings.CollectStats
		if len(enabled) > 0 && !slices.Contains(enabled, aggregator.Name) {
			continue
		}
		a, err := aggregator.New(simcfg)
		if err != nil {
			return &model.SimulationResult{}, err
		}
		aggregators = append(aggregators, a)
	}

	// set up a pool
	respCh := make(chan stats.Result)
	errCh := make(chan error)
	customCh := make(chan T)
	pool := worker.NewWithCustomStats(simcfg.Settings.NumberOfWorkers, respCh, errCh, customCh)
	pool.StopCh = make(chan bool)

	// spin off a go func that will queue jobs for as long as the total queued < iter
	// this should block as queue gets full
	go func() {
		// make all the seeds
		wip := 0
		for wip < simcfg.Settings.Iterations {
			pool.QueueCh <- worker.JobCustomStats[T]{
				Cfg:     simcfg.Copy(),
				Actions: gcsl.Copy(),
				Seed:    CryptoRandSeed(),
				Cstat:   cstat,
			}
			wip++
		}
	}()

	defer close(pool.StopCh)

	// start reading respCh, queueing a new job until wip == number of iterations
	count := 0
	for count < simcfg.Settings.Iterations {
		select {
		case result := <-customCh:
			cagg(result)
		case result := <-respCh:
			for _, a := range aggregators {
				a.Add(result)
			}
			count += 1
		case err := <-errCh:
			// error encountered
			return &model.SimulationResult{}, err
		case <-ctx.Done():
			return &model.SimulationResult{}, ctx.Err()
		}
	}

	result, err := GenerateResult(cfg, simcfg, opts)
	if err != nil {
		return result, err
	}

	// generate final agg results
	stats := &model.SimulationStatistics{}
	for _, a := range aggregators {
		a.Flush(stats)
	}
	result.Statistics = stats

	return result, nil
}

// Note: this generation should be iteration independent (iterations do not change output)
func GenerateResult(cfg string, simcfg *info.ActionList, opts Options) (*model.SimulationResult, error) {
	out := &model.SimulationResult{
		// THIS MUST ALWAYS BE IN SYNC WITH THE VIEWER UPGRADE DIALOG IN UI
		// ONLY CHANGE SCHEMA WHEN THE RESULTS SCHEMA CHANGES. THIS INCLUDES AGG RESULTS CHANGES
		// SemVer spec
		//    Major: increase & reset minor to zero if new schema is backwards incompatible
		//        Ex - changed the location of a critical column (the config file), major refactor
		//    Minor: increase if new schema is backwards compatible with previous
		//        Ex - added new data for new graph on UI. UI still functional if this data is missing
		// Increasing the version will result in the UI flagging all old sims as outdated
		SchemaVersion: &model.Version{Major: "4", Minor: "2"}, // MAKE SURE UI VERSION IS IN SYNC in ui/packages/ui/src/Pages/Viewer/UpgradeDialog.tsx
		SimVersion:    &sha1ver,
		BuildDate:     buildTime,
		Modified:      &modified,
		KeyType:       "NONE",
		SimulatorSettings: &model.SimulatorSettings{
			Duration:          simcfg.Settings.Duration,
			DamageMode:        simcfg.Settings.DamageMode,
			EnableHitlag:      simcfg.Settings.EnableHitlag,
			DefHalt:           simcfg.Settings.DefHalt,
			IgnoreBurstEnergy: simcfg.Settings.IgnoreBurstEnergy,
			NumberOfWorkers:   uint32(simcfg.Settings.NumberOfWorkers),
			Iterations:        uint32(simcfg.Settings.Iterations),
			Delays: &model.Delays{
				Skill:  int32(simcfg.Settings.Delays.Skill),
				Burst:  int32(simcfg.Settings.Delays.Burst),
				Attack: int32(simcfg.Settings.Delays.Attack),
				Charge: int32(simcfg.Settings.Delays.Charge),
				Aim:    int32(simcfg.Settings.Delays.Aim),
				Dash:   int32(simcfg.Settings.Delays.Dash),
				Jump:   int32(simcfg.Settings.Delays.Jump),
				Swap:   int32(simcfg.Settings.Delays.Swap),
			},
		},
		EnergySettings: &model.EnergySettings{
			Active:         simcfg.EnergySettings.Active,
			Once:           simcfg.EnergySettings.Once,
			Start:          int32(simcfg.EnergySettings.Start),
			End:            int32(simcfg.EnergySettings.End),
			Amount:         int32(simcfg.EnergySettings.Amount),
			LastEnergyDrop: int32(simcfg.EnergySettings.LastEnergyDrop),
		},
		Config:           cfg,
		SampleSeed:       strconv.FormatUint(uint64(CryptoRandSeed()), 10),
		InitialCharacter: simcfg.InitialChar.String(),
		TargetDetails:    make([]*model.Enemy, len(simcfg.Targets)),
		PlayerPosition: &model.Coord{
			X: simcfg.InitialPlayerPos.X,
			Y: simcfg.InitialPlayerPos.Y,
			R: simcfg.InitialPlayerPos.R,
		},
	}

	for i, target := range simcfg.Targets {
		resist := make(map[string]float64)
		for k, v := range target.Resist {
			resist[k.String()] = v
		}

		out.TargetDetails[i] = &model.Enemy{
			Level:  int32(target.Level),
			HP:     target.HP,
			Resist: resist,
			Pos: &model.Coord{
				X: target.Pos.X,
				Y: target.Pos.Y,
				R: target.Pos.R,
			},
			ParticleDropThreshold: target.ParticleDropThreshold,
			ParticleDropCount:     target.ParticleDropCount,
			ParticleElement:       target.ParticleElement.String(),
		}
	}

	if simcfg.Settings.DamageMode {
		out.Mode = model.SimMode_TTK_MODE
	}

	charDetails, err := GenerateCharacterDetails(simcfg)
	if err != nil {
		return out, err
	}
	out.CharacterDetails = charDetails

	for i := range simcfg.Characters {
		if !result.IsCharacterComplete(simcfg.Characters[i].Base.Key) {
			out.IncompleteCharacters = append(out.IncompleteCharacters, simcfg.Characters[i].Base.Key.String())
		}
	}

	return out, nil
}

// cryptoRandSeed generates a random seed using crypo rand
func CryptoRandSeed() int64 {
	var b [8]byte
	_, err := rand.Read(b[:])
	if err != nil {
		log.Panic("cannot seed math/rand package with cryptographically secure random number generator")
	}
	return int64(binary.LittleEndian.Uint64(b[:]))
}

var reImport = regexp.MustCompile(`(?m)^import "(.+)"$`)

// readConfig will load and read the config at specified path. Will resolve any import statements
// as well
func ReadConfig(fpath string) (string, error) {
	src, err := os.ReadFile(fpath)
	if err != nil {
		return "", err
	}

	// check for imports
	var data strings.Builder

	rows := strings.Split(strings.ReplaceAll(string(src), "\r\n", "\n"), "\n")
	for _, row := range rows {
		match := reImport.FindStringSubmatch(row)
		if match != nil {
			// read import
			p := path.Join(path.Dir(fpath), match[1])
			src, err = os.ReadFile(p)
			if err != nil {
				return "", err
			}

			data.Write(src)
		} else {
			data.WriteString(row)
			data.WriteString("\n")
		}
	}

	return data.String(), nil
}

package simulator

import (
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/genshinsim/gcsim/pkg/agg"
	"github.com/genshinsim/gcsim/pkg/gcs/ast"
	"github.com/genshinsim/gcsim/pkg/result"
	"github.com/genshinsim/gcsim/pkg/stats"
	"github.com/genshinsim/gcsim/pkg/worker"
)

// Options sets out the settings to run the sim by (such as debug mode, etc..)
type Options struct {
	ResultSaveToPath string // file name (excluding ext) to save the result file; if "" then nothing is saved to file
	GZIPResult       bool   // should the result file be gzipped; only if ResultSaveToPath is not ""
	ConfigPath       string // path to the config file to read
	Version          string
	BuildDate        string
	DebugMinMax      bool // whether to additional include debug logs for min/max-DPS runs
}

var start time.Time

// Run will run the simulation given number of times
func Run(opts Options) (result.Summary, error) {
	start = time.Now()

	cfg, err := ReadConfig(opts.ConfigPath)
	if err != nil {
		return result.Summary{}, err
	}
	parser := ast.New(cfg)
	simcfg, err := parser.Parse()
	if err != nil {
		return result.Summary{}, err
	}
	//check other errors as well
	if len(simcfg.Errors) != 0 {
		fmt.Println("The config has the following errors: ")
		for _, v := range simcfg.Errors {
			fmt.Printf("\t%v\n", v)
		}
		return result.Summary{}, errors.New("sim has errors")
	}
	return RunWithConfig(cfg, simcfg, opts)
}

// Runs the simulation with a given parsed config
func RunWithConfig(cfg string, simcfg *ast.ActionList, opts Options) (result.Summary, error) {
	// initialize aggregators
	var aggregators []agg.Aggregator
	for _, aggregator := range agg.Aggregators() {
		a, err := aggregator(simcfg)
		if err != nil {
			return result.Summary{}, err
		}
		aggregators = append(aggregators, a)
	}

	//set up a pool
	respCh := make(chan stats.Result)
	errCh := make(chan error)
	pool := worker.New(simcfg.Settings.NumberOfWorkers, respCh, errCh)
	pool.StopCh = make(chan bool)

	//spin off a go func that will queue jobs for as long as the total queued < iter
	//this should block as queue gets full
	go func() {
		//make all the seeds
		wip := 0
		for wip < simcfg.Settings.Iterations {
			pool.QueueCh <- worker.Job{
				Cfg:  simcfg.Copy(),
				Seed: CryptoRandSeed(),
			}
			wip++
		}
	}()

	defer close(pool.StopCh)

	//start reading respCh, queueing a new job until wip == number of iterations
	count := 0
	for count < simcfg.Settings.Iterations {
		select {
		case result := <-respCh:
			for _, a := range aggregators {
				a.Add(result, count)
			}
			count += 1
		case err := <-errCh:
			//error encountered
			return result.Summary{}, err
		}
	}

	// generate final agg results
	stats := &agg.Result{}
	for _, a := range aggregators {
		a.Flush(stats)
	}

	result, err := GenerateResult(cfg, simcfg, stats, opts)
	if err != nil {
		return result, err
	}

	//TODO: clean up this code
	if opts.ResultSaveToPath != "" {
		err = result.Save(opts.ResultSaveToPath, opts.GZIPResult)
		if err != nil {
			return result, err
		}
	}

	return result, nil
}

func GenerateResult(cfg string, simcfg *ast.ActionList, stats *agg.Result, opts Options) (result.Summary, error) {
	result := result.Summary{
		V2:            true,
		Version:       opts.Version,
		BuildDate:     opts.BuildDate,
		IsDamageMode:  simcfg.Settings.DamageMode,
		ActiveChar:    simcfg.InitialChar.String(),
		Iterations:    simcfg.Settings.Iterations,
		Runtime:       float64(time.Since(start).Nanoseconds()),
		NumTargets:    len(simcfg.Targets),
		TargetDetails: simcfg.Targets,
		Config:        cfg,
	}
	result.Map(simcfg, stats)
	result.Text = result.PrettyPrint()

	charDetails, err := GenerateCharacterDetails(simcfg)
	if err != nil {
		return result, err
	}
	result.CharDetails = charDetails

	//run one debug
	//debug call will clone before running
	debugOut, err := GenerateDebugLogWithSeed(simcfg, CryptoRandSeed())
	if err != nil {
		return result, err
	}
	result.Debug = debugOut

	// Include debug logs for min/max-DPS runs if requested.
	if opts.DebugMinMax {
		minDPSDebugOut, err := GenerateDebugLogWithSeed(simcfg, int64(result.MinSeed))
		if err != nil {
			return result, err
		}
		result.DebugMinDPSRun = minDPSDebugOut

		maxDPSDebugOut, err := GenerateDebugLogWithSeed(simcfg, int64(result.MaxSeed))
		if err != nil {
			return result, err
		}
		result.DebugMaxDPSRun = maxDPSDebugOut
	}
	return result, nil
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

	src, err := ioutil.ReadFile(fpath)
	if err != nil {
		return "", err
	}

	//check for imports
	var data strings.Builder

	rows := strings.Split(strings.ReplaceAll(string(src), "\r\n", "\n"), "\n")
	for _, row := range rows {
		match := reImport.FindStringSubmatch(row)
		if match != nil {
			//read import
			p := path.Join(path.Dir(fpath), match[1])
			src, err = ioutil.ReadFile(p)
			if err != nil {
				return "", err
			}

			data.WriteString(string(src))

		} else {
			data.WriteString(row)
			data.WriteString("\n")
		}
	}

	return data.String(), nil
}

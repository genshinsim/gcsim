//simhelper provides the methods required to run simulations; the cmd line tools should be a wrapper
//around this
package simulator

import (
	"crypto/rand"
	"encoding/binary"
	"io/ioutil"
	"log"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/parse"
	"github.com/genshinsim/gcsim/pkg/result"
	"github.com/genshinsim/gcsim/pkg/simulation"
	"github.com/genshinsim/gcsim/pkg/worker"
)

//Options sets out the settings to run the sim by (such as debug mode, etc..)
type Options struct {
	ResultSaveToPath string // file name (excluding ext) to save the result file; if "" then nothing is saved to file
	GZIPResult       bool   // should the result file be gzipped; only if ResultSaveToPath is not ""
	ConfigPath       string // path to the config file to read
}

var start time.Time

//Run will run the simulation given number of times
func Run(opts Options) (result.Summary, error) {
	start = time.Now()

	cfg, err := ReadConfig(opts.ConfigPath)
	if err != nil {
		return result.Summary{}, err
	}
	parser := parse.New("single", cfg)
	simcfg, err := parser.Parse()
	if err != nil {
		return result.Summary{}, err
	}
	return RunWithConfig(cfg, simcfg, opts)
}

// Runs the simulation with a given parsed config
func RunWithConfig(cfg string, simcfg core.SimulationConfig, opts Options) (result.Summary, error) {
	//set up a pool
	respCh := make(chan simulation.Result)
	errCh := make(chan error)
	pool := worker.New(simcfg.Settings.NumberOfWorkers, respCh, errCh)

	//spin off a go func that will queue jobs for as long as the total queued < iter
	//this should block as queue gets full
	go func() {
		//make all the seeds
		wip := 0
		for wip < simcfg.Settings.Iterations {
			pool.QueueCh <- worker.Job{
				Cfg:  simcfg.Clone(),
				Seed: cryptoRandSeed(),
			}
			wip++
		}
	}()

	defer close(pool.StopCh)

	var results []simulation.Result

	//start reading respCh, queueing a new job until wip == number of iterations
	count := simcfg.Settings.Iterations
	for count > 0 {
		select {
		case r := <-respCh:
			results = append(results, r)
			count--
		case err := <-errCh:
			//error encountered
			return result.Summary{}, err
		}

	}

	r := aggregateResults(results, simcfg)

	//run one debug
	debugOut, err := GenerateDebugLog(simcfg.Clone())
	if err != nil {
		return r, err
	}
	r.Debug = debugOut
	r.Runtime = time.Since(start)
	r.Config = cfg
	//TODO: clean up this code

	if opts.ResultSaveToPath != "" {
		r.Save(opts.ResultSaveToPath, opts.GZIPResult)
	}

	//all done
	return r, nil
}

func aggregateResults(in []simulation.Result, cfg core.SimulationConfig) result.Summary {
	//aggregate results
	chars := make([]string, len(cfg.Characters.Profile))
	for i, v := range cfg.Characters.Profile {
		chars[i] = v.Base.Key.String()
	}

	r := result.CollectResult(
		in,
		cfg.DamageMode,
		chars,
		true,
		false,
	)

	r.Iterations = cfg.Settings.Iterations
	r.ActiveChar = cfg.Characters.Initial.String()
	if cfg.DamageMode {
		r.Duration.Mean = float64(cfg.Settings.Duration)
		r.Duration.Min = float64(cfg.Settings.Duration)
		r.Duration.Max = float64(cfg.Settings.Duration)
	}

	r.NumTargets = len(cfg.Targets)
	r.CharDetails = in[0].CharDetails
	for i := range r.CharDetails {
		r.CharDetails[i].Stats = cfg.Characters.Profile[i].Stats
	}
	r.TargetDetails = cfg.Targets

	return r
}

//cryptoRandSeed generates a random seed using crypo rand
func cryptoRandSeed() int64 {
	var b [8]byte
	_, err := rand.Read(b[:])
	if err != nil {
		log.Panic("cannot seed math/rand package with cryptographically secure random number generator")
	}
	return int64(binary.LittleEndian.Uint64(b[:]))
}

var reImport = regexp.MustCompile(`(?m)^import "(.+)"$`)

//readConfig will load and read the config at specified path. Will resolve any import statements
//as well
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

package main

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"runtime/pprof"
	"strconv"
	"strings"

	"github.com/genshinsim/gcsim/pkg/model"
	"github.com/genshinsim/gcsim/pkg/optimization"
	"github.com/genshinsim/gcsim/pkg/simulator"
)

var (
	shareKey string
)

type opts struct {
	config           string
	out              string // file result name
	sample           string // file sample name
	sampleMinDps     string // file sample name for the min-DPS run
	sampleMaxDps     string // file sample name for the max-DPS run
	gz               bool
	serve            bool
	nobrowser        bool
	norun            bool
	keepserving      bool
	substatOptim     bool
	substatOptimFull bool
	verbose          bool
	options          string
	cpuprofile       string
	memprofile       string
}

const resultServeFile = "serve_data.json"
const sampleServeFile = "serve_sample.json"
const address = ":8381"

// command line tool; following options are available:
func main() {
	if err := mainImpl(); err != nil {
		log.Fatal(err)
	}
}

func mainImpl() error {
	var opt opts
	var version bool
	flag.BoolVar(&version, "version", false, "check gcsim version (git hash)")
	flag.StringVar(&opt.config, "c", "config.txt", "which profile to use")
	flag.StringVar(&opt.out, "out", "", "output result to file? supply file path (otherwise empty string for disabled). default disabled")
	flag.StringVar(&opt.sample, "sample", "", "create sample result. supply file path (otherwise empty string for disabled). default disabled")
	flag.StringVar(&opt.sampleMinDps, "sampleMinDps", "", "create sample result for the min-DPS run. supply file path (otherwise empty string for disabled). default disabled")
	flag.StringVar(&opt.sampleMaxDps, "sampleMaxDps", "", "create sample result for the max-DPS run. supply file path (otherwise empty string for disabled). default disabled")
	flag.BoolVar(&opt.gz, "gz", false, "gzip json results; require out flag")
	flag.BoolVar(&opt.serve, "s", false, "serve results to viewer (local). default false")
	flag.BoolVar(&opt.norun, "nr", false, "disable running the simulation (useful if you only want to generate a sample")
	flag.BoolVar(&opt.nobrowser, "nb", false, "disable opening default browser")
	flag.BoolVar(&opt.keepserving, "ks", false, "keep serving same results without terminating web server")
	flag.BoolVar(&opt.substatOptim, "substatOptim", false, "Optimize substats according to KQM standards. Set the out flag to output config with optimal substats inserted to a given file path. Alternatively use the substatOptimFull flag to avoid a second config file and second invocation of the sim.")
	flag.BoolVar(&opt.substatOptimFull, "substatOptimFull", false, "Optimize substats according to KQM standards, overwrite the given config with the optimized version and then run the sim on it. Set the out flag and gz flag to save the viewer file. substatOptim flag takes precedence over this flag, so do not use them together.")
	flag.BoolVar(&opt.verbose, "v", false, "Verbose output log (currently only for substat optimization)")
	flag.StringVar(&opt.options, "options", "", `Additional options for substat optimization mode. Currently supports the following flags, set in a semi-colon delimited list (e.g. -options="total_liquid_substats=15;indiv_liquid_cap=8"):
- total_liquid_substats (default = 20): Total liquid substats available to be assigned across all substats
- indiv_liquid_cap (default = 10): Total liquid substats that can be assigned to a single substat
- fixed_substats_count (default = 2): Amount of fixed substats that are assigned to all substats
- fine_tune (default = 1): Set to 0 to disable fine tune step of substat optimizer.`)
	flag.StringVar(&opt.cpuprofile, "cpuprofile", "", `write cpu profile to a file. supply file path (otherwise empty string for disabled). 
can be viewed in the browser via "go tool pprof -http=localhost:3000 cpu.prof" (insert your desired host/port/filename, requires Graphviz)`)
	flag.StringVar(&opt.memprofile, "memprofile", "", `write memory profile to a file. supply file path (otherwise empty string for disabled). 
can be viewed in the browser via "go tool pprof -http=localhost:3000 mem.prof" (insert your desired host/port/filename, requires Graphviz)`)

	flag.Parse()

	_, err := os.Stat(opt.config)
	usedCLI := false
	flag.Visit(func(f *flag.Flag) {
		usedCLI = true
	})
	if errors.Is(err, os.ErrNotExist) && !usedCLI {
		fmt.Printf("The file %s does not exist.\n", opt.config)
		fmt.Println("What is the filepath of the config you would like to run?")
		in := bufio.NewReader(os.Stdin)
		line, _ := in.ReadString('\n')
		opt.config = strings.TrimSpace(line)
		opt.serve = true
	}

	if opt.cpuprofile != "" {
		f, err := os.Create(opt.cpuprofile)
		if err != nil {
			return fmt.Errorf("could not create CPU profile: %w", err)
		}
		defer f.Close() // error handling omitted for example
		if err := pprof.StartCPUProfile(f); err != nil {
			return fmt.Errorf("could not start CPU profile: %w", err)
		}
		defer pprof.StopCPUProfile()
	}

	if version {
		fmt.Println(simulator.Version())
		return nil
	}

	if shareKey == "" {
		shareKey = os.Getenv("GCSIM_SHARE_KEY")
	}

	if opt.serve {
		opt.out = resultServeFile
		opt.sample = sampleServeFile
		opt.gz = true
	}

	simopt := simulator.Options{
		ConfigPath:       opt.config,
		ResultSaveToPath: opt.out,
		GZIPResult:       opt.gz,
	}

	if opt.substatOptim {
		// TODO: Eventually will want to handle verbose/options in some other way.
		// Ideally once documentation is standardized, can move options to a config file, and verbose can also be moved into options or something
		optimization.RunSubstatOptim(simopt, opt.verbose, opt.options)
		return nil
	}

	if opt.substatOptimFull {
		// set output path to input config file so it gets overwritten during substat optimizer run
		simopt.ResultSaveToPath = simopt.ConfigPath
		// run substat optimizer on given config and output optimized config to the same location
		optimization.RunSubstatOptim(simopt, opt.verbose, opt.options)
		// set output path back to given out flag for sim results
		simopt.ResultSaveToPath = opt.out
	}

	// TODO: should perform the config parsing here and then share the parsed results between run & sample
	var res *model.SimulationResult
	var hash string

	if !opt.norun {
		res, err = simulator.Run(context.Background(), simopt)
		if err != nil {
			return err
		}
		hash, _ = res.Sign(shareKey)
		fmt.Println(res.PrettyPrint())

		if simopt.ResultSaveToPath != "" {
			err = res.Save(simopt.ResultSaveToPath, simopt.GZIPResult)
			if err != nil {
				return err
			}
		}
	}

	if opt.sample != "" {
		var err error
		if opt.norun {
			err = writeSample(
				uint64(simulator.CryptoRandSeed()),
				opt.sample,
				opt.config,
				opt.gz,
				simopt,
			)
		} else {
			err = parseStrSeedAndWriteSample(
				res.SampleSeed,
				opt.sample,
				opt.config,
				opt.gz,
				simopt,
			)
		}

		if err != nil {
			return err
		}
	}

	if opt.sampleMinDps != "" {
		err := parseStrSeedAndWriteSample(
			res.Statistics.MinSeed,
			opt.sampleMinDps,
			opt.config,
			opt.gz,
			simopt,
		)

		if err != nil {
			return err
		}
	}

	if opt.sampleMaxDps != "" {
		err := parseStrSeedAndWriteSample(
			res.Statistics.MaxSeed,
			opt.sampleMaxDps,
			opt.config,
			opt.gz,
			simopt,
		)

		if err != nil {
			return err
		}
	}

	if opt.serve && !opt.norun {
		fmt.Println("Serving results & sample to HTTP...")
		idleConnectionsClosed := make(chan struct{})
		serve(idleConnectionsClosed, resultServeFile+".gz", hash, sampleServeFile+".gz", opt.keepserving)

		openBrowser := func() {
			url := "https://gcsim.app/local"
			if opt.nobrowser {
				return
			}

			err := open(url)
			if err == nil {
				return
			}

			// try "xdg-open-wsl"
			err = openWSL(url)
			if err == nil {
				return
			}
			fmt.Printf("Error opening default browser... please visit: %v\n", url)
		}
		openBrowser()

		<-idleConnectionsClosed
	}

	if opt.memprofile != "" {
		f, err := os.Create(fmt.Sprintf(opt.memprofile))
		if err != nil {
			return fmt.Errorf("could not create memory profile: %w", err)
		}
		defer f.Close() // error handling omitted for example
		if err := pprof.WriteHeapProfile(f); err != nil {
			return fmt.Errorf("could not write memory profile: %w", err)
		}
	}

	return nil
}

// open opens the specified URL in the default browser of the user.
func open(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

func openWSL(url string) error {
	cmd := "powershell.exe"
	args := []string{"/c", "start", url}
	return exec.Command(cmd, args...).Start()
}

func parseStrSeedAndWriteSample(seedStr, outputPath, config string, gz bool, simopt simulator.Options) error {
	seed, err := strconv.ParseUint(seedStr, 10, 64)

	if err != nil {
		return err
	}

	return writeSample(seed, outputPath, config, gz, simopt)
}

func writeSample(seed uint64, outputPath, config string, gz bool, simopt simulator.Options) error {
	cfg, err := simulator.ReadConfig(config)
	if err != nil {
		return err
	}

	sample, err := simulator.GenerateSampleWithSeed(cfg, seed, simopt)
	if err != nil {
		return err
	}
	sample.Save(outputPath, gz)
	fmt.Printf("Generated sample with seed %v to %s\n", seed, outputPath)

	return nil
}

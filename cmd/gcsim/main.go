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
	config       string
	out          string //file result name
	sample       string //file sample name
	gz           bool
	serve        bool
	nobrowser    bool
	norun        bool
	keepserving  bool
	substatOptim bool
	verbose      bool
	options      string
	debugMinMax  bool
	cpuprofile   string
	memprofile   string
}

const resultServeFile = "serve_data.json"
const sampleServeFile = "serve_sample.json"
const address = ":8381"

// command line tool; following options are available:
func main() {

	var opt opts
	var version bool
	flag.BoolVar(&version, "version", false, "check gcsim version (git hash)")
	flag.StringVar(&opt.config, "c", "config.txt", "which profile to use")
	flag.StringVar(&opt.out, "out", "", "output result to file? supply file path (otherwise empty string for disabled). default disabled")
	flag.StringVar(&opt.sample, "sample", "", "create sample result. supply file path (otherwise empty string for disabled). default disabled")
	flag.BoolVar(&opt.gz, "gz", false, "gzip json results; require out flag")
	flag.BoolVar(&opt.serve, "s", false, "serve results to viewer (local). default false")
	flag.BoolVar(&opt.norun, "nr", false, "disable running the simulation (useful if you only want to generate a sample")
	flag.BoolVar(&opt.nobrowser, "nb", false, "disable opening default browser")
	flag.BoolVar(&opt.keepserving, "ks", false, "keep serving same results without terminating web server")
	flag.BoolVar(&opt.substatOptim, "substatOptim", false, "optimize substats according to KQM standards. Set the out flag to output config with optimal substats inserted to a given file path")
	flag.BoolVar(&opt.verbose, "v", false, "Verbose output log (currently only for substat optimization)")
	flag.StringVar(&opt.options, "options", "", `Additional options for substat optimization mode. Currently supports the following flags, set in a semi-colon delimited list (e.g. -options="total_liquid_substats=15;indiv_liquid_cap=8"):
- total_liquid_substats (default = 20): Total liquid substats available to be assigned across all substats
- indiv_liquid_cap (default = 10): Total liquid substats that can be assigned to a single substat
- fixed_substats_count (default = 2): Amount of fixed substats that are assigned to all substats
- sim_iter (default = 350): RECOMMENDED TO NOT TOUCH. Number of iterations used when optimizing. Only change (increase) this if you are working with a team with extremely high standard deviation (>25% of mean)
- tol_mean (default = 0.015): RECOMMENDED TO NOT TOUCH. Tolerance of changes in DPS mean used in ER optimization
- tol_sd (default = 0.33): RECOMMENDED TO NOT TOUCH. Tolerance of changes in DPS SD used in ER optimization`)
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
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	if version {
		fmt.Println(simulator.Version())
		return
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
		return
	}

	// TODO: should perform the config parsing here and then share the parsed results between run & sample
	var res *model.SimulationResult
	var hash string

	if !opt.norun {
		res, err = simulator.Run(simopt, context.Background())
		if err != nil {
			log.Println(err)
			return
		}
		hash, _ = res.Sign(shareKey)
		fmt.Println(res.PrettyPrint())
	}

	if opt.sample != "" {
		var seed uint64
		if opt.norun {
			seed = uint64(simulator.CryptoRandSeed())
		} else {
			seed, _ = strconv.ParseUint(res.SampleSeed, 10, 64)
		}

		cfg, err := simulator.ReadConfig(opt.config)
		if err != nil {
			log.Println(err)
			return
		}

		sample, err := simulator.GenerateSampleWithSeed(cfg, seed, simopt)
		if err != nil {
			log.Println(err)
			return
		}
		sample.Save(opt.sample, opt.gz)
		fmt.Printf("Generated sample with seed: %v\n", seed)
	}

	if opt.serve && !opt.norun {
		fmt.Println("Serving results & sample to HTTP...")
		idleConnectionsClosed := make(chan struct{})
		serve(idleConnectionsClosed, resultServeFile+".gz", hash, sampleServeFile+".gz", opt.keepserving)

		url := "https://gcsim.app/local"
		if !opt.nobrowser {
			err := open(url)
			if err != nil {
				//try "xdg-open-wsl"
				err = openWSL(url)
				if err != nil {
					fmt.Printf("Error opening default browser... please visit: %v\n", url)
				}
			}
		}

		<-idleConnectionsClosed
	}

	if opt.memprofile != "" {
		f, err := os.Create(fmt.Sprintf(opt.memprofile))
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		runtime.GC()    // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
	}
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

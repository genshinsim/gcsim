package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"

	"github.com/genshinsim/gcsim/pkg/optimization"
	"github.com/genshinsim/gcsim/pkg/result"
	"github.com/genshinsim/gcsim/pkg/sample"
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
}

const resultServeFile = "serve_data.json"
const sampleServeFile = "serve_sample.json"
const address = ":8381"

// command line tool; following options are available:
func main() {

	var opt opts
	var version bool
	flag.BoolVar(&version, "version", false, "check gcsim version (git hash)")
	flag.StringVar(&opt.config, "c", "config.txt", "which profile to use; default config.txt")
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

	flag.Parse()

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
	var res result.Summary
	var err error

	if !opt.norun {
		res, err = simulator.Run(simopt, context.Background())
		if err != nil {
			log.Println(err)
			return
		}
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

		sample, err := sample.GenerateSampleWithSeed(cfg, seed, simopt)
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
		serve(idleConnectionsClosed, resultServeFile+".gz", sampleServeFile+".gz", opt.keepserving)

		url := "https://gcsim.app/viewer/local"
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

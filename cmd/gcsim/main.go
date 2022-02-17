package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/genshinsim/gcsim/internal/simulator"
	"github.com/genshinsim/gcsim/internal/substatoptimizer"
	"github.com/pkg/profile"
)

type opts struct {
	config       string
	out          string //file result name
	gz           bool
	prof         bool
	substatOptim bool
	verbose      bool
	options      string
}

//command line tool; following options are available:
func main() {

	var opt opts
	flag.StringVar(&opt.config, "c", "config.txt", "which profile to use; default config.txt")
	flag.StringVar(&opt.out, "out", "", "output result to file? supply file path (otherwise empty string for disabled). default disabled")
	flag.BoolVar(&opt.gz, "gz", false, "gzip json results; require out flag")
	flag.BoolVar(&opt.prof, "p", false, "run cpu profile; default false")
	flag.BoolVar(&opt.substatOptim, "substatOptim", false, "optimize substats according to KQM standards. Set the out flag to output config with optimal substats inserted to a given file path")
	flag.BoolVar(&opt.verbose, "v", false, "Verbose output log (currently only for substat optimization)")
	flag.StringVar(&opt.options, "options", "", `Additional options for substat optimization mode. Currently supports the following flags, set in a semi-colon delimited list (e.g. -options="total_liquid_substats=15;indiv_liquid_cap=8"):
- total_liquid_substats (default = 20): Total liquid substats available to be assigned across all substats
- indiv_liquid_cap (default = 10): Total liquid substats that can be assigned to a single substat
- fixed_substats_count (default = 2): Amount of fixed substats that are assigned to all substats
- sim_iter (default = 350): RECOMMENDED TO NOT TOUCH. Number of iterations used when optimizing. Only change (increase) this if you are working with a team with extremely high standard deviation (>25% of mean)
- tol_mean (default = 0.015): RECOMMENDED TO NOT TOUCH. Tolerance of changes in DPS mean used in ER optimization
- tol_sd (default = 0.33): RECOMMENDED TO NOT TOUCH. Tolerance of changes in DPS SD used in ER optimization`)

	if opt.prof {
		defer profile.Start(profile.ProfilePath("./"), profile.CPUProfile).Stop()
	}

	// defer profile.Start(profile.ProfilePath("./mem.pprof"), profile.MemProfileHeap).Stop()

	simopt := simulator.Options{
		ConfigPath:       opt.config,
		ResultSaveToPath: opt.out,
		GZIPResult:       opt.gz,
	}

	if opt.substatOptim {
		// TODO: Eventually will want to handle verbose/options in some other way.
		// Ideally once documentation is standardized, can move options to a config file, and verbose can also be moved into options or something
		substatoptimizer.RunSubstatOptim(simopt, opt.verbose, opt.options)
		return
	}

	res, err := simulator.Run(simopt)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(res.PrettyPrint())
}

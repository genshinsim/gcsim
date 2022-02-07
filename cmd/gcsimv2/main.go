package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/genshinsim/gcsim/internal/simulator"
	"github.com/pkg/profile"
)

type opts struct {
	config string
	out    string //file result name
	gz     bool
	prof   bool
}

//command line tool; following options are available:
func main() {

	var opt opts
	flag.StringVar(&opt.config, "c", "config.txt", "which profile to use; default config.txt")
	flag.StringVar(&opt.out, "out", "", "output result to file? supply file path (otherwise empty string for disabled). default disabled")
	flag.BoolVar(&opt.gz, "gz", false, "gzip json results; require out flag")
	flag.BoolVar(&opt.prof, "p", false, "run cpu profile; default false")

	if opt.prof {
		defer profile.Start(profile.ProfilePath("./"), profile.CPUProfile).Stop()
	}

	// defer profile.Start(profile.ProfilePath("./mem.pprof"), profile.MemProfileHeap).Stop()

	simopt := simulator.Options{
		ConfigPath:       opt.config,
		ResultSaveToPath: opt.out,
		GZIPResult:       opt.gz,
	}
	res, err := simulator.Run(simopt)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(res.PrettyPrint())
}

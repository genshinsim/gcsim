package main

import (
	"fmt"
	"log"

	"github.com/genshinsim/gcsim/internal/simulator"
	"github.com/pkg/profile"
)

func main() {

	defer profile.Start(profile.ProfilePath("./"), profile.CPUProfile).Stop()

	// defer profile.Start(profile.ProfilePath("./mem.pprof"), profile.MemProfileHeap).Stop()

	opts := simulator.Options{
		PrintResultSummaryToScreen: true,
		ConfigPath:                 "./config.txt",
	}
	res, err := simulator.Run(opts)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(res.PrettyPrint())

}

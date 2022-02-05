package main

import (
	"fmt"
	"log"
	"time"

	"github.com/genshinsim/gcsim/internal/simulator"
)

func main() {
	start := time.Now()
	opts := simulator.Options{
		PrintResultSummaryToScreen: true,
		ConfigPath:                 "./config.txt",
	}
	res, err := simulator.Run(opts)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(res.PrettyPrint())
	fmt.Println(time.Since(start))
}

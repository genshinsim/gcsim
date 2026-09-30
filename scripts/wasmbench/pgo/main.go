// Command pgo writes a CPU profile of the per-iteration work the web UI's wasm does, for
// cmd/wasm/default.pgo. It runs natively, since Go can't profile js/wasm, and follows
// cmd/wasm's simulate() and aggregate(): build a core and evaluator from the parsed config, run
// the sim, reduce the result to an agg.Summary, round-trip it through msgpack and add it to the
// aggregators. The GC settings match cmd/wasm/gc.go.
//
// Run it with GOMAXPROCS=1 (see ../pgo.sh): js/wasm has one thread, and with more Ps the
// profile fills up with cross-thread goroutine handoffs and background GC that wasm never does.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"runtime/pprof"
	"slices"
	"time"

	"github.com/genshinsim/gcsim/pkg/agg"
	"github.com/genshinsim/gcsim/pkg/gcs/ast"
	"github.com/genshinsim/gcsim/pkg/gcs/eval"
	"github.com/genshinsim/gcsim/pkg/gcs/parser"
	"github.com/genshinsim/gcsim/pkg/model"
	"github.com/genshinsim/gcsim/pkg/simulation"
	_ "github.com/genshinsim/gcsim/pkg/simulator" // registers characters, weapons and sets like cmd/wasm
)

func main() {
	out := flag.String("cpuprofile", "", "write the CPU profile here")
	secs := flag.Float64("secs", 8, "seconds to run each config")
	flag.Parse()
	if *out == "" || flag.NArg() == 0 {
		log.Fatal("usage: pgo -cpuprofile out.pprof [-secs 8] config.txt...")
	}

	debug.SetGCPercent(400)
	debug.SetMemoryLimit(512 << 20)

	f, err := os.Create(*out)
	if err != nil {
		log.Fatal(err)
	}
	if err := pprof.StartCPUProfile(f); err != nil {
		log.Fatal(err)
	}
	for _, path := range flag.Args() {
		if err := run(path, time.Duration(*secs*float64(time.Second))); err != nil {
			log.Fatalf("%s: %v", path, err)
		}
	}
	pprof.StopCPUProfile()
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

func run(path string, d time.Duration) error {
	raw, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	simcfg, gcsl, err := parser.New(ast.NewFile(), string(raw)).Parse()
	if err != nil {
		return err
	}
	var aggregators []agg.Aggregator
	for _, a := range agg.Aggregators() {
		enabled := simcfg.Settings.CollectStats
		if len(enabled) > 0 && !slices.Contains(enabled, a.Name) {
			continue
		}
		x, err := a.New(simcfg)
		if err != nil {
			return err
		}
		aggregators = append(aggregators, x)
	}

	buffer := make([]byte, 0, 10*1024)
	start := time.Now()
	n := 0
	for ; time.Since(start) < d; n++ {
		cpycfg := simcfg.Copy()
		core, err := simulation.NewCore(int64(n+1), false, cpycfg)
		if err != nil {
			return err
		}
		ev, err := eval.NewEvaluator(ast.NewFile(), gcsl.Copy(), core)
		if err != nil {
			return err
		}
		sim, err := simulation.New(cpycfg, ev, core)
		if err != nil {
			return err
		}
		result, err := sim.Run()
		if err != nil {
			return err
		}
		summary := agg.Summarize(&result)
		if buffer, err = summary.MarshalMsg(buffer[:0]); err != nil {
			return err
		}

		var received agg.Summary
		if _, err := received.UnmarshalMsg(buffer); err != nil {
			return err
		}
		for _, a := range aggregators {
			a.Add(&received)
		}
	}
	stats := &model.SimulationStatistics{}
	for _, a := range aggregators {
		a.Flush(stats)
	}
	fmt.Printf("%s: %d iterations, %.3f ms/iter\n", path, n, float64(time.Since(start).Microseconds())/1000/float64(n))
	return nil
}

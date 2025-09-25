// dbsnapshot pulls all config from db and reruns each config 10x, creating
// a snapshot of the dev debug seeds, and save to a file. Then allows for
// comparison of an older stored snapshot to see if any changes
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/genshinsim/gcsim/pkg/model"
	"github.com/genshinsim/gcsim/pkg/simulator"
	"github.com/schollz/progressbar/v3"
)

type opts struct {
	iters       int
	saveTo      string
	compareFrom string
}

var workers = 5

func main() {
	var opt opts
	flag.IntVar(&workers, "workers", 5, "number of workers to use")
	flag.IntVar(&opt.iters, "iters", 10, "number of iterations to run per config")
	flag.StringVar(&opt.saveTo, "save", "", "file to save snapshot to; if blank will use time stamp")
	flag.StringVar(&opt.compareFrom, "compare", "", "file to compare snapshot from")
	flag.Parse()

	if opt.saveTo == "" {
		opt.saveTo = time.Now().Format("2006-01-02-15-04-05") + ".ss"
	}

	if opt.compareFrom == "" {
		log.Printf("creating snapshot to %v; grabbing from db", opt.saveTo)
		err := createSnapshot(opt.saveTo, opt.iters)
		if err != nil {
			panic(err)
		}
		log.Printf("saved snapshot to %v", opt.saveTo)
		return
	}

	log.Printf("comparing from %v", opt.compareFrom)
	err := compareFromSnapshot(opt.compareFrom, opt.saveTo)
	if err != nil {
		panic(err)
	}
}

func compareFromSnapshot(from, saveTo string) error {
	prev, err := load(from)
	if err != nil {
		return err
	}
	current := &snapshot{}
	var warnings []string
	bar := progressbar.Default(int64(len(prev.results)), "running sims")
	for i, v := range prev.results {
		id := prev.ids[i]
		iters := len(v.Statistics.DevDebug.SeededDps)
		seeds := make([]int64, 0, iters)
		seedRes := make(map[string]float64)

		for _, v := range v.Statistics.DevDebug.SeededDps {
			seedRes[v.Seed] = v.Dps
			s, err := strconv.ParseUint(v.Seed, 10, 64)
			if err != nil {
				return fmt.Errorf("error parsing seed %v for id %v: %w", v.Seed, id, err)
			}
			seeds = append(seeds, int64(s))
		}

		res, err := runSeededSim(context.Background(), v.Config, iters, seeds)
		if err != nil {
			return err
		}
		bar.Add(1)

		current.results = append(current.results, res)
		current.ids = append(current.ids, id)

		// compare result from each seed
		for _, v := range res.Statistics.DevDebug.SeededDps {
			orig, ok := seedRes[v.Seed]
			if !ok {
				warnings = append(warnings, fmt.Sprintf("id %v has unexpected seed %v not found in original result", id, v.Seed))
				continue
			}
			diff := v.Dps - orig
			if diff > 0.01 || diff < -0.01 {
				warnings = append(warnings, fmt.Sprintf("id %v seed %v has different dps: original %v, new %v, diff %v", id, v.Seed, orig, v.Dps, diff))
			}
		}
	}
	if len(warnings) == 0 {
		log.Printf("no differences found")
	} else {
		for _, v := range warnings {
			log.Println(v)
		}
	}
	return current.save(saveTo)
}

func createSnapshot(filename string, iters int) error {
	ctx := context.Background()

	entries, err := getDBEntries()
	if err != nil {
		return err
	}
	log.Printf("fetched %v entries from db", len(entries))
	s := &snapshot{}

	bar := progressbar.Default(int64(len(entries)), "running sims")
	for _, v := range entries {
		simResult, err := runSim(ctx, v.Config, iters)
		if err != nil {
			return err
		}
		bar.Add(1)
		s.results = append(s.results, simResult)
		s.ids = append(s.ids, v.Id)
	}
	return s.save(filename)
}

func runSim(ctx context.Context, config string, iters int) (*model.SimulationResult, error) {
	simcfg, gcsl, err := simulator.Parse(config)
	if err != nil {
		return nil, err
	}
	simcfg.Settings.CollectStats = []string{"overview", "metadata"}
	simcfg.Settings.Iterations = iters
	simcfg.Settings.NumberOfWorkers = workers

	return simulator.RunWithConfig(ctx, config, simcfg, gcsl, simulator.Options{}, time.Now())
}

func runSeededSim(ctx context.Context, config string, iters int, seeds []int64) (*model.SimulationResult, error) {
	simcfg, gcsl, err := simulator.Parse(config)
	if err != nil {
		return nil, err
	}
	simcfg.Settings.CollectStats = []string{"overview", "metadata"}
	simcfg.Settings.Iterations = iters
	simcfg.Settings.NumberOfWorkers = workers

	return simulator.RunWithSeededConfig(ctx, config, simcfg, gcsl, seeds)
}

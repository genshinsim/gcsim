package optstats

import (
	"context"
	"math/rand"
	"slices"

	"github.com/genshinsim/gcsim/pkg/agg"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/gcs/ast"
	"github.com/genshinsim/gcsim/pkg/model"
	"github.com/genshinsim/gcsim/pkg/simulator"
	"github.com/genshinsim/gcsim/pkg/stats"
)

// Runs the simulation with a given parsed config and custom stat collector and aggregator
// TODO: cfg string should be in the action list instead
// TODO: need to add a context here to avoid infinite looping
func RunWithConfigCustomStats[T any](ctx context.Context, cfg string, simcfg *info.ActionList, gcsl ast.Node, opts simulator.Options, seed int64, cstat NewStatsFuncCustomStats[T], cagg func(T)) (*model.SimulationResult, error) {
	// initialize aggregators
	var aggregators []agg.Aggregator
	for _, aggregator := range agg.Aggregators() {
		enabled := simcfg.Settings.CollectStats
		if len(enabled) > 0 && !slices.Contains(enabled, aggregator.Name) {
			continue
		}
		a, err := aggregator.New(simcfg)
		if err != nil {
			return &model.SimulationResult{}, err
		}
		aggregators = append(aggregators, a)
	}

	// set up a pool
	respCh := make(chan stats.Result)
	errCh := make(chan error)
	customCh := make(chan T)
	pool := WorkerNewWithCustomStats(simcfg.Settings.NumberOfWorkers, respCh, errCh, customCh)
	pool.StopCh = make(chan bool)

	// spin off a go func that will queue jobs for as long as the total queued < iter
	// this should block as queue gets full
	go func() {
		src := rand.NewSource(seed)
		// make all the seeds
		wip := 0
		for wip < simcfg.Settings.Iterations {
			pool.QueueCh <- JobCustomStats[T]{
				Cfg:     simcfg.Copy(),
				Actions: gcsl.Copy(),
				Seed:    src.Int63(),
				Cstat:   cstat,
			}
			wip++
		}
	}()

	defer close(pool.StopCh)

	// start reading respCh, queueing a new job until wip == number of iterations
	count := 0
	for count < simcfg.Settings.Iterations {
		select {
		case result := <-customCh:
			cagg(result)
		case result := <-respCh:
			for _, a := range aggregators {
				a.Add(result)
			}
			count += 1
		case err := <-errCh:
			// error encountered
			return &model.SimulationResult{}, err
		case <-ctx.Done():
			return &model.SimulationResult{}, ctx.Err()
		}
	}

	result, err := simulator.GenerateResult(cfg, simcfg)
	if err != nil {
		return result, err
	}

	// generate final agg results
	stats := &model.SimulationStatistics{}
	for _, a := range aggregators {
		a.Flush(stats)
	}
	result.Statistics = stats

	return result, nil
}

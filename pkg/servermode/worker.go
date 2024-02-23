package servermode

import (
	"log/slog"

	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/gcs"
	"github.com/genshinsim/gcsim/pkg/gcs/ast"
	"github.com/genshinsim/gcsim/pkg/model"
	"github.com/genshinsim/gcsim/pkg/simulation"
	"github.com/genshinsim/gcsim/pkg/simulator"
	"github.com/genshinsim/gcsim/pkg/stats"
)

type worker struct {
	// initial stuff
	id     string
	cfg    string
	log    *slog.Logger
	cancel chan bool

	// state
	done   bool                    // if simulation is done
	result *model.SimulationResult // latest result
	err    error                   // any errors running sim
}

type job struct {
	cfg  *info.ActionList
	node ast.Node
	seed int64
}

func (w *worker) handleErr(err error) {
	w.err = err
	w.result = nil
}

func (w *worker) run(workerCount, flushInterval int) {
	w.log.Info("worker run started", "id", w.id)
	// handle panic
	defer func() {
		if r := recover(); r != nil {
			w.handleErr(errorRecover(r))
		}
	}()
	// done is true regardless of reason once this function exits
	defer func() {
		w.done = true
	}()

	simcfg, gcsl, err := parse(w.cfg)
	if err != nil {
		w.log.Info("config parsing failed", "id", w.id, "err", err)
		w.handleErr(err)
		return
	}
	simcfg.Settings.NumberOfWorkers = workerCount
	w.log.Info("parse ok", "id", w.id, "cfg", simcfg)

	aggregators, err := setupAggregators(simcfg)
	if err != nil {
		w.log.Info("aggregator setup failed", "id", w.id, "err", err)
		w.handleErr(err)
		return
	}
	w.log.Info("aggregators ok", "id", w.id)

	// run jobs
	respCh := make(chan stats.Result)
	errCh := make(chan error)
	workChan := make(chan job)
	for i := 0; i < simcfg.Settings.NumberOfWorkers; i++ {
		w.log.Info("spawning worker", "id", w.id, "i", i)
		go w.iter(workChan, respCh, errCh)
	}
	go func() {
		// make all the seeds
		wip := 0
		for wip < simcfg.Settings.Iterations {
			select {
			case <-w.cancel:
				w.log.Info("wip sending ended due to cancel", "id", w.id)
				return
			case workChan <- job{
				cfg:  simcfg.Copy(),
				node: gcsl.Copy(),
				seed: cryptoRandSeed(),
			}:
				wip++
			}
		}
		close(workChan)
	}()

	// setup results
	w.result, err = simulator.GenerateResult(w.cfg, simcfg)
	if err != nil {
		w.log.Info("generate result failed", "id", w.id, "err", err)
		w.handleErr(err)
		return
	}
	w.log.Info("result initialized ok", "id", w.id)

	count := 0
	lastFlush := 0
iters:
	for count < simcfg.Settings.Iterations {
		select {
		case result := <-respCh:
			// w.log.Info("got 1 result", "id", w.id, "count", count)
			for _, a := range aggregators {
				a.Add(result)
			}
			count += 1
		case err := <-errCh:
			// error encountered
			w.log.Info("error running sim", "id", w.id, "err", err)
			w.handleErr(err)
			return
		case <-w.cancel:
			w.log.Info("cancel signal received", "id", w.id)
			// expectation is w.cancel is closed causing all go routines to wrap it up
			break iters
		}
		// flush and update results
		if count-lastFlush > flushInterval {
			w.log.Debug("flushing results", "id", w.id, "count", count, "flush", lastFlush)
			lastFlush = count
			stats := flush(aggregators)
			w.result.Statistics = stats
		}
	}
	w.log.Info("sim done", "id", w.id, "count", count, "flush", lastFlush)
	stats := flush(aggregators)
	w.result.Statistics = stats
}

func (w *worker) iter(work chan job, res chan stats.Result, errChan chan error) {
	for {
		select {
		case <-w.cancel:
			return
		case job, ok := <-work:
			if !ok {
				w.log.Info("work channel closed, iter worker ending", "id", w.id)
				return
			}

			c, err := simulation.NewCore(job.seed, false, job.cfg)
			if err != nil {
				errChan <- err
				return
			}
			eval, err := gcs.NewEvaluator(job.node, c)
			if err != nil {
				errChan <- err
				break
			}
			s, err := simulation.New(job.cfg, eval, c)
			if err != nil {
				errChan <- err
				break
			}
			r, err := s.Run()
			if err != nil {
				errChan <- err
				break
			}
			res <- r
		}
	}
}

package worker

import (
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/gcs"
	"github.com/genshinsim/gcsim/pkg/gcs/ast"
	"github.com/genshinsim/gcsim/pkg/simulation"
	"github.com/genshinsim/gcsim/pkg/stats"
)

type PoolCustomStats[T any] struct {
	respCh   chan stats.Result
	errCh    chan error
	QueueCh  chan JobCustomStats[T]
	customCh chan T
	StopCh   chan bool
}

type JobCustomStats[T any] struct {
	Cfg     *info.ActionList
	Actions ast.Node
	Seed    int64
	Cstat   stats.NewStatsFuncCustomStats[T]
}

// New creates a new Pool. Jobs can be sent to new pool by sending to p.QueueCh
// Closing p.StopCh will cause the pool to stop executing any queued jobs and currently working
// workers will no longer send back responses
func NewWithCustomStats[T any](maxWorker int, respCh chan stats.Result, errCh chan error, customCh chan T) *PoolCustomStats[T] {
	p := &PoolCustomStats[T]{
		respCh:   respCh,
		errCh:    errCh,
		customCh: customCh,
		QueueCh:  make(chan JobCustomStats[T]),
		StopCh:   make(chan bool),
	}
	// create workers
	for i := 0; i < maxWorker; i++ {
		go p.worker()
	}
	return p
}

func (p *PoolCustomStats[T]) worker() {
	for {
		select {
		case job := <-p.QueueCh:
			// fmt.Printf("got job: %s\n", job.Cfg.PrettyPrint())
			c, err := simulation.NewCore(job.Seed, false, job.Cfg)
			if err != nil {
				p.errCh <- err
				break
			}
			t, err := job.Cstat(c)
			if err != nil {
				p.errCh <- err
				break
			}
			eval, err := gcs.NewEvaluator(job.Actions, c)
			if err != nil {
				p.errCh <- err
				break
			}
			s, err := simulation.New(job.Cfg, eval, c)
			if err != nil {
				p.errCh <- err
				break
			}
			res, err := s.Run()
			if err != nil {
				p.errCh <- err
				break
			}
			p.customCh <- t.Flush(c)
			p.respCh <- res

		case _, ok := <-p.StopCh:
			if !ok {
				// stopping
				return
			}
		}
	}
}

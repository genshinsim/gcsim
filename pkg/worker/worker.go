package worker

import (
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/gcs"
	"github.com/genshinsim/gcsim/pkg/gcs/ast"
	"github.com/genshinsim/gcsim/pkg/simulation"
	"github.com/genshinsim/gcsim/pkg/stats"
)

type Pool struct {
	respCh  chan stats.Result
	errCh   chan error
	QueueCh chan Job
	StopCh  chan bool
}

type Job struct {
	Cfg     *info.ActionList
	Actions ast.Node
	Seed    int64
}

// New creates a new Pool. Jobs can be sent to new pool by sending to p.QueueCh
// Closing p.StopCh will cause the pool to stop executing any queued jobs and currently working
// workers will no longer send back responses
func New(maxWorker int, respCh chan stats.Result, errCh chan error) *Pool {

	p := &Pool{
		respCh:  respCh,
		errCh:   errCh,
		QueueCh: make(chan Job),
		StopCh:  make(chan bool),
	}
	//create workers
	for i := 0; i < maxWorker; i++ {
		go p.worker()
	}
	return p
}

func (p *Pool) worker() {
	for {
		select {
		case job := <-p.QueueCh:
			// fmt.Printf("got job: %s\n", job.Cfg.PrettyPrint())
			c, err := simulation.NewCore(job.Seed, false, job.Cfg)
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
			p.respCh <- res

		case _, ok := <-p.StopCh:
			if !ok {
				//stopping
				return
			}
		}

	}
}

package worker

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/simulation"
)

type Pool struct {
	respCh  chan simulation.Result
	errCh   chan error
	QueueCh chan Job
	StopCh  chan bool
}

type Job struct {
	Cfg  core.SimulationConfig
	Seed int64
}

//New creates a new Pool. Jobs can be sent to new pool by sending to p.QueueCh
//Closing p.StopCh will cause the pool to stop executing any queued jobs and currently working
//workers will no longer send back responses
func New(maxWorker int, respCh chan simulation.Result, errCh chan error) *Pool {

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
			c := simulation.NewCore(job.Seed, false, job.Cfg.Settings)
			s, err := simulation.New(job.Cfg, c)
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

		case <-p.StopCh:
			return
		}

	}
}

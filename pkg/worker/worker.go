package worker

import (
	"github.com/genshinsim/gcsim"
	"github.com/genshinsim/gcsim/pkg/core"
)

type Pool struct {
	max     int
	respCh  chan gcsim.Stats
	errCh   chan error
	QueueCh chan Job
	StopCh  chan bool
}

type Job struct {
	Cfg  core.Config
	Opt  core.RunOpt
	Seed int64
	Cust []func(*gcsim.Simulation) error
}

//New creates a new Pool. Jobs can be sent to new pool by sending to p.QueueCh
//Closing p.StopCh will cause the pool to stop executing any queued jobs and currently working
//workers will no longer send back responses
func New(maxWorker int, respCh chan gcsim.Stats, errCh chan error) *Pool {

	p := &Pool{
		max:     maxWorker,
		respCh:  respCh,
		errCh:   errCh,
		QueueCh: make(chan Job, 5),
		StopCh:  make(chan bool),
	}
	go p.run()
	return p
}

func (p *Pool) run() {
	working := 0
	queue := make([]Job, 0, 10)
	done := make(chan bool)
	var j Job
	for {
		select {
		case <-done:
			working--
			if working < 0 {
				working = 0
			}
		case job := <-p.QueueCh:
			queue = append(queue, job)
		case <-p.StopCh:
			//stop all work
			return
		}
		//check if we have any available wokers, if so
		if len(queue) > 0 && working < p.max {
			j, queue = queue[0], queue[1:]
			working++
			go p.worker(j, done)
		}
	}
}

func (p *Pool) worker(job Job, done chan bool) {
	//do stuff
	s, err := gcsim.NewSim(job.Cfg, job.Seed, job.Opt, job.Cust...)

	//make sure we're not supposed to stop
	select {
	case <-p.StopCh:
		return
	default:
	}

	if err != nil {
		p.errCh <- err
		return
	}

	stat, err := s.Run()
	stat.Seed = job.Seed

	//make sure we're not supposed to stop
	select {
	case <-p.StopCh:
		return
	default:
	}

	if err != nil {
		p.errCh <- err
		return
	}

	p.respCh <- stat

}

package simulation

import (
	"context"

	"github.com/genshinsim/gcsim/pkg/gcs"
	"github.com/genshinsim/gcsim/pkg/gcs/ast"
)

func (s *Simulation) Run() (Result, error) {
	//duration
	f := s.cfg.Duration*60 - 1
	stop := false
	var err error

	//setup ast
	s.nextAction = make(chan *ast.ActionStmt)
	s.continueEval = make(chan bool)
	ctx, cancel := context.WithCancel(context.Background())
	s.terminate = cancel
	s.queuer = gcs.Eval{
		AST:  s.cfg.Program,
		Next: s.continueEval,
		Work: s.nextAction,
	}
	go s.queuer.Run(ctx)
	defer s.terminate()

	for !stop {

		err = s.AdvanceFrame()
		if err != nil {
			return s.stats, err
		}

		//TODO: hp mode
		stop = s.C.F == f
	}

	s.stats.Damage = s.C.Combat.TotalDamage
	s.stats.DPS = s.stats.Damage * 60 / float64(s.C.F+1)
	s.stats.Duration = f

	//we're done yay
	return s.stats, nil
}

func (s *Simulation) AdvanceFrame() error {
	var err error
	var done bool
	for !done {
		done, err = s.queueAndExec()
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Simulation) queueAndExec() (bool, error) {
	//check if queue is empty
	if s.queue == nil {
		s.queue = <-s.nextAction
	}

	//we will keep trying to execute this action until it
	//completed successfully
	//TODO: we should do some optimization here to at least skip some frames
	//if we know for sure the action won't be ready

	return true, nil
}

//executes the provided command, returns 2 booleans:
//	- executed if this action was successfully executed and should be purged from queue
//	- frameDone if this action consumes a frame
func (s *Simulation) execCommand(c ast.ActionStmt) (executed, frameDone bool) {

	return
}

func (s *Simulation) QueueNext() {

}

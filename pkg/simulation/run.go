package simulation

import (
	"errors"

	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/gcs"
	"github.com/genshinsim/gcsim/pkg/gcs/ast"
)

func (s *Simulation) Run() (Result, error) {
	//run sim for 90s if no duration set
	if s.cfg.Settings.Duration == 0 {
		// fmt.Println("no duration set, running for 90s")
		s.cfg.Settings.Duration = 90
	}
	//duration
	f := s.cfg.Settings.Duration * 60
	stop := false
	var err error

	//setup ast
	s.nextAction = make(chan *ast.ActionStmt)
	s.continueEval = make(chan bool)
	s.queuer = gcs.Eval{
		AST:  s.cfg.Program,
		Next: s.continueEval,
		Work: s.nextAction,
	}
	go s.queuer.Run()
	defer close(s.continueEval)

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
	s.C.F++
	s.C.Tick()
	s.collectStats()
	err := s.queueAndExec()
	if err != nil {
		return err
	}
	// fmt.Printf("Tick - f = %v\n", s.C.F)
	return nil
}

func (s *Simulation) collectStats() {
	//add char active time
	s.stats.CharActiveTime[s.C.Player.Active()]++
	for i, v := range s.C.Combat.Targets() {
		if t, ok := v.(*enemy.Enemy); ok {
			s.stats.ElementUptime[i][t.AuraType()]++
		}
	}
}

func (s *Simulation) queueAndExec() error {
	//use this to skip some frames as an optimization
	if s.skip > 0 {
		s.skip--
		return nil
	}
	//TODO: this for loops is completely unnecessary
	for {
		if s.queue != nil {
			err := s.C.Player.Exec(s.queue.Action, s.queue.Param)
			switch err {
			case player.ErrActionNotReady:
				//action not ready yet, skipping frame
				//TODO: log something here
				return nil
			case player.ErrPlayerNotReady:
				//player still in animation, skipping frame
				//TODO: log something here
				return nil
			case nil:
				//exeucted successfully
				s.queue = nil
			default:
				//this should now error out
				return err
			}
		}
		//do nothing if no more actions anyways
		if s.noMoreActions {
			//TODO: log here?
			// fmt.Println("no more action")
			return nil
		}
		//check if read to queue first
		if !s.C.Player.CanQueueNextAction() {
			//skip frame if not ready
			return nil
		}
		//check if we can queue an action, if not then skip
		err := s.tryQueueNext()
		switch err {
		case nil:
			//we have an action, continue execute
		case ErrNoMoreActions:
			//make a note no more actions or else <-s.nextAction will block indefinitely
			s.noMoreActions = true
			return nil //do nothing, skip frame
		default:
			//shouldn't really happen??
			return err
		}
	}
}

var ErrNoMoreActions = errors.New("no more actions left")

func (s *Simulation) tryQueueNext() error {
	//tell eval to keep going
	s.continueEval <- true
	//wait for next action
	var ok bool
	s.queue, ok = <-s.nextAction
	if !ok {
		return ErrNoMoreActions
	}
	return nil
}

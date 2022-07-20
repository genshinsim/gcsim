package simulation

import (
	"errors"
	"log"

	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/glog"
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
	f := int(s.cfg.Settings.Duration * 60)
	stop := false
	var err error

	//TODO: enable hp mode?
	s.C.Flags.DamageMode = s.cfg.Settings.DamageMode

	//setup ast
	s.nextAction = make(chan *ast.ActionStmt)
	s.continueEval = make(chan bool)
	s.evalErr = make(chan error)
	s.queuer = gcs.Eval{
		AST:  s.cfg.Program,
		Next: s.continueEval,
		Work: s.nextAction,
		Core: s.C,
		Err:  s.evalErr,
	}
	go s.queuer.Run()
	defer close(s.continueEval)

	//queue up enery tasks
	s.SetupRandEnergyDrop()

	for !stop {
		err = s.AdvanceFrame()
		if err != nil {
			log.Println(err)
			return s.stats, err
		}

		//TODO: hp mode
		stop = s.C.F == f
	}

	s.stats.Seed = s.C.Seed

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
			//handle wait separately
			if s.queue.Action == action.ActionWait {
				//wipe the action here, set skip
				s.skip = s.queue.Param["f"]
				s.C.Log.NewEvent("executed wait", glog.LogActionEvent, s.C.Player.Active()).
					Write("f", s.queue.Param["f"])
				s.queue = nil
				return nil
			} else {
				err := s.C.Player.Exec(s.queue.Action, s.queue.Char, s.queue.Param)
				switch err {
				case player.ErrActionNotReady:
					//action not ready yet, skipping frame
					//TODO: log something here
					return nil
				case player.ErrPlayerNotReady:
					//player still in animation, skipping frame
					//TODO: log something here
					return nil
				case player.ErrActionNoOp:
					//technically the same as nil
					s.C.Log.NewEventBuildMsg(glog.LogActionEvent, s.C.Player.Active(), "noop action: ", s.queue.Action.String())
					s.queue = nil
				case nil:
					//exeucted successfully
					s.queue = nil
				default:
					//this should now error out
					return err
				}
			}
		}
		//do nothing if no more actions anyways
		if s.noMoreActions {
			//TODO: log here?
			// fmt.Println("no more action")
			s.C.Log.NewEvent("no more actions", glog.LogActionEvent, -1)
			return nil
		}
		//check if ready to queue first
		if !s.C.Player.CanQueueNextAction() {
			// s.C.Log.NewEventBuildMsg(glog.LogActionEvent, -1, "action can't be queued yet")
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
			//eval error'd out here
			return err
		}
	}
}

var ErrNoMoreActions = errors.New("no more actions left")

func (s *Simulation) tryQueueNext() error {
	//tell eval to keep going
	s.continueEval <- true
	//eval will either give an action (or keep executing) or error out
	var ok bool
	select {
	case s.queue, ok = <-s.nextAction:
		//wait for next action
		if !ok {
			return ErrNoMoreActions
		}
		return nil
	case err := <-s.evalErr:
		return err
	}
}

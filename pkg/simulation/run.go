package simulation

import (
	"errors"
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/stats"
)

type stateFn func(*Simulation) (stateFn, error)

func (s *Simulation) resFromCurrentState() stats.Result {
	return stats.Result{Seed: uint64(s.C.Seed), Duration: s.C.F + 1}
}

func (s *Simulation) run() (stats.Result, error) {
	// core loop roughly as follows:
	//  - initialize:
	//		- setup
	//		- advance frame by 1
	//		- move to queue phase
	//  - queue phase:
	//		- ask for next action
	//		- move to ready check phase
	//	- ready check phase
	//		- check if action ready (both animation + player); if not ready advance frame until ready
	//		- move to execute action phase
	//	- execute action phase:
	//		- if action has pre-action wait; advance frame until wait is consumed
	//		- execute action and empty queue
	//		- if executed action is no-op, move directly to queue phase
	//		- else advance frame until CanQueueAfter then move to queue phase
	//
	// frame advance will perform the following;
	//	- increment frame counter by 1
	//  - execute any ticks
	//  - check for eneryg procs
	//  - emit OnTick
	//  - perform exit check
	//
	// exit check checks for:
	//	- frame limit
	//  - all enemies dead
	//  - no more actions left

	//TODO: do we need to catch panic here still? or can it be done outside in the worker
	var err error
	for state := initialize; state != nil; {
		state, err = state(s)
		if err != nil {
			return s.resFromCurrentState(), err
		}
	}

	// err = s.eval.Exit()
	// if err != nil {
	// 	fmt.Println("evaluator already closed")
	// 	return s.resFromCurrentState(), err
	// }
	s.eval.Exit()

	err = s.eval.Err()
	if err != nil {
		return s.resFromCurrentState(), err
	}

	return s.gatherResult(), nil
}

func (s *Simulation) gatherResult() stats.Result {
	res := stats.Result{
		Seed:        uint64(s.C.Seed),
		Duration:    s.C.F,
		TotalDamage: s.C.Combat.TotalDamage,
		DPS:         s.C.Combat.TotalDamage * 60 / float64(s.C.F),
		Characters:  make([]stats.CharacterResult, len(s.C.Player.Chars())),
		Enemies:     make([]stats.EnemyResult, s.C.Combat.EnemyCount()),
	}

	for i := range s.cfg.Characters {
		res.Characters[i].Name = s.cfg.Characters[i].Base.Key.String()
	}

	for _, collector := range s.collectors {
		collector.Flush(s.C, &res)
	}

	return res
}

func (s *Simulation) popQueue() int {
	switch len(s.queue) {
	case 0:
	case 1:
		s.queue = s.queue[:0]
	default:
		s.queue = s.queue[1:]
	}
	return len(s.queue)
}

func initialize(s *Simulation) (stateFn, error) {
	go s.eval.Start()
	// run sim for 90s if no duration set
	if s.cfg.Settings.Duration == 0 {
		// fmt.Println("no duration set, running for 90s")
		s.cfg.Settings.Duration = 90
	}
	s.C.Flags.DamageMode = s.cfg.Settings.DamageMode

	return s.advanceFrames(1, queuePhase)
}

func queuePhase(s *Simulation) (stateFn, error) {
	if s.noMoreActions {
		return s.advanceFrames(1, queuePhase)
	}
	s.eval.Continue()
	next, err := s.eval.NextAction()
	if err != nil {
		return nil, err
	}
	// skip a frame and come back to queue phase if eval does not have any more actions
	// relying on advance frame to exit if need be
	if next == nil {
		s.noMoreActions = true
		// we do the same skip here as if eval doesn't have any more ations
		return s.advanceFrames(1, queuePhase)
	}
	// handle sleep here since it's just a frame skip before requeing next
	if next.Action == action.ActionWait {
		return s.handleWait(next)
	}
	// if next action is delay, we can just queue up the action after that right now
	if next.Action == action.ActionDelay {
		// append here because we can have multiple delay chained
		s.preActionDelay += next.Param["f"]
		return queuePhase, nil
	}
	// append swap if called for char is not active
	// check if NoChar incase this is some special action that does not require a character
	if next.Char != keys.NoChar && next.Char != s.C.Player.ActiveChar().Base.Key {
		s.queue = append(s.queue, &action.Eval{
			Char:   next.Char,
			Action: action.ActionSwap,
		})
	}
	s.queue = append(s.queue, next)
	return actionReadyCheckPhase, nil
}

func actionReadyCheckPhase(s *Simulation) (stateFn, error) {
	//TODO: this sanity check is probably not necessary
	if len(s.queue) == 0 {
		return nil, errors.New("unexpected queue length is 0")
	}
	q := s.queue[0]

	//TODO: this loop should be optimized to skip more than 1 frame at a time
	if err := s.C.Player.ReadyCheck(q.Action, q.Char, q.Param); err != nil {
		// repeat this phase until action is ready
		switch {
		case errors.Is(err, player.ErrActionNotReady):
			s.C.Log.NewEvent(fmt.Sprintf("could not execute %v; action not ready", q.Action), glog.LogSimEvent, s.C.Player.Active())
			return s.advanceFrames(1, actionReadyCheckPhase)
		case errors.Is(err, player.ErrPlayerNotReady):
			return s.advanceFrames(1, actionReadyCheckPhase)
		case errors.Is(err, player.ErrActionNoOp):
			// don't do anything here
		default:
			return nil, err
		}
	}

	return executeActionPhase, nil
}

func (s *Simulation) handleWait(q *action.Eval) (stateFn, error) {
	// to maintain existing functionality, wait (alias sleep) is always ready and should cause
	// advanceFrames to be called equal to the param f
	skip := q.Param["f"]
	// log wait(0) differently to make it obvious
	if skip == 0 {
		s.C.Log.NewEvent("executed noop wait(0)", glog.LogActionEvent, s.C.Player.Active()).
			Write("f", skip)
	} else {
		s.C.Log.NewEvent("executed wait", glog.LogActionEvent, s.C.Player.Active()).
			Write("f", skip)
	}
	if l := s.popQueue(); l > 0 {
		// don't go back to queue if there are more actions already queued
		return s.advanceFrames(skip, actionReadyCheckPhase)
	}
	return s.advanceFrames(skip, queuePhase)
}

func executeActionPhase(s *Simulation) (stateFn, error) {
	//TODO: this sanity check is probably not necessary
	if len(s.queue) == 0 {
		return nil, errors.New("unexpected queue length is 0")
	}
	if s.preActionDelay > 0 {
		delay := s.preActionDelay
		s.C.Log.NewEvent("executed pre action delay", glog.LogActionEvent, s.C.Player.Active()).
			Write("f", delay)
		s.preActionDelay = 0
		return s.advanceFrames(delay, executeActionPhase)
	}
	q := s.queue[0]
	err := s.C.Player.Exec(q.Action, q.Char, q.Param)
	if err != nil {
		//TODO: this check probably doesn't do anything
		if errors.Is(err, player.ErrActionNoOp) {
			if l := s.popQueue(); l > 0 {
				// don't go back to queue if there are more actions already queued
				return actionReadyCheckPhase, nil
			}
			return queuePhase, nil
		}
		// this is now unexpected since action should be ready now
		return nil, err
	}
	//TODO: this check here is probably unnecessary
	if l := s.popQueue(); l > 0 {
		// don't go back to queue if there are more actions already queued
		return actionReadyCheckPhase, nil
	}

	return skipUntilCanQueue, nil
}

func skipUntilCanQueue(s *Simulation) (stateFn, error) {
	if !s.C.Player.CanQueueNextAction() {
		return s.advanceFrames(1, skipUntilCanQueue)
	}
	return queuePhase, nil
}

// nextFrame moves up the frame by 1, performing
func (s *Simulation) advanceFrames(f int, next stateFn) (stateFn, error) {
	for i := 0; i < f; i++ {
		done, err := s.nextFrame()
		if err != nil {
			return nil, err
		}
		if done {
			return nil, nil
		}
	}
	return next, nil
}

func (s *Simulation) nextFrame() (bool, error) {
	s.C.F++
	err := s.C.Tick()
	if err != nil {
		return false, err
	}
	s.handleEnergy()
	s.handleHurt()
	s.C.Events.Emit(event.OnTick)
	return s.stopCheck(), nil
}

func (s *Simulation) stopCheck() bool {
	if s.C.Combat.DamageMode {
		// stop if no more actions
		if s.noMoreActions {
			return true
		}
		// stop if all targets are reporting dead
		allDead := true
		for _, t := range s.C.Combat.Enemies() {
			if t.IsAlive() {
				allDead = false
				break
			}
		}
		return allDead
	}
	return s.C.F == int(s.cfg.Settings.Duration*60)
}

// TODO: remove defer in favour of every function actually returning error
//
//nolint:nonamedreturns // not possible to perform the res, err modification without named return
func (s *Simulation) Run() (res stats.Result, err error) {
	defer func() {
		// recover from panic if one occured. Set err to nil otherwise.
		if r := recover(); r != nil {
			res = stats.Result{Seed: uint64(s.C.Seed), Duration: s.C.F + 1}
			err = fmt.Errorf("simulation panic occured: %v", r)
		}
	}()
	res, err = s.run()
	return
}

package simulation

import (
	"errors"
	"fmt"
	"runtime/debug"

	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
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

	// TODO: do we need to catch panic here still? or can it be done outside in the worker
	var err error
	if ev, ok := s.eval.(syncEvaluator); ok {
		err = s.runSync(ev)
	} else if err = s.runStates(initialize); err == nil {
		err = s.exitEval()
	}
	if err != nil {
		return s.resFromCurrentState(), err
	}

	s.C.Events.Emit(event.OnSimEndedSuccessfully)

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
		EndStats:    make([]stats.EndStats, len(s.C.Player.Chars())),
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

// runStates runs the state machine from state until the sim ends or, with a sync evaluator, the
// sim needs the next action
func (s *Simulation) runStates(state stateFn) error {
	var err error
	for state != nil {
		state, err = state(s)
		if err != nil {
			return err
		}
	}
	return nil
}

// exitEval stops the evaluator and returns its error as of now
func (s *Simulation) exitEval() error {
	s.eval.Exit()
	return s.eval.Err()
}

// syncEvaluator runs the program on the calling goroutine, passing each action to exec, which
// returns false to stop the program (see eval.RunSync)
type syncEvaluator interface {
	RunSync(exec func(*action.Eval) bool) error
}

// runSync runs the sim with the evaluator on this goroutine instead of its own: the state machine
// runs until it needs an action, then the program runs until it has one and passes it to exec,
// which runs the state machine again. The two take turns at the same points as over
// Continue/NextAction. It returns the sim's error, or else the evaluator's as of the sim's end.
func (s *Simulation) runSync(ev syncEvaluator) error {
	s.syncEval = true
	if err := s.runStates(initialize); err != nil {
		return err
	}
	if !s.awaitingAction {
		// the sim ended before the program started
		return s.exitEval()
	}
	var err error
	ended := false
	// the program's own errors are read through Err, as without RunSync
	_ = ev.RunSync(func(next *action.Eval) bool {
		err = s.execAction(next)
		if err != nil {
			return false
		}
		if s.awaitingAction {
			return true
		}
		// the sim ended during this action. Read the evaluator's error before the program
		// unwinds, as the goroutine version reads it before the evaluator goroutine resumes
		err = s.exitEval()
		ended = true
		return false
	})
	if err != nil || ended {
		return err
	}
	// the program is done while the sim waits for its next action
	s.awaitingAction = false
	state, err := s.queueAction(nil)
	if err == nil {
		err = s.runStates(state)
	}
	if err != nil {
		return err
	}
	return s.exitEval()
}

// execAction queues next and runs the sim until it needs the next action or ends. A panic is
// returned as the error Run would make of it, since the evaluator frames below would otherwise
// recover it as a program error.
func (s *Simulation) execAction(next *action.Eval) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = panicError(r)
		}
	}()
	s.awaitingAction = false
	state, err := s.queueAction(next)
	if err != nil {
		return err
	}
	return s.runStates(state)
}

func panicError(r any) error {
	return fmt.Errorf("simulation panic occured: %v \n"+string(debug.Stack()), r)
}

func initialize(s *Simulation) (stateFn, error) {
	if !s.syncEval {
		go s.eval.Start()
	}
	// run sim for 90s if no duration set
	if s.cfg.Settings.Duration == 0 {
		// fmt.Println("no duration set, running for 90s")
		s.cfg.Settings.Duration = 90
	}
	s.C.Flags.DamageMode = s.cfg.Settings.DamageMode

	return s.advanceFrames(1, queuePhase)
}

func queuePhase(s *Simulation) (stateFn, error) {
	// noMoreActions stays set, so this runs frames until the sim ends
	for s.noMoreActions {
		if done, err := s.nextFrame(); done || err != nil {
			return nil, err
		}
	}
	if s.syncEval {
		// return to the evaluator, which calls execAction with the next action
		s.awaitingAction = true
		return nil, nil
	}
	s.eval.Continue()
	next, err := s.eval.NextAction()
	if err != nil {
		return nil, err
	}
	return s.queueAction(next)
}

// queueAction handles the evaluator's next action, nil if the program has no more
func (s *Simulation) queueAction(next *action.Eval) (stateFn, error) {
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
		delay := next.Param["f"]
		s.preActionDelay += delay
		s.C.Log.NewEvent(fmt.Sprintf("delay added %v, total: %v", delay, s.preActionDelay), glog.LogActionEvent, s.C.Player.Active()).
			Write("added", delay).
			Write("total", s.preActionDelay)
		return queuePhase, nil
	}
	// IMPORTANT: evaluator should handle adding in implicit swaps if next char is not active
	// we add a sanity check here just to guard against evaluator error
	if next.Char != s.C.Player.ActiveChar().Base.Key && next.Action != action.ActionSwap {
		return nil, fmt.Errorf("internal error: requested next char %v is not active and next action is not swap", next.Char)
	}
	// TODO: consider changing queue to single item. no need for slice without swap in here
	s.queue = append(s.queue, next)
	return actionReadyCheckPhase, nil
}

func actionReadyCheckPhase(s *Simulation) (stateFn, error) {
	// the phases that wait frame by frame loop over nextFrame themselves instead of returning
	// to runStates each frame, which costs two more calls per frame
	for {
		// TODO: this sanity check is probably not necessary
		if len(s.queue) == 0 {
			return nil, errors.New("unexpected queue length is 0")
		}
		q := s.queue[0]

		// check if the next queue item is valid
		// example: most sword characters can't do charge if the previous action was not attack
		char := s.C.Player.ActiveChar()
		if err := char.NextQueueItemIsValid(q.Char, q.Action, q.Param); err != nil {
			switch {
			case errors.Is(err, player.ErrInvalidChargeAction):
				return nil, fmt.Errorf("%v: %w", char.Base.Key, player.ErrInvalidChargeAction)
			default:
				return nil, err
			}
		}

		// TODO: this loop should be optimized to skip more than 1 frame at a time
		// ReadyCheck returns the sentinels unwrapped; comparing with == avoids the interface
		// assertions in errors.Is on every frame the action waits
		//nolint:errorlint // see above
		switch err := s.C.Player.ReadyCheck(q.Action, q.Char, q.Param); err {
		case nil, player.ErrActionNoOp:
			return executeActionPhase, nil
		case player.ErrActionNotReady:
			if s.C.Flags.LogDebug {
				s.C.Log.NewEvent(fmt.Sprintf("could not execute %v; action not ready", q.Action), glog.LogSimEvent, s.C.Player.Active())
			}
		case player.ErrPlayerNotReady:
		default:
			return nil, err
		}
		// repeat this phase on the next frame until action is ready
		if done, err := s.nextFrame(); done || err != nil {
			return nil, err
		}
	}
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

func executeActionDelay(s *Simulation) (stateFn, error) {
	for s.preActionDelay > 0 {
		if !s.C.Player.ActiveChar().FramePausedOnHitlag() {
			s.preActionDelay--
		}
		if done, err := s.nextFrame(); done || err != nil {
			return nil, err
		}
	}
	// go back to the ready check phase in case an action becomes unavailable after delay
	return actionReadyCheckPhase, nil
}

func executeActionPhase(s *Simulation) (stateFn, error) {
	// TODO: this sanity check is probably not necessary
	if len(s.queue) == 0 {
		return nil, errors.New("unexpected queue length is 0")
	}
	if s.preActionDelay > 0 {
		delay := s.preActionDelay
		s.C.Log.NewEvent(fmt.Sprintf("pre action delay: %v", delay), glog.LogActionEvent, s.C.Player.Active()).
			Write("delay", delay)
		return executeActionDelay, nil
	}
	q := s.queue[0]
	err := s.C.Player.Exec(q.Action, q.Char, q.Param)
	if err != nil {
		// TODO: this check probably doesn't do anything
		if errors.Is(err, player.ErrActionNoOp) {
			if l := s.popQueue(); l > 0 {
				// don't go back to queue if there are more actions already queued
				return actionReadyCheckPhase, nil
			}
			return queuePhase, nil
		}
		// this is now unexpected since action should be ready now
		// wrap the error for more context
		return nil, fmt.Errorf("error encountered on %v executing %v: %w", q.Char.String(), q.Action.String(), err)
	}
	// TODO: this check here is probably unnecessary
	if l := s.popQueue(); l > 0 {
		// don't go back to queue if there are more actions already queued
		return actionReadyCheckPhase, nil
	}

	return skipUntilCanQueue, nil
}

func skipUntilCanQueue(s *Simulation) (stateFn, error) {
	for !s.C.Player.CanQueueNextAction() {
		if done, err := s.nextFrame(); done || err != nil {
			return nil, err
		}
	}
	return queuePhase, nil
}

// nextFrame moves up the frame by 1, performing
func (s *Simulation) advanceFrames(f int, next stateFn) (stateFn, error) {
	for range f {
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
	// handleEnergy and handleHurt do nothing unless these hold; checking here saves two
	// calls per frame
	if e := &s.cfg.EnergySettings; e.Active && (e.Once || s.C.F-e.LastEnergyDrop >= e.Start) {
		s.handleEnergy()
	}
	if h := &s.cfg.HurtSettings; h.Active && (h.Once || s.C.F-h.LastHurt >= h.Start) {
		s.handleHurt()
	}
	s.C.Events.Emit(event.OnTick)
	s.yieldToGC()
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
			err = panicError(r)
		}
	}()
	res, err = s.run()
	return res, err
}

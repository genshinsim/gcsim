// package animation provides a simple way of tracking the current
// animation state at any given frame, as well as if the current frame
// is in animation lock or not
package animation

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/task"
)

type AnimationHandler struct { //nolint:revive // cannot just name this Handler because then there is a conflict with Handler in player package
	f      *int
	events event.Eventter
	log    glog.Logger
	tasks  task.Tasker

	char    int
	started int
	lastAct action.Action
	aniEvt  *action.Info

	state       action.AnimationState
	stateExpiry int

	debug bool
}

func New(f *int, debug bool, log glog.Logger, events event.Eventter, tasks task.Tasker) *AnimationHandler {
	h := &AnimationHandler{
		f:      f,
		log:    log,
		events: events,
		tasks:  tasks,
		debug:  debug,
	}
	return h
}

// IsAnimationLocked returns true if the next action can be executed on the
// current frame; false otherwise
func (h *AnimationHandler) IsAnimationLocked(next action.Action) bool {
	if h.aniEvt == nil {
		return false
	}
	// actions are always executed after ticks and right before we advance to
	// the next frame i.e. at the end of a frame
	//
	// if an action should take 20 frames of animation then this action would be ready if
	// f >= s + 20
	//
	// i.e the action lasted 20 frames counting the current frame
	// fmt.Printf("animation check; current frame %v, animation duration %v\n", *h.f, h.info.Frames(next))
	return !h.aniEvt.CanUse(next)
}

// CanQueue returns true if we can start looking for the next action to queue
// on the current frame, false otherwise
func (h *AnimationHandler) CanQueueNextAction() bool {
	if h.aniEvt == nil {
		return true
	}
	return h.aniEvt.CanQueueNext()
}

func (h *AnimationHandler) SetActionUsed(char int, act action.Action, evt *action.Info) {
	// remove previous if still active
	if h.aniEvt != nil {
		if h.aniEvt.OnRemoved != nil {
			h.aniEvt.OnRemoved(evt.State)
		}
		if h.debug {
			h.log.NewEvent(
				fmt.Sprintf("%v from %v ended, time passed: %v (actual: %v)", h.lastAct, h.started, h.aniEvt.TimePassed, h.aniEvt.NormalizedTimePassed),
				glog.LogHitlagEvent,
				h.char,
			)
		}
	}
	// setup next
	h.char = char
	h.started = *h.f
	h.aniEvt = evt
	h.events.Emit(event.OnStateChange, h.state, evt.State)
	h.state = evt.State
	h.stateExpiry = *h.f + evt.AnimationLength
	h.lastAct = act
	if h.debug {
		l := h.log.NewEvent(fmt.Sprintf("%v started", act.String()), glog.LogHitlagEvent, char)
		l.Write("AnimationLength", evt.AnimationLength).
			Write("CanQueueAfter", evt.CanQueueAfter).
			Write("State", evt.State.String())
		for i := action.Action(0); i < action.EndActionType; i++ {
			l.Write(i.String(), evt.Frames(i))
		}
	}
}

func (h *AnimationHandler) CurrentState() action.AnimationState {
	if h.aniEvt == nil {
		return action.Idle
	}
	return h.state
}

func (h *AnimationHandler) CurrentStateStart() int {
	return h.started
}

func (h *AnimationHandler) Tick() {
	if h.aniEvt != nil && h.aniEvt.Tick() {
		h.events.Emit(event.OnStateChange, h.state, action.Idle)
		h.state = action.Idle
		h.aniEvt = nil
	}
}

//package animation provides a simple way of tracking the current
//animation state at any given frame, as well as if the current frame
//is in animation lock or not
package animation

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/task"
)

type AnimationHandler struct {
	f      *int
	events event.Eventter
	log    glog.Logger
	tasks  task.Tasker

	char    int
	started int
	aniEvt  *action.ActionInfo

	state       action.AnimationState
	stateExpiry int
}

func New(f *int, log glog.Logger, events event.Eventter, tasks task.Tasker) *AnimationHandler {
	h := &AnimationHandler{
		f:      f,
		log:    log,
		events: events,
		tasks:  tasks,
	}
	return h
}

//IsAnimationLocked returns true if the next action can be executed on the
//current frame; false otherwise
func (h *AnimationHandler) IsAnimationLocked(next action.Action) bool {
	if h.aniEvt == nil {
		return false
	}
	//actions are always executed after ticks and right before we advance to
	//the next frame i.e. at the end of a frame
	//
	//if an action should take 20 frames of animation then this action would be ready if
	//f >= s + 20
	//
	//i.e the action lasted 20 frames counting the current frame
	// fmt.Printf("animation check; current frame %v, animation duration %v\n", *h.f, h.info.Frames(next))
	return h.aniEvt.CanUse(next)
}

//CanQueue returns true if we can start looking for the next action to queue
//on the current frame, false otherwise
func (h *AnimationHandler) CanQueueNextAction() bool {
	if h.aniEvt == nil {
		return true
	}
	return h.aniEvt.CanQueueNext()
}

func (h *AnimationHandler) SetActionUsed(char int, evt *action.ActionInfo) {
	h.char = char
	h.started = *h.f
	//remove previous if still active
	if h.aniEvt != nil {
		if h.aniEvt.OnRemoved != nil {
			h.aniEvt.OnRemoved()
		}
		pool.Put(h.aniEvt)
	}
	h.aniEvt = evt
	h.events.Emit(event.OnStateChange, h.state, evt.State)
	h.state = evt.State
	h.stateExpiry = *h.f + evt.AnimationLength
}

func (h *AnimationHandler) CurrentState() action.AnimationState {
	if h.aniEvt == nil {
		return action.Idle
	}
	return h.state
}

func (h *AnimationHandler) Tick() {
	if h.aniEvt != nil && h.aniEvt.Tick() {
		h.events.Emit(event.OnStateChange, h.state, action.Idle)
		h.state = action.Idle
		pool.Put(h.aniEvt)
		h.aniEvt = nil
	}
}

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
	info    action.ActionInfo

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
	//actions are always executed after ticks and right before we advance to
	//the next frame i.e. at the end of a frame
	//
	//if an action should take 20 frames of animation then this action would be ready if
	//f >= s + 20
	//
	//i.e the action lasted 20 frames counting the current frame
	return *h.f >= h.info.Frames(next)+h.started
}

//CanQueue returns true if we can start looking for the next action to queue
//on the current frame, false otherwise
func (h *AnimationHandler) CanQueueNextAction() bool {
	return *h.f >= h.info.CanQueueAfter+h.started
}

func (h *AnimationHandler) SetActionUsed(char int, info action.ActionInfo) {
	h.char = char
	h.started = *h.f
	h.info = info
	//update state to current state; drop state to idle if no change otherwise after
	//total animation druation
	h.events.Emit(event.OnStateChange, h.state, info.State)
	h.state = info.State
	h.stateExpiry = *h.f + info.AnimationLength
}

func (h *AnimationHandler) ClearState() {
	h.state = action.Idle
	h.stateExpiry = -1
}

func (h *AnimationHandler) CurrentState() action.AnimationState {
	if h.stateExpiry <= *h.f {
		h.state = action.Idle
		h.stateExpiry = -1
	}

	return h.state
}

func (h *AnimationHandler) Tick() {

}

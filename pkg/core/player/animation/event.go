package animation

import "github.com/genshinsim/gcsim/pkg/core/action"

type AnimationEvent interface {
	//Init is called when this event is first added
	Init(postFunc func())

	//Tick is called every frame and should return true if this animation
	//event has ended
	Tick() bool

	//CanUse returns true if the specified next action can be performed now
	CanUse(next action.Action) bool

	//CanQueueNext returns true if the next action can be queued now
	CanQueueNext() bool

	//State returns the type of animation state we are currently in
	AnimationState() action.AnimationState

	//OnRemoved is called if this animation event is replaced by another one
	OnRemoved()
}

type Action struct {
	F     func()
	Delay float64
}

type Event struct {
	Actions         []Action //sequential list of actions sorted earliest to latest
	Frames          func(action.Action) int
	cachedFrames    []int // cached cancelled frames
	AnimationLength int
	CanQueueAfter   int
	State           action.AnimationState
	OnRemovedCB     func()
	HitlagFactor    func() float64
	timePassed      float64

	//emit PostXX event. this is a bit awkward to handle though
	PostEventAt int
	postFunc    func()
}

func (e *Event) Init(postFunc func()) {
	e.cachedFrames = make([]int, action.EndActionType)
	for i := range e.cachedFrames {
		e.cachedFrames[i] = e.Frames(action.Action(i))
	}
	e.postFunc = postFunc
}

func (e *Event) Tick() bool {
	e.timePassed += 1 * e.HitlagFactor() //1 frame per tick for now

	//execute all action such that timePassed > delay, and then remove from
	//slice
	n := 0
	for i := range e.Actions {
		if e.Actions[i].Delay <= e.timePassed {
			e.Actions[i].F()
		} else {
			n = i
			break
		}
	}
	e.Actions = e.Actions[n:]

	//post even if not done yet
	if e.postFunc != nil && e.timePassed >= float64(e.PostEventAt) {
		//emit event??
		e.postFunc()
		e.postFunc = nil
	}

	//check if animation is over
	if e.timePassed > float64(e.AnimationLength) {
		//handle remove
		e.OnRemovedCB()
		return true
	}

	return false
}

func (e *Event) CanUse(next action.Action) bool {
	return e.timePassed >= float64(e.cachedFrames[next])
}

func (e *Event) CanQueueNext() bool {
	return e.timePassed >= float64(e.CanQueueAfter)
}

func (e *Event) AnimationState() action.AnimationState {
	return e.State
}

func (e *Event) OnRemoved() {
	e.OnRemovedCB()
}

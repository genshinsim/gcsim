package animation

import (
	"sync"

	"github.com/genshinsim/gcsim/pkg/core/action"
)

var pool = sync.Pool{
	New: func() any { return new(Event) },
}

type Action struct {
	F     func()
	Delay float64
}

func NewEvent() *Event {
	evt := pool.Get().(Event)
	evt.Frames = nil
	evt.AnimationLength = 0
	evt.CanQueueAfter = 0
	evt.State = action.Idle
	evt.OnRemovedCB = nil
	evt.HitlagFactor = nil
	evt.timePassed = 0
	evt.postFunc = nil
	evt.Post = 0
	return &evt
}

type Event struct {
	Actions         []Action //sequential list of actions sorted earliest to latest
	Frames          func(action.Action) int
	AnimationLength int
	CanQueueAfter   int
	Post            int
	State           action.AnimationState
	OnRemovedCB     func()
	HitlagFactor    func() float64

	cachedFrames [action.EndActionType]int // cached cancelled frames
	timePassed   float64

	//emit PostXX event. this is a bit awkward to handle though
	postFunc func()
}

func (e *Event) Init(postFunc func()) {
	for i := range e.cachedFrames {
		e.cachedFrames[i] = e.Frames(action.Action(i))
	}
	e.postFunc = postFunc
}

func (e *Event) AddAction(f func(), delay int) {
	e.Actions = append(e.Actions, Action{F: f, Delay: float64(delay)})
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
	if e.postFunc != nil && e.timePassed >= float64(e.Post) {
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

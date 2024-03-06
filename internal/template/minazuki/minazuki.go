// package minazuki provides common implementation for abilities that trigger
// based on normal animation state, i.e. xingqiu burst
package minazuki

import (
	"errors"
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/model"
)

// Watcher watches state change and triggers accordingly
type Watcher struct {
	// mandatory data
	key         keys.Char // name of the watcher; used for keying subscribers
	abil        string
	statusKey   string
	icdKey      string
	triggerFunc func()
	core        *core.Core
	tickerFreq  int

	// other fields including optional overrides
	state       action.AnimationState   // the state change we are watching for
	delayKey    model.AnimationDelayKey // delay key used to check delay func
	shouldDelay func() bool             // function to be called to see if delayed should be applied

	tickSrc int
}

type Config func(w *Watcher) error

func New(cfg ...Config) (*Watcher, error) {
	w := &Watcher{
		// defaults
		state:       action.NormalAttackState,
		shouldDelay: func() bool { return false },
	}
	for _, f := range cfg {
		err := f(w)
		if err != nil {
			return nil, err
		}
	}
	w.stateChangeHook()
	return w, nil
}

func WithMandatory(key keys.Char, abil, statusKey, icdKey string, tickerFreq int, triggerFunc func(), c *core.Core) Config {
	return func(w *Watcher) error {
		if abil == "" {
			return errors.New("ability name cannot be blank")
		}
		if statusKey == "" {
			return errors.New("status key cannot be blank")
		}
		if icdKey == "" {
			return errors.New("icd key cannot be blank")
		}
		if tickerFreq == 0 {
			return errors.New("ticker frequency cannot be 0")
		}
		w.key = key
		w.statusKey = statusKey
		w.icdKey = icdKey
		w.triggerFunc = triggerFunc
		w.tickerFreq = tickerFreq
		w.core = c
		return nil
	}
}

func WithAnimationDelayCheck(key model.AnimationDelayKey, shouldDelay func() bool) Config {
	return func(w *Watcher) error {
		w.delayKey = key
		w.shouldDelay = shouldDelay
		return nil
	}
}

func WithAnimationState(s action.AnimationState) Config {
	return func(w *Watcher) error {
		w.state = s
		return nil
	}
}

func (w *Watcher) Kill() {
	w.tickSrc = -1
}

func (w *Watcher) stateChangeHook() {
	w.core.Events.Subscribe(event.OnStateChange, func(args ...interface{}) bool {
		//TODO: can this ever fail?
		c, _ := w.core.Player.ByKey(w.key)
		// check if buff is up
		if !c.StatusIsActive(w.statusKey) {
			return false
		}
		next := args[1].(action.AnimationState)
		// ignore if it's not the state we are looking for
		if next != w.state {
			return false
		}
		// if we need to check for delay, and there is delay then we skip the rest of
		// this animation check and instead restart a ticker func after delay
		if w.shouldDelay() {
			if delay := w.core.Player.ActiveChar().AnimationStartDelay(w.delayKey); delay > 0 {
				w.core.Log.NewEvent(fmt.Sprintf("%v delay on state change", w.abil), glog.LogCharacterEvent, c.Index).
					Write("delay", delay)
				// it's sufficient to just restart the ticker function since that'll kill any existing
				// and restart the check every 60s from the delay
				w.tickSrc = w.core.F
				c.QueueCharTask(w.tickerFunc(w.core.F), 60) // check every 1sec
				return false
			}
		}

		// otherwise there is no delay and we can proceed as normal
		// ignore if on ICD
		if c.StatusIsActive(w.icdKey) {
			return false
		}
		// this should start a new ticker if not on ICD and state is correct
		w.triggerFunc()
		w.core.Log.NewEventBuildMsg(glog.LogCharacterEvent, c.Index, w.abil, " triggered on state change").
			Write("state", next).
			Write("icd", c.StatusExpiry(w.icdKey))
		w.tickSrc = w.core.F
		// use the hitlag affected queue for this
		c.QueueCharTask(w.tickerFunc(w.core.F), w.tickerFreq) // check every 1sec

		return false
	}, "xq-burst-animation-check")
}

func (w *Watcher) tickerFunc(src int) func() {
	return func() {
		//TODO: can this ever fail?
		c, _ := w.core.Player.ByKey(w.key)
		// check if buff is up
		if !c.StatusIsActive(w.statusKey) {
			return
		}
		if w.tickSrc != src {
			w.core.Log.NewEventBuildMsg(glog.LogCharacterEvent, c.Index, w.abil, " tick check ignored, src diff").
				Write("src", src).
				Write("new src", w.tickSrc)
			return
		}
		// stop if we are no longer in the right animation state
		state := w.core.Player.CurrentState()
		if state != w.state {
			w.core.Log.NewEventBuildMsg(glog.LogCharacterEvent, c.Index, w.abil, " tick check stoped, wrong state").
				Write("src", src).
				Write("state", state)
			return
		}
		// TODO: i THINK this check is not relevant because the ticksrc would have been reset already
		// by the state change watcher
		if w.shouldDelay() {
			// only nesting the if so it's easier to read...
			s := w.core.Player.CurrentStateStart()
			if w.core.F-s < w.core.Player.ActiveChar().AnimationStartDelay(w.delayKey) {
				w.core.Log.NewEventBuildMsg(glog.LogCharacterEvent, c.Index, w.abil, " not triggered; not enough time since normal state start").
					Write("current_state", state).
					Write("state_start", s)
				return
			}
		}
		// if there is a delay check then make sure current frame count is passed the delay
		w.core.Log.NewEventBuildMsg(glog.LogCharacterEvent, c.Index, w.abil, " triggered from ticker").
			Write("src", src).
			Write("state", state).
			Write("icd", c.StatusExpiry(w.statusKey))
		// we can trigger a wave here b/c we're in normal state still and src is still the same
		w.triggerFunc()
		// in theory this should not hit an icd?
		// use the hitlag affected queue for this
		c.QueueCharTask(w.tickerFunc(src), w.tickerFreq) // check every 1sec
	}
}

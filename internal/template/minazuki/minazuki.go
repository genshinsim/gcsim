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
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/model"
)

// Watcher watches state change and triggers accordingly
type Watcher struct {
	// mandatory data
	key         keys.Char // name of the watcher; used for keying subscribers
	caster      *character.CharWrapper
	abil        string
	statusKey   string
	icdKey      string
	triggerFunc func()
	core        *core.Core
	tickerFreq  int

	// other fields including optional overrides
	state        action.AnimationState   // the state change we are watching for
	delayKey     model.AnimationDelayKey // delay key used to check delay func
	shouldDelay  func() bool             // function to be called to see if delayed should be applied
	tickOnActive bool

	tickSrc int
}

type Config func(w *Watcher) error

func New(cfg ...Config) (*Watcher, error) {
	w := &Watcher{
		// defaults
		delayKey: model.InvalidAnimationDelayKey,
		state:    action.NormalAttackState,
	}
	for _, f := range cfg {
		err := f(w)
		if err != nil {
			return nil, err
		}
	}

	caster, ok := w.core.Player.ByKey(w.key)
	if !ok {
		return nil, errors.New("caster key is invalid")
	}
	w.caster = caster

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
		if tickerFreq == 0 {
			return errors.New("ticker frequency cannot be 0")
		}
		w.key = key
		w.abil = abil
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

func WithTickOnActive(v bool) Config {
	return func(w *Watcher) error {
		w.tickOnActive = v
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
		next := args[1].(action.AnimationState)
		// ignore if it's not the state we are looking for
		if next != w.state {
			return false
		}

		// if we need to check for delay, and there is a delay, then we delay the execution
		// of this state check
		// note that this is less performant because we don't actually need to do this check
		// if say the status is not active at all
		// however this just simpler to read so... performance hit shouldn't be that big
		if w.shouldDelay != nil { //TODO: to maintain old implementation equivalent; should be removed
			// if w.shouldDelay() {
			if delay := w.core.Player.ActiveChar().AnimationStartDelay(w.delayKey); delay > 0 {
				c := w.caster
				w.core.Log.NewEventBuildMsg(glog.LogDebugEvent, c.Index, w.abil, " delay on state change").
					Write("delay", delay)
				w.core.Tasks.Add(w.onStateChange(next), delay)
				return false
			}
		}
		w.onStateChange(next)()

		return false
	}, fmt.Sprintf("%v-burst-state-change-hook", w.key.String()))
}

func (w *Watcher) onStateChange(next action.AnimationState) func() {
	return func() {
		c := w.caster
		if !c.StatusIsActive(w.statusKey) {
			return
		}
		if w.icdKey != "" && c.StatusIsActive(w.icdKey) {
			w.core.Log.NewEventBuildMsg(glog.LogCharacterEvent, c.Index, w.abil, " not triggered on state change; on icd").
				Write("icd", c.StatusExpiry(w.icdKey)).
				Write("icd_key", w.icdKey)
			return
		}
		w.triggerFunc()
		w.core.Log.NewEventBuildMsg(glog.LogCharacterEvent, c.Index, w.abil, " triggered on state change").
			Write("state", next).
			Write("icd", c.StatusExpiry(w.icdKey)).
			Write("icd_key", w.icdKey)

		w.tickSrc = w.core.F
		w.queueTick(w.core.F)
	}
}

func (w *Watcher) queueTick(src int) {
	if w.tickerFreq <= 0 {
		return
	}

	c := w.caster
	if w.tickOnActive {
		c = w.core.Player.ActiveChar()
	}
	// use the hitlag affected queue for this
	c.QueueCharTask(w.tickerFunc(src), w.tickerFreq)
}

func (w *Watcher) tickerFunc(src int) func() {
	return func() {
		c := w.caster
		// check if buff is up
		if !c.StatusIsActive(w.statusKey) {
			w.core.Log.NewEventBuildMsg(glog.LogCharacterEvent, c.Index, w.abil, " not triggered on tick; on icd").
				Write("icd", c.StatusExpiry(w.icdKey)).
				Write("icd_key", w.icdKey)
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
			w.core.Log.NewEventBuildMsg(glog.LogCharacterEvent, c.Index, w.abil, " tick check stopped, wrong state").
				Write("src", src).
				Write("state", state)
			return
		}
		if w.shouldDelay != nil && w.shouldDelay() {
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
		w.queueTick(src)
	}
}

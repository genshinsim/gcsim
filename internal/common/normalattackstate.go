package common

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

var percentDelay5 = make([]int, keys.EndCharKeys)
var percentDelay5AltForms = make([]int, keys.EndCharKeys)
var altFormStatusKeys = make([]string, keys.EndCharKeys)

const Unused = -1

func init() {
	for i := range percentDelay5AltForms {
		percentDelay5AltForms[i] = Unused
	}
}

func Get5PercentN0Delay(activeChar *character.CharWrapper) int {
	activeCharKey := activeChar.Base.Key
	// The character doesn't have an alt form
	if percentDelay5AltForms[activeCharKey] == Unused {
		return percentDelay5[activeCharKey]
	}

	if activeChar.StatusIsActive(altFormStatusKeys[activeCharKey]) {
		return percentDelay5AltForms[activeCharKey]
	}

	return percentDelay5[activeCharKey]
}

func Get0PercentN0Delay(activeChar *character.CharWrapper) int {
	// TODO: Collect data for this
	return 0
}

type NAHook struct {
	C           *character.CharWrapper
	AbilName    string
	Core        *core.Core
	AbilKey     string
	AbilProcICD int
	AbilICDKey  string
	DelayFunc   func(*character.CharWrapper) int
	SummonFunc  func()
	abilTickSrc int
	abilHookSrc int
}

func (w *NAHook) Enable() {
	w.Core.Events.Subscribe(event.OnAttack, func(args ...interface{}) bool {
		// check if buff is up
		if !w.checkActive() {
			return false
		}
		w.abilHookSrc = w.Core.F
		delay := w.DelayFunc(w.Core.Player.ActiveChar())

		// This accounts for the delay in n0 timing needed to trigger
		if delay > 0 {
			w.Core.Log.NewEvent(fmt.Sprintf("%v delay on state change", w.AbilName), glog.LogCharacterEvent, w.C.Index).
				Write("delay", delay)
			w.Core.Tasks.Add(w.naStateDelayFuncGen(w.Core.F), delay)
			return false
		}
		// a delay of 0 will actually happen in the next frame, so a seperate conditional is used.

		// Additionally, at the time that OnAttack/OnStateChange events are emitted, the state has not yet changed, so we cannot do an animation check.
		if !w.checkActive() || !w.checkICD() {
			return false
		}
		w.Core.Log.NewEvent(
			fmt.Sprintf("%v triggered on state change", w.AbilName),
			glog.LogCharacterEvent,
			w.C.Index).
			Write("state", w.Core.Player.CurrentState()).
			Write("icd", w.C.StatusExpiry(w.AbilICDKey))
		w.trigger()
		return false
	}, fmt.Sprintf("%v animation check", w.AbilName))
}

func (w *NAHook) naStateDelayFuncGen(src int) func() {
	return func() {
		// ignore if on ICD
		if !w.checkActive() || !w.checkICD() || !w.checkState() || !w.checkSrc(w.abilHookSrc, src) {
			return
		}
		w.Core.Log.NewEvent(
			fmt.Sprintf("%v triggered on state change", w.AbilName),
			glog.LogCharacterEvent,
			w.C.Index).
			Write("state", w.Core.Player.CurrentState()).
			Write("icd", w.C.StatusExpiry(w.AbilICDKey))
		w.trigger()
	}
}

func (w *NAHook) naTickerFunc(src int) func() {
	return func() {
		// check if buff is up
		if !w.checkActive() || !w.checkState() || !w.checkSrc(w.abilTickSrc, src) {
			return
		}
		w.Core.Log.NewEvent(fmt.Sprintf("%v triggered from ticker", w.AbilName), glog.LogCharacterEvent, w.C.Index).
			Write("src", src).
			Write("state", w.Core.Player.CurrentState()).
			Write("icd", w.C.StatusExpiry(w.AbilICDKey))
		w.trigger()
	}
}

func (w *NAHook) trigger() {
	// we can trigger here b/c we're in normal state still and src is still the same
	w.SummonFunc()
	w.C.AddStatus(w.AbilICDKey, w.AbilProcICD, true)
	// in theory this should not hit an icd?
	// use the hitlag affected queue for this
	w.abilTickSrc = w.Core.F
	w.C.QueueCharTask(w.naTickerFunc(w.Core.F), w.AbilProcICD) // check every 1sec
}

func (w *NAHook) checkActive() bool {
	return w.C.StatusIsActive(w.AbilKey)
}

func (w *NAHook) checkICD() bool {
	icd := w.C.StatusIsActive(w.AbilICDKey)
	if icd {
		w.Core.Log.NewEvent(fmt.Sprintf("%v not triggered, on proc ICD", w.AbilName), glog.LogCharacterEvent, w.C.Index).
			Write("icd", w.C.StatusExpiry(w.AbilICDKey))
	}
	return !icd
}

func (w *NAHook) checkState() bool {
	state := w.Core.Player.CurrentState()
	if state != action.NormalAttackState {
		w.Core.Log.NewEvent(fmt.Sprintf("%v not triggered, not in normal state", w.AbilName), glog.LogCharacterEvent, w.C.Index).
			Write("state", state)
		return false
	}
	stateStart := w.Core.Player.CurrentStateStart()
	normCounter := w.Core.Player.ActiveChar().NormalCounter
	if (normCounter == 1) && w.Core.F-stateStart < w.DelayFunc(w.Core.Player.ActiveChar()) {
		w.Core.Log.NewEvent(fmt.Sprintf("%v not triggered, not enough time since normal state start", w.AbilName), glog.LogCharacterEvent, w.C.Index).
			Write("current_state", state).
			Write("state_start", stateStart)
		return false
	}
	return true
}
func (w *NAHook) checkSrc(newSrc, src int) bool {
	if newSrc != src {
		w.Core.Log.NewEvent(fmt.Sprintf("%v not triggered, src diff", w.AbilName), glog.LogCharacterEvent, w.C.Index).
			Write("src", src).
			Write("new src", newSrc)
		return false
	}
	return true
}

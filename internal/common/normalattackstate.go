package common

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/characters/ayato"
	"github.com/genshinsim/gcsim/internal/characters/cyno"
	"github.com/genshinsim/gcsim/internal/characters/tartaglia"
	"github.com/genshinsim/gcsim/internal/characters/wanderer"
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

	percentDelay5[keys.Baizhu] = 8
	percentDelay5[keys.Nahida] = 9
	percentDelay5[keys.Xingqiu] = 7
	percentDelay5[keys.Yelan] = 9
	percentDelay5[keys.Raiden] = 13
	percentDelay5[keys.Bennett] = 7
	percentDelay5[keys.Diluc] = 15
	percentDelay5[keys.Kazuha] = 10
	percentDelay5[keys.Keqing] = 8
	percentDelay5[keys.Xiangling] = 7
	percentDelay5[keys.Albedo] = 9
	percentDelay5[keys.Ayaka] = 7

	percentDelay5[keys.Tartaglia] = 9
	percentDelay5AltForms[keys.Tartaglia] = 12
	altFormStatusKeys[keys.Tartaglia] = tartaglia.MeleeKey

	percentDelay5[keys.Fischl] = 9
	percentDelay5[keys.Ganyu] = 10
	percentDelay5[keys.Jean] = 6

	percentDelay5[keys.Lumine] = 7
	percentDelay5[keys.LumineAnemo] = 7
	percentDelay5[keys.LumineCryo] = 7
	percentDelay5[keys.LumineDendro] = 7
	percentDelay5[keys.LumineElectro] = 7
	percentDelay5[keys.LumineGeo] = 7
	percentDelay5[keys.LumineHydro] = 7
	percentDelay5[keys.LuminePyro] = 7

	percentDelay5[keys.Nilou] = 11
	// I didn't test Nilou E stance, assuming it's the same values for now

	percentDelay5[keys.Venti] = 9
	percentDelay5[keys.Zhongli] = 9
	percentDelay5[keys.Amber] = 8
	percentDelay5[keys.Collei] = 11
	percentDelay5[keys.Diona] = 9
	percentDelay5[keys.Faruzan] = 9
	percentDelay5[keys.Gorou] = 11
	percentDelay5[keys.Heizou] = 10
	percentDelay5[keys.Kaeya] = 6
	percentDelay5[keys.Kuki] = 15
	percentDelay5[keys.Qiqi] = 7
	percentDelay5[keys.Rosaria] = 10
	percentDelay5[keys.Sara] = 14
	percentDelay5[keys.Thoma] = 11
	percentDelay5[keys.Yanfei] = 4
	percentDelay5[keys.Yunjin] = 12

	percentDelay5[keys.Beidou] = 22
	percentDelay5[keys.Kaveh] = 25
	percentDelay5[keys.Chongyun] = 18
	percentDelay5[keys.Dori] = 29
	percentDelay5[keys.Itto] = 27
	percentDelay5[keys.Noelle] = 23
	percentDelay5[keys.Razor] = 18
	percentDelay5[keys.Sayu] = 24
	percentDelay5[keys.Xinyan] = 28

	percentDelay5[keys.Aether] = 8
	percentDelay5[keys.AetherAnemo] = 8
	percentDelay5[keys.AetherCryo] = 8
	percentDelay5[keys.AetherDendro] = 8
	percentDelay5[keys.AetherElectro] = 8
	percentDelay5[keys.AetherGeo] = 8
	percentDelay5[keys.AetherHydro] = 8
	percentDelay5[keys.AetherPyro] = 8

	percentDelay5[keys.Ayato] = 15
	percentDelay5AltForms[keys.Ayato] = 17
	altFormStatusKeys[keys.Ayato] = ayato.SkillBuffKey

	percentDelay5[keys.Candace] = 14
	percentDelay5[keys.Eula] = 22
	percentDelay5[keys.Hutao] = 10
	percentDelay5[keys.Yoimiya] = 17

	percentDelay5[keys.Cyno] = 10
	percentDelay5AltForms[keys.Cyno] = 12
	altFormStatusKeys[keys.Cyno] = cyno.BurstKey

	percentDelay5[keys.Layla] = 12
	percentDelay5[keys.Shenhe] = 12
	percentDelay5[keys.YaeMiko] = 7
	percentDelay5[keys.Alhaitham] = 14
	percentDelay5[keys.Mika] = 10
	percentDelay5[keys.Kirara] = 13
	percentDelay5[keys.Lyney] = 12

	percentDelay5[keys.Wanderer] = 0
	percentDelay5AltForms[keys.Wanderer] = 12
	altFormStatusKeys[keys.Wanderer] = wanderer.SkillKey

	// Technically it's 15 for Left, 5 for Right, and 13 for Twirl
	percentDelay5[keys.Ningguang] = (15 + 5 + 13) / 3

	// jumping/dashing during the NA windup for some catalysts modifies their frames - said by koli
	// thus the current method of NA -> jump to test for N0 timing may not work on them
	percentDelay5[keys.Kokomi] = 0
	percentDelay5[keys.Sucrose] = 0
	percentDelay5[keys.Barbara] = 0
	percentDelay5[keys.Lisa] = 0
	percentDelay5[keys.Mona] = 0
	percentDelay5[keys.Klee] = 0
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
	switch activeChar.Base.Key {
	case keys.Baizhu:
		return 5
	case keys.Kaveh:
		return 23
	case keys.Kirara:
		return 6
	case keys.Lyney:
		return 5
	}
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
		//check if buff is up
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
		w.trigger(w.Core.F)
		return false
	}, fmt.Sprintf("%v animation check", w.AbilName))
}

func (w *NAHook) naStateDelayFuncGen(src int) func() {
	return func() {
		//ignore if on ICD
		if !w.checkActive() || !w.checkICD() || !w.checkState() || !w.checkSrc(w.abilHookSrc, src) {
			return
		}
		w.Core.Log.NewEvent(
			fmt.Sprintf("%v triggered on state change", w.AbilName),
			glog.LogCharacterEvent,
			w.C.Index).
			Write("state", w.Core.Player.CurrentState()).
			Write("icd", w.C.StatusExpiry(w.AbilICDKey))
		w.trigger(src)
	}
}

func (w *NAHook) naTickerFunc(src int) func() {
	return func() {
		//check if buff is up
		if !w.checkActive() || !w.checkState() || !w.checkSrc(w.abilTickSrc, src) {
			return
		}
		w.Core.Log.NewEvent(fmt.Sprintf("%v triggered from ticker", w.AbilName), glog.LogCharacterEvent, w.C.Index).
			Write("src", src).
			Write("state", w.Core.Player.CurrentState()).
			Write("icd", w.C.StatusExpiry(w.AbilICDKey))
		w.trigger(src)
	}
}

func (w *NAHook) trigger(src int) {
	//we can trigger here b/c we're in normal state still and src is still the same
	w.SummonFunc()
	w.C.AddStatus(w.AbilICDKey, w.AbilProcICD, true)
	//in theory this should not hit an icd?
	//use the hitlag affected queue for this
	w.abilTickSrc = w.Core.F
	w.C.QueueCharTask(w.naTickerFunc(w.Core.F), w.AbilProcICD) //check every 1sec
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
	state_start := w.Core.Player.CurrentStateStart()
	norm_counter := w.Core.Player.ActiveChar().NormalCounter
	if (norm_counter == 1) && w.Core.F-state_start < w.DelayFunc(w.Core.Player.ActiveChar()) {
		w.Core.Log.NewEvent(fmt.Sprintf("%v not triggered, not enough time since normal state start", w.AbilName), glog.LogCharacterEvent, w.C.Index).
			Write("current_state", state).
			Write("state_start", state_start)
		return false
	}
	return true

}
func (w *NAHook) checkSrc(newSrc int, src int) bool {
	if newSrc != src {
		w.Core.Log.NewEvent(fmt.Sprintf("%v not triggered, src diff", w.AbilName), glog.LogCharacterEvent, w.C.Index).
			Write("src", src).
			Write("new src", newSrc)
		return false
	}
	return true
}

package common

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/characters/ayato"
	"github.com/genshinsim/gcsim/internal/characters/cyno"
	"github.com/genshinsim/gcsim/internal/characters/tartaglia"
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

	// TODO: Uncomment when Wanderer Implementation is done
	// percentDelay5[keys.Wanderer] = 0
	// percentDelay5AltForms[keys.Wanderer] = 12
	// AltFormStatusKeys[keys.Wanderer] = wanderer.SkillKey

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
	if percentDelay5[activeCharKey] == -1 {
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
	c           *character.CharWrapper
	abilName    string
	Core        *core.Core
	abilKey     string
	abilProcICD int
	abilICDKey  string
	abilTickSrc int
	abilHookSrc int
	delayFunc   func(*character.CharWrapper) int
	summonFunc  func()
}

func NewNAHook(c *character.CharWrapper, Core *core.Core, abilName string, abilKey string, abilProcICD int, abilICDKey string, delayFunc func(*character.CharWrapper) int, summonFunc func()) *NAHook {
	return &NAHook{
		c:           c,
		abilName:    abilName,
		Core:        Core,
		abilKey:     abilKey,
		abilProcICD: abilProcICD,
		abilICDKey:  abilICDKey,
		abilTickSrc: 0,
		abilHookSrc: 0,
		delayFunc:   delayFunc,
		summonFunc:  summonFunc,
	}

}

func (w *NAHook) NAStateHook() {
	w.Core.Events.Subscribe(event.OnAttack, func(args ...interface{}) bool {
		//check if buff is up
		if !w.c.StatusIsActive(w.abilKey) {
			return false
		}
		w.abilHookSrc = w.Core.F
		delay := w.delayFunc(w.Core.Player.ActiveChar())
		w.Core.Log.NewEvent(fmt.Sprintf("%v delay on state change", w.abilName), glog.LogCharacterEvent, w.c.Index).
			Write("delay", delay)
		// This accounts for the delay in n0 timing needed to trigger

		if delay > 0 {

			w.Core.Tasks.Add(w.naStateDelayFuncGen(w.Core.F), delay)
		} else {
			// a delay of 0 will actually happen in the next frame, so a seperate conditional is used.

			// Additionally, at the time that OnAttack/OnStateChange events are emitted, the state has not yet changed, so we cannot do an animation check.
			if w.c.StatusIsActive(w.abilICDKey) {
				w.Core.Log.NewEvent(fmt.Sprintf("%v did not trigger on state change", w.abilName), glog.LogCharacterEvent, w.c.Index).
					Write("state", w.Core.Player.CurrentState()).
					Write("icd", w.c.StatusExpiry(w.abilICDKey))
			} else {
				w.summonFunc()
				w.c.AddStatus(w.abilICDKey, 60, true)
				w.Core.Log.NewEvent(fmt.Sprintf("%v on state change", w.abilName), glog.LogCharacterEvent, w.c.Index).
					Write("state", action.NormalAttackState).
					Write("icd", w.c.StatusExpiry(w.abilICDKey))
				w.abilTickSrc = w.Core.F
				//use the hitlag affected queue for this
				w.c.QueueCharTask(w.naTickerFunc(w.Core.F), w.abilProcICD) //check every `ICD`` frames
			}
		}

		return false
	}, fmt.Sprintf("%v animation check", w.abilName))
}

func (w *NAHook) naStateDelayFuncGen(src int) func() {
	return func() {
		//ignore if on ICD
		if w.c.StatusIsActive(w.abilICDKey) || w.Core.Player.CurrentState() != action.NormalAttackState || w.abilHookSrc != src {
			w.Core.Log.NewEvent(fmt.Sprintf("%v did not trigger on state change", w.abilName), glog.LogCharacterEvent, w.c.Index).
				Write("state", w.Core.Player.CurrentState()).
				Write("icd", w.c.StatusExpiry(w.abilICDKey)).
				Write("ICD check", w.c.StatusIsActive(w.abilICDKey)).
				Write("State check", w.Core.Player.CurrentState() != action.NormalAttackState).
				Write("Source check", w.abilHookSrc != src)

			return
		}
		//this should start a new ticker if not on ICD and state is correct
		w.summonFunc()
		w.c.AddStatus(w.abilICDKey, 60, true)
		w.Core.Log.NewEvent(fmt.Sprintf("%v on state change", w.abilName), glog.LogCharacterEvent, w.c.Index).
			Write("state", action.NormalAttackState).
			Write("icd", w.c.StatusExpiry(w.abilICDKey))
		w.abilTickSrc = w.Core.F
		//use the hitlag affected queue for this
		w.c.QueueCharTask(w.naTickerFunc(w.Core.F), w.abilProcICD) //check every `ICD` frames
	}
}

func (w *NAHook) naTickerFunc(src int) func() {
	return func() {
		//check if buff is up
		if !w.c.StatusIsActive(w.abilKey) {
			return
		}
		if w.abilTickSrc != src {
			w.Core.Log.NewEvent(fmt.Sprintf("%v tick check ignored, src diff", w.abilName), glog.LogCharacterEvent, w.c.Index).
				Write("src", src).
				Write("new src", w.abilTickSrc)
			return
		}
		//stop if we are no longer in normal animation state
		state := w.Core.Player.CurrentState()

		if state != action.NormalAttackState {
			w.Core.Log.NewEvent(fmt.Sprintf("%v tick check stopped, not in normal state", w.abilName), glog.LogCharacterEvent, w.c.Index).
				Write("src", src).
				Write("state", state)
			return
		}
		state_start := w.Core.Player.CurrentStateStart()
		norm_counter := w.Core.Player.ActiveChar().NormalCounter
		if (norm_counter == 1) && w.Core.F-state_start < w.delayFunc(w.Core.Player.ActiveChar()) {
			w.Core.Log.NewEvent(fmt.Sprintf("%v tick check stopped, not enough time since normal state start", w.abilName), glog.LogCharacterEvent, w.c.Index).
				Write("src", src).
				Write("state_start", state_start)
			return
		}

		w.Core.Log.NewEvent(fmt.Sprintf("%v triggered from ticker", w.abilName), glog.LogCharacterEvent, w.c.Index).
			Write("src", src).
			Write("state", state).
			Write("icd", w.c.StatusExpiry(w.abilICDKey))

		//we can trigger here b/c we're in normal state still and src is still the same
		w.summonFunc()
		w.c.AddStatus(w.abilICDKey, 60, true)
		//in theory this should not hit an icd?
		//use the hitlag affected queue for this
		w.abilTickSrc = w.Core.F
		w.c.QueueCharTask(w.naTickerFunc(w.Core.F), w.abilProcICD) //check every 1sec
	}
}

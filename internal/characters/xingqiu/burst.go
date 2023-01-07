package xingqiu

import (
	"github.com/genshinsim/gcsim/internal/characters/raiden"
	"github.com/genshinsim/gcsim/internal/characters/tartaglia"
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

var burstFrames []int
var XQ_N0_delays = make([]int, keys.EndCharKeys)
var XQ_N0_delays_alt_forms = make([]int, keys.EndCharKeys)
var alt_form_status_keys = make([]string, keys.EndCharKeys)

const (
	burstHitmark = 18
	burstKey     = "xingqiuburst"
	burstICDKey  = "xingqiu-burst-icd"
)

func init() {
	burstFrames = frames.InitAbilSlice(40)
	burstFrames[action.ActionAttack] = 33
	burstFrames[action.ActionSkill] = 33
	burstFrames[action.ActionDash] = 33
	burstFrames[action.ActionJump] = 33

	for i := range XQ_N0_delays_alt_forms {
		XQ_N0_delays_alt_forms[i] = -1
	}

	XQ_N0_delays[keys.Nahida] = 9
	XQ_N0_delays[keys.Xingqiu] = 7
	XQ_N0_delays[keys.Yelan] = 9

	XQ_N0_delays[keys.Raiden] = 13
	XQ_N0_delays_alt_forms[keys.Raiden] = 13
	alt_form_status_keys[keys.Raiden] = raiden.BurstKey

	XQ_N0_delays[keys.Bennett] = 7
	XQ_N0_delays[keys.Diluc] = 15
	XQ_N0_delays[keys.Kazuha] = 10
	XQ_N0_delays[keys.Keqing] = 8
	XQ_N0_delays[keys.Xiangling] = 7
	XQ_N0_delays[keys.Albedo] = 9
	XQ_N0_delays[keys.Ayaka] = 7

	XQ_N0_delays[keys.Tartaglia] = 9
	XQ_N0_delays_alt_forms[keys.Tartaglia] = 12
	alt_form_status_keys[keys.Tartaglia] = tartaglia.MeleeKey

	XQ_N0_delays[keys.Fischl] = 9
	XQ_N0_delays[keys.Ganyu] = 10
	XQ_N0_delays[keys.Jean] = 6
	XQ_N0_delays[keys.Lumine] = 7

	XQ_N0_delays[keys.Nilou] = 11
	// I didn't test Nilou E stance, assuming it's the same values for now

	XQ_N0_delays[keys.Venti] = 9
	XQ_N0_delays[keys.Zhongli] = 9
	XQ_N0_delays[keys.Amber] = 8
	XQ_N0_delays[keys.Collei] = 11
	XQ_N0_delays[keys.Diona] = 9
	XQ_N0_delays[keys.Faruzan] = 9
	XQ_N0_delays[keys.Gorou] = 11
	XQ_N0_delays[keys.Heizou] = 10
	XQ_N0_delays[keys.Kaeya] = 6
	XQ_N0_delays[keys.Kuki] = 15
	XQ_N0_delays[keys.Qiqi] = 7
	XQ_N0_delays[keys.Rosaria] = 10
	XQ_N0_delays[keys.Sara] = 14
	XQ_N0_delays[keys.Thoma] = 11
	XQ_N0_delays[keys.Yanfei] = 4
	XQ_N0_delays[keys.Yunjin] = 12

	XQ_N0_delays[keys.Beidou] = 22
	XQ_N0_delays[keys.Chongyun] = 18
	XQ_N0_delays[keys.Dori] = 29
	XQ_N0_delays[keys.Itto] = 27
	XQ_N0_delays[keys.Noelle] = 23
	XQ_N0_delays[keys.Razor] = 18
	XQ_N0_delays[keys.Sayu] = 24
	XQ_N0_delays[keys.Xinyan] = 28

	// Technically it's 15 for Left, 5 for Right, and 13 for Twirl
	XQ_N0_delays[keys.Ningguang] = (15 + 5 + 13) / 3
	// XQ_N0_delays_alt_forms[keys.Ningguang] = 15
	// alt_form_status_keys[keys.Ningguang] = "Left"
	// XQ_N0_delays_alt_forms[keys.Ningguang] = 5
	// alt_form_status_keys[keys.Ningguang] = "Right"

	// jumping/dashing during the NA windup for some catalysts modifies their frames - said by koli
	// thus the current method of NA -> jump to test for N0 timing won't work on them
	XQ_N0_delays[keys.Kokomi] = 0
	XQ_N0_delays[keys.Sucrose] = 0
	XQ_N0_delays[keys.Barbara] = 0
	XQ_N0_delays[keys.Lisa] = 0
	XQ_N0_delays[keys.Mona] = 0

}

/**
The number of Hydro Swords summoned per wave follows a specific pattern, usually alternating between 2 and 3 swords.
At C6, this is upgraded and follows a pattern of 2 → 3 → 5… which then repeats.

There is an approximately 1 second interval between summoned Hydro Sword waves, so that means a theoretical maximum of 15 or 18 waves.

Each wave of Hydro Swords is capable of applying one (1) source of Hydro status, and each individual sword is capable of getting a crit.
**/

func (c *char) Burst(p map[string]int) action.ActionInfo {
	//apply hydro every 3rd hit
	//triggered on normal attack
	//also applies hydro on cast if p=1
	//how we doing that?? trigger 0 dmg?

	/** c2
	Extends the duration of Guhua Sword: Raincutter by 3s.
	Decreases the Hydro RES of opponents hit by sword rain attacks by 15% for 4s.
	**/
	dur := 15
	if c.Base.Cons >= 2 {
		dur += 3
	}
	dur = dur * 60
	c.AddStatus(burstKey, dur+33, true) // add 33f for anim

	orbital, ok := p["orbital"]
	if !ok {
		orbital = 1
	}

	if orbital == 1 {
		c.applyOrbital(dur, burstHitmark)
	}

	c.burstCounter = 0
	c.numSwords = 2
	c.nextRegen = false

	// c.CD[combat.BurstCD] = c.S.F + 20*60
	c.SetCD(action.ActionBurst, 20*60)
	c.ConsumeEnergy(3)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstHitmark,
		State:           action.BurstState,
	}
}

func (c *char) summonSwordWave() {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Guhua Sword: Raincutter",
		AttackTag:  combat.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagElementalBurst,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypePierce,
		Element:    attributes.Hydro,
		Durability: 25,
		Mult:       burst[c.TalentLvlBurst()],
	}

	//only if c.nextRegen is true and first sword
	var c2cb, c6cb func(a combat.AttackCB)
	if c.nextRegen {
		done := false
		c6cb = func(_ combat.AttackCB) {
			if done {
				return
			}
			c.AddEnergy("xingqiu-c6", 3)
			done = true
		}
	}
	if c.Base.Cons >= 2 {
		icd := -1
		c2cb = func(a combat.AttackCB) {
			if c.Core.F < icd {
				return
			}

			e, ok := a.Target.(*enemy.Enemy)
			if !ok {
				return
			}

			icd = c.Core.F + 1
			c.Core.Tasks.Add(func() {
				e.AddResistMod(enemy.ResistMod{
					Base:  modifier.NewBaseWithHitlag("xingqiu-c2", 4*60),
					Ele:   attributes.Hydro,
					Value: -0.15,
				})
			}, 1)
		}
	}

	for i := 0; i < c.numSwords; i++ {
		//TODO: this snapshot timing is off? perhaps should be 0?
		c.Core.QueueAttack(ai, combat.NewCircleHit(c.Core.Combat.PrimaryTarget(), 0.5), 20, 20, c2cb, c6cb)
		c6cb = nil
		c.burstCounter++
	}

	//figure out next wave # of swords
	switch c.numSwords {
	case 2:
		c.numSwords = 3
		c.nextRegen = false
	case 3:
		if c.Base.Cons >= 6 {
			c.numSwords = 5
			c.nextRegen = true
		} else {
			c.numSwords = 2
			c.nextRegen = false
		}
	case 5:
		c.numSwords = 2
		c.nextRegen = false
	}

	c.AddStatus(burstICDKey, 60, true)
}

func (c *char) getN0Delay() int {
	active := c.Core.Player.ActiveChar()
	activeCharKey := active.Base.Key
	// The character doesn't have an alt form
	if XQ_N0_delays_alt_forms[activeCharKey] == -1 {
		return XQ_N0_delays[activeCharKey]
	}

	if active.StatusIsActive(alt_form_status_keys[activeCharKey]) {
		return XQ_N0_delays_alt_forms[activeCharKey]
	}

	return XQ_N0_delays[activeCharKey]
}

func (c *char) burstStateDelayFuncGen(src int) func() {
	return func() {
		//ignore if on ICD
		if c.StatusIsActive(burstICDKey) || c.Core.Player.CurrentState() != action.NormalAttackState || c.burstTickSrc != src {
			return
		}
		//this should start a new ticker if not on ICD and state is correct
		c.summonSwordWave()
		c.Core.Log.NewEvent("xq burst on state change", glog.LogCharacterEvent, c.Index).
			Write("state", action.NormalAttackState).
			Write("icd", c.StatusExpiry(burstICDKey))
		c.burstTickSrc = c.Core.F
		//use the hitlag affected queue for this
		c.QueueCharTask(c.burstTickerFunc(c.Core.F), 60) //check every 1sec
	}
}

func (c *char) burstStateHook() {
	c.Core.Events.Subscribe(event.OnStateChange, func(args ...interface{}) bool {
		//check if buff is up
		if !c.StatusIsActive(burstKey) {
			return false
		}
		next := args[1].(action.AnimationState)
		//ignore if not normal
		if next != action.NormalAttackState {
			return false
		}
		c.burstTickSrc = c.Core.F
		delay := c.getN0Delay()
		c.Core.Log.NewEvent("xq burst delay on state change", glog.LogCharacterEvent, c.Index).
			Write("active", c.Core.Player.ActiveChar().Base.Key.String()).
			Write("delay", delay)
		// This accounts for the delay in n0 timing needed for XQ to trigger rainswords
		c.Core.Tasks.Add(c.burstStateDelayFuncGen(c.Core.F), delay)

		return false
	}, "xq-burst-animation-check")
}

func (c *char) burstTickerFunc(src int) func() {
	return func() {
		//check if buff is up
		if !c.StatusIsActive(burstKey) {
			return
		}
		if c.burstTickSrc != src {
			c.Core.Log.NewEvent("xq burst tick check ignored, src diff", glog.LogCharacterEvent, c.Index).
				Write("src", src).
				Write("new src", c.burstTickSrc)
			return
		}
		//stop if we are no longer in normal animation state
		state := c.Core.Player.CurrentState()

		if state != action.NormalAttackState {
			c.Core.Log.NewEvent("xq burst tick check stopped, not in normal state", glog.LogCharacterEvent, c.Index).
				Write("src", src).
				Write("state", state)
			return
		}
		state_start := c.Core.Player.CurrentStateStart()
		norm_counter := c.Core.Player.ActiveChar().NormalCounter
		if (norm_counter == 1) && c.Core.F-state_start < c.getN0Delay() {
			c.Core.Log.NewEvent("xq burst tick check stopped, not enough time since normal state start", glog.LogCharacterEvent, c.Index).
				Write("src", src).
				Write("state_start", state_start)
			return
		}

		c.Core.Log.NewEvent("xq burst triggered from ticker", glog.LogCharacterEvent, c.Index).
			Write("src", src).
			Write("state", state).
			Write("icd", c.StatusExpiry(burstICDKey))
		//we can trigger a wave here b/c we're in normal state still and src is still the same
		c.summonSwordWave()
		//in theory this should not hit an icd?
		//use the hitlag affected queue for this
		c.QueueCharTask(c.burstTickerFunc(src), 60) //check every 1sec
	}
}

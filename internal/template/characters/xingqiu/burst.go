package xingqiu

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

var burstFrames []int

const burstHitmark = 18

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
	c.Core.Status.Add("xqburst", dur+33) // add 33f for anim
	c.Core.Log.NewEvent("Xingqiu burst activated", glog.LogCharacterEvent, c.Index, "expiry", c.Core.F+dur+33)

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
		Post:            burstHitmark,
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
		Element:    attributes.Hydro,
		Durability: 25,
		Mult:       burst[c.TalentLvlBurst()],
	}

	//only if c.nextRegen is true and first sword
	var c2cb, c6cb func(a combat.AttackCB)
	if c.nextRegen {
		c6cb = func(a combat.AttackCB) {
			c.AddEnergy("xingqiu-c6", 3)
		}
	}
	if c.Base.Cons >= 2 {
		icd := -1
		c2cb = func(a combat.AttackCB) {
			if c.Core.F < icd {
				return
			}

			e, ok := a.Target.(core.Enemy)
			if !ok {
				return
			}

			icd = c.Core.F + 1
			c.Core.Tasks.Add(func() {
				e.AddResistMod("xingqiu-c2", 4*60, attributes.Hydro, -0.15)
			}, 1)
		}
	}

	for i := 0; i < c.numSwords; i++ {
		c.Core.QueueAttack(ai, combat.NewDefCircHit(0.1, false, combat.TargettableEnemy), 20, 20, c2cb, c6cb)
		c6cb = nil
		c.burstCounter++
	}

	//figure out next wave # of swords
	switch c.numSwords {
	case 2:
		c.numSwords = 3
		c.nextRegen = false
	case 3:
		if c.Base.Cons == 6 {
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

	c.burstSwordICD = c.Core.F + 60
}

func (c *char) burstStateHook() {
	c.Core.Events.Subscribe(event.OnStateChange, func(args ...interface{}) bool {
		//check if buff is up
		if c.Core.Status.Duration("xqburst") <= 0 {
			return false
		}
		next := args[1].(action.AnimationState)
		//ignore if not normal
		if next != action.NormalAttackState {
			return false
		}
		//ignore if on ICD
		if c.burstSwordICD > c.Core.F {
			return false
		}
		//this should start a new ticker if not on ICD and state is correct
		c.summonSwordWave()
		c.Core.Log.NewEvent("xq burst on state change", glog.LogCharacterEvent, c.Index, "state", next, "icd", c.burstSwordICD)
		c.burstTickSrc = c.Core.F
		c.Core.Tasks.Add(c.burstTickerFunc(c.Core.F), 60) //check every 1sec

		return false
	}, "xq-burst-animation-check")
}

func (c *char) burstTickerFunc(src int) func() {
	return func() {
		//check if buff is up
		if c.Core.Status.Duration("xqburst") <= 0 {
			return
		}
		if c.burstTickSrc != src {
			c.Core.Log.NewEvent("xq burst tick check ignored, src diff", glog.LogCharacterEvent, c.Index, "src", src, "new src", c.burstTickSrc)
			return
		}
		//stop if we are no longer in normal animation state
		state := c.Core.Player.CurrentState()
		if state != action.NormalAttackState {
			c.Core.Log.NewEvent("xq burst tick check stopped, not normal state", glog.LogCharacterEvent, c.Index, "src", src, "state", state)
			return
		}
		c.Core.Log.NewEvent("xq burst triggered from ticker", glog.LogCharacterEvent, c.Index, "src", src, "state", state, "icd", c.burstSwordICD)
		//we can trigger a wave here b/c we're in normal state still and src is still the same
		c.summonSwordWave()
		//in theory this should not hit an icd?
		c.Core.Tasks.Add(c.burstTickerFunc(src), 60) //check every 1sec
	}
}

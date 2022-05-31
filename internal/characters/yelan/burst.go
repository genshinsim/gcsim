package yelan

import "github.com/genshinsim/gcsim/pkg/core"

func (c *char) summonExquisiteThrow() {
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Exquisite Throw",
		AttackTag:  core.AttackTagElementalBurst,
		ICDTag:     core.ICDTagYelanBurst,
		ICDGroup:   core.ICDGroupYelanBurst,
		Element:    core.Hydro,
		Durability: 25,
		Mult:       0,
		FlatDmg:    burstDice[c.TalentLvlBurst()] * c.MaxHP(),
	}
	for i := 0; i < 3; i++ {
		//TODO: frames timing on this?
		c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(0.1, false, core.TargettableEnemy), 22+i*6, 22+i*6)
	}
	if c.Base.Cons >= 2 && c.c2icd <= c.Core.F {
		ai.Abil = "Yelan C2 Proc"
		ai.FlatDmg = 12.0 / 100 * c.MaxHP()
		c.c2icd = c.Core.F + 1.6*60
		//TODO: frames timing on this?
		c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(0.1, false, core.TargettableEnemy), 22+4*6, 22+4*6)
	}

	c.burstDiceICD = c.Core.F + 60
}

func (c *char) burstStateHook() {
	c.Core.Events.Subscribe(core.OnStateChange, func(args ...interface{}) bool {
		//check if buff is up
		if c.Core.Status.Duration(burstStatus) <= 0 {
			return false
		}
		next := args[1].(core.AnimationState)
		//ignore if not normal
		if next != core.NormalAttackState {
			return false
		}
		//ignore if on ICD
		if c.burstDiceICD > c.Core.F {
			return false
		}
		//this should start a new ticker if not on ICD and state is correct
		c.summonExquisiteThrow()
		c.Core.Log.NewEvent("yelan burst on state change", core.LogCharacterEvent, c.Index, "state", next, "icd", c.burstDiceICD)
		c.burstTickSrc = c.Core.F
		c.AddTask(c.burstTickerFunc(c.Core.F), "yelan-ticker", 60) //check every 1sec

		return false
	}, "yelan-burst-animation-check")
}

func (c *char) burstTickerFunc(src int) func() {
	return func() {
		//check if buff is up
		if c.Core.Status.Duration(burstStatus) <= 0 {
			return
		}
		if c.burstTickSrc != src {
			c.Core.Log.NewEvent("yelan burst tick check ignored, src diff", core.LogCharacterEvent, c.Index, "src", src, "new src", c.burstTickSrc)
			return
		}
		//stop if we are no longer in normal animation state
		state := c.Core.State()
		if state != core.NormalAttackState {
			c.Core.Log.NewEvent("yelan burst tick check stopped, not normal state", core.LogCharacterEvent, c.Index, "src", src, "state", state)
			return
		}
		c.Core.Log.NewEvent("yelan burst triggered from ticker", core.LogCharacterEvent, c.Index, "src", src, "state", state, "icd", c.burstDiceICD)
		//we can trigger a wave here b/c we're in normal state still and src is still the same
		c.summonExquisiteThrow()
		//in theory this should not hit an icd?
		c.AddTask(c.burstTickerFunc(src), "yelan-ticker", 60) //check every 1sec
	}
}

package amber

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func (c *char) Attack(p map[string]int) (int, int) {
	travel, ok := p["travel"]
	if !ok {
		travel = 20
	}

	f, a := c.ActionFrames(core.ActionAttack, p)

	ai := core.AttackInfo{
		Abil:       fmt.Sprintf("Normal %v", c.NormalCounter),
		ActorIndex: c.Index,
		AttackTag:  core.AttackTagNormal,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Physical,
		Durability: 25,
		Mult:       attack[c.NormalCounter][c.TalentLvlAttack()],
	}

	c.Core.Combat.QueueAttack(ai, core.NewDefSingleTarget(1, core.TargettableEnemy), f, f+travel)

	c.AdvanceNormalIndex()

	return f, a
}

func (c *char) Aimed(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionAim, p)

	travel, ok := p["travel"]
	if !ok {
		travel = 20
	}

	b := p["bunny"]

	if c.Base.Cons >= 2 && b != 0 {
		//explode the first bunny
		c.AddTask(func() {
			c.manualExplode()
		}, "bunny", travel+f)

		//also don't do any dmg since we're shooting at bunny

		return f, a
	}

	ai := core.AttackInfo{
		Abil:         "Aim (Charged)",
		ActorIndex:   c.Index,
		AttackTag:    core.AttackTagNormal,
		ICDTag:       core.ICDTagNone,
		ICDGroup:     core.ICDGroupDefault,
		Element:      core.Pyro,
		Durability:   50,
		Mult:         aim[c.TalentLvlAttack()],
		HitWeakPoint: true,
	}

	// d.AnimationFrames = f

	//add 15% since 360noscope
	cb := func(a core.AttackCB) {
		if a.AttackEvent.Info.HitWeakPoint {
			c.AddMod(core.CharStatMod{
				Key: "a2",
				Amount: func(a core.AttackTag) ([]float64, bool) {
					val := make([]float64, core.EndStatType)
					val[core.ATKP] = 0.15
					return val, true
				},
				Expiry: c.Core.F + 600,
			})
		}

	}

	c.Core.Combat.QueueAttack(ai, core.NewDefSingleTarget(1, core.TargettableEnemy), f, f+travel, cb)

	if c.Base.Cons >= 1 {
		ai.Mult = .2 * ai.Mult
		c.Core.Combat.QueueAttack(ai, core.NewDefSingleTarget(1, core.TargettableEnemy), f, f+travel)
	}

	return f, a
}

func (c *char) Skill(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionSkill, p)
	hold := p["hold"]

	c.AddTask(func() {
		c.makeBunny()
	}, "new-bunny", f+hold)

	c.overloadExplode()

	if c.Base.Cons < 4 {
		c.SetCD(core.ActionSkill, 900)
		return f + hold, a + hold
	}

	switch c.eCharge {
	case c.eChargeMax:
		c.Core.Log.Debugw("amber bunny at max charge, queuing next recovery", "frame", c.Core.F, "event", core.LogCharacterEvent, "recover at", c.Core.F+721)
		c.eNextRecover = c.Core.F + 721
		c.AddTask(c.recoverCharge(c.Core.F), "charge", 720)
		c.eTickSrc = c.Core.F
	case 1:
		c.SetCD(core.ActionSkill, c.eNextRecover)
	}
	c.eCharge--

	return f + hold, a + hold
}

func (c *char) recoverCharge(src int) func() {
	return func() {
		if c.eTickSrc != src {
			c.Core.Log.Debugw("amber bunny recovery function ignored, src diff", "frame", c.Core.F, "event", core.LogCharacterEvent, "src", src, "new src", c.eTickSrc)
			return
		}
		c.eCharge++
		c.Core.Log.Debugw("amber bunny recovering a charge", "frame", c.Core.F, "event", core.LogCharacterEvent, "src", src, "total charge", c.eCharge)
		c.SetCD(core.ActionSkill, 0)
		if c.eCharge >= c.eChargeMax {
			//fully charged
			return
		}
		//other wise restore another charge
		c.Core.Log.Debugw("amber bunny queuing next recovery", "frame", c.Core.F, "event", core.LogCharacterEvent, "src", src, "recover at", c.Core.F+720)
		c.eNextRecover = c.Core.F + 721
		c.AddTask(c.recoverCharge(src), "charge", 720)

	}
}

func (c *char) Burst(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionBurst, p)

	ai := core.AttackInfo{
		Abil:       "Fiery Rain",
		ActorIndex: c.Index,
		AttackTag:  core.AttackTagElementalBurst,
		ICDTag:     core.ICDTagElementalBurst,
		ICDGroup:   core.ICDGroupAmber,
		Element:    core.Pyro,
		Durability: 25,
		Mult:       burstTick[c.TalentLvlSkill()],
	}
	snap := c.Snapshot(&ai)

	//2sec duration, tick every .4 sec in zone 1
	//2sec duration, tick every .6 sec in zone 2
	//2sec duration, tick every .2 sec in zone 3

	//TODO: properly implement random hits and hit box range. right now everything is just radius 3
	for i := f + 24; i < 120+f; i += 24 {
		c.Core.Combat.QueueAttackWithSnap(ai, snap, core.NewDefCircHit(3, false, core.TargettableEnemy), i)
	}

	for i := f + 36; i < 120+f; i += 36 {
		c.Core.Combat.QueueAttackWithSnap(ai, snap, core.NewDefCircHit(3, false, core.TargettableEnemy), i)
	}

	for i := f + 12; i < 120+f; i += 12 {
		c.Core.Combat.QueueAttackWithSnap(ai, snap, core.NewDefCircHit(3, false, core.TargettableEnemy), i)
	}

	if c.Base.Cons == 6 {
		for _, active := range c.Core.Chars {
			val := make([]float64, core.EndStatType)
			val[core.ATKP] = 0.15
			active.AddMod(core.CharStatMod{
				Key:    "amber-c6",
				Amount: func(a core.AttackTag) ([]float64, bool) { return val, true },
				Expiry: c.Core.F + 900,
			})
			c.Core.Log.Debugw("c6 - adding atk %", "frame", c.Core.F, "event", core.LogCharacterEvent, "character", c.Name())
		}
	}

	c.ConsumeEnergy(64)
	c.SetCD(core.ActionBurst, 720)
	return f, a
}

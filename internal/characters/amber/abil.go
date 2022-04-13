package amber

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func (c *char) Attack(p map[string]int) (int, int) {
	travel, ok := p["travel"]
	if !ok {
		travel = 10
	}

	f, a := c.ActionFrames(core.ActionAttack, p)

	ai := core.AttackInfo{
		Abil:       fmt.Sprintf("Normal %v", c.NormalCounter),
		ActorIndex: c.Index,
		AttackTag:  core.AttackTagNormal,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupAmber,
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
		travel = 10
	}
	weakspot := p["weakspot"]

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
		AttackTag:    core.AttackTagExtra,
		ICDTag:       core.ICDTagExtraAttack,
		ICDGroup:     core.ICDGroupAmber,
		Element:      core.Pyro,
		Durability:   50,
		Mult:         aim[c.TalentLvlAttack()],
		HitWeakPoint: weakspot == 1,
	}
	c.Core.Combat.QueueAttack(ai, core.NewDefSingleTarget(1, core.TargettableEnemy), f, f+travel, c.a4)

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

	c.SetCD(core.ActionSkill, 720)

	return f + hold, a + hold
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
				Amount: func() ([]float64, bool) { return val, true },
				Expiry: c.Core.F + 900,
			})
			c.Core.Log.NewEvent("c6 - adding atk %", core.LogCharacterEvent, c.Index, "character", c.Name())
		}
	}

	c.ConsumeEnergy(64)
	c.SetCDWithDelay(core.ActionBurst, 720, 64)
	return f, a
}

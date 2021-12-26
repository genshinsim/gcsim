package ayaka

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func (c *char) Attack(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionAttack, p)

	ai := core.AttackInfo{
		Abil:       fmt.Sprintf("Normal %v", c.NormalCounter),
		ActorIndex: c.Index,
		AttackTag:  core.AttackTagNormal,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Physical,
		Durability: 25,
	}

	for i, mult := range attack[c.NormalCounter] {
		ai.Mult = mult[c.TalentLvlAttack()]
		c.Core.Combat.QueueAttack(ai, core.NewDefSingleTarget(1, core.TargettableEnemy), 0, f-5+i)
	}

	c.AdvanceNormalIndex()

	return f, a
}

func (c *char) ChargeAttack(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionCharge, p)

	ai := core.AttackInfo{
		Abil:       "Charge",
		ActorIndex: c.Index,
		AttackTag:  core.AttackTagExtra,
		ICDTag:     core.ICDTagExtraAttack,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Physical,
		Durability: 25,
		Mult:       ca[c.TalentLvlAttack()],
	}

	for i := 0; i < 3; i++ {
		c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(2, false, core.TargettableEnemy), 0, f-3+i)
	}

	return f, a
}

func (c *char) Dash(p map[string]int) (int, int) {
	f, ok := p["f"]
	if !ok {
		f = 36
	}
	//no dmg attack at end of dash
	ai := core.AttackInfo{
		Abil:       "Dash",
		ActorIndex: c.Index,
		AttackTag:  core.AttackTagNone,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Cryo,
		Durability: 25,
	}

	//restore on hit, once per attack
	once := false
	cb := func(a core.AttackCB) {
		if once {
			return
		}

		c.Core.RestoreStam(10)
		val := make([]float64, core.EndStatType)
		val[core.CryoP] = 0.18
		//a2 increase normal + ca dmg by 30% for 6s
		c.AddMod(core.CharStatMod{
			Key:    "ayaka-a4",
			Expiry: c.Core.F + 600,
			Amount: func(a core.AttackTag) ([]float64, bool) {
				return val, true
			},
		})
		once = true
	}
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(2, false, core.TargettableEnemy), 0, f, cb)
	//add cryo infuse
	c.AddWeaponInfuse(core.WeaponInfusion{
		Key:    "ayaka-dash",
		Ele:    core.Cryo,
		Tags:   []core.AttackTag{core.AttackTagNormal, core.AttackTagExtra, core.AttackTagPlunge},
		Expiry: c.Core.F + 300,
	})
	return f, f
}

func (c *char) Skill(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionSkill, p)
	ai := core.AttackInfo{
		Abil:       "Hyouka",
		ActorIndex: c.Index,
		AttackTag:  core.AttackTagElementalArt,
		ICDTag:     core.ICDTagElementalArt,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Cryo,
		Durability: 50,
		Mult:       skill[c.TalentLvlSkill()],
	}

	//2 or 3 1:1 ratio
	count := 4
	if c.Core.Rand.Float64() < 0.5 {
		count = 5
	}
	c.QueueParticle("ayaka", count, core.Cryo, f+100)

	//a2 increase normal + ca dmg by 30% for 6s
	c.AddMod(core.CharStatMod{
		Key:    "ayaka-a2",
		Expiry: c.Core.F + 360,
		Amount: func(a core.AttackTag) ([]float64, bool) {

			val := make([]float64, core.EndStatType)
			val[core.DmgP] = 0.3
			return val, a == core.AttackTagNormal || a == core.AttackTagExtra
		},
	})

	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(4, false, core.TargettableEnemy), 0, f)

	c.SetCD(core.ActionSkill, 600)
	return f, a

}

func (c *char) Burst(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionBurst, p)

	ai := core.AttackInfo{
		Abil:       "Soumetsu",
		ActorIndex: c.Index,
		AttackTag:  core.AttackTagElementalBurst,
		ICDTag:     core.ICDTagElementalBurst,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Cryo,
		Durability: 25,
	}

	//5 second, 20 ticks, so once every 15 frames, bloom after 5 seconds
	ai.Mult = burstBloom[c.TalentLvlBurst()]
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(5, false, core.TargettableEnemy), f, f+300)

	ai.Mult = burstCut[c.TalentLvlBurst()]
	for i := 0; i < 19; i++ {
		c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(5, false, core.TargettableEnemy), f, f+i*15)
	}

	c.SetCD(core.ActionBurst, 20*60)
	c.ConsumeEnergy(13)

	return f, a
}

package itto

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func (c *char) Attack(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionAttack, p)

	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       fmt.Sprintf("Normal %v", c.NormalCounter),
		AttackTag:  core.AttackTagNormal,
		ICDTag:     core.ICDTagNormalAttack,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeBlunt,
		Element:    core.Physical,
		Durability: 25,
		Mult:       attack[c.NormalCounter][c.TalentLvlAttack()],
	}
	// guess i don't have to worry about missing because you can't miss in gcsim?
	if c.Core.Status.Duration("unga-bunga") > 0 {
		ai.Element = core.Geo //is this overridable?
		if c.NormalCounter == 0 || c.NormalCounter == 2 {
			if c.skillStacks < 5 {
				c.skillStacks++
			}
		}
	}
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(1, false, core.TargettableEnemy), f, f)
	c.AdvanceNormalIndex()

	return f, a
}

func (c *char) ChargeAttack(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionCharge, p)
	chargedIndex := 0
	if c.skillStacks == 1 {
		chargedIndex = 2
		c.skillStacks--
	} else if c.skillStacks > 1 {
		chargedIndex = 1
		c.skillStacks--
	}

	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Charge",
		AttackTag:  core.AttackTagExtra,
		ICDTag:     core.ICDTagExtraAttack,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeBlunt,
		Element:    core.Physical,
		Durability: 25,
		Mult:       charged[chargedIndex][c.TalentLvlAttack()],
	}
	if c.Core.Status.Duration("unga-bunga") > 0 {
		ai.Element = core.Geo //is this overridable?
	}
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(1, false, core.TargettableEnemy), f, f)
	return f, a
}

func (c *char) Skill(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionSkill, p)

	// no idea
	orb := 1
	if c.Core.Rand.Float64() < 0.33 {
		orb = 2
	}
	c.QueueParticle("Itto", orb, core.Geo, f+60)

	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Masatsu Zetsugi: Akaushi Burst!",
		AttackTag:  core.AttackTagElementalArt,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeBlunt,
		Element:    core.Geo,
		Durability: 25,
		Mult:       skill[c.TalentLvlSkill()],
	}

	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(2, false, core.TargettableEnemy), f-5, f-5)
	c.throwCow()

	c.SetCD(core.ActionSkill, 10*60)
	return f, a
}

func (c *char) throwCow() {

	dur := 6 * 60

	//create a construct
	c.Core.Constructs.New(c.newCow(dur), true) //30 seconds

	c.Core.Log.Debugw("Ushi thrown", "frame", c.Core.F, "event", core.LogCharacterEvent, "expected end", c.Core.F+dur)

	c.Core.Status.AddStatus("ittoUshi", dur)

}

func (c *char) Burst(p map[string]int) (int, int) {
	dur := 11 * 60

	c.Core.Status.AddStatus("unga-bunga", dur)
	f, a := c.ActionFrames(core.ActionBurst, p)

	//enhance weapon for 12 seconds
	// Infusion starts when burst starts and ends when burst comes off CD - check any diluc video
	c.AddWeaponInfuse(core.WeaponInfusion{
		Key:    "itto-geo-weapon",
		Ele:    core.Geo,
		Tags:   []core.AttackTag{core.AttackTagNormal, core.AttackTagExtra, core.AttackTagPlunge},
		Expiry: c.Core.F + dur, //with a4
	})

	// decrease res
	val := make([]float64, core.EndStatType)
	val[core.ATK] = burstpp[c.TalentLvlBurst()] * c.Stat(core.DEF)
	c.AddMod(core.CharStatMod{
		Key:    "itto-unga-bunga",
		Amount: func() ([]float64, bool) { return val, true },
		Expiry: c.Core.F + dur,
	})

	c.ConsumeEnergy(24)
	c.SetCD(core.ActionBurst, 720)
	return f, a
}

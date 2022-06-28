package traveleranemo

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
		StrikeType: core.StrikeTypeSlash,
		Element:    core.Physical,
		Durability: 25,
		Mult:       attack[c.NormalCounter][c.TalentLvlAttack()],
	}

	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(0.3, false, core.TargettableEnemy), f-1, f-1)

	if c.NormalCounter == c.NormalHitNum-1 {
		//add 60% as anemo dmg
		ai := core.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "A1",
			AttackTag:  core.AttackTagNormal,
			ICDTag:     core.ICDTagNone,
			ICDGroup:   core.ICDGroupDefault,
			StrikeType: core.StrikeTypeDefault, // TODO: I don't know what strike type this is?
			Element:    core.Anemo,
			Durability: 25,
			Mult:       0.6,
		}
		c.Core.Combat.QueueAttack(ai, core.NewDefBoxHit(1, 3, false, core.TargettableEnemy), f-1, f-1)
	}

	c.AdvanceNormalIndex()

	return f, a
}

func (c *char) Skill(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionSkill, p)

	holdTicks := 0
	if p["hold"] == 1 {
		holdTicks = 6
	}
	if 0 < p["hold_ticks"] && p["hold_ticks"] <= 6 {
		holdTicks = p["hold_ticks"]
	}
	if holdTicks == 0 {
		hitmark := 34
		ai := core.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Palm Vortex (Tap)",
			AttackTag:  core.AttackTagElementalArt,
			ICDTag:     core.ICDTagElementalArt,
			ICDGroup:   core.ICDGroupDefault,
			Element:    core.Anemo,
			Durability: 25,
			Mult:       skillInitialStorm[c.TalentLvlSkill()],
		}

		c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(2, false, core.TargettableEnemy), hitmark, hitmark)

		c.QueueParticle(c.Name(), 2, core.Anemo, hitmark+90)
		c.SetCDWithDelay(core.ActionSkill, 5*60, hitmark-5)
	} else {
		c.eInfuse = core.NoElement
		c.eICDTag = core.ICDTagNone
		aiCut := core.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Palm Vortex Initial Cutting (Hold)",
			AttackTag:  core.AttackTagElementalArt,
			ICDTag:     core.ICDTagElementalArt,
			ICDGroup:   core.ICDGroupDefault,
			Element:    core.Anemo,
			Durability: 25,
			Mult:       skillInitialCutting[c.TalentLvlSkill()],
		}
		aiCutAbs := core.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Palm Vortex Initial Cutting Absorbed (Hold)",
			AttackTag:  core.AttackTagElementalArt,
			ICDTag:     core.ICDTagNone,
			ICDGroup:   core.ICDGroupDefault,
			Element:    core.NoElement,
			Durability: 25,
			Mult:       skillInitialCuttingAbsorb[c.TalentLvlSkill()],
		}

		aiMaxCutAbs := core.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Palm Vortex Max Cutting Absorbed (Hold)",
			AttackTag:  core.AttackTagElementalArt,
			ICDTag:     core.ICDTagNone,
			ICDGroup:   core.ICDGroupDefault,
			Element:    core.NoElement,
			Durability: 25,
			Mult:       skillMaxCuttingAbsorb[c.TalentLvlSkill()],
		}
		hitmark := 31
		for i := 0; i < holdTicks; i += 1 {

			c.Core.Combat.QueueAttack(aiCut, core.NewDefCircHit(1, false, core.TargettableEnemy), hitmark, hitmark)
			if i > 1 {
				c.AddTask(func() {
					if c.eInfuse != core.NoElement {
						aiMaxCutAbs.Element = c.eInfuse
						aiMaxCutAbs.ICDTag = c.eICDTag
						c.Core.Combat.QueueAttack(aiMaxCutAbs, core.NewDefCircHit(1.5, false, core.TargettableEnemy), 0, 0)
					}
					//check if infused
				}, "amc-e-cutting-absorb", hitmark)
			} else {
				c.AddTask(func() {
					if c.eInfuse != core.NoElement {
						aiCutAbs.Element = c.eInfuse
						aiCutAbs.ICDTag = c.eICDTag
						c.Core.Combat.QueueAttack(aiCutAbs, core.NewDefCircHit(1.5, false, core.TargettableEnemy), 0, 0)
					}
					//check if infused
				}, "amc-e-cutting-absorb", hitmark)
			}

			hitmark += 15
			if i == 1 {
				aiCut.Mult = skillMaxCutting[c.TalentLvlSkill()]
				aiCut.Abil = "Palm Vortex Max Cutting (Hold)"

				// there is a 5 frame delay when it shifts from initial to max
				hitmark += 5
			}
		}
		// move the hitmark back by 1 tick (15f) then forward by 5f for the Storm damage
		hitmark = hitmark - 15 + 5
		aiStorm := core.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Palm Vortex Initial Storm (Hold)",
			AttackTag:  core.AttackTagElementalArt,
			ICDTag:     core.ICDTagElementalArt,
			ICDGroup:   core.ICDGroupDefault,
			Element:    core.Anemo,
			Durability: 25,
			Mult:       skillInitialStorm[c.TalentLvlSkill()],
		}
		aiStormAbs := core.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Palm Vortex Initial Storm Absorbed (Hold)",
			AttackTag:  core.AttackTagElementalArt,
			ICDTag:     core.ICDTagNone,
			ICDGroup:   core.ICDGroupDefault,
			Element:    core.Anemo,
			Durability: 25,
			Mult:       skillInitialStormAbsorb[c.TalentLvlSkill()],
		}

		// it does max storm when there are 2 or more ticks
		if holdTicks >= 2 {
			aiStorm.Mult = skillMaxStorm[c.TalentLvlSkill()]
			aiStorm.Abil = "Palm Vortex Max Storm (Hold)"

			aiStormAbs.Mult = skillMaxStormAbsorb[c.TalentLvlSkill()]
			aiStormAbs.Abil = "Palm Vortex Max Storm Absorbed (Hold)"

			count := 3
			if c.Core.Rand.Float64() < 0.33 {
				count = 4
			}
			c.QueueParticle(c.Name(), count, core.Anemo, hitmark+90)
			c.SetCDWithDelay(core.ActionSkill, 8*60, hitmark-5)
		} else {
			c.QueueParticle(c.Name(), 2, core.Anemo, hitmark+90)
			c.SetCDWithDelay(core.ActionSkill, 5*60, hitmark-5)
		}

		c.Core.Combat.QueueAttack(aiCut, core.NewDefCircHit(2, false, core.TargettableEnemy), hitmark, hitmark)
		c.AddTask(func() {
			if c.eInfuse != core.NoElement {
				aiStormAbs.Element = c.eInfuse
				aiStormAbs.ICDTag = c.eICDTag
				c.Core.Combat.QueueAttack(aiStormAbs, core.NewDefCircHit(1.5, false, core.TargettableEnemy), 0, 0)
			}
			//check if infused
		}, "amc-e-storm-absorb", hitmark)

		// starts absorbing after the first tick?
		c.AddTask(c.absorbCheckE(c.Core.F, 0, int((hitmark)/18)), "amc-e-absorb-check", 32)
	}

	return f, a
}

func (c *char) Burst(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionBurst, p)

	c.qInfuse = core.NoElement
	c.qICDTag = core.ICDTagNone
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Gust Surge",
		AttackTag:  core.AttackTagElementalBurst,
		ICDTag:     core.ICDTagVentiBurstAnemo,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Anemo,
		Durability: 25,
		Mult:       burstDot[c.TalentLvlBurst()],
	}

	aiAbsorb := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Gust Surge (Absorbed)",
		AttackTag:  core.AttackTagElementalBurst,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.NoElement,
		Durability: 50,
		Mult:       burstAbsorbDot[c.TalentLvlBurst()],
	}

	// snapshot is on cast?
	var snap core.Snapshot
	c.AddTask(func() {
		snap = c.Snapshot(&ai)
	}, "amc-q-snapshot", 1)

	var cb core.AttackCBFunc
	if c.Base.Cons >= 6 {
		cb = c6cb(core.Anemo)
	}

	//First hit 94f, then 30f between hits. Max 9 hits total
	for i := 0; i < 9; i += 1 {
		c.Core.Combat.QueueAttackWithSnap(ai, snap, core.NewDefCircHit(2, false, core.TargettableEnemy), 94+30*i, cb)

		c.AddTask(func() {
			if c.qInfuse != core.NoElement {
				aiAbsorb.Element = c.qInfuse
				aiAbsorb.ICDTag = c.qICDTag
				var cbAbs core.AttackCBFunc
				if c.Base.Cons >= 6 {
					cbAbs = c6cb(c.qInfuse)
				}
				c.Core.Combat.QueueAttackWithSnap(aiAbsorb, snap, core.NewDefCircHit(1.5, false, core.TargettableEnemy), 0, cbAbs)
			}
			//check if infused
		}, "amc-q-absorb", 94+30*i)
	}

	//it absorbs before the first hit.
	c.AddTask(c.absorbCheckQ(c.Core.F, 0, int((94+8*30)/18)), "amc-q-absorb-check", 39)

	c.SetCDWithDelay(core.ActionBurst, 15*60, 2)
	c.ConsumeEnergy(8)
	return f, a
}

func (c *char) absorbCheckQ(src, count, max int) func() {
	return func() {
		if count == max {
			return
		}
		c.qInfuse = c.Core.AbsorbCheck(c.infuseCheckLocation, core.Cryo, core.Pyro, core.Hydro, core.Electro)
		switch c.qInfuse {
		case core.Cryo:
			c.qICDTag = core.ICDTagVentiBurstCryo
		case core.Pyro:
			c.qICDTag = core.ICDTagVentiBurstPyro
		case core.Electro:
			c.qICDTag = core.ICDTagVentiBurstElectro
		case core.Hydro:
			c.qICDTag = core.ICDTagVentiBurstHydro
		case core.NoElement:
			c.AddTask(c.absorbCheckQ(src, count+1, max), "amc-q-absorb-check", 18)
		}
	}
}

func (c *char) absorbCheckE(src, count, max int) func() {
	return func() {
		if count == max {
			return
		}
		c.eInfuse = c.Core.AbsorbCheck(c.infuseCheckLocation, core.Cryo, core.Pyro, core.Hydro, core.Electro)
		switch c.eInfuse {
		case core.Cryo:
			c.eICDTag = core.ICDTagSayuSkillCryo
		case core.Pyro:
			c.eICDTag = core.ICDTagSayuSkillPyro
		case core.Electro:
			c.eICDTag = core.ICDTagSayuSkillElectro
		case core.Hydro:
			c.eICDTag = core.ICDTagSayuSkillHydro
		case core.NoElement:
			c.AddTask(c.absorbCheckE(src, count+1, max), "amc-e-absorb-check", 18)
		}
	}
}

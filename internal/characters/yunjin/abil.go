package yunjin

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/tmpl/shield"
	"github.com/genshinsim/gcsim/pkg/core"
)

// Normal attack damage queue generator
// Very standard
func (c *char) Attack(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionAttack, p)
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       fmt.Sprintf("Normal %v", c.NormalCounter),
		AttackTag:  core.AttackTagNormal,
		ICDTag:     core.ICDTagNormalAttack,
		ICDGroup:   core.ICDGroupPole,
		Element:    core.Physical,
		Durability: 25,
	}

	for _, mult := range attack[c.NormalCounter] {
		ai.Mult = mult[c.TalentLvlAttack()]
		c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(0.1, false, core.TargettableEnemy), f, f)
	}

	c.AdvanceNormalIndex()

	return f, a
}

// Charge attack damage queue generator
// Very standard - consistent with other characters like Xiangling
func (c *char) ChargeAttack(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionCharge, p)

	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Charge",
		AttackTag:  core.AttackTagExtra,
		ICDTag:     core.ICDTagExtraAttack,
		ICDGroup:   core.ICDGroupPole,
		Element:    core.Physical,
		Durability: 25,
		Mult:       charge[c.TalentLvlAttack()],
	}
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(0.1, false, core.TargettableEnemy), f, f)

	//return animation cd
	return f, a
}

// Skill - modelled after Beidou E
// Has two parameters:
// perfect = 1 if you are doing a perfect counter
// hold = 1 or 2 for regular charging up to level 1 or 2
func (c *char) Skill(p map[string]int) (int, int) {
	// Hold parameter gets used in action frames to get earliest possible release frame
	f, a := c.ActionFrames(core.ActionSkill, p)

	chargeLevel := 0
	if p["perfect"] == 1 {
		chargeLevel = 2
	} else {
		chargeLevel = p["hold"]
	}

	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Opening Flourish Press (E)",
		AttackTag:  core.AttackTagElementalArt,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeBlunt,
		Element:    core.Geo,
		Durability: 50,
		Mult:       skillDmg[chargeLevel][c.TalentLvlSkill()],
	}

	ai.UseDef = true

	// TODO: Fix hit frames when known
	// Particle should spawn after hit
	hitDelay := f
	switch chargeLevel {
	case 0:
		c.QueueParticle("yunjin", 2, core.Geo, 100+hitDelay)
	case 1:
		// Currently believed to be 2-3 particles with the ratio 3:2
		if c.Core.Rand.Float64() < .6 {
			c.QueueParticle("yunjin", 2, core.Geo, 100+hitDelay)
		} else {
			c.QueueParticle("yunjin", 3, core.Geo, 100+hitDelay)
		}
		ai.Abil = "Opening Flourish Level 1 (E)"
	case 2:
		c.QueueParticle("yunjin", 3, core.Geo, 100+hitDelay)
		ai.Durability = 100
		ai.Abil = "Opening Flourish Level 2 (E)"
	}

	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(1, false, core.TargettableEnemy), hitDelay, hitDelay)

	// Add shield until skill unleashed (treated as frame when attack hits)
	c.Core.Shields.Add(&shield.Tmpl{
		Src:        c.Core.F,
		ShieldType: core.ShieldYunjinSkill,
		HP:         skillShieldPct[c.TalentLvlSkill()]*c.MaxHP() + skillShieldFlat[c.TalentLvlSkill()],
		Ele:        core.Geo,
		Expires:    c.Core.F + f,
	})

	if c.Base.Cons >= 1 {
		// 18% doesn't result in a whole number - 442.8 frames. We round up
		c.SetCD(core.ActionSkill, 443)
	} else {
		c.SetCD(core.ActionSkill, 9*60)
	}

	return f, a
}

// Burst - The main buff effects are handled in a separate function
func (c *char) Burst(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionBurst, p)

	// AoE Geo damage
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Cliffbreaker's Banner",
		AttackTag:  core.AttackTagElementalBurst,
		ICDTag:     core.ICDTagElementalBurst,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeBlunt,
		Element:    core.Geo,
		Durability: 50,
		Mult:       burstDmg[c.TalentLvlBurst()],
	}
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(1, false, core.TargettableEnemy), f, f)

	c.Core.Status.AddStatus("yunjinburst", 12*60)

	// Reset number of burst triggers to 30
	for i := range c.burstTriggers {
		c.burstTriggers[i] = 30
		c.updateBuffTags()
	}

	// TODO: Need to obtain exact timing of the 12s. Currently assume that it starts when burst is used
	if c.Base.Cons >= 2 {
		val := make([]float64, core.EndStatType)
		val[core.DmgP] = .15
		for _, char := range c.Core.Chars {
			char.AddPreDamageMod(core.PreDamageMod{
				Key: "yunjin-c2",
				Amount: func(ae *core.AttackEvent, t core.Target) ([]float64, bool) {
					if ae.Info.AttackTag == core.AttackTagNormal {
						return val, true
					}
					return nil, false
				},
				Expiry: c.Core.F + 12*60,
			})
		}
	}

	if c.Base.Cons >= 6 {
		val := make([]float64, core.EndStatType)
		val[core.AtkSpd] = .12
		for _, char := range c.Core.Chars {
			char.AddMod(core.CharStatMod{
				Key:    "yunjin-c6",
				Expiry: c.Core.F + 12*60,
				Amount: func() ([]float64, bool) {
					if c.Core.Status.Duration("yunjinburst") == 0 {
						return nil, false
					}
					return val, true
				},
			})
		}
	}

	c.ConsumeEnergy(8)
	c.SetCDWithDelay(core.ActionBurst, 15*60, 8)

	return f, a
}

func (c *char) burstProc() {
	// Add Flying Cloud Flag Formation as a pre-damage hook
	c.Core.Events.Subscribe(core.OnAttackWillLand, func(args ...interface{}) bool {
		ae := args[1].(*core.AttackEvent)

		if ae.Info.AttackTag != core.AttackTagNormal {
			return false
		}
		if c.Core.Status.Duration("yunjinburst") == 0 || c.burstTriggers[ae.Info.ActorIndex] == 0 {
			return false
		}

		finalBurstBuff := burstBuff[c.TalentLvlBurst()]
		if c.partyElementalTypes == 4 {
			finalBurstBuff += .115
		} else {
			finalBurstBuff += 0.025 * float64(c.partyElementalTypes)
		}

		// ai := core.AttackInfo{
		// 	Abil:      "Yunjin Burst Buff",
		// 	AttackTag: core.AttackTagNone,
		// }
		stats, _ := c.SnapshotStats()
		dmgAdded := (c.Base.Def*(1+stats[core.DEFP]) + stats[core.DEF]) * finalBurstBuff
		ae.Info.FlatDmg += dmgAdded

		c.burstTriggers[ae.Info.ActorIndex]--
		c.updateBuffTags()

		c.Core.Log.NewEvent("yunjin burst adding damage", core.LogPreDamageMod, ae.Info.ActorIndex, "damage_added", dmgAdded, "stacks_remaining_for_char", c.burstTriggers[ae.Info.ActorIndex], "burst_def_pct", finalBurstBuff)

		return false
	}, "yunjin-burst")
}

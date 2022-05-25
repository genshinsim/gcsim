package tartaglia

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

// Normal attack
// Perform up to 6 consecutive shots with a bow.
func (c *char) Attack(p map[string]int) (int, int) {
	travel, ok := p["travel"]
	if !ok {
		travel = 10
	}

	f, a := c.ActionFrames(core.ActionAttack, p)

	if c.Core.Status.Duration("tartagliamelee") > 0 {
		return c.meleeAttack(f, a)
	}
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       fmt.Sprintf("Normal %v", c.NormalCounter),
		AttackTag:  core.AttackTagNormal,
		ICDTag:     core.ICDTagNormalAttack,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypePierce,
		Element:    core.Physical,
		Durability: 25,
		Mult:       attack[c.NormalCounter][c.TalentLvlAttack()],
	}
	// TODO - double check this snapshotDelay
	c.Core.Combat.QueueAttack(ai, core.NewDefSingleTarget(1, core.TargettableEnemy), f, f+travel)

	c.AdvanceNormalIndex()

	return f, a
}

// random delayed numbers
var meleeDelayOffset = [][]int{
	{0},
	{0},
	{0},
	{0},
	{0},
	{1, 0},
}

// Melee stance attack.
// Perform up to 6 consecutive Hydro strikes.
func (c *char) meleeAttack(f, a int) (int, int) {
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       fmt.Sprintf("Normal %v", c.NormalCounter),
		AttackTag:  core.AttackTagNormal,
		ICDTag:     core.ICDTagNormalAttack,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeSlash,
		Element:    core.Hydro,
		Durability: 25,
	}
	for i, mult := range eAttack[c.NormalCounter] {
		ai.Mult = mult[c.TalentLvlSkill()]
		delay := f - meleeDelayOffset[c.NormalCounter][i]
		c.Core.Combat.QueueAttack(
			ai,
			core.NewDefCircHit(.5, false, core.TargettableEnemy),
			delay,
			delay,
			//TODO: what's the ordering on these 2 callbacks?
			c.meleeApplyRiptide, //call back for applying riptide
			c.rtSlashCallback,   //call back for triggering slash
		)
	}

	c.AdvanceNormalIndex()
	return f, a
}

//Once fully charged, deal Hydro DMG and apply the Riptide status.
func (c *char) Aimed(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionAim, p)

	travel, ok := p["travel"]
	if !ok {
		travel = 10
	}
	weakspot, ok := p["weakspot"]

	ai := core.AttackInfo{
		ActorIndex:   c.Index,
		Abil:         "Aim (Charged)",
		AttackTag:    core.AttackTagExtra,
		ICDTag:       core.ICDTagNone,
		ICDGroup:     core.ICDGroupDefault,
		StrikeType:   core.StrikeTypePierce,
		Element:      core.Hydro,
		Durability:   25,
		Mult:         aim[c.TalentLvlAttack()],
		HitWeakPoint: weakspot == 1,
	}

	c.Core.Combat.QueueAttack(
		ai,
		core.NewDefSingleTarget(1, core.TargettableEnemy),
		f,
		f+travel,
		//TODO: what's the ordering on these 2 callbacks?
		c.rtFlashCallback,   //call back for triggering slash
		c.aimedApplyRiptide, //call back for applying riptide
	)

	return f, a
}

var meleeChargeDelayOffset = []int{
	2, 0,
}

// since E is aoe, so this should be considered aoe too
// hitWeakPoint: tartaglia can proc Prototype Cresent's Passive on Geovishap's weakspots.
// Evidence: https://youtu.be/oOfeu5pW0oE
func (c *char) ChargeAttack(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionCharge, p)

	if c.Core.Status.Duration("tartagliamelee") == 0 {
		return f, a
	}

	hitWeakPoint, ok := p["hitWeakPoint"]
	if !ok {
		hitWeakPoint = 0
	}

	ai := core.AttackInfo{
		ActorIndex:   c.Index,
		Abil:         "Charged Attack",
		AttackTag:    core.AttackTagExtra,
		ICDTag:       core.ICDTagExtraAttack,
		ICDGroup:     core.ICDGroupDefault,
		StrikeType:   core.StrikeTypeSlash,
		Element:      core.Hydro,
		Durability:   25,
		HitWeakPoint: hitWeakPoint != 0,
	}
	for i, mult := range eCharge {
		ai.Mult = mult[c.TalentLvlSkill()]
		c.Core.Combat.QueueAttack(
			ai,
			core.NewDefCircHit(1, false, core.TargettableEnemy),
			f-meleeChargeDelayOffset[i],
			f-meleeChargeDelayOffset[i],
			//TODO: what's the ordering on these 2 callbacks?
			c.meleeApplyRiptide, //call back for applying riptide
			c.rtSlashCallback,   //call back for triggering slash
		)
	}
	return f, a
}

//Cast: AoE strong hydro damage
//Melee Stance: infuse NA/CA to hydro damage
func (c *char) Skill(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionSkill, p)

	if c.Core.Status.Duration("tartagliamelee") > 0 {
		c.onExitMeleeStance()
		c.ResetNormalCounter()
		return f, a
	}

	c.eCast = c.Core.F
	c.Core.Status.AddStatus("tartagliamelee", 30*60)
	c.Core.Log.NewEvent("Foul Legacy activated", core.LogCharacterEvent, c.Index, "rtexpiry", c.Core.F+30*60)

	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Foul Legacy: Raging Tide",
		AttackTag:  core.AttackTagElementalArt,
		ICDTag:     core.ICDTagNormalAttack,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeDefault,
		Element:    core.Hydro,
		Durability: 50,
		Mult:       skill[c.TalentLvlSkill()],
	}

	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(1, false, core.TargettableEnemy), f, f)

	src := c.eCast
	c.AddTask(func() {
		if src == c.eCast && c.Core.Status.Duration("tartagliamelee") > 0 {
			c.onExitMeleeStance()
			c.ResetNormalCounter()
		}
	}, "tartagliamelee-cd", 30*60)
	c.SetCD(core.ActionSkill, 60)

	return f, a
}

func (c *char) onExitMeleeStance() {
	// Precise skill CD from Risuke:
	// Aligns with separate table on wiki except the 4 second duration one
	// https://discord.com/channels/763583452762734592/851428030094114847/899416824117084210
	// https://media.discordapp.net/attachments/778615842916663357/781978094495727646/unknown-20.png

	skillCD := 0

	switch timeInMeleeStance := c.Core.F - c.eCast; {
	case timeInMeleeStance < 2*60:
		skillCD = 7 * 60
	case 2*60 <= timeInMeleeStance && timeInMeleeStance < 4*60:
		skillCD = 8 * 60
	case 4*60 <= timeInMeleeStance && timeInMeleeStance < 5*60:
		skillCD = 9 * 60
	case 5*60 <= timeInMeleeStance && timeInMeleeStance < 8*60:
		skillCD = 5*60 + timeInMeleeStance
	case 8*60 <= timeInMeleeStance && timeInMeleeStance < 30*60:
		skillCD = 6*60 + timeInMeleeStance
	case timeInMeleeStance >= 30*60:
		skillCD = 45 * 60
	}

	if c.Base.Cons >= 1 {
		skillCD = int(float64(skillCD) * 0.8)
	}

	if c.mlBurstUsed {
		c.ResetActionCooldown(core.ActionSkill)
		c.mlBurstUsed = false
	} else {
		c.SetCD(core.ActionSkill, skillCD)
	}
	c.Core.Status.DeleteStatus("tartagliamelee")
}

//Performs a different attack depending on the stance in which it is cast.
//Ranged Stance: dealing AoE Hydro DMG. Apply Riptide status to enemies hit. Returns 20 Energy after use
//Melee Stance: dealing AoE Hydro DMG. Triggers Riptide Blast (clear riptide after triggering riptide blast)
func (c *char) Burst(p map[string]int) (int, int) {
	mult := burst[c.TalentLvlBurst()]

	f, a := c.ActionFrames(core.ActionBurst, p)

	skillName := "Ranged Stance: Flash of Havoc"
	if c.Core.Status.Duration("tartagliamelee") > 0 {
		skillName = "Melee Stance: Light of Obliteration"
		mult = meleeBurst[c.TalentLvlBurst()]
	}

	c.AddTask(func() {
		ai := core.AttackInfo{
			ActorIndex: c.Index,
			Abil:       skillName,
			AttackTag:  core.AttackTagElementalBurst,
			ICDTag:     core.ICDTagNone,
			ICDGroup:   core.ICDGroupDefault,
			StrikeType: core.StrikeTypeDefault,
			Element:    core.Hydro,
			Durability: 50,
			Mult:       mult,
		}
		var cb core.AttackCBFunc
		if c.Core.Status.Duration("tartagliamelee") > 0 {
			cb = c.rtBlastCallback
		} else {
			cb = c.rangedBurstApplyRiptide
		}

		c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(5, false, core.TargettableEnemy), 0, 0, cb)
		if c.Core.Status.Duration("tartagliamelee") > 0 {
			if c.Base.Cons >= 6 {
				c.mlBurstUsed = true
			}
		} else {
			c.AddTask(func() {
				c.AddEnergy("tartaglia-ranged-burst-refund", 20)
			}, "tartaglia-ranged-burst-energy-refund", 9)
		}
	}, "tartaglia-burst-clear", f-5) //random 5 frame

	if c.Core.Status.Duration("tartagliamelee") == 0 {
		c.ConsumeEnergy(8)
		c.SetCDWithDelay(core.ActionBurst, 900, 8)
	} else {
		c.ConsumeEnergy(75)
		c.SetCDWithDelay(core.ActionBurst, 900, 75)
	}

	return f, a
}

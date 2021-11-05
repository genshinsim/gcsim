package tartaglia

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/core"
)

// Normal attack
// Perform up to 6 consecutive shots with a bow.
func (c *char) Attack(p map[string]int) (int, int) {
	travel, ok := p["travel"]
	if !ok {
		travel = 20
	}

	f, a := c.ActionFrames(core.ActionAttack, p)

	if c.Core.Status.Duration("tartagliamelee") > 0 {
		return c.meleeAttack(f, a)
	}

	d := c.Snapshot(
		fmt.Sprintf("Normal %v", c.NormalCounter),
		core.AttackTagNormal,
		core.ICDTagNormalAttack,
		core.ICDGroupDefault,
		core.StrikeTypePierce,
		core.Physical,
		25,
		attack[c.NormalCounter][c.TalentLvlAttack()],
	)

	c.QueueDmg(&d, f+travel)

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
	for i, mult := range eAttack[c.NormalCounter] {
		c.AddTask(func() {
			d := c.Snapshot(
				fmt.Sprintf("Normal %v", c.NormalCounter),
				core.AttackTagNormal,
				core.ICDTagNormalAttack,
				core.ICDGroupDefault,
				core.StrikeTypeSlash,
				core.Hydro,
				25,
				mult[c.TalentLvlSkill()],
			)
			d.OnHitCallback = c.rtSlashCallback
			c.Core.Combat.ApplyDamage(&d)
		}, "tartaglia-attack", f-meleeDelayOffset[c.NormalCounter][i])
	}

	c.AdvanceNormalIndex()
	return f, a
}

//Once fully charged, deal Hydro DMG and apply the Riptide status.
func (c *char) Aimed(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionAim, p)

	travel, ok := p["travel"]
	if !ok {
		travel = 20
	}
	hitWeakPoint, ok := p["hitWeakPoint"]
	if !ok {
		hitWeakPoint = 0
	}

	d := c.Snapshot(
		"Aim (Charged)",
		core.AttackTagExtra,
		core.ICDTagNone,
		core.ICDGroupDefault,
		core.StrikeTypePierce,
		core.Hydro,
		25,
		aim[c.TalentLvlAttack()],
	)

	if hitWeakPoint != 0 {
		d.HitWeakPoint = true
	}
	d.AnimationFrames = f
	d.OnHitCallback = c.rtFlashCallback

	c.QueueDmg(&d, travel+f)

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

	// simply set 2 instances of CA to be snapshotted
	for i, mult := range eCharge {
		d := c.Snapshot(
			"Charged Attack",
			core.AttackTagExtra,
			core.ICDTagExtraAttack,
			core.ICDGroupDefault,
			core.StrikeTypeSlash,
			core.Hydro,
			25,
			mult[c.TalentLvlSkill()],
		)
		if hitWeakPoint != 0 {
			d.HitWeakPoint = true
		}
		d.Targets = core.TargetAll
		d.OnHitCallback = c.rtSlashCallback
		c.QueueDmg(&d, f-meleeChargeDelayOffset[i])
		// c.AddTask(func() {
		// 	c.Core.Combat.ApplyDamage(&d)
		// }, "tartaglia-charge-attack", f-meleeChargeDelayOffset[i])
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
	c.Core.Log.Debugw("Foul Legacy acivated", "frame", c.Core.F, "event", core.LogCharacterEvent, "rtexpiry", c.Core.F+30*60)

	c.AddTask(func() {
		d := c.Snapshot(
			"Foul Legacy: Raging Tide",
			core.AttackTagElementalArt,
			core.ICDTagNormalAttack,
			core.ICDGroupDefault,
			core.StrikeTypeDefault,
			core.Hydro,
			50,
			skill[c.TalentLvlSkill()],
		)
		d.Targets = core.TargetAll
		c.Core.Combat.ApplyDamage(&d)
	}, "tartaglia-skill", f)

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
		skillCD = (5 + timeInMeleeStance) * 60
	case 8*60 <= timeInMeleeStance && timeInMeleeStance < 30*60:
		skillCD = (6 + timeInMeleeStance) * 60
	case timeInMeleeStance >= 30*60:
		skillCD = 45 * 60
	}

	if c.Base.Cons >= 1 {
		skillCD = int(float64(skillCD) * 0.8)
	}

	if c.mlBurstUsed {
		c.SetCD(core.ActionSkill, 0)
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
		d := c.Snapshot(
			skillName,
			core.AttackTagElementalBurst,
			core.ICDTagNone,
			core.ICDGroupDefault,
			core.StrikeTypeDefault,
			core.Hydro,
			50,
			mult,
		)
		d.Targets = core.TargetAll
		if c.Core.Status.Duration("tartagliamelee") > 0 {
			d.OnHitCallback = c.rtBlastCallback
		}
		c.Core.Combat.ApplyDamage(&d)
		if c.Core.Status.Duration("tartagliamelee") > 0 {
			if c.Base.Cons >= 4 {
				c.mlBurstUsed = true
			}
		} else {
			c.AddEnergy(20)
			c.Core.Log.Debugw("Tartaglia ranged burst restoring 20 energy", "frame", c.Core.F, "event", core.LogEnergyEvent, "new energy", c.Energy)
		}
	}, "tartaglia-burst-clear", f-5) //random 5 frame

	c.Energy = 0
	c.SetCD(core.ActionBurst, 900)
	return f, a
}

func (c *char) rtFlashCallback(t core.Target) {
	if c.rtExpiry[t.Index()] > c.Core.F {
		if c.rtFlashICD[t.Index()] > c.Core.F {
			return
		}

		c.AddTask(func() {
			d := c.Snapshot(
				"Riptide Flash",
				core.AttackTagNormal,
				core.ICDTagTartagliaRiptideFlash,
				core.ICDGroupDefault,
				core.StrikeTypeDefault,
				core.Hydro,
				25,
				rtFlash[0][c.TalentLvlAttack()],
			)
			d.Targets = core.TargetAll

			//proc 3 hits
			for i := 1; i <= 3; i++ {
				x := d.Clone()
				c.QueueDmg(&x, i)
			}
			c.Core.Log.Debugw("Riptide Flash ticked", "frame", c.Core.F, "event", core.LogCharacterEvent, "dur",
				c.Core.Status.Duration("tartagliamelee"), "target", t.Index(), "flashICD", c.rtFlashICD[t.Index()], "rtExpiry", c.rtExpiry[t.Index()])

			if c.rtParticleICD < c.Core.F {
				c.rtParticleICD = c.Core.F + 180 //3 sec
				c.QueueParticle("tartaglia", 1, core.Hydro, 100)
			}
			c.rtFlashICD[t.Index()] = c.Core.F + 42 //0.7s icd
		}, "Riptide Flash", 5)
	}
}

func (c *char) rtSlashCallback(t core.Target) {
	if c.rtExpiry[t.Index()] > c.Core.F {
		if c.rtSlashICD[t.Index()] > c.Core.F {
			return
		}

		c.AddTask(func() {
			d := c.Snapshot(
				"Riptide Slash",
				core.AttackTagElementalArt,
				core.ICDTagNone,
				core.ICDGroupDefault,
				core.StrikeTypeDefault,
				core.Hydro,
				25,
				rtSlash[c.TalentLvlSkill()],
			)
			d.Targets = core.TargetAll

			c.Core.Combat.ApplyDamage(&d)
			c.Core.Log.Debugw("Riptide Slash ticked", "frame", c.Core.F, "event", core.LogCharacterEvent, "dur",
				c.Core.Status.Duration("tartagliamelee"), "target", t.Index(), "slashICD", c.rtSlashICD[t.Index()], "rtExpiry", c.rtExpiry[t.Index()])

			if c.rtParticleICD < c.Core.F {
				c.rtParticleICD = c.Core.F + 180 //3 sec
				c.QueueParticle("tartaglia", 1, core.Hydro, 100)
			}
		}, "Riptide Slash", 5)
		c.rtSlashICD[t.Index()] = c.Core.F + 90 //1.5s icd
	}
}

func (c *char) rtBlastCallback(t core.Target) {
	if c.rtExpiry[t.Index()] > c.Core.F {
		if c.rtSlashICD[t.Index()] > c.Core.F {
			return
		}

		c.AddTask(func() {
			d := c.Snapshot(
				"Riptide Blast",
				core.AttackTagElementalBurst,
				core.ICDTagNone,
				core.ICDGroupDefault,
				core.StrikeTypeDefault,
				core.Hydro,
				50,
				rtBlast[c.TalentLvlBurst()],
			)
			d.Targets = core.TargetAll

			c.Core.Combat.ApplyDamage(&d)
			// triggering riptide blast will clear riptide status
			c.rtExpiry[t.Index()] = 0
			c.Core.Log.Debugw("Riptide Blast ticked", "frame", c.Core.F, "event", core.LogCharacterEvent, "dur",
				c.Core.Status.Duration("tartagliamelee"), "target", t.Index(), "rtExpiry", c.rtExpiry[t.Index()])
		}, "Riptide Blast", 5)
	}
}

//apply riptide status to enemy hit
func (c *char) applyRT() {
	c.Core.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		ds := args[1].(*core.Snapshot)
		t := args[0].(core.Target)
		crit := args[3].(bool)

		if ds.AttackTag == core.AttackTagNormal || ds.AttackTag == core.AttackTagExtra {
			if c.Core.Status.Duration("tartagliamelee") > 0 {
				if !crit {
					return false
				}

				//dont log if it just refreshes riptide status
				if c.rtExpiry[t.Index()] <= c.Core.F {
					c.Core.Log.Debugw("Tartaglia applied riptide", "frame", c.Core.F, "event", core.LogCharacterEvent, "target", t.Index(), "rtExpiry", c.Core.F+rtA1)
				}

				c.rtExpiry[t.Index()] = c.Core.F + rtA1
			} else {
				if ds.AttackTag == core.AttackTagNormal {
					return false
				}
				//aim mode
				//dont log if it just refreshes riptide status
				if c.rtExpiry[t.Index()] <= c.Core.F {
					c.Core.Log.Debugw("Tartaglia applied riptide", "frame", c.Core.F, "event", core.LogCharacterEvent, "target", t.Index(), "rtExpiry", c.Core.F+rtA1)
				}

				c.rtExpiry[t.Index()] = c.Core.F + rtA1
			}
		}

		if ds.AttackTag == core.AttackTagElementalBurst {
			if c.Core.Status.Duration("tartagliamelee") == 0 {
				//ranged burst
				//dont log if it just refreshes riptide status
				if c.rtExpiry[t.Index()] <= c.Core.F {
					c.Core.Log.Debugw("Tartaglia applied riptide", "frame", c.Core.F, "event", core.LogCharacterEvent, "target", t.Index(), "rtExpiry", c.Core.F+rtA1)
				}

				c.rtExpiry[t.Index()] = c.Core.F + rtA1
			}
		}
		return false
	}, "tartaglia-apply-riptide")
}

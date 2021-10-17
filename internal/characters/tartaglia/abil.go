package tartaglia

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/core"
)

// TODO: Global comment for all callbacks - maybe it would be efficient to clone a snapshot for each riptide effect?
// Current method is wasteful when there's a bunch happening all on the same frame, and it functionally shouldn't matter

// Regular normal bow attacks - largely uninteresting
func (c *char) Attack(p map[string]int) (int, int) {

	travel, ok := p["travel"]
	if !ok {
		travel = 20
	}

	f, a := c.ActionFrames(core.ActionAttack, p)

	if c.Core.Status.Duration("tartagliamelee") > 0 {
		return c.meleeAttack(f, a)
	}

	c.QueueDmgDynamic(func() *core.Snapshot {
		d := c.Snapshot(
			fmt.Sprintf("Normal %v", c.NormalCounter),
			core.AttackTagNormal,
			core.ICDTagNone,
			core.ICDGroupDefault,
			core.StrikeTypePierce,
			core.Physical,
			25,
			attack[c.NormalCounter][c.TalentLvlAttack()],
		)
		return &d
	}, travel+f)

	c.AdvanceNormalIndex()

	return f, a
}

// Need to space out the multihits by at least 2 frames so riptide procs are handled correctly
var meleeAttackDelayOffset = [][]int{
	{0},
	{0},
	{0},
	{0},
	{0},
	{2, 0},
}

func (c *char) meleeAttack(f int, a int) (int, int) {
	for i, mult := range attackMelee[c.NormalCounter] {
		c.QueueDmgDynamic(func() *core.Snapshot {
			d := c.Snapshot(
				fmt.Sprintf("Melee Normal %v", c.NormalCounter),
				core.AttackTagNormal,
				core.ICDTagNormalAttack,
				core.ICDGroupDefault,
				core.StrikeTypePierce,
				core.Hydro,
				25,
				mult[c.TalentLvlSkill()],
			)
			d.OnHitCallback = c.riptideSlashCallback
			// TODO: Targeting?

			return &d
		}, f-meleeAttackDelayOffset[c.NormalCounter][i])
	}

	c.AdvanceNormalIndex()

	return f, a
}

func (c *char) Aimed(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionAim, p)

	travel, ok := p["travel"]
	if !ok {
		travel = 20
	}

	// Bow attacks are snapshot upon release
	c.AddTask(func() {
		d := c.Snapshot(
			"Aimed Shot",
			core.AttackTagExtra,
			core.ICDTagNone,
			core.ICDGroupDefault,
			core.StrikeTypePierce,
			core.Hydro,
			25,
			chargeRange[c.TalentLvlAttack()],
		)
		d.OnHitCallback = c.riptideFlashCallback

		c.QueueDmg(&d, travel)
	}, "tartaglia-aim", f)

	return f, a
}

// On hit callback that procs damage for Riptide Flash (from Ranged CAs)
// Hitting an opponent affected by Riptide with a fully charged aimed shot deals consecutive bouts of AoE DMG. Can occur once every 0.7s.
func (c *char) riptideFlashCallback(t core.Target) {
	if c.riptideStatusLastProc[t.Index()]+18*60 < c.Core.F {
		return
	}

	if c.riptideFlashLastProc[t.Index()]+int(0.7*60) > c.Core.F {
		return
	}

	// TODO: AoE range on this is so small that it might as well be single target?
	d := c.Snapshot(
		"Riptide Flash",
		core.AttackTagExtra,
		core.ICDTagTartagliaRiptideFlash,
		core.ICDGroupDefault,
		core.StrikeTypePierce,
		core.Hydro,
		25,
		riptideFlash[c.TalentLvlAttack()],
	)
	d.Targets = t.Index()

	// Apply damage 3 times
	// TODO: Not sure if this is one frame after another or how snapshot works
	// For now for simplicity, assume 1 frame apart and that it's snapshot on the initial hit
	for i := 0; i < 3; i++ {
		x := d.Clone()
		c.QueueDmg(&x, 1+i)
	}

	// Riptide generates 1 particle upon proc, with an ICD of 3 seconds
	if c.riptideParticleLastProc+180 <= c.Core.F {
		c.QueueParticle("tartaglia", 1, core.Hydro, 100)
		c.riptideParticleLastProc = c.Core.F
	}
}

// Need to space out the multihits by at least 2 frames so riptide procs are handled correctly
var chargeAttackDelayOffset = []int{2, 0}

func (c *char) ChargeAttack(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionCharge, p)

	for i, mult := range chargeMelee {

		c.QueueDmgDynamic(func() *core.Snapshot {
			d := c.Snapshot(
				"Melee Charge Attack",
				core.AttackTagExtra,
				core.ICDTagExtraAttack,
				core.ICDGroupDefault,
				core.StrikeTypePierce,
				core.Hydro,
				25,
				mult[c.TalentLvlSkill()],
			)
			d.Targets = core.TargetAll
			d.OnHitCallback = c.riptideSlashCallback

			return &d
		}, f-chargeAttackDelayOffset[i])
	}

	return f, a
}

// Handles stance change
func (c *char) Skill(p map[string]int) (int, int) {

	// Enters melee stance
	if c.Core.Status.Duration("tartagliamelee") == 0 {

		f, a := c.ActionFrames(core.ActionSkill, p)

		c.QueueDmgDynamic(func() *core.Snapshot {
			d := c.Snapshot(
				"Stance Change",
				core.AttackTagElementalArt,
				core.ICDTagNormalAttack,
				core.ICDGroupDefault,
				core.StrikeTypePierce,
				core.Hydro,
				50,
				skill[c.TalentLvlSkill()],
			)
			return &d
		}, f)

		// Enter melee stance
		c.Core.Status.AddStatus("tartagliamelee", 30*60)
		c.skillStartFrame = c.Core.F

		return f, a
	}

	// Otherwise exits melee stance
	c.onExitMeleeStance()

	// TODO: May actually be some kind of frame delay as he has that bow spinning animation?
	return 0, 0
}

// Removes melee status and sets correct cooldown. Used when manually removing stance, on switch, and when status runs out
// TODO: Have not added a check when status runs out. Too lazy since no one should ever reasonably do it
func (c *char) onExitMeleeStance() {
	// Precise skill CD from Risuke:
	// Aligns with separate table on wiki except the 4 second duration one
	// https://discord.com/channels/763583452762734592/851428030094114847/899416824117084210
	// https://media.discordapp.net/attachments/778615842916663357/781978094495727646/unknown-20.png

	skillCD := 0

	switch timeInMeleeStance := c.Core.F - c.skillStartFrame; {
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

	if c.c6BurstUsed {
		c.SetCD(core.ActionSkill, 0)
		c.c6BurstUsed = false
	} else {
		c.SetCD(core.ActionSkill, skillCD)
	}
	c.Core.Status.DeleteStatus("tartagliamelee")
}

func (c *char) riptideSlashCallback(t core.Target) {
	// Proc riptide if able
	if c.riptideStatusLastProc[t.Index()]+18*60 < c.Core.F {
		return
	}
	if c.riptideSlashLastProc[t.Index()]+90 > c.Core.F {
		return
	}
	c.QueueDmgDynamic(func() *core.Snapshot {
		d := c.Snapshot(
			"Riptide Slash",
			core.AttackTagElementalArt,
			core.ICDTagNone,
			core.ICDGroupDefault,
			core.StrikeTypePierce,
			core.Hydro,
			25,
			riptideSlash[c.TalentLvlSkill()],
		)
		d.Targets = core.TargetAll

		// Riptide generates 1 particle upon proc, with an ICD of 3 seconds
		if c.riptideParticleLastProc+180 <= c.Core.F {
			c.QueueParticle("tartaglia", 1, core.Hydro, 100)
			c.riptideParticleLastProc = c.Core.F
		}

		c.riptideSlashLastProc[t.Index()] = c.Core.F

		return &d
	}, 1)
}

func (c *char) Burst(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionBurst, p)

	if c.Core.Status.Duration("tartagliamelee") == 0 {

		c.QueueDmgDynamic(func() *core.Snapshot {
			d := c.Snapshot(
				"Ranged Burst: Flash of Havoc",
				core.AttackTagElementalBurst,
				core.ICDTagElementalBurst,
				core.ICDGroupDefault,
				core.StrikeTypePierce,
				core.Hydro,
				50,
				burstRanged[c.TalentLvlBurst()],
			)
			d.Targets = core.TargetAll

			return &d
		}, f)

		c.Energy = 20
		c.SetCD(core.ActionBurst, 15*60)

		return f, a
	}
	// Melee
	c.QueueDmgDynamic(func() *core.Snapshot {
		d := c.Snapshot(
			"Melee Burst: Light of Obliteration",
			core.AttackTagElementalBurst,
			core.ICDTagElementalBurst,
			core.ICDGroupDefault,
			core.StrikeTypePierce,
			core.Hydro,
			50,
			burstMelee[c.TalentLvlBurst()],
		)
		d.Targets = core.TargetAll
		d.OnHitCallback = c.riptideBlastCallback

		c.c6BurstUsed = true

		return &d
	}, f)

	c.Energy = 0
	c.SetCD(core.ActionBurst, 15*60)

	return f, a
}

func (c *char) riptideBlastCallback(t core.Target) {
	if c.riptideStatusLastProc[t.Index()]+18*60 < c.Core.F {
		return
	}

	// "Clears" riptide
	c.riptideStatusLastProc[t.Index()] = -9999
	c.QueueDmgDynamic(func() *core.Snapshot {
		d := c.Snapshot(
			"Riptide Blast",
			core.AttackTagElementalBurst,
			core.ICDTagNone,
			core.ICDGroupDefault,
			core.StrikeTypePierce,
			core.Hydro,
			50,
			riptideBlast[c.TalentLvlBurst()],
		)
		d.Targets = core.TargetAll

		return &d
	}, 1)
}

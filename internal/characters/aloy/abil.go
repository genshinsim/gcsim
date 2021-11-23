package aloy

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

// Standard attack - infusion mechanics are handled as part of the skill
func (c *char) Attack(p map[string]int) (int, int) {

	travel, ok := p["travel"]
	if !ok {
		travel = 20
	}

	f, a := c.ActionFrames(core.ActionAttack, p)

	for i, mult := range attack[c.NormalCounter] {
		d := c.Snapshot(
			fmt.Sprintf("Normal %v", c.NormalCounter),
			core.AttackTagNormal,
			core.ICDTagNormalAttack,
			core.ICDGroupDefault,
			core.StrikeTypeDefault,
			core.Physical,
			25,
			mult[c.TalentLvlAttack()],
		)
		c.QueueDmg(&d, f+i+travel)
	}

	c.AdvanceNormalIndex()

	// return animation cd
	return f, a
}

// Skill
func (c *char) Skill(p map[string]int) (int, int) {

	bomblets, ok := p["bomblets"]
	if !ok {
		bomblets = 2
	}

	bombletCoilStacks, ok := p["bomblet_coil_stacks"]
	if !ok {
		bombletCoilStacks = 2
	}

	delay, ok := p["delay"]
	if !ok {
		delay = 0
	}

	f, a := c.ActionFrames(core.ActionSkill, p)

	// Initial damage
	c.QueueDmgDynamic(func() *core.Snapshot {
		// TODO: Not 100% sure about ICD
		d := c.Snapshot(
			"Freeze Bomb",
			core.AttackTagElementalArt,
			core.ICDTagNone,
			core.ICDGroupDefault,
			core.StrikeTypeDefault,
			core.Cryo,
			25,
			skillMain[c.TalentLvlSkill()],
		)
		d.Targets = core.TargetAll

		c.coilStacks()

		return &d
	}, f)

	// Bomblets snapshot on cast
	dBomblets := c.Snapshot(
		"Chillwater Bomblets",
		core.AttackTagElementalArt,
		core.ICDTagElementalArt,
		core.ICDGroupDefault,
		core.StrikeTypeDefault,
		core.Cryo,
		25,
		skillBomblets[c.TalentLvlSkill()],
	)

	// Queue up bomblets
	for i := 0; i < bomblets; i++ {
		x := dBomblets.Clone()

		c.QueueDmg(&x, f+delay+((i+1)*6))
	}

	// Queue up bomblet coil stacks
	for i := 0; i < bombletCoilStacks; i++ {
		c.AddTask(func() {
			c.coilStacks()
		}, "aloy-bomblet-coil-stacks", f+delay+((i+1)*6))
	}

	c.QueueParticle("aloy", 5, core.Cryo, f+100)
	c.SetCD(core.ActionSkill, 20*60)

	return f, a
}

// Handles coil stacking and associated effects, including triggering rushing ice
func (c *char) coilStacks() {
	if c.coilICDExpiry > c.Core.F {
		return
	}
	// Can't gain coil stacks while in rushing ice
	if c.Core.Status.Duration("aloyrushingice") > 0 {
		return
	}
	c.Tags["coil_stacks"]++
	c.coilICDExpiry = c.Core.F + 6

	// A1
	// When Aloy receives the Coil effect from Frozen Wilds, her ATK is increased by 16%, while nearby party members' ATK is increased by 8%. This effect lasts 10s.
	valA1 := make([]float64, core.EndStatType)
	for _, char := range c.Core.Chars {
		valA1[core.ATKP] = .08
		if char.CharIndex() == c.Index {
			valA1[core.ATKP] = .16
		}

		char.AddMod(core.CharStatMod{
			Key:    "aloy-a1",
			Expiry: c.Core.F + 600,
			Amount: func(a core.AttackTag) ([]float64, bool) {
				return valA1, true
			},
		})
	}

	if c.Tags["coil_stacks"] == 4 {
		c.rushingIce()
	}
}

// Handles rushing ice state
func (c *char) rushingIce() {
	c.Core.Status.AddStatus("aloyrushingice", 600)

	c.AddWeaponInfuse(core.WeaponInfusion{
		Key:    "aloy-rushing-ice",
		Ele:    core.Cryo,
		Tags:   []core.AttackTag{core.AttackTagNormal},
		Expiry: c.Core.F + 600,
	})

	// Rushing ice NA bonus
	val := make([]float64, core.EndStatType)
	val[core.DmgP] = skillRushingIceNABonus[c.TalentLvlSkill()]
	c.AddMod(core.CharStatMod{
		Key:    "aloy-rushing-ice",
		Expiry: c.Core.F + 600,
		Amount: func(a core.AttackTag) ([]float64, bool) {
			if a == core.AttackTagNormal {
				return val, true
			}
			return nil, false
		},
	})

	// A4 cryo damage increase
	valA4 := make([]float64, core.EndStatType)
	stacks := 1
	c.AddMod(core.CharStatMod{
		Key:    "aloy-strong-strike",
		Expiry: c.Core.F + 600,
		Amount: func(a core.AttackTag) ([]float64, bool) {
			if stacks > 10 {
				stacks = 10
			}
			valA4[core.CryoP] = float64(stacks) * 0.035
			return valA4, true
		},
	})
	for i := 0; i < 10; i++ {
		c.AddTask(func() { stacks++ }, "aloy-strone-strike-stack", 60*(1+i))
	}
}

// Burst - doesn't do much other than damage, so fairly straightforward
func (c *char) Burst(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionBurst, p)

	// TODO: Assuming dynamic
	c.QueueDmgDynamic(func() *core.Snapshot {
		d := c.Snapshot(
			"Prophecies of Dawn",
			core.AttackTagElementalBurst,
			core.ICDTagElementalBurst,
			core.ICDGroupDefault,
			core.StrikeTypeDefault,
			core.Cryo,
			50,
			burst[c.TalentLvlBurst()],
		)
		d.Targets = core.TargetAll
		return &d
	}, f)

	c.SetCD(core.ActionBurst, 12*60)
	// TODO: Not sure when energy drain happens
	c.Energy = 0

	return f, a
}

package aloy

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

// Standard attack - infusion mechanics are handled as part of the skill
func (c *char) Attack(p map[string]int) (int, int) {

	travel, ok := p["travel"]
	if !ok {
		travel = 10
	}

	f, a := c.ActionFrames(core.ActionAttack, p)

	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       fmt.Sprintf("Normal %v", c.NormalCounter),
		AttackTag:  core.AttackTagNormal,
		ICDTag:     core.ICDTagNormalAttack,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Physical,
		Durability: 25,
	}

	for i, mult := range attack[c.NormalCounter] {
		ai.Mult = mult[c.TalentLvlAttack()]
		c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(0.1, false, core.TargettableEnemy), f+i, f+i+travel)
	}

	c.AdvanceNormalIndex()

	// return animation cd
	return f, a
}

// Standard aimed attack
func (c *char) Aimed(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionAim, p)

	travel, ok := p["travel"]
	if !ok {
		travel = 10
	}
	weakspot, ok := p["weakspot"]

	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Charge Shot",
		// TODO: Not sure about CA ICD
		AttackTag:    core.AttackTagExtra,
		ICDTag:       core.ICDTagExtraAttack,
		ICDGroup:     core.ICDGroupDefault,
		Element:      core.Cryo,
		Durability:   25,
		Mult:         aim[c.TalentLvlAttack()],
		HitWeakPoint: weakspot == 1,
	}
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(0.1, false, core.TargettableEnemy), f, f+travel)

	return f, a
}

// Skill - Handles main damage, bomblet, and coil effects
// Has 3 parameters, "bomblets" = Number of bomblets that hit
// "bomblet_coil_stacks" = Number of coil stacks gained
// "delay" - Delay in frames before bomblets go off and coil stacks get added
// Too many potential bomblet hit variations to keep syntax short, so we simplify how they can be handled here
func (c *char) Skill(p map[string]int) (int, int) {

	bomblets, ok := p["bomblets"]
	if !ok {
		bomblets = 2
	}

	bombletCoilStacks, ok := p["bomblet_coil_stacks"]
	if !ok {
		bombletCoilStacks = 2
	}

	delay, ok := p["bomb_delay"]
	if !ok {
		delay = 0
	}

	f, a := c.ActionFrames(core.ActionSkill, p)

	c.Core.Tasks.Add(func() {
		// TODO: Not 100% sure about ICD
		ai := core.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Freeze Bomb",
			AttackTag:  core.AttackTagElementalArt,
			ICDTag:     core.ICDTagNone,
			ICDGroup:   core.ICDGroupDefault,
			Element:    core.Cryo,
			Durability: 25,
			Mult:       skillMain[c.TalentLvlSkill()],
		}
		c.coilStacks()
		c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(3, false, core.TargettableEnemy), 0, 0)
	}, f)

	// Bomblets snapshot on cast
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Chillwater Bomblets",
		AttackTag:  core.AttackTagElementalArt,
		ICDTag:     core.ICDTagElementalArt,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Cryo,
		Durability: 25,
		Mult:       skillBomblets[c.TalentLvlSkill()],
	}

	// Queue up bomblets
	for i := 0; i < bomblets; i++ {
		c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(0.1, false, core.TargettableEnemy), 0, f+delay+((i+1)*6))
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
	for _, char := range c.Core.Chars {
		valA1 := make([]float64, core.EndStatType)
		valA1[core.ATKP] = .08
		if char.CharIndex() == c.Index {
			valA1[core.ATKP] = .16
		}
		char.AddMod(core.CharStatMod{
			Key:    "aloy-a1",
			Expiry: c.Core.F + 600,
			Amount: func() ([]float64, bool) {
				return valA1, true
			},
		})
	}

	if c.Tags["coil_stacks"] == 4 {
		c.Tags["coil_stacks"] = 0
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
	c.AddPreDamageMod(core.PreDamageMod{
		Key:    "aloy-rushing-ice",
		Expiry: c.Core.F + 600,
		Amount: func(atk *core.AttackEvent, t core.Target) ([]float64, bool) {
			if atk.Info.AttackTag == core.AttackTagNormal {
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
		Amount: func() ([]float64, bool) {
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
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Prophecies of Dawn",
		AttackTag:  core.AttackTagElementalBurst,
		ICDTag:     core.ICDTagElementalBurst,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Cryo,
		Durability: 50,
		Mult:       burst[c.TalentLvlBurst()],
	}
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(3, false, core.TargettableEnemy), f, f)

	c.SetCDWithDelay(core.ActionBurst, 12*60, 8)
	c.ConsumeEnergy(8)

	return f, a
}

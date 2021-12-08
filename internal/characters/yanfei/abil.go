package yanfei

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/shield"
)

// Standard attack function with seal handling
func (c *char) Attack(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionAttack, p)

	travel, ok := p["travel"]
	if !ok {
		travel = 20
	}

	d := c.Snapshot(
		fmt.Sprintf("Normal %v", c.NormalCounter),
		core.AttackTagNormal,
		core.ICDTagNormalAttack,
		core.ICDGroupDefault,
		core.StrikeTypeDefault,
		core.Pyro,
		25,
		attack[c.NormalCounter][c.TalentLvlAttack()],
	)

	c.AddTask(func() {
		// Technically seals are earned on hitting the enemy, but we just keep it here instead of an event
		if c.Tags["seal"] < c.maxTags {
			c.Tags["seal"]++
		}
		c.sealExpiry = c.Core.F + 600
		c.Core.Log.Debugw("yanfei gained a seal from normal attack", "frame", c.Core.F, "event", core.LogCharacterEvent, "char", c.Index, "current_seals", c.Tags["seal"], "expiry", c.sealExpiry)

		c.Core.Combat.ApplyDamage(&d)
	}, "yanfei-attack", f+travel)

	c.AdvanceNormalIndex()

	// return animation cd
	return f, a
}

// Charge attack function - handles seal use
func (c *char) ChargeAttack(p map[string]int) (int, int) {

	//check for seal stacks
	if c.Core.F > c.sealExpiry {
		c.Tags["seal"] = 0
	}
	stacks := c.Tags["seal"]

	//a1
	// When Yan Fei's Charged Attack consumes Scarlet Seals, each Scarlet Seal consumed will increase her Pyro DMG by 5% for 6 seconds. When this effect is repeatedly triggered it will overwrite the oldest bonus first.
	// The Pyro DMG bonus from Proviso is applied before charged attack damage is calculated.
	var m [core.EndStatType]float64
	m[core.PyroP] = float64(stacks) * 0.05
	c.AddMod(core.CharStatMod{
		Key: "yanfei-a1",
		Amount: func(a core.AttackTag) ([core.EndStatType]float64, bool) {
			return m, true
		},
		Expiry: c.Core.F + 360,
	})

	f, a := c.ActionFrames(core.ActionCharge, p)

	d := c.Snapshot(
		"Charge Attack",
		core.AttackTagExtra,
		core.ICDTagNone,
		core.ICDGroupDefault,
		core.StrikeTypeBlunt,
		core.Pyro,
		25,
		charge[stacks][c.TalentLvlAttack()],
	)
	d.Targets = core.TargetAll

	// TODO: Not sure of snapshot timing
	c.QueueDmg(&d, f)

	c.Core.Log.Debugw("yanfei charge attack consumed seals", "frame", c.Core.F, "event", core.LogCharacterEvent, "char", c.Index, "current_seals", c.Tags["seal"], "expiry", c.sealExpiry)

	// Clear the seals next frame just in case for some reason we call stam check late
	c.AddTask(func() {
		c.Tags["seal"] = 0
		c.sealExpiry = c.Core.F - 1
	}, "clear-seals", 1)

	return f, a
}

// Yanfei skill - Straightforward as it has little interactions with the rest of her kit
// Summons flames that deal AoE Pyro DMG. Opponents hit by the flames will grant Yanfei the maximum number of Scarlet Seals.
func (c *char) Skill(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionSkill, p)

	c.AddTask(func() {
		d := c.Snapshot(
			"Signed Edict",
			core.AttackTagElementalArt,
			core.ICDTagNone,
			core.ICDGroupDefault,
			core.StrikeTypeBlunt,
			core.Pyro,
			25,
			skill[c.TalentLvlSkill()],
		)
		d.Targets = core.TargetAll

		// Create max seals on hit
		if c.Tags["seal"] < c.maxTags {
			c.Tags["seal"] = c.maxTags
		}
		c.sealExpiry = c.Core.F + 600

		c.Core.Log.Debugw("yanfei gained max seals", "frame", c.Core.F, "event", core.LogCharacterEvent, "char", c.Index, "current_seals", c.Tags["seal"], "expiry", c.sealExpiry)

		c.Core.Combat.ApplyDamage(&d)

	}, "yanfei-skill", f)

	c.QueueParticle("yanfei", 3, core.Pyro, f+100)

	c.SetCD(core.ActionSkill, 540)

	return f, a
}

// Burst - Deals burst damage and adds status for charge attack bonus
func (c *char) Burst(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionBurst, p)

	// +1 is to make sure the scarlet seal grant works correctly on the last frame
	// TODO: Not 100% sure whether this adds a seal at the exact moment the burst ends or not
	c.Core.Status.AddStatus("yanfeiburst", 15*60+1)

	var m [core.EndStatType]float64
	m[core.DmgP] = burstBonus[c.TalentLvlBurst()]
	c.AddMod(core.CharStatMod{
		Key: "yanfei-burst",
		Amount: func(a core.AttackTag) ([core.EndStatType]float64, bool) {
			if a == core.AttackTagExtra {
				return m, true
			}
			return nil, false
		},
		Expiry: c.Core.F + 15*60,
	})

	c.AddTask(func() {
		d := c.Snapshot(
			"Done Deal",
			core.AttackTagElementalBurst,
			core.ICDTagNone,
			core.ICDGroupDefault,
			core.StrikeTypeBlunt,
			core.Pyro,
			50,
			burst[c.TalentLvlBurst()],
		)
		d.Targets = core.TargetAll

		// Create max seals on hit
		if c.Tags["seal"] < c.maxTags {
			c.Tags["seal"] = c.maxTags
		}
		c.sealExpiry = c.Core.F + 600

		c.Core.Log.Debugw("yanfei gained max seals", "frame", c.Core.F, "event", core.LogCharacterEvent, "char", c.Index, "current_seals", c.Tags["seal"], "expiry", c.sealExpiry)

		c.Core.Combat.ApplyDamage(&d)
	}, "yanfei-burst", f)

	c.AddTask(c.burstAddSealHook(), "burst-add-seals-task", 60)

	c.c4()

	c.SetCD(core.ActionBurst, 20*60)
	c.Energy = 0

	return f, a
}

// Handles C4 shield creation
// When Done Deal is used:
// Creates a shield that absorbs up to 45% of Yan Fei's Max HP for 15s
// This shield absorbs Pyro DMG 250% more effectively
func (c *char) c4() {
	if c.Base.Cons >= 4 {
		c.Core.Shields.Add(&shield.Tmpl{
			Src:        c.Core.F,
			ShieldType: core.ShieldYanfeiC4,
			HP:         c.HPMax * .45,
			Ele:        core.Pyro,
			Expires:    c.Core.F + 15*60,
		})
	}
}

// Recurring task to add seals every second while burst is up
func (c *char) burstAddSealHook() func() {
	return func() {
		if c.Core.Status.Duration("yanfeiburst") == 0 {
			return
		}
		if c.Tags["seal"] < c.maxTags {
			c.Tags["seal"]++
		}
		c.sealExpiry = c.Core.F + 600

		c.Core.Log.Debugw("yanfei gained seal from burst", "frame", c.Core.F, "event", core.LogCharacterEvent, "char", c.Index, "current_seals", c.Tags["seal"], "expiry", c.sealExpiry)

		c.AddTask(c.burstAddSealHook(), "burst-add-seals", 60)
	}
}

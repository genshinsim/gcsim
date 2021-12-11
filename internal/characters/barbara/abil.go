package barbara

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

// Standard attack function with seal handling
func (c *char) Attack(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionAttack, p)

	travel, ok := p["travel"]
	if !ok {
		travel = 20
	}

	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       fmt.Sprintf("Normal %v", c.NormalCounter),
		AttackTag:  core.AttackTagNormal,
		ICDTag:     core.ICDTagNormalAttack,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Pyro,
		Durability: 25,
		Mult:       attack[c.NormalCounter][c.TalentLvlAttack()],
	}
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(0.1, false, core.TargettableEnemy), 0, f+travel)

	c.AdvanceNormalIndex()

	// return animation cd
	return f, a
}

// Charge attack function - handles seal use
func (c *char) ChargeAttack(p map[string]int) (int, int) {

	//a1
	// When Yan Fei's Charged Attack consumes Scarlet Seals, each Scarlet Seal consumed will increase her Pyro DMG by 5% for 6 seconds. When this effect is repeatedly triggered it will overwrite the oldest bonus first.
	// The Pyro DMG bonus from Proviso is applied before charged attack damage is calculated.
	var m [core.EndStatType]float64
	c.AddMod(core.CharStatMod{
		Key: "barbara-a1",
		Amount: func(a core.AttackTag) ([core.EndStatType]float64, bool) {
			return m, true
		},
		Expiry: c.Core.F + 360,
	})

	f, a := c.ActionFrames(core.ActionCharge, p)

	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Charge Attack",
		AttackTag:  core.AttackTagExtra,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeBlunt,
		Element:    core.Pyro,
		Durability: 25,
	}
	// TODO: Not sure of snapshot timing
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(2, false, core.TargettableEnemy), 0, f)

	c.Core.Log.Debugw("barbara charge attack consumed seals", "frame", c.Core.F, "event", core.LogCharacterEvent, "char", c.Index, "current_seals", c.Tags["seal"], "expiry")

	// Clear the seals next frame just in case for some reason we call stam check late
	c.AddTask(func() {
		c.Tags["seal"] = 0
	}, "clear-seals", 1)

	return f, a
}

// barbara skill - Straightforward as it has little interactions with the rest of her kit
// Summons flames that deal AoE Pyro DMG. Opponents hit by the flames will grant barbara the maximum number of Scarlet Seals.
func (c *char) Skill(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionSkill, p)

	done := false
	addSeal := func(t core.Target, ae *core.AttackEvent) {
		if done {
			return
		}
		// Create max seals on hit
		c.Core.Log.Debugw("barbara gained max seals", "frame", c.Core.F, "event", core.LogCharacterEvent, "char", c.Index)
		done = true
	}

	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Signed Edict",
		AttackTag:  core.AttackTagElementalArt,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeBlunt,
		Element:    core.Pyro,
		Durability: 25,
		Mult:       skill[c.TalentLvlSkill()],
	}
	// TODO: Not sure of snapshot timing
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(2, false, core.TargettableEnemy), 0, f, addSeal)

	c.QueueParticle("barbara", 3, core.Pyro, f+100)

	c.SetCD(core.ActionSkill, 540)

	return f, a
}

// Burst - Deals burst damage and adds status for charge attack bonus
func (c *char) Burst(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionBurst, p)

	// +1 is to make sure the scarlet seal grant works correctly on the last frame
	// TODO: Not 100% sure whether this adds a seal at the exact moment the burst ends or not
	c.Core.Status.AddStatus("barbaraburst", 15*60+1)

	var m [core.EndStatType]float64
	m[core.DmgP] = burstBonus[c.TalentLvlBurst()]
	c.AddMod(core.CharStatMod{
		Key: "barbara-burst",
		Amount: func(a core.AttackTag) ([core.EndStatType]float64, bool) {
			if a == core.AttackTagExtra {
				return m, true
			}
			return m, false
		},
		Expiry: c.Core.F + 15*60,
	})

	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Signed Edict",
		AttackTag:  core.AttackTagElementalBurst,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeBlunt,
		Element:    core.Pyro,
		Durability: 50,
		Mult:       burst[c.TalentLvlBurst()],
	}

	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(2, false, core.TargettableEnemy), 0, f)

	c.AddTask(c.burstAddSealHook(), "burst-add-seals-task", 60)

	c.SetCD(core.ActionBurst, 20*60)
	c.Energy = 0

	return f, a
}

// Recurring task to add seals every second while burst is up
func (c *char) burstAddSealHook() func() {
	return func() {
		if c.Core.Status.Duration("barbaraburst") == 0 {
			return
		}

		c.Core.Log.Debugw("barbara gained seal from burst", "frame", c.Core.F, "event", core.LogCharacterEvent, "char", c.Index)

		c.AddTask(c.burstAddSealHook(), "burst-add-seals", 60)
	}
}

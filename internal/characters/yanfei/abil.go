package yanfei

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/tmpl/shield"
	"github.com/genshinsim/gcsim/pkg/core"
)

// Standard attack function with seal handling
func (c *char) Attack(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionAttack, p)

	travel, ok := p["travel"]
	if !ok {
		travel = 10
	}

	done := false
	addSeal := func(a core.AttackCB) {
		if done {
			return
		}
		if c.Tags["seal"] < c.maxTags {
			c.Tags["seal"]++
		}
		c.sealExpiry = c.Core.F + 600
		c.Core.Log.NewEvent("yanfei gained a seal from normal attack", core.LogCharacterEvent, c.Index, "current_seals", c.Tags["seal"], "expiry", c.sealExpiry)
		done = true
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
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(0.1, false, core.TargettableEnemy), 0, f+travel, addSeal)

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
	m := make([]float64, core.EndStatType)
	m[core.PyroP] = float64(stacks) * 0.05
	c.AddMod(core.CharStatMod{
		Key: "yanfei-a1",
		Amount: func() ([]float64, bool) {
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
		Mult:       charge[stacks][c.TalentLvlAttack()],
	}
	// TODO: Not sure of snapshot timing
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(2, false, core.TargettableEnemy), 0, f)

	c.Core.Log.NewEvent("yanfei charge attack consumed seals", core.LogCharacterEvent, c.Index, "current_seals", c.Tags["seal"], "expiry", c.sealExpiry)

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

	done := false
	addSeal := func(a core.AttackCB) {
		if done {
			return
		}
		// Create max seals on hit
		if c.Tags["seal"] < c.maxTags {
			c.Tags["seal"] = c.maxTags
		}
		c.sealExpiry = c.Core.F + 600
		c.Core.Log.NewEvent("yanfei gained max seals", core.LogCharacterEvent, c.Index, "current_seals", c.Tags["seal"], "expiry", c.sealExpiry)
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

	m := make([]float64, core.EndStatType)
	m[core.DmgP] = burstBonus[c.TalentLvlBurst()]
	c.AddPreDamageMod(core.PreDamageMod{
		Key: "yanfei-burst",
		Amount: func(atk *core.AttackEvent, t core.Target) ([]float64, bool) {
			if atk.Info.AttackTag == core.AttackTagExtra {
				return m, true
			}
			return nil, false
		},
		Expiry: c.Core.F + 15*60,
	})

	done := false
	addSeal := func(a core.AttackCB) {
		if done {
			return
		}
		// Create max seals on hit
		if c.Tags["seal"] < c.maxTags {
			c.Tags["seal"] = c.maxTags
		}
		c.sealExpiry = c.Core.F + 600
		c.Core.Log.NewEvent("yanfei gained max seals", core.LogCharacterEvent, c.Index, "current_seals", c.Tags["seal"], "expiry", c.sealExpiry)
		done = true
	}

	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Done Deal",
		AttackTag:  core.AttackTagElementalBurst,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeBlunt,
		Element:    core.Pyro,
		Durability: 50,
		Mult:       burst[c.TalentLvlBurst()],
	}

	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(2, false, core.TargettableEnemy), 0, f, addSeal)

	c.AddTask(c.burstAddSealHook(), "burst-add-seals-task", 60)

	c.c4()

	c.SetCDWithDelay(core.ActionBurst, 20*60, 8)
	c.ConsumeEnergy(8)

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
			Name:       "Yanfei C4",
			HP:         c.MaxHP() * .45,
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

		c.Core.Log.NewEvent("yanfei gained seal from burst", core.LogCharacterEvent, c.Index, "current_seals", c.Tags["seal"], "expiry", c.sealExpiry)

		c.AddTask(c.burstAddSealHook(), "burst-add-seals", 60)
	}
}

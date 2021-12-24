package kokomi

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

// Standard attack damage function
// Has "travel" parameter, used to set the number of frames that the projectile is in the air (default = 20)
func (c *char) Attack(p map[string]int) (int, int) {
	travel, ok := p["travel"]
	if !ok {
		travel = 20
	}

	f, a := c.ActionFrames(core.ActionAttack, p)

	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       fmt.Sprintf("Normal %v", c.NormalCounter),
		AttackTag:  core.AttackTagNormal,
		ICDTag:     core.ICDTagNormalAttack,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Hydro,
		Durability: 25,
		Mult:       attack[c.NormalCounter][c.TalentLvlAttack()],
	}
	ai.FlatDmg = c.burstDmgBonus(ai.AttackTag)

	c.Core.Combat.QueueAttack(ai, core.NewDefSingleTarget(1, core.TargettableEnemy), 0, f+travel)
	// TODO: Assume that this is not dynamic (snapshot on projectile release)

	if c.NormalCounter == c.NormalHitNum-1 {
		c.c1(f)
	}

	c.AdvanceNormalIndex()

	return f, a
}

func (c *char) c1(f int) {
	if c.Base.Cons == 0 {
		return
	}
	if c.Core.Status.Duration("kokomiburst") == 0 {
		return
	}

	// TODO: Assume that these are 1A (not specified in library)
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       fmt.Sprintf("Normal %v", c.NormalCounter),
		AttackTag:  core.AttackTagNone,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Hydro,
		Durability: 25,
		Mult:       0,
	}
	ai.FlatDmg = 0.3 * c.HPMax

	// TODO: Frames not in library - Think it's 7 frames based on a rough count
	// TODO: Is this snapshotted/dynamic?
	c.Core.Combat.QueueAttack(ai, core.NewDefSingleTarget(1, core.TargettableEnemy), 0, f+7)
}

// Standard charge attack
func (c *char) ChargeAttack(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionCharge, p)

	// CA has no travel time

	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Charge",
		AttackTag:  core.AttackTagExtra,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Hydro,
		Durability: 25,
		Mult:       charge[c.TalentLvlAttack()],
	}
	ai.FlatDmg = c.burstDmgBonus(ai.AttackTag)

	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(2, false, core.TargettableEnemy), f, f)
	return f, a
}

// Skill handling - Handles primary damage instance
// Deals Hydro DMG to surrounding opponents and heal nearby active characters once every 2s. This healing is based on Kokomi's Max HP.
// TODO: Have not handled the fact that you can snapshot burst bonus onto skill if you switch immediately after casting burst
func (c *char) Skill(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionSkill, p)

	// Plus 1 to avoid same frame issues with skill ticks
	c.Core.Status.AddStatus("kokomiskill", 12*60+1)

	d := c.createSkillSnapshot()

	// You get 1 tick immediately, then 1 tick every 2 seconds for a total of 7 ticks
	c.AddTask(func() { c.skillTick(d) }, "kokomi-e-tick", 1)

	c.AddTask(c.skillTickTask(d, c.Core.F), "kokomi-e-ticks", 120)

	c.skillLastUsed = c.Core.F - 1
	c.SetCD(core.ActionSkill, 20*60)

	return f, a
}

// Helper function since this needs to be created both on skill use and burst use
func (c *char) createSkillSnapshot() *core.AttackEvent {

	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Bake-Kurage",
		AttackTag:  core.AttackTagElementalArt,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Hydro,
		Durability: 25,
		Mult:       skillDmg[c.TalentLvlSkill()],
	}
	ai.FlatDmg = c.burstDmgBonus(ai.AttackTag)
	snap := c.Snapshot(&ai)

	return (&core.AttackEvent{
		Info:        ai,
		Pattern:     core.NewDefCircHit(5, false, core.TargettableEnemy),
		SourceFrame: c.Core.F,
		Snapshot:    snap,
	})

}

// Helper function that handles damage, healing, and particle components of every tick of her E
func (c *char) skillTick(d *core.AttackEvent) {

	c.Core.Combat.QueueAttackEvent(d, 0)
	c.Core.Health.HealActive(c.Index, skillHealPct[c.TalentLvlSkill()]*c.HPMax+skillHealFlat[c.TalentLvlSkill()])

	// Particles are 0~1 (1:2) on every damage instance
	if c.Core.Rand.Float64() < .6667 {
		c.QueueParticle("kokomi", 1, core.Hydro, 100)
	}

	c.skillLastTick = c.Core.F

	// C2 handling - believe this is an additional instance of flat healing
	// Sangonomiya Kokomi gains the following Healing Bonuses with regard to characters with 50% or less HP via the following methods:
	// Kurage's Oath Bake-Kurage: 4.5% of Kokomi's Max HP.
	if c.Base.Cons >= 2 {
		active := c.Core.Chars[c.Core.ActiveChar]
		if active.HP()/active.MaxHP() <= .5 {
			c.Core.Health.HealActive(c.Index, 0.045*c.HPMax)
		}
	}
}

// Handles repeating skill damage ticks. Split into a separate function as you can only have 1 jellyfish on field at once
// Skill snapshots, so inputs into the function are the originating snapshot
func (c *char) skillTickTask(originalSnapshot *core.AttackEvent, src int) func() {
	return func() {
		c.Core.Log.Debugw("Skill Tick Debug", "frame", c.Core.F, "event", core.LogCharacterEvent, "current dur", c.Core.Status.Duration("kokomiskill"), "skilllastused", c.skillLastUsed, "src", src)
		if c.Core.Status.Duration("kokomiskill") == 0 {
			return
		}

		// Basically stops "old" casts of E from working, and also stops further ticks from that source
		if c.skillLastUsed > src {

			return
		}

		c.skillTick(originalSnapshot)

		c.AddTask(c.skillTickTask(originalSnapshot, src), "kokomi-skill-tick", 120)
	}
}

// Burst - This function only handles initial damage and status setting
// Damage bonus modification is handled in a separate function based on status
/* The might of Watatsumi descends, dealing Hydro DMG to surrounding opponents, before robing Kokomi in a Ceremonial Garment made from the flowing waters of Sangonomiya.
Ceremonial Garment
Sangonomiya Kokomi's Normal Attack, Charged Attack and Bake-Kurage DMG are increased based on her Max HP.When her Normal and Charged Attacks hit opponents, Kokomi will restore HP for all nearby party members, and the amount restored is based on her Max HP.Increases Sangonomiya Kokomi's resistance to interruption and allows her to move on the water's surface.
These effects will be cleared once Sangonomiya Kokomi leaves the field.
*/
func (c *char) Burst(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionBurst, p)

	// TODO: Snapshot timing is not yet known. Assume it's dynamic for now
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Nereid's Ascension",
		AttackTag:  core.AttackTagElementalBurst,
		ICDTag:     core.ICDTagElementalBurst,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Hydro,
		Durability: 50,
		Mult:       0,
	}
	ai.FlatDmg = burstDmg[c.TalentLvlBurst()] * c.HPMax

	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(5, false, core.TargettableEnemy), f, f)

	c.Core.Status.AddStatus("kokomiburst", 10*60)

	// Ascension 1 - reset duration of E Skill and also resnapshots it
	// Should not activate HoD consistent with in game since it is not a skill usage
	if c.Core.Status.Duration("kokomiskill") > 0 {
		// +1 to avoid same frame expiry issues with skill tick
		c.Core.Status.AddStatus("kokomiskill", 12*60+1)
		c.skillLastUsed = c.Core.F - 1
		d := c.createSkillSnapshot()
		// Tick intervals stay the same if duration is refreshed
		c.AddTask(c.skillTickTask(d, c.Core.F), "kokomi-e-ticks", 120-(c.Core.F-c.skillLastTick))
	}

	// C4 attack speed buff
	if c.Base.Cons >= 4 {
		c.AddMod(core.CharStatMod{
			Key: "kokomi-c4",
			Amount: func(a core.AttackTag) ([]float64, bool) {
				val := make([]float64, core.EndStatType)
				val[core.AtkSpd] = 0.1
				if c.Core.Status.Duration("kokomiburst") > 0 {
					return val, true
				}
				return nil, false
			},
			Expiry: c.Core.F + 10*60,
		})
	}

	// Cannot be prefed particles
	c.AddTask(func() {
		c.ConsumeEnergy(57)
	}, "kokomi-q-energy-drain", f)

	c.SetCD(core.ActionBurst, 18*60)
	return f, a
}

// Helper function for determining whether burst damage bonus should apply
// TODO: Technically A4 cannot be snapshotted so it has to be pulled out into an event subscription...
func (c *char) burstDmgBonus(a core.AttackTag) float64 {
	if c.Core.Status.Duration("kokomiburst") == 0 {
		return 0
	}
	a4Bonus := c.Stats[core.Heal] * 0.15
	switch a {
	case core.AttackTagNormal:
		return (burstBonusNormal[c.TalentLvlBurst()] + a4Bonus) * c.HPMax
	case core.AttackTagExtra:
		return (burstBonusCharge[c.TalentLvlBurst()] + a4Bonus) * c.HPMax
	case core.AttackTagElementalArt:
		return burstBonusSkill[c.TalentLvlBurst()] * c.HPMax
	default:
		return 0
	}
}

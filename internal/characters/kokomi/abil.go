package kokomi

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

// Standard attack damage function
// Has "travel" parameter, used to set the number of frames that the projectile is in the air (default = 10)
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
		Element:    core.Hydro,
		Durability: 25,
		Mult:       attack[c.NormalCounter][c.TalentLvlAttack()],
	}
	ai.FlatDmg = c.burstDmgBonus(ai.AttackTag)

	c.Core.Combat.QueueAttack(ai, core.NewDefSingleTarget(1, core.TargettableEnemy), f, f+travel)
	// TODO: Assume that this is not dynamic (snapshot on projectile release)

	if c.NormalCounter == c.NormalHitNum-1 {
		c.c1(f, travel)
	}

	c.AdvanceNormalIndex()

	return f, a
}

func (c *char) c1(f, travel int) {
	if c.Base.Cons == 0 {
		return
	}
	if c.Core.Status.Duration("kokomiburst") == 0 {
		return
	}

	// TODO: Assume that these are 1A (not specified in library)
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "At Water's Edge (C1)",
		AttackTag:  core.AttackTagNone,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Hydro,
		Durability: 25,
		Mult:       0,
	}
	ai.FlatDmg = 0.3 * c.MaxHP()

	// TODO: Is this snapshotted/dynamic?
	c.Core.Combat.QueueAttack(ai, core.NewDefSingleTarget(1, core.TargettableEnemy), f, f+travel)
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
func (c *char) Skill(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionSkill, p)

	// skill duration is ~12.5s
	// Plus 1 to avoid same frame issues with skill ticks
	c.Core.Status.AddStatus("kokomiskill", 12*60+30+1)

	d := c.createSkillSnapshot()

	// You get 1 tick immediately, then 1 tick every 2 seconds for a total of 7 ticks
	c.AddTask(func() { c.skillTick(d) }, "kokomi-e-tick", 24)
	c.AddTask(c.skillTickTask(d, c.Core.F), "kokomi-e-ticks", 24+126)

	c.skillLastUsed = c.Core.F
	c.SetCDWithDelay(core.ActionSkill, 20*60, 20)

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

	// check if skill has burst bonus snapshot
	// max swap frame should be 40 frame before 2nd tick
	if c.swapEarlyF > c.skillLastUsed && c.swapEarlyF < (c.skillLastUsed+120-40) {
		d.Info.FlatDmg = c.skillFlatDmg
	} else {
		d.Info.FlatDmg = c.burstDmgBonus(d.Info.AttackTag)
	}

	maxhp := c.MaxHP()

	c.Core.Combat.QueueAttackEvent(d, 0)
	c.Core.Health.Heal(core.HealInfo{
		Caller:  c.Index,
		Target:  c.Core.ActiveChar,
		Message: "Bake-Kurage",
		Src:     skillHealPct[c.TalentLvlSkill()]*maxhp + skillHealFlat[c.TalentLvlSkill()],
		Bonus:   d.Snapshot.Stats[core.Heal],
	})

	// Particles are 0~1 (1:2) on every damage instance
	if c.Core.Rand.Float64() < .6667 {
		c.QueueParticle("kokomi", 1, core.Hydro, 100)
	}

	// C2 handling - believe this is an additional instance of flat healing
	// Sangonomiya Kokomi gains the following Healing Bonuses with regard to characters with 50% or less HP via the following methods:
	// Kurage's Oath Bake-Kurage: 4.5% of Kokomi's Max HP.
	if c.Base.Cons >= 2 {
		active := c.Core.Chars[c.Core.ActiveChar]
		if active.HP()/active.MaxHP() <= .5 {
			c.Core.Health.Heal(core.HealInfo{
				Caller:  c.Index,
				Target:  c.Core.ActiveChar,
				Message: "The Clouds Like Waves Rippling",
				Src:     0.045 * maxhp,
				Bonus:   c.Stat(core.Heal),
			})
		}
	}
}

// Handles repeating skill damage ticks. Split into a separate function as you can only have 1 jellyfish on field at once
// Skill snapshots, so inputs into the function are the originating snapshot
func (c *char) skillTickTask(originalSnapshot *core.AttackEvent, src int) func() {
	return func() {
		c.Core.Log.NewEvent("Skill Tick Debug", core.LogCharacterEvent, c.Index, "current dur", c.Core.Status.Duration("kokomiskill"), "skilllastused", c.skillLastUsed, "src", src)
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
	ai.FlatDmg = burstDmg[c.TalentLvlBurst()] * c.MaxHP()

	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(5, false, core.TargettableEnemy), f, f)

	c.Core.Status.AddStatus("kokomiburst", 10*60)

	// Ascension 1 - reset duration of E Skill and also resnapshots it
	// Should not activate HoD consistent with in game since it is not a skill usage
	if c.Core.Status.Duration("kokomiskill") > 0 {
		// +1 to avoid same frame expiry issues with skill tick
		c.Core.Status.AddStatus("kokomiskill", 12*60+1)
	}

	// C4 attack speed buff
	if c.Base.Cons >= 4 {
		m := make([]float64, core.EndStatType)
		m[core.AtkSpd] = 0.1
		c.AddMod(core.CharStatMod{
			Key:    "kokomi-c4",
			Expiry: c.Core.F + 10*60,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})
	}

	// Cannot be prefed particles
	c.ConsumeEnergy(57)
	c.SetCDWithDelay(core.ActionBurst, 18*60, 46)
	return f, a
}

// Helper function for determining whether burst damage bonus should apply
func (c *char) burstDmgBonus(a core.AttackTag) float64 {
	if c.Core.Status.Duration("kokomiburst") == 0 {
		return 0
	}
	switch a {
	case core.AttackTagNormal:
		return burstBonusNormal[c.TalentLvlBurst()] * c.MaxHP()
	case core.AttackTagExtra:
		return burstBonusCharge[c.TalentLvlBurst()] * c.MaxHP()
	case core.AttackTagElementalArt:
		return burstBonusSkill[c.TalentLvlBurst()] * c.MaxHP()
	default:
		return 0
	}
}

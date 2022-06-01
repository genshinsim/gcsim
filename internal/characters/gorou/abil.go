package gorou

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

// Normal attack damage queue generator
// relatively standard with no major differences versus other bow characters
// Has "travel" parameter, used to set the number of frames that the arrow is in the air (default = 10)
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
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypePierce,
		Element:    core.Physical,
		Durability: 25,
		Mult:       attack[c.NormalCounter][c.TalentLvlAttack()],
	}
	c.Core.Combat.QueueAttack(ai, core.NewDefSingleTarget(1, core.TargettableEnemy), f, f+travel)

	c.AdvanceNormalIndex()

	return f, a
}

// Aimed charge attack damage queue generator
func (c *char) Aimed(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionAim, p)

	travel, ok := p["travel"]
	if !ok {
		travel = 10
	}
	weakspot, ok := p["weakspot"]

	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Aim Charge Attack",
		AttackTag:  core.AttackTagExtra,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeBlunt,
		Element:    core.Geo,
		//confirmed 25
		Durability:   25,
		Mult:         aimed[c.TalentLvlAttack()],
		HitWeakPoint: weakspot == 1,
	}
	// d.AnimationFrames = f
	c.Core.Combat.QueueAttack(ai, core.NewDefSingleTarget(1, core.TargettableEnemy), f, f+travel)

	return f, a
}

/**
Provides up to 3 buffs to active characters within the skill's AoE based on the number of Geo characters in
the party at the time of casting:
• 1 Geo character: Adds "Standing Firm" - DEF Bonus.
• 2 Geo characters: Adds "Impregnable" - Increased resistance to interruption.
• 3 Geo characters: Adds "Crunch" - Geo DMG Bonus.
Gorou can deploy only 1 General's War Banner on the field at any one time. Characters can only benefit from
1 General's War Banner at a time. When a party member leaves the field, the active buff will last for 2s.
**/
func (c *char) Skill(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionSkill, p)

	c.Core.Tasks.Add(func() {
		ai := core.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Inuzaka All-Round Defense",
			AttackTag:  core.AttackTagElementalArt,
			ICDTag:     core.ICDTagNone,
			ICDGroup:   core.ICDGroupDefault,
			StrikeType: core.StrikeTypeBlunt,
			Element:    core.Geo,
			Durability: 25,
			Mult:       skill[c.TalentLvlSkill()],
		}
		//Inuzaka All-Round Defense: Skill DMG increased by 156% of DEF
		snap := c.Snapshot(&ai)
		ai.FlatDmg = (snap.BaseDef*snap.Stats[core.DEFP] + snap.Stats[core.DEF]) * 1.56

		c.Core.Combat.QueueAttackWithSnap(
			ai,
			snap,
			core.NewDefCircHit(5, false, core.TargettableEnemy),
			//TODO: skill damage frames
			0,
		)
	}, f+10)

	//2 particles apparently
	//TODO: particle frames
	c.QueueParticle(c.Name(), 2, core.Geo, f+100)

	//c6 check
	if c.Base.Cons == 6 {
		c.c6()
	}

	//so it looks like gorou fields works much the same was as bennett field
	//however e field cant be placed if q field still active
	if c.Core.Status.Duration(generalGloryKey) == 0 {

		//TODO: when does ticks start?
		c.eFieldSrc = c.Core.F
		c.Core.Tasks.Add(c.gorouSkillBuffField(c.Core.F), 59) //59 so we get one last tick

		//add a status for general's banner, 10 seconds
		c.Core.Status.AddStatus(generalWarBannerKey, 600)

		if c.Base.Cons >= 4 && c.geoCharCount > 1 {
			//TODO: not sure if this actually snapshots stats
			// ai := core.AttackInfo{
			// 	Abil:      "Inuzaka All-Round Defense C4",
			// 	AttackTag: core.AttackTagNone,
			// }
			stats, _ := c.SnapshotStats()
			c.Core.Tasks.Add(c.gorouSkillHealField(c.Core.F, stats[:]), 90)
		}
	}

	//10s coold down
	c.SetCD(core.ActionSkill, 600)
	return f, a
}

//recursive function for queueing up ticks
func (c *char) gorouSkillBuffField(src int) func() {
	return func() {
		//do nothing if this has been overwritten
		if c.eFieldSrc != src {
			return
		}
		//do nothing if both field expired
		if c.Core.Status.Duration(generalWarBannerKey) == 0 && c.Core.Status.Duration(generalGloryKey) == 0 {
			return
		}
		//do nothing if expired
		//add buff to active char based on number of geo chars
		//ok to overwrite existing mod
		active := c.Core.Chars[c.Core.ActiveChar]
		active.AddMod(core.CharStatMod{
			Key:    defenseBuffKey,
			Expiry: c.Core.F + 126, //2.1s
			Amount: func() ([]float64, bool) {
				return c.gorouBuff, true
			},
		})

		//tick again every second
		c.Core.Tasks.Add(c.gorouSkillBuffField(src), 60)
	}
}

func (c *char) gorouSkillHealField(src int, stats []float64) func() {
	return func() {
		//do nothing if this has been overwritten
		if c.eFieldHealSrc != src {
			return
		}
		//do nothing if field expired
		if c.Core.Status.Duration(generalWarBannerKey) == 0 {
			return
		}
		//When General's Glory is in the "Impregnable" or "Crunch" states, it will also heal active characters
		//within its AoE by 50% of Gorou's own DEF every 1.5s.
		amt := c.Base.Def*(1+stats[core.DEFP]) + stats[core.DEF]
		c.Core.Health.Heal(core.HealInfo{
			Caller:  c.Index,
			Target:  c.Core.ActiveChar,
			Message: "Lapping Hound: Warm as Water",
			Src:     0.5 * amt,
			Bonus:   c.Stat(core.Heal),
		})

		//tick every 1.5s
		c.Core.Tasks.Add(c.gorouSkillBuffField(src), 90)
	}
}

func (c *char) Burst(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionBurst, p)

	c.Core.Tasks.Add(func() {
		ai := core.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Juuga: Forward Unto Victory",
			AttackTag:  core.AttackTagElementalBurst,
			ICDTag:     core.ICDTagNone,
			ICDGroup:   core.ICDGroupDefault,
			Element:    core.Geo,
			StrikeType: core.StrikeTypeBlunt,
			//TODO: don't know the gauge of this
			Durability: 25,
			Mult:       burst[c.TalentLvlSkill()],
		}
		//Juuga: Forward Unto Victory: Skill DMG and Crystal Collapse DMG increased by 15.6% of DEF
		snap := c.Snapshot(&ai)
		ai.FlatDmg = (snap.BaseDef*snap.Stats[core.DEFP] + snap.Stats[core.DEF]) * 0.156

		c.Core.Combat.QueueAttackWithSnap(
			ai,
			snap,
			core.NewDefCircHit(5, false, core.TargettableEnemy),
			//TODO: skill damage frames
			0,
		)
	}, f+10)

	//Like the General's War Banner created by Inuzaka All-Round Defense, provides buffs to active characters
	//within the skill's AoE based on the number of Geo characters in the party. Also moves together with
	//your active character.
	c.eFieldSrc = c.Core.F
	c.Core.Tasks.Add(c.gorouSkillBuffField(c.Core.F), 59) //59 so we get one last tick

	//If a General's War Banner created by Gorou currently exists on the field when this ability is used,
	//it will be destroyed. In addition, for the duration of General's Glory, Gorou's
	//Elemental Skill "Inuzaka All-Round Defense" will not create the General's War Banner.
	c.Core.Status.DeleteStatus(generalWarBannerKey)
	c.Core.Status.AddStatus(generalGloryKey, generalGloryDuration)

	//Generates 1 Crystal Collapse every 1.5s that deals AoE Geo DMG to 1 opponent within the skill's AoE.
	//Pulls 1 elemental shard in the skill's AoE to your active character's position every 1.5s (elemental
	//shards are created by Crystallize reactions).
	c.qFieldSrc = c.Core.F
	c.Core.Tasks.Add(c.gorouCrystalCollapse(c.Core.F), 90) //every 90s?

	//TODO:  If Gorou falls, the effects of General's Glory will be cleared.

	//A1: After using Juuga: Forward Unto Victory, all nearby party members' DEF is increased by 25% for 12s.
	val := make([]float64, core.EndStatType)
	val[core.DEFP] = .25
	for _, char := range c.Core.Chars {
		char.AddMod(core.CharStatMod{
			Key:    heedlessKey,
			Expiry: c.Core.F + 720, //12s
			Amount: func() ([]float64, bool) {
				return val, true
			},
		})
	}

	//c6 check
	if c.Base.Cons == 6 {
		c.c6()
	}

	c.c2Extension = 0

	c.SetCDWithDelay(core.ActionBurst, 20*60, 8)
	c.ConsumeEnergy(8)
	return f, a
}

//recursive function for dealing damage
func (c *char) gorouCrystalCollapse(src int) func() {

	return func() {
		//do nothing if this has been overwritten
		if c.qFieldSrc != src {
			return
		}
		//do nothing if field expired
		if c.Core.Status.Duration(generalGloryKey) == 0 {
			return
		}
		//trigger damage
		ai := core.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Crystal Collapse",
			AttackTag:  core.AttackTagElementalBurst,
			ICDTag:     core.ICDTagElementalBurst,
			ICDGroup:   core.ICDGroupDefault,
			Element:    core.Geo,
			StrikeType: core.StrikeTypeBlunt,
			//TODO: don't know the gauge of this
			Durability: 25,
			Mult:       burstTick[c.TalentLvlSkill()],
		}
		//Juuga: Forward Unto Victory: Skill DMG and Crystal Collapse DMG increased by 15.6% of DEF
		snap := c.Snapshot(&ai)
		ai.FlatDmg = (snap.BaseDef*snap.Stats[core.DEFP] + snap.Stats[core.DEF]) * 0.156

		c.Core.Combat.QueueAttackWithSnap(
			ai,
			snap,
			core.NewDefCircHit(5, false, core.TargettableEnemy),
			//TODO: skill damage frames
			1,
		)

		//tick every 1.5s
		c.Core.Tasks.Add(c.gorouCrystalCollapse(src), 90)
	}
}

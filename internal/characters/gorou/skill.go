package gorou

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player"
)

var skillFrames []int

const skillHitmark = 35

func init() {
	skillFrames = frames.InitAbilSlice(35)
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
func (c *char) Skill(p map[string]int) action.ActionInfo {
	c.Core.Tasks.Add(func() {
		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Inuzaka All-Round Defense",
			AttackTag:  combat.AttackTagElementalArt,
			ICDTag:     combat.ICDTagNone,
			ICDGroup:   combat.ICDGroupDefault,
			StrikeType: combat.StrikeTypeBlunt,
			Element:    attributes.Geo,
			Durability: 25,
			Mult:       skill[c.TalentLvlSkill()],
		}
		//Inuzaka All-Round Defense: Skill DMG increased by 156% of DEF
		snap := c.Snapshot(&ai)
		ai.FlatDmg = (snap.BaseDef*snap.Stats[attributes.DEFP] + snap.Stats[attributes.DEF]) * 1.56

		c.Core.QueueAttackWithSnap(
			ai,
			snap,
			combat.NewDefCircHit(5, false, combat.TargettableEnemy),
			//TODO: skill damage frames
			0,
		)
	}, skillHitmark+10)

	//2 particles apparently
	//TODO: particle frames
	c.Core.QueueParticle("gorou", 2, attributes.Geo, skillHitmark+100)

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
		c.Core.Status.Add(generalWarBannerKey, 600)

		if c.Base.Cons >= 4 && c.geoCharCount > 1 {
			//TODO: not sure if this actually snapshots stats
			// ai := combat.AttackInfo{
			// 	Abil:      "Inuzaka All-Round Defense C4",
			// 	AttackTag: combat.AttackTagNone,
			// }
			stats, _ := c.Stats()
			c.Core.Tasks.Add(c.gorouSkillHealField(c.Core.F, stats[:]), 90)
		}
	}

	//10s coold down
	c.SetCD(action.ActionSkill, 600)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillHitmark,
		State:           action.SkillState,
	}
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
		active := c.Core.Player.ActiveChar()
		active.AddStatMod(defenseBuffKey, 126, attributes.NoStat, func() ([]float64, bool) {
			return c.gorouBuff, true
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
		amt := c.Base.Def*(1+stats[attributes.DEFP]) + stats[attributes.DEF]
		c.Core.Player.Heal(player.HealInfo{
			Caller:  c.Index,
			Target:  c.Core.Player.Active(),
			Message: "Lapping Hound: Warm as Water",
			Src:     0.5 * amt,
			Bonus:   c.Stat(attributes.Heal),
		})

		//tick every 1.5s
		c.Core.Tasks.Add(c.gorouSkillBuffField(src), 90)
	}
}

package gorou

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

var skillFrames []int

const skillHitmark = 34

func init() {
	skillFrames = frames.InitAbilSlice(47) // E -> N1/Q
	skillFrames[action.ActionDash] = 33    // E -> D
	skillFrames[action.ActionJump] = 33    // E -> J
	skillFrames[action.ActionSwap] = 46    // E -> Swap
}

/*
*
Provides up to 3 buffs to active characters within the skill's AoE based on the number of Geo characters in
the party at the time of casting:
• 1 Geo character: Adds "Standing Firm" - DEF Bonus.
• 2 Geo characters: Adds "Impregnable" - Increased resistance to interruption.
• 3 Geo characters: Adds "Crunch" - Geo DMG Bonus.
Gorou can deploy only 1 General's War Banner on the field at any one time. Characters can only benefit from
1 General's War Banner at a time. When a party member leaves the field, the active buff will last for 2s.
*
*/
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

		// A1 Part 1
		// Inuzaka All-Round Defense: Skill DMG increased by 156% of DEF
		snap := c.Snapshot(&ai)
		ai.FlatDmg = (snap.BaseDef*snap.Stats[attributes.DEFP] + snap.Stats[attributes.DEF]) * 1.56

		c.eFieldArea = combat.NewCircleHitOnTarget(c.Core.Combat.Player(), combat.Point{Y: 2}, 8)
		c.Core.QueueAttackWithSnap(
			ai,
			snap,
			combat.NewCircleHitOnTarget(c.eFieldArea.Shape.Pos(), nil, 5),
			0,
		)

		// E
		// so it looks like gorou fields works much the same was as bennett field
		// however e field cant be placed if q field still active
		if c.Core.Status.Duration(generalGloryKey) == 0 {

			c.eFieldSrc = c.Core.F
			c.Core.Tasks.Add(c.gorouSkillBuffField(c.Core.F), 17) // 17 so we get one last tick

			// add a status for general's banner, 10 seconds
			c.Core.Status.Add(generalWarBannerKey, 600)
		}

		// C6
		if c.Base.Cons == 6 {
			c.c6()
		}
	}, skillHitmark)

	// 2 particles apparently
	// TODO: particle frames
	c.Core.QueueParticle("gorou", 2, attributes.Geo, skillHitmark+c.ParticleDelay)

	// 10s cooldown
	c.SetCDWithDelay(action.ActionSkill, 600, 32)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}
}

// recursive function for queueing up ticks
func (c *char) gorouSkillBuffField(src int) func() {
	return func() {
		//do nothing if this has been overwritten
		if c.eFieldSrc != src {
			return
		}
		//do nothing if both field expired
		eActive := c.Core.Status.Duration(generalWarBannerKey) > 0
		qActive := c.Core.Status.Duration(generalGloryKey) > 0
		if !eActive && !qActive {
			return
		}
		// do nothing if only e is up and player is outside of the field area
		// if q is up then the player is always inside of the field area
		if eActive && !qActive && !combat.TargetIsWithinArea(c.Core.Combat.Player(), c.eFieldArea) {
			return
		}

		//add buff to active char based on number of geo chars
		//ok to overwrite existing mod
		active := c.Core.Player.ActiveChar()
		active.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag(defenseBuffKey, 120), // looks like it lasts 2 seconds
			AffectedStat: attributes.NoStat,
			Amount: func() ([]float64, bool) {
				return c.gorouBuff, true
			},
		})

		//looks like tick every 0.3s
		c.Core.Tasks.Add(c.gorouSkillBuffField(src), 18)
	}
}

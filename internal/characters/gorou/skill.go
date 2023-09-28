package gorou

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/targets"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

var skillFrames []int

const (
	skillHitmark   = 34
	particleICDKey = "gorou-particle-icd"
)

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
func (c *char) Skill(p map[string]int) (action.Info, error) {
	c.Core.Tasks.Add(func() {
		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Inuzaka All-Round Defense",
			AttackTag:  attacks.AttackTagElementalArt,
			ICDTag:     attacks.ICDTagNone,
			ICDGroup:   attacks.ICDGroupDefault,
			StrikeType: attacks.StrikeTypeBlunt,
			Element:    attributes.Geo,
			Durability: 25,
			Mult:       skill[c.TalentLvlSkill()],
			FlatDmg:    c.a4Skill(),
		}
		c.eFieldArea = combat.NewCircleHitOnTarget(c.Core.Combat.Player(), geometry.Point{Y: 2}, 8)
		c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.eFieldArea.Shape.Pos(), nil, 5), 0, 0, c.particleCB)

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

	// 10s cooldown
	c.SetCDWithDelay(action.ActionSkill, 600, 32)

	return action.Info{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}, nil
}

func (c *char) particleCB(a combat.AttackCB) {
	if a.Target.Type() != targets.TargettableEnemy {
		return
	}
	if c.StatusIsActive(particleICDKey) {
		return
	}
	c.AddStatus(particleICDKey, 0.2*60, true)
	c.Core.QueueParticle(c.Base.Key.String(), 2, attributes.Geo, c.ParticleDelay)
}

// recursive function for queueing up ticks
func (c *char) gorouSkillBuffField(src int) func() {
	return func() {
		// do nothing if this has been overwritten
		if c.eFieldSrc != src {
			return
		}
		// do nothing if both field expired
		eActive := c.Core.Status.Duration(generalWarBannerKey) > 0
		qActive := c.Core.Status.Duration(generalGloryKey) > 0
		if !eActive && !qActive {
			return
		}
		// do nothing if only e is up and player is outside of the field area
		// if q is up then the player is always inside of the field area
		if eActive && !qActive && !c.Core.Combat.Player().IsWithinArea(c.eFieldArea) {
			return
		}

		// add buff to active char based on number of geo chars
		// ok to overwrite existing mod
		active := c.Core.Player.ActiveChar()
		active.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag(defenseBuffKey, 120), // looks like it lasts 2 seconds
			AffectedStat: attributes.NoStat,
			Amount: func() ([]float64, bool) {
				return c.gorouBuff, true
			},
		})

		// looks like tick every 0.3s
		c.Core.Tasks.Add(c.gorouSkillBuffField(src), 18)
	}
}

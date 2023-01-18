package faruzan

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var skillFrames []int

const (
	VortexAbilName = "Pressurized Collapse"
	skillHitmark   = 14
	skillKey       = "faruzan-e"
	particleICDKey = "faruzan-particle-icd"
	vortexHitmark  = 33
)

func init() {
	skillFrames = frames.InitAbilSlice(35)
	skillFrames[action.ActionSkill] = 34 // E -> E
	skillFrames[action.ActionBurst] = 34 // E -> Q
	skillFrames[action.ActionDash] = 28  // E -> N1
	skillFrames[action.ActionJump] = 27  // E -> J
	skillFrames[action.ActionWalk] = 34  // E -> J
	skillFrames[action.ActionSwap] = 33  // E -> Swap
}

// Faruzan deploys a polyhedron that deals AoE Anemo DMG to nearby opponents.
// She will also enter the Manifest Gale state. While in the Manifest Gale
// state, Faruzan's next fully charged shot will consume this state and will
// become a Hurricane Arrow that deals Anemo DMG to opponents hit. This DMG
// will be considered Charged Attack DMG.
//
// Pressurized Collapse
// The Hurricane Arrow will create a Pressurized Collapse effect at its point
// of impact, applying the Pressurized Collapse effect to the opponent or
// character hit. This effect will be removed after a short delay, creating a
// vortex that deals AoE Anemo DMG and pulls nearby objects and opponents in.
// The vortex DMG is considered Elemental Skill DMG.
func (c *char) Skill(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Wind Realm of Nasamjnin (E)",
		AttackTag:  combat.AttackTagElementalArt,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeDefault,
		Element:    attributes.Anemo,
		Durability: 25,
		Mult:       skill[c.TalentLvlSkill()],
	}
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), combat.Point{Y: 3}, 3),
		skillHitmark,
		skillHitmark,
	)

	// C1: Faruzan can fire off a maximum of 2 Hurricane
	// Arrows using fully charged Aimed Shots while under a
	// single Wind Realm of Nasamjnin effect.
	c.hurricaneCount = 1
	if c.Base.Cons >= 1 {
		c.hurricaneCount = 2
	}

	c.Core.Tasks.Add(func() {
		c.AddStatus(skillKey, 1080, true)
		c.SetCD(action.ActionSkill, 360)
	}, 12)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionJump], // earliest cancel
		State:           action.SkillState,
	}
}

func (c *char) pressurizedCollapse(pos combat.Point) {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       VortexAbilName,
		AttackTag:  combat.AttackTagElementalArt,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeDefault,
		Element:    attributes.Anemo,
		Durability: 25,
		Mult:       vortexDmg[c.TalentLvlSkill()],
	}
	done := false
	particleCb := func(a combat.AttackCB) {
		if done {
			return
		}
		if c.StatusIsActive(particleICDKey) {
			return
		}
		c.Core.QueueParticle("faruzan", 2.0, attributes.Anemo, c.ParticleDelay)
		c.AddStatus(particleICDKey, 330, true)
		done = true
	}
	snap := c.Snapshot(&ai)
	c.Core.Tasks.Add(func() {
		c.Core.QueueAttackWithSnap(
			ai,
			snap,
			combat.NewCircleHitOnTarget(pos, nil, 6),
			0,
			c.makeC4Callback(),
			applyBurstShredCb,
			particleCb)
	}, vortexHitmark)
}

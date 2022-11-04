package faruzan

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var skillFrames []int

const (
	skillHitmark  = 20
	skillKey      = "faruzan-e"
	vortexHitmark = 40
)

func init() {
	skillFrames = frames.InitAbilSlice(52) // E -> D
	skillFrames[action.ActionAttack] = 29  // E -> N1
	skillFrames[action.ActionAim] = 30     // E -> CA
	skillFrames[action.ActionBurst] = 32   // E -> Q
	skillFrames[action.ActionJump] = 51    // E -> J
	skillFrames[action.ActionSwap] = 50    // E -> Swap
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
		ICDTag:     combat.ICDTagElementalArt,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypePierce,
		Element:    attributes.Anemo,
		Durability: 25,
		Mult:       skill[c.TalentLvlSkill()],
	}
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHit(c.Core.Combat.Player(), 2),
		skillHitmark,
		skillHitmark,
		c.a4,
	) // TODO: hitmark and size

	c.AddStatus(skillKey, 1080, false)
	c.SetCDWithDelay(action.ActionSkill, 360, 7) // TODO: check cooldown delay

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionAttack], // earliest cancel
		State:           action.SkillState,
	}
}

func (c *char) hurricaneArrow(travel int, weakspot bool) {
	ai := combat.AttackInfo{
		ActorIndex:           c.Index,
		Abil:                 "Hurricane Arrow",
		AttackTag:            combat.AttackTagExtra,
		ICDTag:               combat.ICDTagNone, // TODO: check ICD
		ICDGroup:             combat.ICDGroupDefault,
		StrikeType:           combat.StrikeTypePierce,
		Element:              attributes.Anemo,
		Durability:           25,
		Mult:                 hurricane[c.TalentLvlAttack()],
		HitWeakPoint:         weakspot,
		HitlagHaltFrames:     .12 * 60, // TODO: check hitlag for special hurricane arrow
		HitlagOnHeadshotOnly: true,
		IsDeployable:         true,
	}

	done := false
	vortexCb := func(a combat.AttackCB) {
		if done {
			return
		}
		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Pressurized Collapse",
			AttackTag:  combat.AttackTagElementalArt,
			ICDTag:     combat.ICDTagElementalArt, // TODO: check ICD
			ICDGroup:   combat.ICDGroupDefault,
			StrikeType: combat.StrikeTypePierce,
			Element:    attributes.Anemo,
			Durability: 25,
			Mult:       hurricane[c.TalentLvlSkill()],
		}
		c.Core.QueueAttack(ai, combat.NewCircleHit(a.Target, 2), vortexHitmark, vortexHitmark) // TODO: hitmark and size
		done = true
	}

	c.Core.QueueAttack(ai, combat.NewDefSingleTarget(c.Core.Combat.DefaultTarget), 0, travel, vortexCb)
}

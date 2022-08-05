package itto

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var skillFrames []int

const skillHitmark = 14

func init() {
	skillFrames = frames.InitAbilSlice(44) // CA0 frames
	skillFrames[action.ActionAttack] = 42
	skillFrames[action.ActionBurst] = 42
	skillFrames[action.ActionDash] = 28
	skillFrames[action.ActionJump] = 28
	skillFrames[action.ActionSwap] = 41
}

func (c *char) Skill(p map[string]int) action.ActionInfo {
	// Added "travel" parameter for future, since Ushi is thrown and takes 12 frames to hit the ground from a press E
	travel, ok := p["travel"]
	if !ok {
		travel = 4
	}

	//deal damage when created
	ai := combat.AttackInfo{
		ActorIndex:       c.Index,
		Abil:             "Masatsu Zetsugi: Akaushi Burst!",
		AttackTag:        combat.AttackTagElementalArt,
		ICDTag:           combat.ICDTagElementalArt,
		ICDGroup:         combat.ICDGroupDefault,
		StrikeType:       combat.StrikeTypeBlunt,
		Element:          attributes.Geo,
		Durability:       25,
		HitlagHaltFrames: 0.02 * 60,
		HitlagFactor:     0.01,
		IsDeployable:     true,
		Mult:             skill[c.TalentLvlSkill()],
	}

	// assume ushi spawn on hitmark not including travel
	c.Core.Tasks.Add(func() {
		c.Core.Constructs.New(c.newUshi(360), true) // 6 seconds from hit/land
	}, skillHitmark)

	done := false
	hitcb := func(a combat.AttackCB) {
		if done {
			return
		}
		done = true

		var count float64 = 3
		if c.Core.Rand.Float64() < 0.33 {
			count = 4
		}
		c.Core.QueueParticle("itto", count, attributes.Geo, c.Core.Flags.ParticleDelay)
	}
	// TODO: snapshot timing
	c.Core.QueueAttack(ai, combat.NewCircleHit(c.Core.Combat.Player(), 1, false, combat.TargettableEnemy), skillHitmark, skillHitmark+travel, hitcb)

	// Assume that Ushi always hits for a stack
	c.addStrStack(1)

	c.SetCDWithDelay(action.ActionSkill, 10*60, 14)

	return action.ActionInfo{
		Frames: func(next action.Action) int {
			// if next is CA1/CA2
			if next == action.ActionCharge && c.Tags[strStackKey] > 0 {
				return 28
			}
			return skillFrames[next]
		},
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}
}

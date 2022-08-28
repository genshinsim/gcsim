package collei

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

const skillKey = "collei-skill"

var (
	skillHitmarks = []int{50, 100} // TODO actual hitmarks
	skillFrames   []int
)

func init() {
	skillFrames = frames.InitAbilSlice(50) // TODO: actual frames
}

func (c *char) Skill(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Floral Brush",
		AttackTag:  combat.AttackTagElementalArt,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeDefault,
		Element:    attributes.Dendro,
		Durability: 25,
		Mult:       skill[c.TalentLvlSkill()],
	}
	snap := c.Snapshot(&ai) // TODO: snapshot timing
	c.Core.QueueAttackWithSnap(
		ai,
		snap,
		combat.NewCircleHit(c.Core.Combat.Player(), 2, false, combat.TargettableEnemy),
		skillHitmarks[0],
	)
	c.AddStatus(skillKey, 480, false) // TODO: find boomerang return frames

	// 50% chance of 3 orbs
	count := 2.0
	if c.Core.Rand.Float64() < .50 {
		count = 3.0
	}
	c.Core.QueueParticle("collei", count, attributes.Cryo, skillHitmarks[0]+c.ParticleDelay)

	c.Core.QueueAttackWithSnap(
		ai,
		snap,
		combat.NewCircleHit(c.Core.Combat.Player(), 2, false, combat.TargettableEnemy),
		skillHitmarks[1],
	)

	c.SetCDWithDelay(action.ActionSkill, 720, 10) // TODO: cd delay

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.InvalidAction], // TODO: fix earliest cancel
		State:           action.SkillState,
	}
}

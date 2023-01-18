package jean

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var skillFrames []int

const skillHitmark = 21

func init() {
	skillFrames = frames.InitAbilSlice(46)
	skillFrames[action.ActionDash] = 28
	skillFrames[action.ActionJump] = 28
	skillFrames[action.ActionSwap] = 45
}

func (c *char) Skill(p map[string]int) action.ActionInfo {
	hold := p["hold"]
	//hold for p up to 5 seconds
	if hold > 300 {
		hold = 300
	}
	hitmark := skillHitmark + hold

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Gale Blade",
		AttackTag:  combat.AttackTagElementalArt,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeDefault,
		Element:    attributes.Anemo,
		Durability: 50,
		Mult:       skill[c.TalentLvlSkill()],
	}
	snap := c.Snapshot(&ai)
	if c.Base.Cons >= 1 && p["hold"] >= 60 {
		c.c1(&snap)
	}

	c.Core.QueueAttackWithSnap(
		ai,
		snap,
		combat.NewBoxHitOnTarget(c.Core.Combat.Player(), nil, 4, 4.1),
		hitmark,
	)

	var count float64 = 2
	if c.Core.Rand.Float64() < 2.0/3.0 {
		count++
	}
	c.Core.QueueParticle("jean", count, attributes.Anemo, hitmark+c.ParticleDelay)

	c.SetCDWithDelay(action.ActionSkill, 360, hitmark-2)

	return action.ActionInfo{
		Frames:          func(next action.Action) int { return skillFrames[next] + hold },
		AnimationLength: skillFrames[action.InvalidAction] + hold,
		CanQueueAfter:   skillFrames[action.ActionDash] + hold, // earliest cancel
		State:           action.SkillState,
	}
}

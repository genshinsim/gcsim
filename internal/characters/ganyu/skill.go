package ganyu

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var skillFrames []int

func init() {
	skillFrames = frames.InitAbilSlice(28)
	skillFrames[action.ActionSwap] = 27
}

func (c *char) Skill(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Ice Lotus",
		AttackTag:  combat.AttackTagElementalArt,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeDefault,
		Element:    attributes.Cryo,
		Durability: 25,
		Mult:       lotus[c.TalentLvlSkill()],
	}

	snap := c.Snapshot(&ai)
	ap := combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 4)
	//flower damage immediately
	c.Core.QueueAttackWithSnap(ai, snap, ap, 13)
	//we get the orbs right away
	c.Core.QueueParticle("ganyu", 2, attributes.Cryo, c.ParticleDelay)

	//flower damage is after 6 seconds
	c.Core.QueueAttackWithSnap(ai, snap, ap, 373)
	c.Core.QueueParticle("ganyu", 2, attributes.Cryo, 373+c.ParticleDelay)

	//add cooldown to sim
	// c.CD[charge] = c.Core.F + 10*60

	if c.Base.Cons == 6 {
		c.Core.Status.Add(c6Key, 1800)
	}

	c.SetCDWithDelay(action.ActionSkill, 600, 10)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionSwap], // earliest cancel
		State:           action.SkillState,
	}
}

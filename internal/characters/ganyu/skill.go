package ganyu

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
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
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Cryo,
		Durability: 25,
		Mult:       lotus[c.TalentLvlSkill()],
	}

	snap := c.Snapshot(&ai)
	ap := combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 4)
	// flower damage immediately
	c.Core.QueueAttackWithSnap(ai, snap, ap, 13, c.makeParticleCB())
	// flower explosion is after 6 seconds
	c.Core.QueueAttackWithSnap(ai, snap, ap, 373, c.makeParticleCB())

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

func (c *char) makeParticleCB() combat.AttackCBFunc {
	done := false
	return func(a combat.AttackCB) {
		if a.Target.Type() != combat.TargettableEnemy {
			return
		}
		if done {
			return
		}
		done = true
		c.Core.QueueParticle(c.Base.Key.String(), 2, attributes.Cryo, c.ParticleDelay)
	}
}

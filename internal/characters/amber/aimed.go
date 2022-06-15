package amber

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var aimedFrames []int

const aimedHitmark = 94

func init() {
	aimedFrames = frames.InitAbilSlice(94)
	aimedFrames[action.ActionDash] = aimedHitmark
	aimedFrames[action.ActionJump] = aimedHitmark
}

func (c *char) Aimed(p map[string]int) action.ActionInfo {
	travel, ok := p["travel"]
	if !ok {
		travel = 10
	}
	weakspot := p["weakspot"]

	b := p["bunny"]

	if c.Base.Cons >= 2 && b != 0 {
		//explode the first bunny
		c.Core.Tasks.Add(func() {
			c.manualExplode()
		}, aimedHitmark+travel)

		//also don't do any dmg since we're shooting at bunny
		return action.ActionInfo{
			Frames:          frames.NewAbilFunc(aimedFrames),
			AnimationLength: aimedFrames[action.InvalidAction],
			CanQueueAfter:   aimedHitmark,

			State: action.AimState,
		}
	}

	ai := combat.AttackInfo{
		Abil:         "Aim (Charged)",
		ActorIndex:   c.Index,
		AttackTag:    combat.AttackTagExtra,
		ICDTag:       combat.ICDTagExtraAttack,
		ICDGroup:     combat.ICDGroupAmber,
		Element:      attributes.Pyro,
		Durability:   50,
		Mult:         aim[c.TalentLvlAttack()],
		HitWeakPoint: weakspot == 1,
	}
	c.Core.QueueAttack(ai, combat.NewDefSingleTarget(1, combat.TargettableEnemy), aimedHitmark, aimedHitmark+travel, c.a4)

	if c.Base.Cons >= 1 {
		ai.Mult = .2 * ai.Mult
		c.Core.QueueAttack(ai, combat.NewDefSingleTarget(1, combat.TargettableEnemy), aimedHitmark, aimedHitmark+travel)
	}

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(aimedFrames),
		AnimationLength: aimedFrames[action.InvalidAction],
		CanQueueAfter:   aimedHitmark,

		State: action.AimState,
	}
}

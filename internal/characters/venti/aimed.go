package venti

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var aimedFrames []int

const aimedHitmark = 86

func init() {
	// TODO: get separate counts for each cancel, currently using generic frames for all of them
	aimedFrames = frames.InitAbilSlice(94)
	aimedFrames[action.ActionDash] = aimedHitmark
	aimedFrames[action.ActionJump] = aimedHitmark
}

func (c *char) Aimed(p map[string]int) action.ActionInfo {
	travel, ok := p["travel"]
	if !ok {
		travel = 10
	}
	weakspot, ok := p["weakspot"]

	ai := combat.AttackInfo{
		ActorIndex:   c.Index,
		Abil:         "Aim (Charged)",
		AttackTag:    combat.AttackTagExtra,
		ICDTag:       combat.ICDTagNone,
		ICDGroup:     combat.ICDGroupDefault,
		Element:      attributes.Anemo,
		Durability:   25,
		Mult:         aim[c.TalentLvlAttack()],
		HitWeakPoint: weakspot == 1,
	}

	c.Core.QueueAttack(ai, combat.NewDefCircHit(.1, false, combat.TargettableEnemy), aimedHitmark, aimedHitmark+travel)
	if c.Base.Cons >= 1 {
		ai.Abil = "Aim (Charged) C1"
		ai.Mult = ai.Mult / 3.0
		c.Core.QueueAttack(ai, combat.NewDefCircHit(.1, false, combat.TargettableEnemy), aimedHitmark, aimedHitmark+travel)
		c.Core.QueueAttack(ai, combat.NewDefCircHit(.1, false, combat.TargettableEnemy), aimedHitmark, aimedHitmark+travel)
	}

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(aimedFrames),
		AnimationLength: aimedFrames[action.InvalidAction],
		CanQueueAfter:   aimedHitmark,
		State:           action.AimState,
	}
}

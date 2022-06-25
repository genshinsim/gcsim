package gorou

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

// Aimed charge attack damage queue generator
func (c *char) Aimed(p map[string]int) action.ActionInfo {
	travel, ok := p["travel"]
	if !ok {
		travel = 10
	}
	weakspot, ok := p["weakspot"]

	ai := combat.AttackInfo{
		ActorIndex:   c.Index,
		Abil:         "Aim Charge Attack",
		AttackTag:    combat.AttackTagExtra,
		ICDTag:       combat.ICDTagNone,
		ICDGroup:     combat.ICDGroupDefault,
		StrikeType:   combat.StrikeTypeBlunt,
		Element:      attributes.Geo,
		Durability:   25,
		Mult:         aimed[c.TalentLvlAttack()],
		HitWeakPoint: weakspot == 1,
	}

	c.Core.QueueAttack(ai, combat.NewDefSingleTarget(1, combat.TargettableEnemy), aimedHitmark, aimedHitmark+travel)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(aimedFrames),
		AnimationLength: aimedFrames[action.InvalidAction],
		CanQueueAfter:   aimedHitmark,
		State:           action.AimState,
	}
}

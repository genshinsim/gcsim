package diona

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var aimedFrames []int
var aimedC4Frames []int

const aimedHitmark = 86
const aimedC4Hitmark = 50

func init() {
	aimedFrames = frames.InitAbilSlice(94)
	aimedFrames[action.ActionDash] = aimedHitmark
	aimedFrames[action.ActionJump] = aimedHitmark

	aimedC4Frames = frames.InitAbilSlice(58)
	aimedC4Frames[action.ActionDash] = aimedC4Hitmark
	aimedC4Frames[action.ActionJump] = aimedC4Hitmark
}

func (c *char) Aimed(p map[string]int) action.ActionInfo {
	travel, ok := p["travel"]
	if !ok {
		travel = 10
	}
	weakspot := p["weakspot"]

	ai := combat.AttackInfo{
		ActorIndex:           c.Index,
		Abil:                 "Aim (Charged)",
		AttackTag:            combat.AttackTagExtra,
		ICDTag:               combat.ICDTagExtraAttack,
		ICDGroup:             combat.ICDGroupDefault,
		StrikeType:           combat.StrikeTypePierce,
		Element:              attributes.Cryo,
		Durability:           25,
		Mult:                 aim[c.TalentLvlAttack()],
		HitWeakPoint:         weakspot == 1,
		HitlagHaltFrames:     0.12 * 60,
		HitlagFactor:         0.01,
		HitlagOnHeadshotOnly: true,
		IsDeployable:         true,
	}

	var a action.ActionInfo

	if c.Base.Cons >= 4 && c.Core.Status.Duration("diona-q") > 0 && combat.TargetIsWithinArea(c.Core.Combat.Player(), c.burstBuffArea) {
		a = action.ActionInfo{
			Frames:          frames.NewAbilFunc(aimedC4Frames),
			AnimationLength: aimedC4Frames[action.InvalidAction],
			CanQueueAfter:   aimedC4Hitmark,
			State:           action.AimState,
		}
	} else {
		a = action.ActionInfo{
			Frames:          frames.NewAbilFunc(aimedFrames),
			AnimationLength: aimedFrames[action.InvalidAction],
			CanQueueAfter:   aimedHitmark,
			State:           action.AimState,
		}

	}

	c.Core.QueueAttack(
		ai,
		combat.NewBoxHit(
			c.Core.Combat.Player(),
			c.Core.Combat.PrimaryTarget(),
			combat.Point{Y: -0.5},
			0.1,
			1,
		),
		a.CanQueueAfter,
		a.CanQueueAfter+travel,
	)

	return a
}

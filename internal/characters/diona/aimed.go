package diona

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
)

var aimedFrames [][]int
var aimedC4Frames []int

var aimedHitmarks = []int{15, 86}

const aimedC4Hitmark = 50

func init() {
	aimedFrames = make([][]int, 2)

	// Aimed Shot
	aimedFrames[0] = frames.InitAbilSlice(23)
	aimedFrames[0][action.ActionDash] = aimedHitmarks[0]
	aimedFrames[0][action.ActionJump] = aimedHitmarks[0]

	// Fully-Charged Aimed Shot
	aimedFrames[1] = frames.InitAbilSlice(94)
	aimedFrames[1][action.ActionDash] = aimedHitmarks[1]
	aimedFrames[1][action.ActionJump] = aimedHitmarks[1]

	// Fully-Charged Aimed Shot (C4)
	aimedC4Frames = frames.InitAbilSlice(58)
	aimedC4Frames[action.ActionDash] = aimedC4Hitmark
	aimedC4Frames[action.ActionJump] = aimedC4Hitmark
}

func (c *char) Aimed(p map[string]int) (action.Info, error) {
	hold, ok := p["hold"]
	if !ok {
		hold = attacks.AimParamLv1
	}
	switch hold {
	case attacks.AimParamPhys:
	case attacks.AimParamLv1:
	default:
		return action.Info{}, fmt.Errorf("invalid hold param supplied, got %v", hold)
	}
	travel, ok := p["travel"]
	if !ok {
		travel = 10
	}
	weakspot := p["weakspot"]

	ai := combat.AttackInfo{
		ActorIndex:           c.Index,
		Abil:                 "Fully-Charged Aimed Shot",
		AttackTag:            attacks.AttackTagExtra,
		ICDTag:               attacks.ICDTagExtraAttack,
		ICDGroup:             attacks.ICDGroupDefault,
		StrikeType:           attacks.StrikeTypePierce,
		Element:              attributes.Cryo,
		Durability:           25,
		Mult:                 fullaim[c.TalentLvlAttack()],
		HitWeakPoint:         weakspot == 1,
		HitlagHaltFrames:     0.12 * 60,
		HitlagFactor:         0.01,
		HitlagOnHeadshotOnly: true,
		IsDeployable:         true,
	}
	if hold < attacks.AimParamLv1 {
		ai.Abil = "Aimed Shot"
		ai.Element = attributes.Physical
		ai.Mult = aim[c.TalentLvlAttack()]
	}

	var a action.Info

	if c.Base.Cons >= 4 && c.Core.Status.Duration("diona-q") > 0 && c.Core.Combat.Player().IsWithinArea(c.burstBuffArea) && hold == attacks.AimParamLv1 {
		ai.Abil += " (C4)"
		a = action.Info{
			Frames:          frames.NewAbilFunc(aimedC4Frames),
			AnimationLength: aimedC4Frames[action.InvalidAction],
			CanQueueAfter:   aimedC4Hitmark,
			State:           action.AimState,
		}
	} else {
		a = action.Info{
			Frames:          frames.NewAbilFunc(aimedFrames[hold]),
			AnimationLength: aimedFrames[hold][action.InvalidAction],
			CanQueueAfter:   aimedHitmarks[hold],
			State:           action.AimState,
		}
	}

	c.Core.QueueAttack(
		ai,
		combat.NewBoxHit(
			c.Core.Combat.Player(),
			c.Core.Combat.PrimaryTarget(),
			geometry.Point{Y: -0.5},
			0.1,
			1,
		),
		a.CanQueueAfter,
		a.CanQueueAfter+travel,
	)

	return a, nil
}

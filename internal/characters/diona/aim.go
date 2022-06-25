package diona

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var aimedFrames []int
var aimedC4Frames []int

const aimAnimationDuration = 84
const aimC4AnimationDuration = 34

func init() {
	aimedFrames = frames.InitAbilSlice(aimAnimationDuration)
	aimedC4Frames = frames.InitAbilSlice(aimC4AnimationDuration)
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
		ICDTag:       combat.ICDTagExtraAttack,
		ICDGroup:     combat.ICDGroupDefault,
		StrikeType:   combat.StrikeTypePierce,
		Element:      attributes.Cryo,
		Durability:   25,
		Mult:         aim[c.TalentLvlAttack()],
		HitWeakPoint: weakspot == 1,
	}

	var a action.ActionInfo

	if c.Base.Cons >= 4 && c.Core.Status.Duration("dionaburst") > 0 {
		a = action.ActionInfo{
			Frames:          frames.NewAbilFunc(aimedC4Frames),
			AnimationLength: aimC4AnimationDuration,
			CanQueueAfter:   aimC4AnimationDuration,
			State:           action.AimState,
		}
	} else {
		a = action.ActionInfo{
			Frames:          frames.NewAbilFunc(aimedFrames),
			AnimationLength: aimAnimationDuration,
			CanQueueAfter:   aimAnimationDuration,
			State:           action.AimState,
		}

	}

	c.Core.QueueAttack(ai,
		combat.NewDefSingleTarget(1, combat.TargettableEnemy),
		a.AnimationLength,
		a.AnimationLength+travel,
	)

	return a
}

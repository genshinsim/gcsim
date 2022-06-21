package yelan

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

var aimedFrames []int
var aimedBarbFrames []int

const aimedHitmark = 86
const aimedBarbHitmark = 32

func init() {
	// TODO: confirm that aim->x is the same for all cancels
	aimedFrames = frames.InitAbilSlice(96)
	aimedFrames[action.ActionDash] = aimedHitmark
	aimedFrames[action.ActionJump] = aimedHitmark

	aimedBarbFrames = frames.InitAbilSlice(42)
	aimedBarbFrames[action.ActionDash] = aimedHitmark
	aimedBarbFrames[action.ActionJump] = aimedHitmark
}

// Aimed charge attack damage queue generator
func (c *char) Aimed(p map[string]int) action.ActionInfo {
	travel, ok := p["travel"]
	if !ok {
		travel = 10
	}
	weakspot := p["weakspot"]

	if c.Tag(breakthroughStatus) > 0 {
		c.RemoveTag(breakthroughStatus)
		c.Core.Log.NewEvent("breakthrough state deleted", glog.LogCharacterEvent, c.Index)

		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Breakthrough Barb",
			AttackTag:  combat.AttackTagExtra,
			ICDTag:     combat.ICDTagYelanBreakthrough,
			ICDGroup:   combat.ICDGroupYelanBreakthrough,
			Element:    attributes.Hydro,
			Durability: 25,
			FlatDmg:    barb[c.TalentLvlAttack()] * c.MaxHP(),
		}
		c.Core.QueueAttack(ai, combat.NewDefCircHit(1, false, combat.TargettableEnemy), aimedBarbHitmark, aimedBarbHitmark+travel)

		return action.ActionInfo{
			Frames:          frames.NewAbilFunc(aimedBarbFrames),
			AnimationLength: aimedBarbFrames[action.InvalidAction],
			CanQueueAfter:   aimedBarbHitmark,
			State:           action.AimState,
		}
	}

	ai := combat.AttackInfo{
		ActorIndex:   c.Index,
		Abil:         "Aim Charge Attack",
		AttackTag:    combat.AttackTagExtra,
		ICDTag:       combat.ICDTagNone,
		ICDGroup:     combat.ICDGroupDefault,
		Element:      attributes.Hydro,
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

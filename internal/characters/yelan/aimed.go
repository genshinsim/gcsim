package yelan

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

var aimedFrames [][]int
var aimedBarbFrames []int

var aimedHitmarks = []int{15 - 12, 15, 86}

const aimedBarbHitmark = 32

func init() {
	aimedFrames = make([][]int, 3)

	// Aimed Shot (ARCC)
	aimedFrames[0] = frames.InitAbilSlice(25 - 12)
	aimedFrames[0][action.ActionDash] = aimedHitmarks[0]
	aimedFrames[0][action.ActionJump] = aimedHitmarks[0]

	// Aimed Shot
	aimedFrames[1] = frames.InitAbilSlice(25)
	aimedFrames[1][action.ActionDash] = aimedHitmarks[1]
	aimedFrames[1][action.ActionJump] = aimedHitmarks[1]

	// Fully-Charged Aimed Shot
	aimedFrames[2] = frames.InitAbilSlice(96)
	aimedFrames[2][action.ActionDash] = aimedHitmarks[2]
	aimedFrames[2][action.ActionJump] = aimedHitmarks[2]

	// Breakthrough Barb
	aimedBarbFrames = frames.InitAbilSlice(42)
	aimedBarbFrames[action.ActionDash] = aimedBarbHitmark
	aimedBarbFrames[action.ActionJump] = aimedBarbHitmark
}

// Aimed charge attack damage queue generator
func (c *char) Aimed(p map[string]int) action.ActionInfo {
	hold, ok := p["hold"]
	if !ok || hold < 0 {
		hold = 2
	}
	if hold > 2 {
		hold = 2
	}
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
		c.Core.QueueAttack(ai, combat.NewCircleHit(c.Core.Combat.Player(), 1, false, combat.TargettableEnemy), aimedBarbHitmark, aimedBarbHitmark+travel)

		return action.ActionInfo{
			Frames:          frames.NewAbilFunc(aimedBarbFrames),
			AnimationLength: aimedBarbFrames[action.InvalidAction],
			CanQueueAfter:   aimedBarbHitmark,
			State:           action.AimState,
		}
	}

	ai := combat.AttackInfo{
		ActorIndex:   c.Index,
		Abil:         "Fully-Charged Aimed Shot",
		AttackTag:    combat.AttackTagExtra,
		ICDTag:       combat.ICDTagNone,
		ICDGroup:     combat.ICDGroupDefault,
		Element:      attributes.Hydro,
		Durability:   25,
		Mult:         fullaim[c.TalentLvlAttack()],
		HitWeakPoint: weakspot == 1,
	}
	if hold < 2 {
		ai.Abil = "Aimed Shot"
		if hold == 0 {
			ai.Abil += " (ARCC)"
		}
		ai.Element = attributes.Physical
		ai.Mult = aim[c.TalentLvlAttack()]
	}
	c.Core.QueueAttack(ai, combat.NewDefSingleTarget(c.Core.Combat.DefaultTarget, combat.TargettableEnemy), aimedHitmarks[hold], aimedHitmarks[hold]+travel)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(aimedFrames[hold]),
		AnimationLength: aimedFrames[hold][action.InvalidAction],
		CanQueueAfter:   aimedHitmarks[hold],
		State:           action.AimState,
	}
}

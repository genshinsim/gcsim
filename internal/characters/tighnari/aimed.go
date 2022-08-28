package tighnari

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var aimedFrames []int
var aimedWreathFrames []int

const aimedHitmark = 86
const aimedWreathHitmark = 175

func init() {
	aimedFrames = frames.InitAbilSlice(94)
	aimedFrames[action.ActionDash] = aimedHitmark
	aimedFrames[action.ActionJump] = aimedHitmark

	aimedWreathFrames = frames.InitAbilSlice(167)
	aimedWreathFrames[action.ActionDash] = aimedWreathHitmark
	aimedWreathFrames[action.ActionJump] = aimedWreathHitmark
}

func (c *char) Aimed(p map[string]int) action.ActionInfo {
	level, ok := p["level"]
	if !ok {
		level = 0
	}

	if c.StatusIsActive(vijnanasuffusionStatus) {
		level = 1
	}
	if level == 1 {
		return c.WreathAimed(p)
	}

	travel, ok := p["travel"]
	if !ok {
		travel = 10
	}
	weakspot := p["weakspot"]

	ai := combat.AttackInfo{
		ActorIndex:   c.Index,
		Abil:         "Aim (Charged)",
		AttackTag:    combat.AttackTagExtra,
		ICDTag:       combat.ICDTagNone,
		ICDGroup:     combat.ICDGroupDefault,
		Element:      attributes.Dendro,
		Durability:   25,
		Mult:         aim[c.TalentLvlAttack()],
		HitWeakPoint: weakspot == 1,
	}

	c.Core.QueueAttack(ai, combat.NewCircleHit(c.Core.Combat.Player(), .1, false, combat.TargettableEnemy), aimedHitmark, aimedHitmark+travel)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(aimedFrames),
		AnimationLength: aimedFrames[action.InvalidAction],
		CanQueueAfter:   aimedHitmark,
		State:           action.AimState,
	}
}

func (c *char) WreathAimed(p map[string]int) action.ActionInfo {
	travel, ok := p["travel"]
	if !ok {
		travel = 10
	}
	wreathTravel, ok := p["wreath"]
	if !ok {
		wreathTravel = 40
	}
	weakspot := p["weakspot"]

	skip := 0
	if c.StatusIsActive(vijnanasuffusionStatus) {
		skip = 2.4 * 60

		arrows := c.Tag(wreatharrows) - 1
		c.SetTag(wreatharrows, arrows)
		if arrows == 0 {
			c.DeleteStatus(vijnanasuffusionStatus)
		}
	}

	c.a1()

	ai := combat.AttackInfo{
		ActorIndex:   c.Index,
		Abil:         "Wreath Arrow",
		AttackTag:    combat.AttackTagExtra,
		ICDTag:       combat.ICDTagNone,
		ICDGroup:     combat.ICDGroupDefault,
		Element:      attributes.Dendro,
		Durability:   25,
		Mult:         wreath[c.TalentLvlAttack()],
		HitWeakPoint: weakspot == 1,
	}
	c.Core.QueueAttack(ai, combat.NewCircleHit(c.Core.Combat.Player(), .1, false, combat.TargettableEnemy), aimedWreathHitmark-skip, aimedWreathHitmark+travel-skip)

	ai = combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Clusterbloom Arrow",
		AttackTag:  combat.AttackTagExtra,
		ICDTag:     combat.ICDTagExtraAttack,
		ICDGroup:   combat.ICDGroupDefault, // combat.ICDGroupTighnari, but need to fix dendro application
		Element:    attributes.Dendro,
		Durability: 25,
		Mult:       clusterbloom[c.TalentLvlAttack()],
	}
	snap := c.Snapshot(&ai)
	for i := 0; i < 4; i++ {
		ai.HitWeakPoint = c.Core.Rand.Float64() < .5 // random
		c.Core.QueueAttackWithSnap(
			ai,
			snap,
			combat.NewCircleHit(c.Core.Combat.Player(), 0.1, false, combat.TargettableEnemy),
			aimedWreathHitmark+travel+wreathTravel-skip,
		)
	}

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(aimedWreathFrames),
		AnimationLength: aimedWreathFrames[action.InvalidAction],
		CanQueueAfter:   aimedWreathHitmark - skip,
		State:           action.AimState,
	}
}

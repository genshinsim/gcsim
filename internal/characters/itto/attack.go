package itto

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var attackFrames [][]int
var attackHitmarks = []int{23, 25, 16, 48}
var attackHitlagHaltFrames = []float64{0.08, 0.08, 0.10, 0.10}

const normalHitNum = 4

func init() {
	attackFrames = make([][]int, normalHitNum)

	// ActionCharge is CA1/CAF frames, while InvalidAction is CA0 expect for the last NA

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0], 41) // CA0 frames
	attackFrames[0][action.ActionAttack] = 33
	attackFrames[0][action.ActionCharge] = 23

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1], 51) // CA0 frames
	attackFrames[1][action.ActionAttack] = 36
	attackFrames[1][action.ActionCharge] = 27

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2], 57) // CA0 frames
	attackFrames[2][action.ActionAttack] = 43
	attackFrames[2][action.ActionCharge] = 21

	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3], 83) // NA frames
	attackFrames[3][action.ActionCharge] = 51
}

func (c *char) Attack(p map[string]int) action.ActionInfo {

	// Additionally, Itto's Normal Attack combo does not immediately reset after sprinting or using his Elemental Skill
	switch c.Core.Player.CurrentState() {
	case action.DashState, action.SkillState:
		c.NormalCounter = c.dasshuCount
	}

	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               fmt.Sprintf("Normal %v", c.NormalCounter),
		AttackTag:          combat.AttackTagNormal,
		ICDTag:             combat.ICDTagNormalAttack,
		ICDGroup:           combat.ICDGroupDefault,
		StrikeType:         combat.StrikeTypeBlunt,
		Element:            attributes.Physical,
		Durability:         25,
		HitlagHaltFrames:   attackHitlagHaltFrames[c.NormalCounter] * 60,
		HitlagFactor:       0.01,
		CanBeDefenseHalted: true,
		Mult:               attack[c.NormalCounter][c.TalentLvlAttack()],
	}

	// Check burst status
	r := 1.0
	if c.StatModIsActive(burstBuffKey) {
		r = 2
	}
	// TODO: hitmark is not getting adjusted for atk speed
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHit(c.Core.Combat.Player(), r, false, combat.TargettableEnemy),
		attackHitmarks[c.NormalCounter],
		attackHitmarks[c.NormalCounter],
	)

	// TODO: assume NAs always hit. since it is not possible to know if the next CA is CA0 or CA1/CAF when deciding what CA frames to return.
	// Add superlative strength stacks on damage
	n := c.NormalCounter // needed for the frames func
	if n == 1 {
		c.addStrStack(1)
	} else if n == 3 {
		c.addStrStack(2)
	}
	if c.StatModIsActive(burstBuffKey) && (n == 0 || n == 2) {
		c.addStrStack(1)
	}

	defer func() {
		c.AdvanceNormalIndex()
		c.dasshuCount = c.NormalCounter
	}()

	return action.ActionInfo{
		Frames: func(next action.Action) int {
			// check if next is CA0. NA4->CA0 doesn't exist
			if next == action.ActionCharge && InvalidSlash.Next(c.Tags[strStackKey]) == SaichiSlash {
				// assume InvalidAction is CA0 frames
				next = action.InvalidAction
				// CA0 after the last NA is illegal. so return 500
				if n == c.NormalHitNum-1 {
					return 500
				}
			}
			return frames.AtkSpdAdjust(attackFrames[n][next], c.Stat(attributes.AtkSpd))
		},
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackHitmarks[c.NormalCounter],
		State:           action.NormalAttackState,
	}
}

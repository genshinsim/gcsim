package klee

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var attackFrames [][]int
var attackHitmarks = []int{16, 23, 37}

const normalHitNum = 3

func init() {
	attackFrames = make([][]int, normalHitNum)
	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0], 34)
	attackFrames[0][action.ActionAttack] = 31
	attackFrames[0][action.ActionCharge] = 23
	attackFrames[0][action.ActionSkill] = 6
	attackFrames[0][action.ActionBurst] = 6
	attackFrames[0][action.ActionDash] = 7
	attackFrames[0][action.ActionJump] = 7
	attackFrames[0][action.ActionWalk] = 34
	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1], 41)
	attackFrames[1][action.ActionAttack] = 38
	attackFrames[1][action.ActionCharge] = 32
	attackFrames[1][action.ActionSkill] = 2
	attackFrames[1][action.ActionBurst] = 2
	attackFrames[1][action.ActionDash] = 2
	attackFrames[1][action.ActionJump] = 2
	attackFrames[1][action.ActionWalk] = 41
	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2], 77)
	attackFrames[2][action.ActionCharge] = 49
	attackFrames[2][action.ActionWalk] = 72
}

func (c *char) Attack(p map[string]int) action.ActionInfo {
	travel, ok := p["travel"]
	if !ok {
		travel = 10
	}

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       fmt.Sprintf("Normal %v", c.NormalCounter),
		AttackTag:  combat.AttackTagNormal,
		ICDTag:     combat.ICDTagKleeFireDamage,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeBlunt,
		Element:    attributes.Pyro,
		Durability: 25,
		Mult:       attack[c.NormalCounter][c.TalentLvlAttack()],
	}

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHit(c.Core.Combat.Player(), 1, false, combat.TargettableEnemy),
		attackHitmarks[c.NormalCounter],
		attackHitmarks[c.NormalCounter]+travel,
		c.a1,
	)
	c.c1(attackHitmarks[c.NormalCounter] + travel)

	defer c.AdvanceNormalIndex()

	return action.ActionInfo{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackHitmarks[c.NormalCounter],
		State:           action.NormalAttackState,
	}
}

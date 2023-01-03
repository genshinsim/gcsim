package wanderer

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var (
	attackFramesNormal   [][]int
	attackFramesE        [][]int
	attackHitmarksNormal = [][]int{{11}, {6}, {32, 41}}
	attackHitmarksE      = [][]int{{15}, {3}, {32, 40}}
	attackRadiusNormal   = []float64{1, 1, 1}
	attackRadiusE        = []float64{2.5, 2.5, 3}
)

const normalHitNum = 3

func init() {
	//TODO: Release = Hitmark? (No Travel Time)
	attackFramesNormal = make([][]int, normalHitNum)

	attackFramesNormal[0] = frames.InitNormalCancelSlice(attackHitmarksNormal[0][0], 35)
	attackFramesNormal[0][action.ActionAttack] = 26
	attackFramesNormal[0][action.ActionCharge] = 24
	attackFramesNormal[0][action.ActionSkill] = 12
	attackFramesNormal[0][action.ActionBurst] = 12
	attackFramesNormal[0][action.ActionDash] = 12

	attackFramesNormal[1] = frames.InitNormalCancelSlice(attackHitmarksNormal[1][0], 39)
	attackFramesNormal[1][action.ActionAttack] = 18
	attackFramesNormal[1][action.ActionCharge] = 27
	attackFramesNormal[1][action.ActionSkill] = 5
	attackFramesNormal[1][action.ActionBurst] = 5
	attackFramesNormal[1][action.ActionDash] = 6
	attackFramesNormal[1][action.ActionJump] = 5
	attackFramesNormal[1][action.ActionSwap] = 5

	attackFramesNormal[2] = frames.InitNormalCancelSlice(attackHitmarksNormal[2][0], 76)
	attackFramesNormal[2][action.ActionAttack] = 64
	attackFramesNormal[2][action.ActionCharge] = 50
	attackFramesNormal[2][action.ActionSkill] = 33
	attackFramesNormal[2][action.ActionBurst] = 33
	attackFramesNormal[2][action.ActionDash] = 34
	attackFramesNormal[2][action.ActionJump] = 34
	attackFramesNormal[2][action.ActionSwap] = 33

	attackFramesE = make([][]int, normalHitNum)

	attackFramesE[0] = frames.InitNormalCancelSlice(attackHitmarksE[0][0], 43)
	attackFramesE[0][action.ActionAttack] = 30
	attackFramesE[0][action.ActionCharge] = 31

	attackFramesE[1] = frames.InitNormalCancelSlice(attackHitmarksE[1][0], 34)
	attackFramesE[1][action.ActionAttack] = 17
	attackFramesE[1][action.ActionCharge] = 23
	attackFramesE[1][action.ActionSkill] = 4
	attackFramesE[1][action.ActionBurst] = 6
	attackFramesE[1][action.ActionDash] = 5
	attackFramesE[1][action.ActionJump] = 5

	attackFramesE[2] = frames.InitNormalCancelSlice(attackHitmarksE[2][0], 70)
	attackFramesE[2][action.ActionAttack] = 54
	attackFramesE[2][action.ActionCharge] = 53
	attackFramesE[2][action.ActionSkill] = 33
	attackFramesE[2][action.ActionBurst] = 33
	attackFramesE[2][action.ActionJump] = 33
}

func (c *char) Attack(p map[string]int) action.ActionInfo {
	delay := c.checkForSkillEnd()

	if c.StatusIsActive(skillKey) {
		// Can only occur if delay == 0, so it can be disregarded
		return c.WindfavoredAttack(p)
	}

	windup := c.attackWindupNormal()

	currentNormalCounter := c.NormalCounter

	for i, mult := range attack[c.NormalCounter] {
		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       fmt.Sprintf("Normal %v", c.NormalCounter),
			AttackTag:  combat.AttackTagNormal,
			ICDTag:     combat.ICDTagNormalAttack,
			ICDGroup:   combat.ICDGroupDefault,
			StrikeType: combat.StrikeTypeDefault,
			Element:    attributes.Anemo,
			Durability: 25,
			Mult:       mult[c.TalentLvlAttack()],
		}
		radius := attackRadiusNormal[c.NormalCounter]

		c.Core.QueueAttack(
			ai,
			combat.NewCircleHit(c.Core.Combat.PrimaryTarget(), radius),
			delay,
			delay+windup+attackHitmarksNormal[c.NormalCounter][i],
			c.makeC6Callback(),
		)
	}

	defer c.AdvanceNormalIndex()

	return action.ActionInfo{
		Frames: func(next action.Action) int {
			return delay + windup +
				frames.AtkSpdAdjust(attackFramesNormal[currentNormalCounter][next], c.Stat(attributes.AtkSpd))
		},
		AnimationLength: delay + windup + attackFramesNormal[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   delay + windup + attackHitmarksNormal[c.NormalCounter][len(attackHitmarksNormal[c.NormalCounter])-1],
		State:           action.NormalAttackState,
	}

}

func (c *char) WindfavoredAttack(p map[string]int) action.ActionInfo {
	windup := c.attackWindupE()

	currentNormalCounter := c.NormalCounter

	for i, mult := range attack[c.NormalCounter] {
		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       fmt.Sprintf("Normal %v (Windfavored)", c.NormalCounter),
			AttackTag:  combat.AttackTagNormal,
			ICDTag:     combat.ICDTagNormalAttack,
			ICDGroup:   combat.ICDGroupDefault,
			StrikeType: combat.StrikeTypeDefault,
			Element:    attributes.Anemo,
			Durability: 25,
			Mult:       mult[c.TalentLvlAttack()],
		}
		radius := attackRadiusE[c.NormalCounter]

		c.Core.QueueAttack(
			ai,
			combat.NewCircleHit(c.Core.Combat.PrimaryTarget(), radius),
			0,
			windup+attackHitmarksE[c.NormalCounter][i],
			c.makeC6Callback(),
		)
	}

	defer c.AdvanceNormalIndex()

	return action.ActionInfo{
		Frames: func(next action.Action) int {
			return windup +
				frames.AtkSpdAdjust(attackFramesE[currentNormalCounter][next], c.Stat(attributes.AtkSpd))
		},
		AnimationLength: windup + attackFramesE[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   windup + attackHitmarksE[c.NormalCounter][len(attackHitmarksE[c.NormalCounter])-1],
		State:           action.NormalAttackState,
	}
}

func (c *char) attackWindupNormal() int {
	switch c.Core.Player.LastAction.Type {
	case action.ActionAttack:
		if c.NormalCounter == 0 {
			return 3
		}
		return 0
	case action.ActionCharge,
		action.ActionBurst:
		return -2
	default:
		return 0
	}
}

func (c *char) attackWindupE() int {
	switch c.Core.Player.LastAction.Type {
	case action.ActionDash:
		return -3
	case action.ActionCharge,
		action.ActionJump:
		return -2
	default:
		return 0
	}
}

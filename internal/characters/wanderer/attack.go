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
	attackRadius         = []float64{1.8, 1.8, 2.2}
)

const normalHitNum = 3

func init() {
	//TODO: Release = Hitmark? (No Travel Time)
	attackFramesNormal = make([][]int, normalHitNum)

	attackFramesNormal[0] = frames.InitNormalCancelSlice(attackHitmarksNormal[0][0], 35)
	attackFramesNormal[0][action.ActionAttack] = 25
	attackFramesNormal[0][action.ActionCharge] = 25
	attackFramesNormal[0][action.ActionSkill] = 12
	attackFramesNormal[0][action.ActionBurst] = 12
	attackFramesNormal[0][action.ActionDash] = 12
	attackFramesNormal[0][action.ActionJump] = 11
	attackFramesNormal[0][action.ActionSwap] = 11

	attackFramesNormal[1] = frames.InitNormalCancelSlice(attackHitmarksNormal[1][0], 40)
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
	attackFramesE[0][action.ActionSkill] = 15
	attackFramesE[0][action.ActionBurst] = 15
	attackFramesE[0][action.ActionDash] = 15
	attackFramesE[0][action.ActionJump] = 15

	attackFramesE[1] = frames.InitNormalCancelSlice(attackHitmarksE[1][0], 34)
	attackFramesE[1][action.ActionAttack] = 17
	attackFramesE[1][action.ActionCharge] = 23
	attackFramesE[1][action.ActionSkill] = 4
	attackFramesE[1][action.ActionBurst] = 6
	attackFramesE[1][action.ActionDash] = 5
	attackFramesE[1][action.ActionJump] = 34

	attackFramesE[2] = frames.InitNormalCancelSlice(attackHitmarksE[2][0], 70)
	attackFramesE[2][action.ActionAttack] = 54
	attackFramesE[2][action.ActionCharge] = 53
	attackFramesE[2][action.ActionSkill] = 33
	attackFramesE[2][action.ActionBurst] = 33
	attackFramesE[2][action.ActionDash] = 32
	attackFramesE[2][action.ActionJump] = 33
}

func (c *char) Attack(p map[string]int) action.ActionInfo {
	delay := c.checkForSkillEnd()
	frameChange := 0

	currentNormalCounter := c.NormalCounter

	relevantHitmarks := attackHitmarksNormal
	relevantFrames := attackFramesNormal

	if c.Core.Player.LastAction.Char == c.Index &&
		c.Core.Player.LastAction.Type == action.ActionAttack &&
		currentNormalCounter == 0 {
		frameChange = 3
	}

	if c.Core.Player.LastAction.Char == c.Index &&
		(c.Core.Player.LastAction.Type == action.ActionCharge || c.Core.Player.LastAction.Type == action.ActionBurst) &&
		currentNormalCounter == 0 {
		frameChange = -2
	}

	if c.StatusIsActive(skillKey) {
		relevantHitmarks = attackHitmarksE
		relevantFrames = attackFramesE

		if c.Core.Player.LastAction.Char == c.Index &&
			c.Core.Player.LastAction.Type == action.ActionDash &&
			currentNormalCounter == 0 {
			frameChange = -3
		}

		if c.Core.Player.LastAction.Char == c.Index &&
			(c.Core.Player.LastAction.Type == action.ActionCharge || c.Core.Player.LastAction.Type == action.ActionJump) &&
			currentNormalCounter == 0 {
			frameChange = -2
		}
	}

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
		radius := attackRadius[c.NormalCounter]

		c.Core.QueueAttack(
			ai,
			combat.NewCircleHit(c.Core.Combat.Player(), radius),
			delay,
			delay+frameChange+relevantHitmarks[c.NormalCounter][i],
		)
	}

	defer c.AdvanceNormalIndex()

	return action.ActionInfo{
		Frames: func(next action.Action) int {
			return delay + frameChange +
				frames.AtkSpdAdjust(relevantFrames[currentNormalCounter][next], c.Stat(attributes.AtkSpd))
		},
		AnimationLength: delay + frameChange + relevantFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   delay + frameChange + relevantHitmarks[c.NormalCounter][len(relevantHitmarks[c.NormalCounter])-1],
		State:           action.NormalAttackState,
	}

}

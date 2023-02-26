package wanderer

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var (
	attackFramesNormal  [][]int
	attackFramesE       [][]int
	attackReleaseNormal = [][]int{{11}, {6}, {32, 41}}
	attackReleaseE      = [][]int{{15}, {3}, {32, 40}}
	attackRadiusE       = []float64{2.5, 2.5, 3}
)

const normalHitNum = 3

func init() {
	attackFramesNormal = make([][]int, normalHitNum)

	attackFramesNormal[0] = frames.InitNormalCancelSlice(attackReleaseNormal[0][0], 35)
	attackFramesNormal[0][action.ActionAttack] = 26
	attackFramesNormal[0][action.ActionCharge] = 24
	attackFramesNormal[0][action.ActionSkill] = 12
	attackFramesNormal[0][action.ActionBurst] = 12
	attackFramesNormal[0][action.ActionDash] = 12

	attackFramesNormal[1] = frames.InitNormalCancelSlice(attackReleaseNormal[1][0], 39)
	attackFramesNormal[1][action.ActionAttack] = 18
	attackFramesNormal[1][action.ActionCharge] = 27
	attackFramesNormal[1][action.ActionSkill] = 5
	attackFramesNormal[1][action.ActionBurst] = 5
	attackFramesNormal[1][action.ActionDash] = 6
	attackFramesNormal[1][action.ActionJump] = 5
	attackFramesNormal[1][action.ActionSwap] = 5

	attackFramesNormal[2] = frames.InitNormalCancelSlice(attackReleaseNormal[2][0], 76)
	attackFramesNormal[2][action.ActionAttack] = 64
	attackFramesNormal[2][action.ActionCharge] = 50
	attackFramesNormal[2][action.ActionSkill] = 33
	attackFramesNormal[2][action.ActionBurst] = 33
	attackFramesNormal[2][action.ActionDash] = 34
	attackFramesNormal[2][action.ActionJump] = 34
	attackFramesNormal[2][action.ActionSwap] = 33

	attackFramesE = make([][]int, normalHitNum)

	attackFramesE[0] = frames.InitNormalCancelSlice(attackReleaseE[0][0], 43)
	attackFramesE[0][action.ActionAttack] = 30
	attackFramesE[0][action.ActionCharge] = 31

	attackFramesE[1] = frames.InitNormalCancelSlice(attackReleaseE[1][0], 34)
	attackFramesE[1][action.ActionAttack] = 17
	attackFramesE[1][action.ActionCharge] = 23
	attackFramesE[1][action.ActionSkill] = 4
	attackFramesE[1][action.ActionBurst] = 6
	attackFramesE[1][action.ActionDash] = 5
	attackFramesE[1][action.ActionJump] = 5

	attackFramesE[2] = frames.InitNormalCancelSlice(attackReleaseE[2][0], 70)
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

	windup := c.attackWindupNormal() + delay

	currentNormalCounter := c.NormalCounter

	travel, ok := p["travel"]
	if !ok {
		travel = 5
	}

	for i, mult := range attack[c.NormalCounter] {
		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       fmt.Sprintf("Normal %v", c.NormalCounter),
			AttackTag:  attacks.AttackTagNormal,
			ICDTag:     attacks.ICDTagNormalAttack,
			ICDGroup:   attacks.ICDGroupDefault,
			StrikeType: attacks.StrikeTypeDefault,
			Element:    attributes.Anemo,
			Durability: 25,
			Mult:       mult[c.TalentLvlAttack()],
		}

		release := windup + attackReleaseNormal[c.NormalCounter][i]

		c.Core.QueueAttack(
			ai,
			combat.NewCircleHit(c.Core.Combat.Player(), c.Core.Combat.PrimaryTarget(), nil, 1),
			release,
			release+travel,
			c.makeA4CB(),
			c.makeA1ElectroCB(),
			c.makeC6Callback(),
			c.particleCB,
		)
	}

	defer c.AdvanceNormalIndex()

	return action.ActionInfo{
		Frames: func(next action.Action) int {
			return windup +
				frames.AtkSpdAdjust(attackFramesNormal[currentNormalCounter][next], c.Stat(attributes.AtkSpd))
		},
		AnimationLength: windup + attackFramesNormal[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   windup + attackReleaseNormal[c.NormalCounter][len(attackReleaseNormal[c.NormalCounter])-1],
		State:           action.NormalAttackState,
	}

}

func (c *char) WindfavoredAttack(p map[string]int) action.ActionInfo {
	// TODO: E can expire during N3, not implemented yet

	windup := c.attackWindupE()

	currentNormalCounter := c.NormalCounter

	travel, ok := p["travel"]
	if !ok {
		travel = 5
	}

	for i, mult := range attack[c.NormalCounter] {
		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       fmt.Sprintf("Normal %v (Windfavored)", c.NormalCounter),
			AttackTag:  attacks.AttackTagNormal,
			ICDTag:     attacks.ICDTagNormalAttack,
			ICDGroup:   attacks.ICDGroupDefault,
			StrikeType: attacks.StrikeTypeDefault,
			Element:    attributes.Anemo,
			Durability: 25,
			Mult:       skillNABonus[c.TalentLvlSkill()] * mult[c.TalentLvlAttack()],
		}
		radius := attackRadiusE[c.NormalCounter]

		release := windup + attackReleaseE[c.NormalCounter][i]

		c.Core.QueueAttack(
			ai,
			combat.NewCircleHit(c.Core.Combat.Player(), c.Core.Combat.PrimaryTarget(), nil, radius),
			release,
			release+travel,
			c.makeA4CB(),
			c.makeA1ElectroCB(),
			c.makeC6Callback(),
			c.particleCB,
		)
	}

	defer c.AdvanceNormalIndex()

	return action.ActionInfo{
		Frames: func(next action.Action) int {
			return windup +
				frames.AtkSpdAdjust(attackFramesE[currentNormalCounter][next], c.Stat(attributes.AtkSpd))
		},
		AnimationLength: windup + attackFramesE[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   windup + attackReleaseE[c.NormalCounter][len(attackReleaseE[c.NormalCounter])-1],
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

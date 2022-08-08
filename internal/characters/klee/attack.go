package klee

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
)

var (
	attackFrames    [][]int
	attackHitmarks  = []int{16, 23, 37}
	attackCancelKey = "klee-attack-cancel"
)

const normalHitNum = 3

func init() {
	attackFrames = make([][]int, normalHitNum)
	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0], 34)
	attackFrames[0][action.ActionAttack] = 31
	attackFrames[0][action.ActionCharge] = 23
	attackFrames[0][action.ActionSkill] = 6
	attackFrames[0][action.ActionBurst] = 6
	attackFrames[0][action.ActionDash] = 6
	attackFrames[0][action.ActionJump] = 6
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

	doDamage := func() {
		c.Core.QueueAttack(
			ai,
			combat.NewCircleHit(c.Core.Combat.Player(), 1, false, combat.TargettableEnemy),
			0,
			travel,
			c.a1,
		)
		c.c1(travel)
	}
	earlyTrigger := false
	c.Core.Events.Subscribe(event.OnStateChange, func(args ...interface{}) bool {
		if earlyTrigger {
			return false
		}
		switch args[1].(action.AnimationState) {
		case action.SkillState,
			action.BurstState,
			action.DashState,
			action.JumpState:
			doDamage()
			earlyTrigger = true
		}
		return false
	}, attackCancelKey)
	animationLag := func() int {
		lastAction := &c.Core.Player.LastAction
		if lastAction.Char == c.Index {
			switch lastAction.Type { // if Klee does either of these, N1 will take 9f longer
			case action.ActionDash,
				action.ActionSkill,
				action.ActionBurst:
				return 9
			}
		}
		return 0
	}()
	c.Core.Tasks.Add(func() {
		c.Core.Events.Unsubscribe(event.OnStateChange, attackCancelKey)
		if earlyTrigger {
			return
		}
		doDamage()
	}, attackHitmarks[c.NormalCounter]+animationLag)

	defer c.AdvanceNormalIndex()

	adjustedFrames := attackFrames
	if animationLag > 0 {
		adjustedFrames = make([][]int, len(attackFrames))
		for i := range attackFrames {
			adjustedFrames[i] = make([]int, len(attackFrames[i]))
			copy(adjustedFrames[i], attackFrames[i])
		}
		for i := range attackFrames[0] {
			switch action.Action(i) {
			case action.ActionBurst,
				action.ActionDash,
				action.ActionJump,
				action.ActionSkill:
			default:
				adjustedFrames[0][i] += animationLag
			}
		}
	}

	return action.ActionInfo{
		Frames:          frames.NewAttackFunc(c.Character, adjustedFrames),
		AnimationLength: adjustedFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   0,
		State:           action.NormalAttackState,
	}
}

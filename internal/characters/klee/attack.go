package klee

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var (
	attackFrames          [][]int
	attackFramesWithLag   [][]int
	attackHitmarks        = []int{16, 23, 37}
	attackHitmarksWithLag []int
	attackCancelKey       = "klee-attack-cancel"
)

const normalHitNum = 3

func init() {
	attackHitmarksWithLag = make([]int, len(attackHitmarks))
	copy(attackHitmarksWithLag, attackHitmarks)
	for i := range attackHitmarksWithLag {
		attackHitmarksWithLag[i] += 9
	}
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
	attackFramesWithLag = make([][]int, len(attackFrames))
	for i := range attackFrames {
		attackFramesWithLag[i] = make([]int, len(attackFrames[i]))
		copy(attackFramesWithLag[i], attackFrames[i])
	}
	add9FrameLag(attackFramesWithLag[0])
}

// klee has 9f lag on normals when using after skill/dash/burst
func add9FrameLag(frames []int) {
	for i := range frames {
		switch action.Action(i) {
		case action.ActionBurst,
			action.ActionDash,
			action.ActionJump,
			action.ActionSkill:
		default:
			frames[i] += 9
		}
	}
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

	performAttack := func() {
		c.Core.QueueAttack(
			ai,
			combat.NewCircleHit(c.Core.Combat.Player(), 1, false, combat.TargettableEnemy),
			0,
			travel,
			c.a1,
		)
		c.c1(travel)
	}

	defer c.AdvanceNormalIndex()

	adjustedFrames := attackFrames
	adjustedHitmarks := attackHitmarks
	switch c.Core.Player.CurrentState() {
	case action.DashState,
		action.SkillState,
		action.BurstState:
		adjustedFrames = attackFramesWithLag
		adjustedHitmarks = attackHitmarksWithLag
	}

	actionInfo := action.ActionInfo{
		Frames:          frames.NewAttackFunc(c.Character, adjustedFrames),
		AnimationLength: adjustedFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   0,
		State:           action.NormalAttackState,
		OnRemoved: func(next action.AnimationState) {
			switch next {
			case action.SkillState,
				action.BurstState,
				action.DashState,
				action.JumpState:
				performAttack()
			}
		},
	}
	actionInfo.QueueAction(performAttack, adjustedHitmarks[c.NormalCounter])
	return actionInfo
}

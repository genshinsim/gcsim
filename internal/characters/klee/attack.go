package klee

import (
	"fmt"
	"math"

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
	attackRadius          = []float64{1, 1, 1.5}
)

const normalHitNum = 3

func init() {
	attackHitmarksWithLag = make([]int, len(attackHitmarks))
	copy(attackHitmarksWithLag, attackHitmarks)
	for i := range attackHitmarksWithLag {
		attackHitmarksWithLag[i] += 9
	}
	attackFrames = make([][]int, normalHitNum)

	// N1 -> x
	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0], 34)
	attackFrames[0][action.ActionAttack] = 31
	attackFrames[0][action.ActionCharge] = 23
	attackFrames[0][action.ActionSkill] = 6
	attackFrames[0][action.ActionBurst] = 6
	attackFrames[0][action.ActionDash] = 6
	attackFrames[0][action.ActionJump] = 6
	attackFrames[0][action.ActionWalk] = 34

	// N2 -> x
	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1], 41)
	attackFrames[1][action.ActionAttack] = 38
	attackFrames[1][action.ActionCharge] = 32
	attackFrames[1][action.ActionSkill] = 2
	attackFrames[1][action.ActionBurst] = 2
	attackFrames[1][action.ActionDash] = 2
	attackFrames[1][action.ActionJump] = 2
	attackFrames[1][action.ActionWalk] = 41

	// N3 -> x
	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2], 77)
	attackFrames[2][action.ActionCharge] = 49
	attackFrames[2][action.ActionWalk] = 72

	// N1 -> x (9f lag)
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

	done := false
	tryPerformAttack := func() {
		if done {
			return
		}
		c.Core.QueueAttack(
			ai,
			combat.NewCircleHit(
				c.Core.Combat.Player(),
				c.Core.Combat.PrimaryTarget(),
				nil,
				attackRadius[c.NormalCounter],
			),
			0,
			travel,
			c.makeA1CB(),
		)
		c.c1(travel)
		done = true
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

	canQueueAfter := math.MaxInt
	for _, f := range adjustedFrames[c.NormalCounter] {
		if f < canQueueAfter {
			canQueueAfter = f
		}
	}
	actionInfo := action.ActionInfo{
		Frames:          frames.NewAttackFunc(c.Character, adjustedFrames),
		AnimationLength: adjustedFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   canQueueAfter,
		State:           action.NormalAttackState,
		OnRemoved: func(next action.AnimationState) {
			switch next {
			case action.SkillState,
				action.BurstState,
				action.DashState,
				action.JumpState:
				tryPerformAttack()
			}
		},
	}
	actionInfo.QueueAction(tryPerformAttack, adjustedHitmarks[c.NormalCounter])
	return actionInfo
}

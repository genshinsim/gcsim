package raiden

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

const normalHitNum = 5

var (
	attackFrames   [][]int
	attackHitmarks = [][]int{{14}, {9}, {14}, {14, 27}, {34}}
	// same between polearm and burst attacks so just use these arrays for both
	attackHitlagHaltFrame = [][]float64{{0.02}, {0.02}, {0.02}, {0, 0}, {0.02}}
	attackDefHalt         = [][]bool{{true}, {true}, {true}, {false, false}, {true}}
	attackRadius          = []float64{1.66, 2.5, 2.19, 2.8, 3}
)

func init() {
	// NA cancels (polearm)
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0][0], 24)
	attackFrames[0][action.ActionAttack] = 18

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1][0], 26)
	attackFrames[1][action.ActionAttack] = 13

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2][0], 36)
	attackFrames[2][action.ActionAttack] = 26

	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3][1], 57)
	attackFrames[3][action.ActionAttack] = 41

	attackFrames[4] = frames.InitNormalCancelSlice(attackHitmarks[4][0], 50)
	attackFrames[4][action.ActionCharge] = 500 //TODO: this action is illegal; need better way to handle it
}

func (c *char) Attack(p map[string]int) action.ActionInfo {
	if c.StatusIsActive(BurstKey) {
		return c.swordAttack(p)
	}

	for i, mult := range attack[c.NormalCounter] {
		ai := combat.AttackInfo{
			ActorIndex:         c.Index,
			Abil:               fmt.Sprintf("Normal %v", c.NormalCounter),
			AttackTag:          combat.AttackTagNormal,
			ICDTag:             combat.ICDTagNormalAttack,
			ICDGroup:           combat.ICDGroupDefault,
			StrikeType:         combat.StrikeTypeSlash,
			Element:            attributes.Physical,
			Durability:         25,
			HitlagHaltFrames:   attackHitlagHaltFrame[c.NormalCounter][i] * 60,
			HitlagFactor:       0.01,
			CanBeDefenseHalted: attackDefHalt[c.NormalCounter][i],
		}
		ai.Mult = mult[c.TalentLvlAttack()]
		radius := attackRadius[c.NormalCounter]
		c.QueueCharTask(func() {
			c.Core.QueueAttack(ai, combat.NewCircleHit(c.Core.Combat.Player(), radius), 0, 0)
		}, attackHitmarks[c.NormalCounter][i])
	}

	defer c.AdvanceNormalIndex()

	return action.ActionInfo{
		Frames:              frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength:     attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:       attackHitmarks[c.NormalCounter][len(attackHitmarks[c.NormalCounter])-1],
		State:               action.NormalAttackState,
		FramePausedOnHitlag: c.FramePausedOnHitlag,
	}
}

var (
	swordFrames   [][]int
	swordHitmarks = [][]int{{12}, {13}, {11}, {22, 33}, {33}}
	swordRadius   = [][]float64{{3.06}, {5.39}, {2.68}, {6.41, 4.8}, {6.15}}
)

func init() {
	// NA cancels (burst)
	swordFrames = make([][]int, normalHitNum)

	swordFrames[0] = frames.InitNormalCancelSlice(swordHitmarks[0][0], 24)
	swordFrames[0][action.ActionAttack] = 19

	swordFrames[1] = frames.InitNormalCancelSlice(swordHitmarks[1][0], 26)
	swordFrames[1][action.ActionAttack] = 16

	swordFrames[2] = frames.InitNormalCancelSlice(swordHitmarks[2][0], 34)
	swordFrames[2][action.ActionAttack] = 16

	swordFrames[3] = frames.InitNormalCancelSlice(swordHitmarks[3][1], 67)
	swordFrames[3][action.ActionAttack] = 44

	swordFrames[4] = frames.InitNormalCancelSlice(swordHitmarks[4][0], 59)
	swordFrames[4][action.ActionCharge] = 500 //TODO: this action is illegal; need better way to handle it
}

func (c *char) swordAttack(p map[string]int) action.ActionInfo {
	for i, mult := range attackB[c.NormalCounter] {
		ai := combat.AttackInfo{
			ActorIndex:         c.Index,
			Abil:               fmt.Sprintf("Musou Isshin %v", c.NormalCounter),
			AttackTag:          combat.AttackTagElementalBurst,
			ICDTag:             combat.ICDTagNormalAttack,
			ICDGroup:           combat.ICDGroupDefault,
			StrikeType:         combat.StrikeTypeSlash,
			Element:            attributes.Electro,
			Durability:         25,
			HitlagHaltFrames:   attackHitlagHaltFrame[c.NormalCounter][i] * 60,
			HitlagFactor:       0.01,
			CanBeDefenseHalted: attackDefHalt[c.NormalCounter][i],
		}
		if c.NormalCounter == 2 {
			ai.StrikeType = combat.StrikeTypeSpear
		}
		// Sword hits are dynamic - group snapshots with damage proc
		ai.Mult = mult[c.TalentLvlBurst()]
		ai.Mult += resolveBonus[c.TalentLvlBurst()] * c.stacksConsumed
		if c.Base.Cons >= 2 {
			ai.IgnoreDefPercent = .6
		}
		radius := swordRadius[c.NormalCounter][i]
		c.QueueCharTask(func() {
			c.Core.QueueAttack(ai, combat.NewCircleHit(c.Core.Combat.Player(), radius), 0, 0, c.burstRestorefunc, c.c6)
		}, swordHitmarks[c.NormalCounter][i])
	}

	defer c.AdvanceNormalIndex()

	return action.ActionInfo{
		Frames:          frames.NewAttackFunc(c.Character, swordFrames),
		AnimationLength: swordFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   swordHitmarks[c.NormalCounter][len(swordHitmarks[c.NormalCounter])-1],
		State:           action.NormalAttackState,
	}
}

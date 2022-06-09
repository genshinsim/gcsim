package raiden

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var attackFrames [][]int
var attackHitmarks = [][]int{{14}, {9}, {14}, {14, 27}, {34}}

func initAttackFrames() {
	// NA cancels
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
	if c.Core.Status.Duration("raidenburst") > 0 {
		return c.swordAttack(p)
	}

	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               fmt.Sprintf("Normal %v", c.NormalCounter),
		AttackTag:          combat.AttackTagNormal,
		ICDTag:             combat.ICDTagNormalAttack,
		ICDGroup:           combat.ICDGroupDefault,
		Element:            attributes.Physical,
		Durability:         25,
		HitlagHaltFrames:   10, //all raiden normals have 0.02s hitlag
		HitlagFactor:       0.01,
		CanBeDefenseHalted: true,
	}

	act := action.ActionInfo{
		Frames:              frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength:     attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:       attackHitmarks[c.NormalCounter][len(attackHitmarks[c.NormalCounter])-1],
		Post:                attackHitmarks[c.NormalCounter][len(attackHitmarks[c.NormalCounter])-1],
		State:               action.NormalAttackState,
		FramePausedOnHitlag: c.FramePausedOnHitlag,
	}

	for i, mult := range attack[c.NormalCounter] {
		ai.Mult = mult[c.TalentLvlAttack()]
		act.QueueAction(func() {
			c.Core.QueueAttack(ai, combat.NewDefCircHit(0.5, false, combat.TargettableEnemy), 0, 0)
		}, attackHitmarks[c.NormalCounter][i])
	}

	defer c.AdvanceNormalIndex()

	return act
}

var swordFrames [][]int
var swordHitmarks = [][]int{{12}, {13}, {11}, {22, 33}, {33}}

func initSwordFrames() {
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
	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               fmt.Sprintf("Musou Isshin %v", c.NormalCounter),
		AttackTag:          combat.AttackTagElementalBurst,
		ICDTag:             combat.ICDTagNormalAttack,
		ICDGroup:           combat.ICDGroupDefault,
		Element:            attributes.Electro,
		Durability:         25,
		HitlagHaltFrames:   12, //all raiden normals have 0.2s hitlag
		HitlagFactor:       0.01,
		CanBeDefenseHalted: true,
	}
	act := action.ActionInfo{
		Frames:          frames.NewAttackFunc(c.Character, swordFrames),
		AnimationLength: swordFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   swordHitmarks[c.NormalCounter][len(swordHitmarks[c.NormalCounter])-1],
		Post:            swordHitmarks[c.NormalCounter][len(swordHitmarks[c.NormalCounter])-1],
		State:           action.NormalAttackState,
	}

	for i, mult := range attackB[c.NormalCounter] {
		// Sword hits are dynamic - group snapshots with damage proc
		ai.Mult = mult[c.TalentLvlBurst()]
		ai.Mult += resolveBonus[c.TalentLvlBurst()] * c.stacksConsumed
		if c.Base.Cons >= 2 {
			ai.IgnoreDefPercent = .6
		}
		act.QueueAction(func() {
			c.Core.QueueAttack(ai, combat.NewDefCircHit(2, false, combat.TargettableEnemy), 0, 0, c.burstRestorefunc, c.c6)
		}, swordHitmarks[c.NormalCounter][i])
	}

	defer c.AdvanceNormalIndex()

	return act
}

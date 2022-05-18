package raiden

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var attackFrames [][]int
var hitmarks = [][]int{{14}, {9}, {14}, {14, 27}, {34}}

func (c *char) attackFrameFunc(next action.Action) int {
	//back out what last attack was
	n := c.NormalCounter - 1
	if n < 0 {
		n = c.NormalHitNum - 1
	}
	//atkspd
	f := attackFrames[n][next]
	f = int(-0.5 * float64(f) * c.Stat(attributes.AtkSpd))
	return f
}

func (c *char) initNormalCancels() {

	//normal cancels
	attackFrames = make([][]int, c.NormalHitNum) //should be 5

	//n1 animations
	frames.InitNormalCancelSlice(&attackFrames, 0, hitmarks[0][0], 24)
	attackFrames[0][action.ActionAttack] = 18
	attackFrames[0][action.ActionCharge] = 24

	//n2 animations
	frames.InitNormalCancelSlice(&attackFrames, 1, hitmarks[1][0], 26)
	attackFrames[1][action.ActionAttack] = 13
	attackFrames[1][action.ActionCharge] = 26

	//n3 animations
	frames.InitNormalCancelSlice(&attackFrames, 2, hitmarks[2][0], 36)
	attackFrames[2][action.ActionAttack] = 26
	attackFrames[2][action.ActionCharge] = 36

	//n4 animations
	frames.InitNormalCancelSlice(&attackFrames, 3, hitmarks[3][1], 57)
	attackFrames[3][action.ActionAttack] = 41
	attackFrames[3][action.ActionCharge] = 57

	//n5 animations
	frames.InitNormalCancelSlice(&attackFrames, 4, hitmarks[4][0], 50)
	attackFrames[4][action.ActionAttack] = 41
	attackFrames[4][action.ActionCharge] = 100 //TODO: this action is illegal; need better way to handle it

}

func (c *char) Attack(p map[string]int) action.ActionInfo {
	if c.Core.Status.Duration("raidenburst") > 0 {
		return c.swordAttack()
	}

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       fmt.Sprintf("Normal %v", c.NormalCounter),
		AttackTag:  combat.AttackTagNormal,
		ICDTag:     combat.ICDTagNormalAttack,
		ICDGroup:   combat.ICDGroupDefault,
		Element:    attributes.Physical,
		Durability: 25,
	}

	for i, mult := range attack[c.NormalCounter] {
		ai.Mult = mult[c.TalentLvlAttack()]
		c.Core.QueueAttack(
			ai,
			combat.NewDefCircHit(0.5, false, combat.TargettableEnemy),
			hitmarks[c.NormalCounter][i],
			hitmarks[c.NormalCounter][i],
		)
	}

	defer c.AdvanceNormalIndex()

	return action.ActionInfo{
		Frames:          c.attackFrameFunc,
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   hitmarks[c.NormalCounter][len(hitmarks[c.NormalCounter])-1],
		Post:            hitmarks[c.NormalCounter][len(hitmarks[c.NormalCounter])-1],
		State:           action.NormalAttackState,
	}
}

var swordFrames [][]int
var swordHitmarks = [][]int{{12}, {13}, {11}, {22, 33}, {33}}

func (c *char) swordAttackFramesFunc(next action.Action) int {
	//back out what last attack was
	n := c.NormalCounter - 1
	if n < 0 {
		n = c.NormalHitNum - 1
	}
	//atkspd
	f := swordFrames[n][next]
	f = int(-0.5 * float64(f) * c.Stat(attributes.AtkSpd))
	return f
}

func (c *char) initBurstAttackCancels() {

	//normal cancels
	swordFrames = make([][]int, c.NormalHitNum) //should be 5

	//n1 animations
	frames.InitNormalCancelSlice(&swordFrames, 0, swordHitmarks[0][0], 24)
	swordFrames[0][action.ActionAttack] = 19
	swordFrames[0][action.ActionCharge] = 24

	//n2 animations
	frames.InitNormalCancelSlice(&swordFrames, 1, swordHitmarks[1][0], 26)
	swordFrames[1][action.ActionAttack] = 16
	swordFrames[1][action.ActionCharge] = 26

	//n3 animations
	frames.InitNormalCancelSlice(&swordFrames, 2, swordHitmarks[2][0], 34)
	swordFrames[2][action.ActionAttack] = 16
	swordFrames[2][action.ActionCharge] = 34

	//n4 animations
	frames.InitNormalCancelSlice(&swordFrames, 3, swordHitmarks[3][1], 67)
	swordFrames[3][action.ActionAttack] = 44
	swordFrames[3][action.ActionCharge] = 67

	//n5 animations
	frames.InitNormalCancelSlice(&swordFrames, 4, swordHitmarks[4][0], 83)
	swordFrames[4][action.ActionAttack] = 59
	swordFrames[4][action.ActionCharge] = 83
}
func (c *char) swordAttack() action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       fmt.Sprintf("Musou Isshin %v", c.NormalCounter),
		AttackTag:  combat.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagNormalAttack,
		ICDGroup:   combat.ICDGroupDefault,
		Element:    attributes.Electro,
		Durability: 25,
	}

	for i, mult := range attackB[c.NormalCounter] {
		// Sword hits are dynamic - group snapshots with damage proc
		ai.Mult = mult[c.TalentLvlBurst()]
		ai.Mult += resolveBonus[c.TalentLvlBurst()] * c.stacksConsumed
		if c.Base.Cons >= 2 {
			ai.IgnoreDefPercent = .6
		}
		c.Core.QueueAttack(
			ai,
			combat.NewDefCircHit(2, false, combat.TargettableEnemy),
			swordHitmarks[c.NormalCounter][i],
			swordHitmarks[c.NormalCounter][i],
			c.burstRestorefunc,
			c.c6(),
		)
	}

	defer c.AdvanceNormalIndex()

	return action.ActionInfo{
		Frames:          c.swordAttackFramesFunc,
		AnimationLength: swordFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   swordHitmarks[c.NormalCounter][len(swordHitmarks[c.NormalCounter])-1],
		Post:            swordHitmarks[c.NormalCounter][len(swordHitmarks[c.NormalCounter])-1],
		State:           action.NormalAttackState,
	}
}

package kazuha

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

const normalHitNum = 5

var attackFrames [][]int
var attackAnimation = []int{18, 25, 35, 40, 71}
var attackHitmarks = [][]int{
	{12},         //n1
	{11},         //n2
	{16, 25},     //n3
	{15},         //n4
	{15, 23, 31}, //n5
}

func init() {
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0][0], 12)
	attackFrames[0][action.ActionAttack] = 16
	attackFrames[0][action.ActionCharge] = 16

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1][0], 11)
	attackFrames[1][action.ActionAttack] = 20
	attackFrames[1][action.ActionCharge] = 25

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2][1], 25)
	attackFrames[2][action.ActionAttack] = 30
	attackFrames[2][action.ActionCharge] = 35

	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3][0], 15)
	attackFrames[3][action.ActionAttack] = 40
	attackFrames[3][action.ActionCharge] = 36

	attackFrames[4] = frames.InitNormalCancelSlice(attackHitmarks[4][2], 31)
	attackFrames[4][action.ActionAttack] = 71
	attackFrames[4][action.ActionCharge] = 71 //TODO: missing frame for n5 -> charge
}

func (c *char) Attack(p map[string]int) action.ActionInfo {

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       fmt.Sprintf("Normal %v", c.NormalCounter),
		AttackTag:  combat.AttackTagNormal,
		ICDTag:     combat.ICDTagNormalAttack,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeSlash,
		Element:    attributes.Physical,
		Durability: 25,
	}

	act := action.ActionInfo{
		Frames:              frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength:     attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:       attackHitmarks[c.NormalCounter][len(attackHitmarks[c.NormalCounter])-1],
		State:               action.NormalAttackState,
		FramePausedOnHitlag: c.FramePausedOnHitlag,
	}

	for i, mult := range attack[c.NormalCounter] {
		ai.Mult = mult[c.TalentLvlAttack()]
		c.Core.QueueAttack(
			ai,
			combat.NewDefCircHit(0.3, false, combat.TargettableEnemy),
			attackHitmarks[c.NormalCounter][i],
			attackHitmarks[c.NormalCounter][i],
		)
	}

	defer c.AdvanceNormalIndex()

	return act
}

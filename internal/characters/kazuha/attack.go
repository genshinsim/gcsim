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

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2][0], 25)
	attackFrames[1][action.ActionAttack] = 20
	attackFrames[1][action.ActionCharge] = 25
}

func (c *char) Attack(p map[string]int) action.ActionInfo {

	f, a := c.ActionFrames(action.ActionAttack, p)

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

	for i, mult := range attack[c.NormalCounter] {
		ai.Mult = mult[c.TalentLvlAttack()]
		c.Core.QueueAttack(ai, combat.NewDefCircHit(0.3, false, combat.TargettableEnemy), attackHitmarks[c.NormalCounter][i], attackHitmarks[c.NormalCounter][i])
	}

	defer c.AdvanceNormalIndex()

	return f, a
}

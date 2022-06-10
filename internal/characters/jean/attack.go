package jean

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player"
)

var attackFrames [][]int
var attackHitmarks = []int{13, 6, 17, 37, 25}

const normalHitNum = 5

func init() {
	// NA cancels
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0], 25)
	attackFrames[0][action.ActionAttack] = 22

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1], 20)
	attackFrames[1][action.ActionAttack] = 14

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2], 31)
	attackFrames[2][action.ActionAttack] = 28

	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3], 49)
	attackFrames[3][action.ActionAttack] = 44

	attackFrames[4] = frames.InitNormalCancelSlice(attackHitmarks[4], 68)
	attackFrames[4][action.ActionCharge] = 500 //TODO: this action is illegal; need better way to handle it
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
		Mult:       auto[c.NormalCounter][c.TalentLvlAttack()],
	}

	c.Core.Tasks.Add(func() {
		snap := c.Snapshot(&ai)
		c.Core.QueueAttackWithSnap(ai, snap, combat.NewDefCircHit(0.4, false, combat.TargettableEnemy), 0)

		//check for healing
		if c.Core.Rand.Float64() < 0.5 {
			heal := 0.15 * (snap.BaseAtk*(1+snap.Stats[attributes.ATKP]) + snap.Stats[attributes.ATK])
			c.Core.Player.Heal(player.HealInfo{
				Caller:  c.Index,
				Target:  -1,
				Message: "Wind Companion",
				Src:     heal,
				Bonus:   c.Stat(attributes.Heal),
			})
		}
	}, attackHitmarks[c.NormalCounter])

	defer c.AdvanceNormalIndex()

	return action.ActionInfo{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackHitmarks[c.NormalCounter],
		Post:            attackHitmarks[c.NormalCounter],
		State:           action.NormalAttackState,
	}
}

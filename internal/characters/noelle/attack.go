package noelle

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var (
	attackFrames          [][]int
	attackHitmarks        = []int{28, 25, 20, 42}
	attackHitlagHaltFrame = []float64{0.10, 0.10, 0.09, 0.15}
	attackRadius          = [][]float64{{2, 2, 2, 1.8}, {5.2, 5.2, 5.2, 3.51}}
)

const normalHitNum = 4

func init() {
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0], 38)
	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1], 46)
	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2], 31)
	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3], 107)
}

func (c *char) Attack(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               fmt.Sprintf("Normal %v", c.NormalCounter),
		AttackTag:          combat.AttackTagNormal,
		ICDTag:             combat.ICDTagNormalAttack,
		ICDGroup:           combat.ICDGroupDefault,
		StrikeType:         combat.StrikeTypeBlunt,
		Element:            attributes.Physical,
		Durability:         25,
		Mult:               attack[c.NormalCounter][c.TalentLvlAttack()],
		HitlagFactor:       0.01,
		HitlagHaltFrames:   attackHitlagHaltFrame[c.NormalCounter] * 60,
		CanBeDefenseHalted: true,
	}
	burstIndex := 0
	if c.StatModIsActive(burstBuffKey) {
		burstIndex = 1
		if c.NormalCounter == 2 {
			//q-n3 has different hit lag
			ai.HitlagHaltFrames = 0.1 * 60
		}
	}
	// TODO: don't forget this when implementing her CA
	done := false
	cb := c.skillHealCB(done)
	radius := attackRadius[burstIndex][c.NormalCounter]
	// need char queue because of potential hitlag from C4
	c.QueueCharTask(func() {
		c.Core.QueueAttack(
			ai,
			combat.NewCircleHit(c.Core.Combat.Player(), radius),
			0,
			0,
			cb,
		)
	}, attackHitmarks[c.NormalCounter])

	defer c.AdvanceNormalIndex()

	c.a4Counter++
	if c.a4Counter == 4 {
		c.a4Counter = 0
		if c.Cooldown(action.ActionSkill) > 0 {
			c.ReduceActionCooldown(action.ActionSkill, 60)
		}
	}

	return action.ActionInfo{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackHitmarks[c.NormalCounter],
		State:           action.NormalAttackState,
	}
}

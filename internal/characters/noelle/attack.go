package noelle

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
)

var (
	attackFrames          [][]int
	attackHitmarks        = []int{28, 25, 20, 42}
	attackPoiseDMG        = [][]float64{{105.8, 98.1, 115.34, 151.68}, {132.25, 122.82, 143.865, 189.75}}
	attackHitlagHaltFrame = []float64{0.10, 0.10, 0.09, 0.15}
	attackHitboxes        = [][][]float64{{{2}, {2}, {2}, {2, 3}}, {{5.2}, {5.2}, {5.2}, {3.3, 6.2}}}
	attackOffsets         = []float64{1, 1, 1, -1}
)

const normalHitNum = 4

func init() {
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0], 38)
	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1], 46)
	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2], 31)
	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3], 107)
}

func (c *char) Attack(p map[string]int) (action.Info, error) {
	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               fmt.Sprintf("Normal %v", c.NormalCounter),
		AttackTag:          attacks.AttackTagNormal,
		ICDTag:             attacks.ICDTagNormalAttack,
		ICDGroup:           attacks.ICDGroupDefault,
		StrikeType:         attacks.StrikeTypeBlunt,
		Element:            attributes.Physical,
		Durability:         25,
		Mult:               attack[c.NormalCounter][c.TalentLvlAttack()],
		HitlagFactor:       0.01,
		HitlagHaltFrames:   attackHitlagHaltFrame[c.NormalCounter] * 60,
		CanBeDefenseHalted: true,
	}

	counter := c.NormalCounter
	// need char queue because of potential hitlag from C4
	c.QueueCharTask(func() {
		burstIndex := 0
		if c.StatModIsActive(burstBuffKey) {
			burstIndex = 1
			if counter == 2 {
				// q-n3 has different hit lag
				ai.HitlagHaltFrames = 0.1 * 60
			}
			ai.ICDTag = attacks.ICDTagNone
		}
		ai.PoiseDMG = attackPoiseDMG[burstIndex][counter]

		ap := combat.NewCircleHitOnTarget(
			c.Core.Combat.Player(),
			geometry.Point{Y: attackOffsets[counter]},
			attackHitboxes[burstIndex][counter][0],
		)
		if counter == 3 {
			ap = combat.NewBoxHitOnTarget(
				c.Core.Combat.Player(),
				geometry.Point{Y: attackOffsets[counter]},
				attackHitboxes[burstIndex][counter][0],
				attackHitboxes[burstIndex][counter][1],
			)
		}

		c.Core.QueueAttack(ai, ap, 0, 0, c.skillHealCB(), c.makeA4CB())
	}, attackHitmarks[c.NormalCounter])

	defer c.AdvanceNormalIndex()

	return action.Info{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackHitmarks[c.NormalCounter],
		State:           action.NormalAttackState,
	}, nil
}

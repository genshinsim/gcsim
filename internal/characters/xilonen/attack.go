package xilonen

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
	attackFrames           [][]int
	attackHitmarks         = [][]int{{18}, {16, 32}, {22}}
	attackHitlagHaltFrames = []float64{0.03, 0.03, 0.06}
	attackHitboxes         = [][]float64{{2.1}, {2.5, 2.5}, {2.8}}
	attackFanAngles        = [][]float64{{0}, {150, 150}, {160}}
	attackOffsets          = [][]float64{{1.0}, {-0.5, -0.5}, {-0.5}}

	rollerFrames           [][]int
	rollerHitmarks         = []int{17, 17, 22, 32}
	rollerHitlagHaltFrames = []float64{0.03, 0.03, 0.03, 0.06}
	rollerHitboxes         = [][]float64{{2.7, 270}, {5, 3.7}, {3.5, 4.5}, {3.5, 270}}
	rollerOffsets          = []float64{0.7, -0.5, -0.5, 0.6}
	rollerPoiseDMG         = []float64{52.8, 51.9, 62.1, 81.1}
)

const normalHitNum = 3
const rollerHitNum = 4

func init() {
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0][0], 34)
	attackFrames[0][action.ActionAttack] = 20
	attackFrames[0][action.ActionCharge] = 26

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1][1], 57)
	attackFrames[1][action.ActionAttack] = 46
	attackFrames[1][action.ActionCharge] = 39

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2][0], 70)
	attackFrames[2][action.ActionAttack] = 57
	attackFrames[2][action.ActionCharge] = 500 //TODO: this action is illegal; need better way to handle it

	rollerFrames = make([][]int, rollerHitNum)

	rollerFrames[0] = frames.InitNormalCancelSlice(rollerHitmarks[0], 44)
	rollerFrames[0][action.ActionAttack] = 20

	rollerFrames[1] = frames.InitNormalCancelSlice(rollerHitmarks[1], 48)
	rollerFrames[1][action.ActionAttack] = 28

	rollerFrames[2] = frames.InitNormalCancelSlice(rollerHitmarks[2], 50)
	rollerFrames[2][action.ActionAttack] = 30

	rollerFrames[3] = frames.InitNormalCancelSlice(rollerHitmarks[3], 69)
	rollerFrames[3][action.ActionAttack] = 68 // TODO: this action is illegal; need better way to handle it
}

func (c *char) Attack(p map[string]int) (action.Info, error) {
	if c.nightsoulState.HasBlessing() {
		return c.nightsoulAttack(), nil
	}

	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		AttackTag:          attacks.AttackTagNormal,
		ICDTag:             attacks.ICDTagNormalAttack,
		ICDGroup:           attacks.ICDGroupDefault,
		StrikeType:         attacks.StrikeTypeSlash,
		Element:            attributes.Physical,
		Durability:         25,
		HitlagHaltFrames:   attackHitlagHaltFrames[c.NormalCounter] * 60,
		HitlagFactor:       0.01,
		CanBeDefenseHalted: true,
	}

	for i, mult := range attack[c.NormalCounter] {
		ax := ai
		ax.Abil = fmt.Sprintf("Normal %v", c.NormalCounter)
		ax.Mult = mult[c.TalentLvlAttack()]

		var ap combat.AttackPattern
		if attackFanAngles[c.NormalCounter][i] > 0 {
			ap = combat.NewCircleHitOnTargetFanAngle(
				c.Core.Combat.Player(),
				geometry.Point{Y: attackOffsets[c.NormalCounter][i]},
				attackHitboxes[c.NormalCounter][i],
				attackFanAngles[c.NormalCounter][i],
			)
		} else {
			ap = combat.NewCircleHitOnTarget(
				c.Core.Combat.Player(),
				geometry.Point{Y: attackOffsets[c.NormalCounter][i]},
				attackHitboxes[c.NormalCounter][i],
			)
		}

		c.QueueCharTask(func() {
			c.Core.QueueAttack(ax, ap, 0, 0)
		}, attackHitmarks[c.NormalCounter][i])
	}

	defer c.AdvanceNormalIndex()

	return action.Info{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackHitmarks[c.NormalCounter][len(attackHitmarks[c.NormalCounter])-1],
		State:           action.NormalAttackState,
	}, nil
}

func (c *char) nightsoulAttack() action.Info {
	c.c6()
	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		AttackTag:          attacks.AttackTagNormal,
		AdditionalTags:     []attacks.AdditionalTag{attacks.AdditionalTagNightsoul},
		ICDTag:             attacks.ICDTagXilonenSkate,
		ICDGroup:           attacks.ICDGroupDefault,
		StrikeType:         attacks.StrikeTypeBlunt,
		PoiseDMG:           rollerPoiseDMG[c.NormalCounter],
		Element:            attributes.Geo,
		Durability:         25,
		HitlagHaltFrames:   rollerHitlagHaltFrames[c.NormalCounter] * 60,
		HitlagFactor:       0.01,
		CanBeDefenseHalted: c.NormalCounter == 0, // only N1 can be defhalted
		UseDef:             true,
		IgnoreInfusion:     true,
	}

	ax := ai
	ax.Abil = fmt.Sprintf("Blade Roller %v", c.NormalCounter)
	ax.Mult = attackE[c.NormalCounter][c.TalentLvlAttack()] + c.c6DmgMult()
	var ap combat.AttackPattern
	if c.NormalCounter == 0 || c.NormalCounter == 3 {
		ap = combat.NewCircleHitOnTargetFanAngle(
			c.Core.Combat.Player(),
			geometry.Point{Y: rollerOffsets[c.NormalCounter]},
			rollerHitboxes[c.NormalCounter][0],
			rollerHitboxes[c.NormalCounter][1],
		)
	} else {
		ap = combat.NewBoxHitOnTarget(
			c.Core.Combat.Player(),
			geometry.Point{Y: rollerOffsets[c.NormalCounter]},
			rollerHitboxes[c.NormalCounter][0],
			rollerHitboxes[c.NormalCounter][1],
		)
	}

	c.QueueCharTask(func() {
		c.Core.QueueAttack(ax, ap, 0, 0, c.a1cb)
	}, rollerHitmarks[c.NormalCounter])

	defer c.AdvanceNormalIndex()

	return action.Info{
		Frames:          frames.NewAttackFunc(c.Character, rollerFrames),
		AnimationLength: rollerFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   rollerHitmarks[c.NormalCounter],
		State:           action.NormalAttackState,
	}
}

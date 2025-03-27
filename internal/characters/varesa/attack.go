package varesa

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
)

// TODO: update frames
var (
	// based on heizou frames
	attackFrames          [][]int
	attackHitmarks        = []int{21, 8, 32}
	attackHitlagHaltFrame = []float64{0.03, 0.03, 0.06}
	attackHitboxes        = [][]float64{{2, 3}, {2, 3}, {2.2}}
	attackOffsets         = []float64{-0.2, -0.2, 1.1}

	// based on wriothesley frames
	fieryAttackFrames          [][]int
	fieryAttackHitmarks        = []int{20, 7, 31}
	fieryAttackHitlagFactor    = []float64{0, 0, 0.01}
	fieryAttackHitlagHaltFrame = []float64{0, 0, 0.03}
	fieryAttackHitboxes        = [][]float64{{2, 3}, {2, 3}, {2.5, 3}}
)

const normalHitNum = 3

func init() {
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0], 35)
	attackFrames[0][action.ActionAttack] = 20
	attackFrames[0][action.ActionCharge] = 21

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1], 30)
	attackFrames[1][action.ActionAttack] = 17
	attackFrames[1][action.ActionCharge] = 21

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2], 56)
	attackFrames[2][action.ActionAttack] = 45
	attackFrames[2][action.ActionCharge] = 46

	fieryAttackFrames = make([][]int, normalHitNum)

	fieryAttackFrames[0] = frames.InitNormalCancelSlice(fieryAttackHitmarks[0], 27)
	fieryAttackFrames[0][action.ActionAttack] = 14
	fieryAttackFrames[0][action.ActionCharge] = 23

	fieryAttackFrames[1] = frames.InitNormalCancelSlice(fieryAttackHitmarks[0], 25)
	fieryAttackFrames[1][action.ActionAttack] = 13
	fieryAttackFrames[1][action.ActionCharge] = 20

	fieryAttackFrames[2] = frames.InitNormalCancelSlice(fieryAttackHitmarks[0], 41)
	fieryAttackFrames[2][action.ActionAttack] = 24
	fieryAttackFrames[2][action.ActionCharge] = 19
}

func (c *char) Attack(p map[string]int) (action.Info, error) {
	if c.StatusIsActive(skillStatus) {
		// TODO: or c.Core.Player.Exec(action.ActionCharge, c.Base.Key, nil)
		return c.ChargeAttack(p)
	}
	if c.nightsoulState.HasBlessing() {
		return c.fieryAttack(), nil
	}

	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               fmt.Sprintf("Normal %v", c.NormalCounter),
		AttackTag:          attacks.AttackTagNormal,
		ICDTag:             attacks.ICDTagNormalAttack,
		ICDGroup:           attacks.ICDGroupDefault,
		StrikeType:         attacks.StrikeTypeDefault,
		Element:            attributes.Electro,
		Durability:         25,
		Mult:               attack[c.NormalCounter][c.TalentLvlAttack()],
		HitlagFactor:       0.01,
		HitlagHaltFrames:   attackHitlagHaltFrame[c.NormalCounter] * 60,
		CanBeDefenseHalted: true,
	}
	ap := combat.NewCircleHitOnTarget(
		c.Core.Combat.Player(),
		geometry.Point{Y: attackOffsets[c.NormalCounter]},
		attackHitboxes[c.NormalCounter][0],
	)
	if c.NormalCounter == 0 || c.NormalCounter == 1 {
		ap = combat.NewBoxHitOnTarget(
			c.Core.Combat.Player(),
			geometry.Point{Y: attackOffsets[c.NormalCounter]},
			attackHitboxes[c.NormalCounter][0],
			attackHitboxes[c.NormalCounter][1],
		)
	}
	c.Core.QueueAttack(ai, ap, attackHitmarks[c.NormalCounter], attackHitmarks[c.NormalCounter])

	defer c.AdvanceNormalIndex()

	return action.Info{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackHitmarks[c.NormalCounter],
		State:           action.NormalAttackState,
	}, nil
}

func (c *char) fieryAttack() action.Info {
	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               fmt.Sprintf("Fiery Passion %v", c.NormalCounter),
		AttackTag:          attacks.AttackTagNormal,
		ICDTag:             attacks.ICDTagNormalAttack,
		ICDGroup:           attacks.ICDGroupDefault,
		StrikeType:         attacks.StrikeTypeDefault,
		Element:            attributes.Electro,
		Durability:         25,
		Mult:               fieryAttack[c.NormalCounter][c.TalentLvlAttack()],
		HitlagFactor:       fieryAttackHitlagFactor[c.NormalCounter],
		HitlagHaltFrames:   fieryAttackHitlagHaltFrame[c.NormalCounter] * 60,
		CanBeDefenseHalted: true,
	}

	ap := combat.NewBoxHitOnTarget(
		c.Core.Combat.Player(),
		nil,
		fieryAttackHitboxes[c.NormalCounter][0],
		fieryAttackHitboxes[c.NormalCounter][1],
	)
	c.Core.QueueAttack(ai, ap, fieryAttackHitmarks[c.NormalCounter], fieryAttackHitmarks[c.NormalCounter])

	defer c.AdvanceNormalIndex()

	return action.Info{
		Frames:          frames.NewAttackFunc(c.Character, fieryAttackFrames),
		AnimationLength: fieryAttackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   fieryAttackHitmarks[c.NormalCounter],
		State:           action.NormalAttackState,
	}
}

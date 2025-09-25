package varesa

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var (
	attackFrames          [][]int
	attackHitmarks        = []int{23, 7, 33}
	attackHitlagHaltFrame = []float64{0, 0, 0.03}
	attackHitlagFactor    = []float64{0, 0, 0.01}
	attackHitboxes        = [][]float64{{2.5, 2.5}, {2.8, 3.5}, {2.5}}
	attackOffsets         = []float64{-0.2, -0.2, 0}

	fieryAttackFrames   [][]int
	fieryAttackHitmarks = []int{17, 29, 37}
	fieryAttackHitboxes = [][]float64{{2.5}, {4, 4}, {4, 4.5}}
	fieryAttackOffsets  = []float64{1, -0.5, -0.5}
)

const normalHitNum = 3

func init() {
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0], 49) // N1 -> Walk
	attackFrames[0][action.ActionAttack] = 33
	attackFrames[0][action.ActionCharge] = 23

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1], 30) // N2 -> N3
	attackFrames[1][action.ActionCharge] = 17
	attackFrames[1][action.ActionWalk] = 28

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2], 59) // N3 -> Walk
	attackFrames[2][action.ActionAttack] = 45
	attackFrames[2][action.ActionCharge] = 32

	fieryAttackFrames = make([][]int, normalHitNum)

	fieryAttackFrames[0] = frames.InitNormalCancelSlice(fieryAttackHitmarks[0], 39) // N1 -> Walk
	fieryAttackFrames[0][action.ActionAttack] = 29
	fieryAttackFrames[0][action.ActionCharge] = 31

	fieryAttackFrames[1] = frames.InitNormalCancelSlice(fieryAttackHitmarks[1], 47) // N2 -> Walk
	fieryAttackFrames[1][action.ActionAttack] = 39
	fieryAttackFrames[1][action.ActionCharge] = 25

	fieryAttackFrames[2] = frames.InitNormalCancelSlice(fieryAttackHitmarks[2], 63) // N3 -> N4/Walk
	fieryAttackFrames[2][action.ActionCharge] = 37
}

func (c *char) Attack(p map[string]int) (action.Info, error) {
	// OnRemoved is sometimes called after the next action is executed. so we need to exit nightsoul here too
	c.clearNightsoulCB(action.NormalAttackState)

	if c.StatusIsActive(skillStatus) {
		// TODO: or c.Core.Player.Exec(action.ActionCharge, c.Base.Key, nil)
		return c.ChargeAttack(p)
	}
	if c.nightsoulState.HasBlessing() {
		return c.fieryAttack(), nil
	}

	windup := 0
	if c.NormalCounter == 0 {
		if c.Core.Player.CurrentState() == action.BurstState && c.usedShortBurst {
			windup = 4
		} else {
			windup = 6
		}
	}

	ai := info.AttackInfo{
		ActorIndex:         c.Index(),
		Abil:               fmt.Sprintf("Normal %v", c.NormalCounter),
		AttackTag:          attacks.AttackTagNormal,
		ICDTag:             attacks.ICDTagNormalAttack,
		ICDGroup:           attacks.ICDGroupDefault,
		StrikeType:         attacks.StrikeTypeDefault,
		Element:            attributes.Electro,
		Durability:         25,
		Mult:               attack[c.NormalCounter][c.TalentLvlAttack()],
		HitlagFactor:       attackHitlagFactor[c.NormalCounter],
		HitlagHaltFrames:   attackHitlagHaltFrame[c.NormalCounter] * 60,
		CanBeDefenseHalted: true,
	}
	var ap info.AttackPattern
	switch len(attackHitboxes[c.NormalCounter]) {
	case 1: // circle
		ap = combat.NewCircleHitOnTarget(
			c.Core.Combat.Player(),
			info.Point{Y: attackOffsets[c.NormalCounter]},
			attackHitboxes[c.NormalCounter][0],
		)
	case 2: // box
		ap = combat.NewBoxHitOnTarget(
			c.Core.Combat.Player(),
			info.Point{Y: attackOffsets[c.NormalCounter]},
			attackHitboxes[c.NormalCounter][0],
			attackHitboxes[c.NormalCounter][1],
		)
	}
	c.Core.QueueAttack(ai, ap, attackHitmarks[c.NormalCounter]-windup, attackHitmarks[c.NormalCounter]-windup)

	defer c.AdvanceNormalIndex()

	return action.Info{
		Frames:          func(next action.Action) int { return attackFrames[c.NormalCounter][next] - windup },
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction] - windup,
		CanQueueAfter:   attackHitmarks[c.NormalCounter] - windup,
		State:           action.NormalAttackState,
	}, nil
}

func (c *char) fieryAttack() action.Info {
	ai := info.AttackInfo{
		ActorIndex:         c.Index(),
		Abil:               fmt.Sprintf("Fiery Passion %v", c.NormalCounter),
		AttackTag:          attacks.AttackTagNormal,
		ICDTag:             attacks.ICDTagNormalAttack,
		ICDGroup:           attacks.ICDGroupDefault,
		StrikeType:         attacks.StrikeTypeDefault,
		Element:            attributes.Electro,
		Durability:         25,
		Mult:               fieryAttack[c.NormalCounter][c.TalentLvlAttack()],
		HitlagFactor:       0.01,
		HitlagHaltFrames:   0.03 * 60,
		CanBeDefenseHalted: true,
	}

	windup := 6
	switch c.Core.Player.CurrentState() {
	case action.NormalAttackState:
		windup = 0
	case action.BurstState:
		if c.usedShortBurst {
			windup = 0
		} else {
			windup = 4
		}
	}

	var ap info.AttackPattern
	switch len(fieryAttackHitboxes[c.NormalCounter]) {
	case 1: // circle
		ap = combat.NewCircleHitOnTarget(
			c.Core.Combat.Player(),
			info.Point{Y: fieryAttackOffsets[c.NormalCounter]},
			fieryAttackHitboxes[c.NormalCounter][0],
		)
	case 2: // box
		ap = combat.NewBoxHitOnTarget(
			c.Core.Combat.Player(),
			info.Point{Y: fieryAttackOffsets[c.NormalCounter]},
			fieryAttackHitboxes[c.NormalCounter][0],
			fieryAttackHitboxes[c.NormalCounter][1],
		)
	}
	c.Core.QueueAttack(ai, ap, fieryAttackHitmarks[c.NormalCounter]-windup, fieryAttackHitmarks[c.NormalCounter]-windup)

	defer c.AdvanceNormalIndex()

	return action.Info{
		Frames:          func(next action.Action) int { return fieryAttackFrames[c.NormalCounter][next] - windup },
		AnimationLength: fieryAttackFrames[c.NormalCounter][action.InvalidAction] - windup,
		CanQueueAfter:   fieryAttackHitmarks[c.NormalCounter] - windup,
		State:           action.NormalAttackState,
	}
}

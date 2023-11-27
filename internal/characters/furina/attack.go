package furina

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
	attackFrames    [][]int
	attackHitmarks  = []int{15, 12, 21, 27}
	attackOffsets   = []float64{-1.4, -0.85, -0.95, -3}
	attackOffsetsC6 = []float64{-1.5, -1.15, -1.1, -3}

	// these ones should be correct
	attackHitlagHaltFrame = []float64{0.01, 0.01, 0.02, 0.02}
	attackHitboxes        = [][]float64{{1.5, 2.8}, {1.7}, {1.9}, {5, 6}}
	attackHitboxesC6      = [][]float64{{1.5, 3}, {2.3}, {2.2}, {5, 6}}
	attackStrikeType      = []attacks.StrikeType{attacks.StrikeTypePierce, attacks.StrikeTypeSlash, attacks.StrikeTypeSlash, attacks.StrikeTypePierce}

	arkheIcdKeys     = []string{"spiritbreath-thorn-icd", "surging-blade-icd"}
	arkhePrettyPrint = []string{"Spiritbreath Thorn", "Surging Blade"}
)

const normalHitNum = 4

func init() {
	// NA cancels
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0], 60) // N1 -> Walk
	attackFrames[0][action.ActionAttack] = 34
	attackFrames[0][action.ActionCharge] = 31

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1], 60) // N2 -> Walk
	attackFrames[1][action.ActionAttack] = 23
	attackFrames[1][action.ActionCharge] = 28

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2], 90) // N3 -> Walk
	attackFrames[2][action.ActionAttack] = 36
	attackFrames[2][action.ActionCharge] = 45

	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3], 120) // N4 -> Walk
	attackFrames[3][action.ActionAttack] = 53
	attackFrames[3][action.ActionCharge] = 58
}

func (c *char) arkheCB(a combat.AttackCB) {
	if c.StatusIsActive(arkheIcdKeys[c.arkhe]) {
		return
	}

	c.AddStatus(arkheIcdKeys[c.arkhe], 6*60, true)

	c.QueueCharTask(func() {
		ai := combat.AttackInfo{
			ActorIndex:     c.Index,
			Abil:           arkhePrettyPrint[c.arkhe] + " (" + c.Base.Key.Pretty() + ")",
			AttackTag:      attacks.AttackTagNormal,
			ICDTag:         attacks.ICDTagNone,
			ICDGroup:       attacks.ICDGroupDefault,
			StrikeType:     attacks.StrikeTypeSlash,
			Element:        attributes.Hydro,
			Durability:     0,
			Mult:           arkhe[c.TalentLvlAttack()],
			IgnoreInfusion: true,
		}
		// https://www.youtube.com/watch?v=sbKIEzelynE
		// Furina's 18% Max HP boost applies to her Arkhe attacks
		if c.Base.Cons >= 6 && c.StatusIsActive(c6Key) {
			ai.FlatDmg = c.c6BonusDMGArkhe()
		}
		ap := combat.NewBoxHitOnTarget(
			a.Target,
			nil,
			1.2,
			4.5,
		)
		c.Core.QueueAttack(ai, ap, 0, 0)
	}, 42)
}
func (c *char) Attack(p map[string]int) (action.Info, error) {
	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               fmt.Sprintf("Normal %v", c.NormalCounter),
		AttackTag:          attacks.AttackTagNormal,
		ICDTag:             attacks.ICDTagNormalAttack,
		ICDGroup:           attacks.ICDGroupDefault,
		StrikeType:         attackStrikeType[c.NormalCounter],
		Element:            attributes.Physical,
		Durability:         25,
		Mult:               attack[c.NormalCounter][c.TalentLvlAttack()],
		HitlagFactor:       0.01,
		HitlagHaltFrames:   attackHitlagHaltFrame[c.NormalCounter] * 60,
		CanBeDefenseHalted: true,
	}

	var ap combat.AttackPattern

	if c.Base.Cons >= 6 && c.StatusIsActive(c6Key) {
		ai.Element = attributes.Hydro
		ai.IgnoreInfusion = true
		ai.FlatDmg = c.c6BonusDMG()
		switch c.NormalCounter {
		case 0, 3:
			ap = combat.NewBoxHitOnTarget(
				c.Core.Combat.Player(),
				geometry.Point{Y: attackOffsetsC6[c.NormalCounter]},
				attackHitboxesC6[c.NormalCounter][0],
				attackHitboxesC6[c.NormalCounter][1],
			)
		case 1, 2:
			ap = combat.NewCircleHitOnTarget(
				c.Core.Combat.Player(),
				geometry.Point{Y: attackOffsetsC6[c.NormalCounter]},
				attackHitboxesC6[c.NormalCounter][0],
			)
		}
		c.Core.QueueAttack(ai, ap, attackHitmarks[c.NormalCounter], attackHitmarks[c.NormalCounter], c.arkheCB, c.c6cb)
	} else {
		switch c.NormalCounter {
		case 0, 3:
			ap = combat.NewBoxHitOnTarget(
				c.Core.Combat.Player(),
				geometry.Point{Y: attackOffsets[c.NormalCounter]},
				attackHitboxes[c.NormalCounter][0],
				attackHitboxes[c.NormalCounter][1],
			)
		case 1, 2:
			ap = combat.NewCircleHitOnTarget(
				c.Core.Combat.Player(),
				geometry.Point{Y: attackOffsets[c.NormalCounter]},
				attackHitboxes[c.NormalCounter][0],
			)
		}
		c.Core.QueueAttack(ai, ap, attackHitmarks[c.NormalCounter], attackHitmarks[c.NormalCounter], c.arkheCB)
	}

	defer c.AdvanceNormalIndex()

	return action.Info{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackHitmarks[c.NormalCounter],
		State:           action.NormalAttackState,
	}, nil
}

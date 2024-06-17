package arlecchino

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/targets"
)

var (
	attackFrames          [][]int
	attackHitmarks        = [][]int{{11}, {16}, {17}, {24, 35}, {21}, {44}}
	attackHitlagHaltFrame = [][]float64{{0.02}, {0.02}, {0.02}, {0, 0}, {0}, {0.02}}
	attackDefHalt         = [][]bool{{true}, {true}, {true}, {true, false}, {false}, {true}}
	attackStrikeTypes     = [][]attacks.StrikeType{
		{attacks.StrikeTypeSlash},
		{attacks.StrikeTypeSlash},
		{attacks.StrikeTypeSlash},
		{attacks.StrikeTypeSlash, attacks.StrikeTypeSlash},
		{attacks.StrikeTypeSpear},
		{attacks.StrikeTypeSlash},
	}
	attackHitboxes = [][][][]float64{
		{
			{{1.9, 3}},     // box
			{{2.6}},        // fan
			{{1.9, 4}},     // box
			{{2.8}, {2.8}}, // circle, circle
			{{2.5}},        // circle
			{{3}},          // circle
		},
		{
			{{1.9, 4.2}},   // box
			{{3.1}},        // fan
			{{3.4, 5.6}},   // box
			{{3.3}, {3.3}}, // circle, circle
			{{2.8}},        // circle
			{{3.7}},        // circle
		},
	}
	attackOffsets = [][][]float64{
		{{0, -0.15}},
		{{0, 0.5}},
		{{0, -1.2}},
		{{-0.5, 0.7}, {-0.5, 0.7}},
		{{0, 2.4}},
		{{0, 2.5}},
	}

	attackFanAngles = [][]float64{{360}, {300}, {360}, {360, 360}, {360}, {360}}
)

const naBuffKey = "masque-of-the-red-death"
const bondConsumeICDKey = "bond-consume-icd"
const normalHitNum = 6

func init() {
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0][0], 24)
	attackFrames[0][action.ActionAttack] = 11
	attackFrames[0][action.ActionCharge] = 17

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1][0], 31)
	attackFrames[1][action.ActionAttack] = 21
	attackFrames[1][action.ActionCharge] = 26

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2][0], 39)
	attackFrames[2][action.ActionAttack] = 31
	attackFrames[2][action.ActionCharge] = 33

	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3][1], 55)
	attackFrames[3][action.ActionAttack] = 49
	attackFrames[3][action.ActionCharge] = 49

	attackFrames[4] = frames.InitNormalCancelSlice(attackHitmarks[4][0], 43)
	attackFrames[4][action.ActionAttack] = 30
	attackFrames[4][action.ActionCharge] = 33

	attackFrames[5] = frames.InitNormalCancelSlice(attackHitmarks[5][0], 59)
	attackFrames[5][action.ActionAttack] = 58
	attackFrames[5][action.ActionCharge] = 55
}

func (c *char) naBuff() {
	c.Core.Events.Subscribe(event.OnHPDebt, func(args ...interface{}) bool {
		target := args[0].(int)
		if target != c.Index {
			return false
		}
		// TODO: Remove when BoL changes get logged for all characters
		c.Core.Log.NewEvent("Bond of Life changed", glog.LogCharacterEvent, c.Index).
			Write("arle_hp_debt", c.CurrentHPDebt()).
			Write("arle_hp_debt%", c.CurrentHPDebt()/c.MaxHP())
		if c.CurrentHPDebt() >= c.MaxHP()*0.3 {
			// can't use negative duration or else `if .arlecchino.status.masque-of-the-red-death` won't work
			c.AddStatus(naBuffKey, 999999, false)
		} else {
			c.DeleteStatus(naBuffKey)
		}
		return false
	}, "arlechinno-bol-hook")
}

func (c *char) Attack(p map[string]int) (action.Info, error) {
	counter := c.NormalCounter
	for i, mult := range attack[counter] {
		// clone the values into another variable so that it won't be changed when the queued task executes
		i := i
		mult := mult
		c.QueueCharTask(func() {
			ai := combat.AttackInfo{
				ActorIndex:         c.Index,
				Abil:               fmt.Sprintf("Normal %v", counter),
				AttackTag:          attacks.AttackTagNormal,
				ICDTag:             attacks.ICDTagNormalAttack,
				ICDGroup:           attacks.ICDGroupDefault,
				StrikeType:         attackStrikeTypes[counter][i],
				Element:            attributes.Physical,
				Durability:         25,
				Mult:               mult[c.TalentLvlAttack()],
				HitlagFactor:       0.01,
				HitlagHaltFrames:   attackHitlagHaltFrame[counter][i] * 60,
				CanBeDefenseHalted: attackDefHalt[counter][i],
			}
			if c.NormalCounter == 3 && i == 0 {
				ai.HitlagFactor = 0
			}
			naIndex := 0
			if c.StatusIsActive(naBuffKey) {
				naIndex = 1
				ai.Element = attributes.Pyro
				ai.IgnoreInfusion = true
				ai.FlatDmg += c.bondBonus()
			}

			var ap combat.AttackPattern
			if len(attackHitboxes[naIndex][counter][i]) == 1 { // circle or fan
				ap = combat.NewCircleHitOnTargetFanAngle(
					c.Core.Combat.Player(),
					geometry.Point{X: attackOffsets[counter][i][0], Y: attackOffsets[counter][i][1]},
					attackHitboxes[naIndex][counter][i][0],
					attackFanAngles[counter][i],
				)
			} else { // box
				ap = combat.NewBoxHitOnTarget(
					c.Core.Combat.Player(),
					geometry.Point{X: attackOffsets[counter][i][0], Y: attackOffsets[counter][i][1]},
					attackHitboxes[naIndex][counter][i][0],
					attackHitboxes[naIndex][counter][i][1],
				)
			}

			c.Core.QueueAttack(ai, ap, 0, 0, c.bondConsumeCB)
		}, attackHitmarks[counter][i])
	}

	defer c.AdvanceNormalIndex()

	return action.Info{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackHitmarks[c.NormalCounter][len(attackHitmarks[c.NormalCounter])-1],
		State:           action.NormalAttackState,
	}, nil
}

func (c *char) bondBonus() float64 {
	c1Bonus := 0.0
	if c.Base.Cons >= 1 {
		c1Bonus = 1.0
	}
	amt := (masque[c.TalentLvlAttack()] + c1Bonus) * c.CurrentHPDebt() / c.MaxHP() * c.TotalAtk()
	return amt
}

func (c *char) bondConsumeCB(a combat.AttackCB) {
	if a.Target.Type() != targets.TargettableEnemy {
		return
	}
	if !c.StatusIsActive(naBuffKey) {
		return
	}

	if c.StatusIsActive(bondConsumeICDKey) {
		return
	}

	// 0.03*60 = 1.8 rounded to 2 frames
	c.AddStatus(bondConsumeICDKey, 2, true)

	amt := -0.075 * c.CurrentHPDebt()

	c.ModifyHPDebtByAmount(amt)

	c.ReduceActionCooldown(action.ActionSkill, 0.8*60)
}

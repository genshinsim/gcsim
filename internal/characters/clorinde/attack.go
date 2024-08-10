package clorinde

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
	attackHitmarks        = [][]int{{18}, {12}, {23, 32}, {12, 18, 23}, {21}}
	attackHitlagHaltFrame = [][]float64{{0.03}, {0.03}, {0.03, 0.03}, {0.02, 0.02, 0.02}, {0.03}}
	attackHitlagFactor    = [][]float64{{0.01}, {0.01}, {0.01, 0.01}, {0.05, 0.05, 0.05}, {0.05}}
	attackDefHalt         = [][]bool{{true}, {true}, {true, true}, {true, true, true}, {true}}
	attackHitboxes        = [][][]float64{{{1.7}}, {{1.9}}, {{2.1}, {2.1}}, {{2, 3.5}, {2, 3}, {2, 3}}, {{2.5}}} // n4 is a box
	attackOffsets         = []float64{1.1, 1.3, 1.2, 1.3, 1.4}
)

var (
	skillAttackFrames   [][]int
	skillAttackHitmarks = []int{10, 10, 11}
)

const (
	normalHitNum = 5
	skillHitNum  = 3

	arkheHitmark = 42
	arkheICDKey  = "clorinde-arkhe-icd"
)

func init() {
	// Normal attack
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0][0], 24) // N1 -> CA
	attackFrames[0][action.ActionAttack] = 19

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1][0], 27)
	attackFrames[1][action.ActionAttack] = 12

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2][1], 42)
	attackFrames[2][action.ActionAttack] = 40

	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3][2], 35)
	attackFrames[3][action.ActionAttack] = 32

	attackFrames[4] = frames.InitNormalCancelSlice(attackHitmarks[4][0], 60)

	// Skill attack
	skillAttackFrames = make([][]int, skillHitNum)

	skillAttackFrames[0] = frames.InitNormalCancelSlice(skillAttackHitmarks[0], 19)
	skillAttackFrames[0][action.ActionSkill] = 11
	skillAttackFrames[0][action.ActionBurst] = 10

	skillAttackFrames[1] = frames.InitNormalCancelSlice(skillAttackHitmarks[1], 17)
	skillAttackFrames[1][action.ActionSkill] = 10
	skillAttackFrames[1][action.ActionBurst] = 10

	skillAttackFrames[2] = frames.InitNormalCancelSlice(skillAttackHitmarks[2], 20)
	skillAttackFrames[2][action.ActionSkill] = 11
	skillAttackFrames[2][action.ActionBurst] = 11
}

func (c *char) Attack(p map[string]int) (action.Info, error) {
	if c.StatusIsActive(skillStateKey) {
		return c.skillAttack(p)
	}

	for i, mult := range attack[c.NormalCounter] {
		ai := combat.AttackInfo{
			ActorIndex:         c.Index,
			Abil:               fmt.Sprintf("Normal %v", c.NormalCounter),
			AttackTag:          attacks.AttackTagNormal,
			ICDTag:             attacks.ICDTagNormalAttack,
			ICDGroup:           attacks.ICDGroupDefault,
			StrikeType:         attacks.StrikeTypeSlash,
			Element:            attributes.Physical,
			Durability:         25,
			Mult:               mult[c.TalentLvlAttack()],
			HitlagFactor:       attackHitlagFactor[c.NormalCounter][i],
			HitlagHaltFrames:   attackHitlagHaltFrame[c.NormalCounter][i] * 60,
			CanBeDefenseHalted: attackDefHalt[c.NormalCounter][i],
		}

		ap := combat.NewCircleHitOnTarget(
			c.Core.Combat.Player(),
			geometry.Point{Y: attackOffsets[c.NormalCounter]},
			attackHitboxes[c.NormalCounter][i][0],
		)
		if c.NormalCounter == 3 {
			ap = combat.NewBoxHitOnTarget(
				c.Core.Combat.Player(),
				geometry.Point{Y: attackOffsets[c.NormalCounter]},
				attackHitboxes[c.NormalCounter][i][0],
				attackHitboxes[c.NormalCounter][i][1],
			)
		}

		c.Core.QueueAttack(ai, ap, attackHitmarks[c.NormalCounter][i], attackHitmarks[c.NormalCounter][i])
	}

	defer c.AdvanceNormalIndex()

	return action.Info{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackHitmarks[c.NormalCounter][len(attackHitmarks[c.NormalCounter])-1],
		State:           action.NormalAttackState,
	}, nil
}

func (c *char) skillAttack(_ map[string]int) (action.Info, error) {
	ai := combat.AttackInfo{
		ActorIndex:     c.Index,
		Abil:           fmt.Sprintf("Swift Hunt (Piercing Shot) %d", c.normalSCounter),
		AttackTag:      attacks.AttackTagNormal,
		ICDTag:         attacks.ICDTagNormalAttack,
		ICDGroup:       attacks.ICDGroupDefault,
		StrikeType:     attacks.StrikeTypeDefault,
		Element:        attributes.Electro,
		Durability:     25,
		Mult:           skillEnhancedNA[c.TalentLvlSkill()],
		IgnoreInfusion: true,
	}

	t := c.Core.Combat.PrimaryTarget()
	gainBOL := true
	var ap combat.AttackPattern
	if c.currentHPDebtRatio() < 1 {
		// TODO: assume this is just a big rectangle center on target
		ap = combat.NewBoxHitOnTarget(t, nil, 2, 14)
	} else {
		ai.Abil = fmt.Sprintf("Swift Hunt (Normal shot) %d", c.normalSCounter)
		ai.Mult = skillNA[c.TalentLvlSkill()]
		ap = combat.NewCircleHitOnTarget(t, nil, 0.6)
		gainBOL = false
	}

	// TODO: assume no snapshotting on this
	c.QueueCharTask(func() {
		c.Core.QueueAttack(ai, ap, 0, 0, c.particleCB)
		c.arkheAttack()
		if gainBOL {
			c.gainBOLOnAttack() // Bond of Life timing is ping dependent
		}
	}, skillAttackHitmarks[c.normalSCounter])

	defer c.AdvanceNormalIndex()

	return action.Info{
		Frames:          frames.NewAbilFunc(skillAttackFrames[c.normalSCounter]),
		AnimationLength: skillAttackFrames[c.normalSCounter][action.InvalidAction],
		CanQueueAfter:   skillAttackFrames[c.normalSCounter][action.ActionBurst],
		State:           action.NormalAttackState,
	}, nil
}

func (c *char) arkheAttack() {
	if c.StatusIsActive(arkheICDKey) {
		return
	}
	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               "Surging Blade",
		AttackTag:          attacks.AttackTagElementalArt,
		ICDTag:             attacks.ICDTagNone,
		ICDGroup:           attacks.ICDGroupDefault,
		StrikeType:         attacks.StrikeTypeSpear,
		Element:            attributes.Electro,
		Durability:         0,
		Mult:               arkheDamage[c.TalentLvlSkill()],
		HitlagFactor:       0.01,
		CanBeDefenseHalted: true,
	}
	c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 4.5), arkheHitmark, arkheHitmark)
	c.AddStatus(arkheICDKey, int(arkheCD[c.TalentLvlSkill()]*60), true)
}

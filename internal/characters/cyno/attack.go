package cyno

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var (
	attackFrames          [][]int
	attackHitmarks        = [][]int{{14}, {17}, {13, 22}, {27}}
	attackHitlagHaltFrame = [][]float64{{0.01}, {0.06}, {0, 0.02}, {0.04}}
	attackDefHalt         = [][]bool{{false}, {true}, {false, true}, {true}}
	attackHitboxes        = [][]float64{{1.8}, {1.8, 2.7}, {2.2, 3.6}, {2.3}}
	attackOffsets         = []float64{0.5, -0.2, 0, 1}
	attackFanAngles       = []float64{270, 360, 360, 360}
)

const normalHitNum = 4

func init() {
	attackFrames = make([][]int, normalHitNum)
	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0][0], 28) // N1 -> CA
	attackFrames[0][action.ActionAttack] = 15

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1][0], 23) // N2 -> CA
	attackFrames[1][action.ActionAttack] = 22

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2][1], 27) // N3 -> N4
	attackFrames[2][action.ActionCharge] = 26

	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3][0], 58) // N4 -> N1
	attackFrames[3][action.ActionCharge] = 500                               // impossible action
}

func (c *char) Attack(p map[string]int) action.ActionInfo {
	if c.StatusIsActive(BurstKey) {
		return c.attackB(p) // go to burst mode attacks
	}
	for i, mult := range attack[c.NormalCounter] {
		ai := combat.AttackInfo{
			ActorIndex:         c.Index,
			Abil:               fmt.Sprintf("Normal %v", c.NormalCounter),
			Mult:               mult[c.TalentLvlAttack()],
			AttackTag:          combat.AttackTagNormal,
			ICDTag:             combat.ICDTagNormalAttack,
			ICDGroup:           combat.ICDGroupDefault,
			StrikeType:         combat.StrikeTypeSlash,
			Element:            attributes.Physical,
			Durability:         25,
			HitlagFactor:       0.01,
			HitlagHaltFrames:   attackHitlagHaltFrame[c.NormalCounter][i],
			CanBeDefenseHalted: attackDefHalt[c.NormalCounter][i],
		}
		if c.NormalCounter == 2 {
			ai.StrikeType = combat.StrikeTypeSpear
		}
		ap := combat.NewCircleHitOnTargetFanAngle(
			c.Core.Combat.Player(),
			combat.Point{Y: attackOffsets[c.NormalCounter]},
			attackHitboxes[c.NormalCounter][0],
			attackFanAngles[c.NormalCounter],
		)
		if c.NormalCounter == 1 || c.NormalCounter == 2 {
			ap = combat.NewBoxHitOnTarget(
				c.Core.Combat.Player(),
				combat.Point{Y: attackOffsets[c.NormalCounter]},
				attackHitboxes[c.NormalCounter][0],
				attackHitboxes[c.NormalCounter][1],
			)
		}
		c.Core.QueueAttack(ai, ap, attackHitmarks[c.NormalCounter][i], attackHitmarks[c.NormalCounter][i])
	}

	defer c.AdvanceNormalIndex()

	return action.ActionInfo{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackHitmarks[c.NormalCounter][len(attackHitmarks[c.NormalCounter])-1],
		State:           action.NormalAttackState,
	}
}

const burstHitNum = 5

var (
	attackBFrames          [][]int
	attackBHitmarks        = [][]int{{12}, {14}, {18}, {5, 14}, {40}}
	attackBHitlagHaltFrame = [][]float64{{0.01}, {0.01}, {0.03}, {0.01, 0.03}, {0.05}}
	attackBDefHalt         = [][]bool{{false}, {false}, {false}, {false, false}, {true}}
	attackBHitboxes        = [][]float64{{2}, {2}, {2.5, 6}, {2.5}, {3.5}}
	attackBOffsets         = []float64{1, 1, -0.2, 0.8, 1.5}
)

func init() {
	// NA cancels (burst)
	attackBFrames = make([][]int, burstHitNum)
	attackBFrames[0] = frames.InitNormalCancelSlice(attackBHitmarks[0][0], 28) // N1 -> CA
	attackBFrames[0][action.ActionAttack] = 16

	attackBFrames[1] = frames.InitNormalCancelSlice(attackBHitmarks[1][0], 35) // N2 -> CA
	attackBFrames[1][action.ActionAttack] = 31

	attackBFrames[2] = frames.InitNormalCancelSlice(attackBHitmarks[2][0], 41) // N3 -> N4
	attackBFrames[2][action.ActionCharge] = 39

	attackBFrames[3] = frames.InitNormalCancelSlice(attackBHitmarks[3][0], 36) // N4 -> CA
	attackBFrames[3][action.ActionAttack] = 27

	attackBFrames[4] = frames.InitNormalCancelSlice(attackBHitmarks[4][0], 62) // N5 -> N1
	attackBFrames[4][action.ActionCharge] = 500                                // illegal action
}

func (c *char) attackB(p map[string]int) action.ActionInfo {
	c.tryBurstPPSlide(attackBHitmarks[c.normalBCounter][len(attackBHitmarks[c.normalBCounter])-1])

	for i, mult := range attackB[c.normalBCounter] {
		ai := combat.AttackInfo{
			ActorIndex:         c.Index,
			Abil:               fmt.Sprintf("Pactsworn Pathclearer %v", c.normalBCounter),
			AttackTag:          combat.AttackTagNormal,
			ICDTag:             combat.ICDTagNormalAttack,
			ICDGroup:           combat.ICDGroupDefault,
			StrikeType:         combat.StrikeTypeSlash,
			Element:            attributes.Electro,
			Durability:         25,
			HitlagFactor:       0.01,
			HitlagHaltFrames:   attackBHitlagHaltFrame[c.normalBCounter][i],
			CanBeDefenseHalted: attackBDefHalt[c.normalBCounter][i],
			Mult:               mult[c.TalentLvlBurst()],
			FlatDmg:            c.Stat(attributes.EM) * 1.5, // this is A4
			IgnoreInfusion:     true,
		}
		if c.normalBCounter == 2 || c.normalBCounter == 4 {
			ai.StrikeType = combat.StrikeTypeBlunt
		}
		ap := combat.NewCircleHitOnTarget(
			c.Core.Combat.Player(),
			combat.Point{Y: attackBOffsets[c.normalBCounter]},
			attackBHitboxes[c.normalBCounter][0],
		)
		if c.normalBCounter == 2 {
			ap = combat.NewBoxHitOnTarget(
				c.Core.Combat.Player(),
				combat.Point{Y: attackBOffsets[c.normalBCounter]},
				attackBHitboxes[c.normalBCounter][0],
				attackBHitboxes[c.normalBCounter][1],
			)
		}
		c.QueueCharTask(func() {
			c.Core.QueueAttack(ai, ap, 0, 0)
		}, attackBHitmarks[c.normalBCounter][i])
	}

	defer c.AdvanceNormalIndex()

	return action.ActionInfo{
		Frames: func(next action.Action) int {
			n := c.normalBCounter - 1
			if n < 0 {
				n = burstHitNum - 1
			}
			return frames.AtkSpdAdjust(attackBFrames[n][next], c.Stat(attributes.AtkSpd))
		},
		AnimationLength: attackBFrames[c.normalBCounter][action.InvalidAction],
		CanQueueAfter:   attackBHitmarks[c.normalBCounter][len(attackBHitmarks[c.normalBCounter])-1],
		State:           action.NormalAttackState,
	}
}

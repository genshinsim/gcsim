package alyosha

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
	attackHitmarks        = [][]int{{13}, {14}, {17, 17 + 15}, {17}}
	attackHitlagHaltFrame = [][]float64{{0.06}, {0.06}, {0, 0.00}, {0.00}}
	attackDefHalt         = [][]bool{{true}, {true}, {false, false}, {false}}
	attackHitboxes        = [][][]float64{{{2}}, {{2, 4}}, {{1.8}, {2, 4}}, {{2.4, 6}}}
	attackOffsets         = [][][]float64{{{0, 0.6}}, {{0.2, 0}}, {{0, 0.8}, {0, -0.5}}, {{0, 0.5}}}
)

const normalHitNum = 4

func init() {
	attackFrames = make([][]int, normalHitNum)
	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0][0], 34) // N1 -> W
	attackFrames[0][action.ActionAttack] = 18
	attackFrames[0][action.ActionCharge] = 24

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1][0], 38) // N2 -> W
	attackFrames[1][action.ActionAttack] = 23
	attackFrames[1][action.ActionCharge] = 26

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2][1], 65) // N3 -> W
	attackFrames[2][action.ActionAttack] = 53
	attackFrames[2][action.ActionCharge] = 54

	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3][0], 62) // N4 -> W
	attackFrames[3][action.ActionAttack] = 56
}

func (c *char) Attack(p map[string]int) (action.Info, error) {
	for i, mult := range attack[c.NormalCounter] {
		ai := info.AttackInfo{
			ActorIndex:         c.Index(),
			Abil:               fmt.Sprintf("Normal %v", c.NormalCounter),
			Mult:               mult[c.TalentLvlAttack()],
			AttackTag:          attacks.AttackTagNormal,
			ICDTag:             attacks.ICDTagNormalAttack,
			ICDGroup:           attacks.ICDGroupDefault,
			StrikeType:         attacks.StrikeTypeSlash,
			Element:            attributes.Physical,
			Durability:         25,
			HitlagFactor:       0.01,
			HitlagHaltFrames:   attackHitlagHaltFrame[c.NormalCounter][i] * 60,
			CanBeDefenseHalted: attackDefHalt[c.NormalCounter][i],
		}

		offset := info.Point{X: attackOffsets[c.NormalCounter][i][0], Y: attackOffsets[c.NormalCounter][i][1]}
		ap := combat.NewCircleHitOnTarget(
			c.Core.Combat.Player(),
			offset,
			attackHitboxes[c.NormalCounter][i][0],
		)

		if c.NormalCounter == 1 || (c.NormalCounter == 2 && i == 1) || c.NormalCounter == 3 {
			ai.StrikeType = attacks.StrikeTypeSpear
			ap = combat.NewBoxHitOnTarget(
				c.Core.Combat.Player(),
				offset,
				attackHitboxes[c.NormalCounter][i][0],
				attackHitboxes[c.NormalCounter][i][1],
			)
		}

		var cb info.AttackCBFunc
		if c.NormalCounter == normalHitNum-1 {
			cb = c.triggerSkillMarkCB(true)
		}

		c.QueueCharTask(func() {
			c.Core.QueueAttack(
				ai,
				ap,
				0,
				0,
				cb,
			)
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

package keqing

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
	attackHitmarks        = [][]int{{10}, {10}, {14}, {11, 21}, {22}}
	attackHitlagHaltFrame = [][]float64{{.03}, {.03}, {.06}, {0, .03}, {0}}
	attackDefHalt         = [][]bool{{true}, {true}, {true}, {false, true}, {false}}
	attackHitboxes        = [][]float64{{1.5, 2.2}, {1.5}, {1.8}, {1.5}, {0.8}}
	attackOffsets         = []float64{0, 0.5, 1, 0.6, 0}
)

const normalHitNum = 5

func init() {
	// NA cancels
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0][0], 19)
	attackFrames[0][action.ActionAttack] = 15

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1][0], 24)
	attackFrames[1][action.ActionAttack] = 16

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2][0], 36)
	attackFrames[2][action.ActionAttack] = 27

	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3][1], 58)
	attackFrames[3][action.ActionAttack] = 31

	attackFrames[4] = frames.InitNormalCancelSlice(attackHitmarks[4][0], 66)
	attackFrames[4][action.ActionCharge] = 500 //TODO: this action is illegal; need better way to handle it
}

func (c *char) Attack(p map[string]int) (action.Info, error) {
	centerTarget := c.Core.Combat.Player()
	if c.NormalCounter == 4 {
		centerTarget = c.Core.Combat.PrimaryTarget() // N5 is a bullet
	}
	c2CB := c.makeC2CB()
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
			HitlagFactor:       0.01,
			HitlagHaltFrames:   attackHitlagHaltFrame[c.NormalCounter][i] * 60,
			CanBeDefenseHalted: attackDefHalt[c.NormalCounter][i],
		}
		ap := combat.NewCircleHit(
			c.Core.Combat.Player(),
			centerTarget,
			geometry.Point{Y: attackOffsets[c.NormalCounter]},
			attackHitboxes[c.NormalCounter][0],
		)
		if c.NormalCounter == 0 {
			ap = combat.NewBoxHit(
				c.Core.Combat.Player(),
				centerTarget,
				geometry.Point{Y: attackOffsets[c.NormalCounter]},
				attackHitboxes[c.NormalCounter][0],
				attackHitboxes[c.NormalCounter][1],
			)
		}
		c.QueueCharTask(func() {
			c.Core.QueueAttack(ai, ap, 0, 0, c2CB)
		}, attackHitmarks[c.NormalCounter][i])
	}

	if c.Base.Cons >= 6 {
		c.c6("attack")
	}

	defer c.AdvanceNormalIndex()

	act := action.Info{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackHitmarks[c.NormalCounter][len(attackHitmarks[c.NormalCounter])-1],
		State:           action.NormalAttackState,
	}

	if c.NormalCounter == 0 {
		act.UseNormalizedTime = func(next action.Action) bool {
			return next == action.ActionCharge
		}
	}

	return act, nil
}

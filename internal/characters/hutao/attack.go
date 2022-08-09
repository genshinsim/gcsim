package hutao

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var attackFrames [][]int
var attackHitmarks = [][]int{{12}, {9}, {17}, {23}, {16, 26}, {27}}
var attackHitlagHaltFrame = [][]float64{{0.01}, {0.01}, {0.01}, {0.02}, {0.02, 0.02}, {0.04}}
var attackDefHalt = [][]bool{{true}, {true}, {true}, {true}, {false, true}, {true}}

const normalHitNum = 6

func init() {
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0][0], 20)
	attackFrames[0][action.ActionAttack] = 14

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1][0], 17)
	attackFrames[1][action.ActionAttack] = 12

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2][0], 26)
	attackFrames[2][action.ActionCharge] = 23

	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3][0], 31)
	attackFrames[3][action.ActionAttack] = 29

	attackFrames[4] = frames.InitNormalCancelSlice(attackHitmarks[4][1], 48)
	attackFrames[4][action.ActionAttack] = 36

	attackFrames[5] = frames.InitNormalCancelSlice(attackHitmarks[5][0], 72)
	attackFrames[5][action.ActionCharge] = 500 //TODO: this action is illegal; need better way to handle it
}

func (c *char) Attack(p map[string]int) action.ActionInfo {
	if c.StatModIsActive(paramitaBuff) {
		return c.ppAttack(p)
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
			HitlagHaltFrames:   attackHitlagHaltFrame[c.NormalCounter][i] * 60,
			CanBeDefenseHalted: attackDefHalt[c.NormalCounter][i],
		}
		c.QueueCharTask(func() {
			c.Core.QueueAttack(
				ai,
				combat.NewCircleHit(c.Core.Combat.Player(), 0.1, false, combat.TargettableEnemy),
				0,
				0,
			)
		}, attackHitmarks[c.NormalCounter][i])
	}

	defer c.AdvanceNormalIndex()

	return action.ActionInfo{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackHitmarks[c.NormalCounter][len(attackHitmarks[c.NormalCounter])-1],
		State:           action.NormalAttackState,
	}
}

var ppAttackHitlagHaltFrame = [][]float64{{0.01}, {0.01}, {0.01}, {0.02}, {0.02, 0.02}, {0.04}}
var ppAttackDefHalt = [][]bool{{true}, {true}, {true}, {true}, {false, true}, {true}}

func (c *char) ppAttack(p map[string]int) action.ActionInfo {

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
			HitlagHaltFrames:   ppAttackHitlagHaltFrame[c.NormalCounter][i] * 60,
			CanBeDefenseHalted: ppAttackDefHalt[c.NormalCounter][i],
		}
		c.QueueCharTask(func() {
			c.Core.QueueAttack(
				ai,
				combat.NewCircleHit(c.Core.Combat.Player(), 0.1, false, combat.TargettableEnemy),
				0,
				0,
				c.ppParticles,
			)
		}, attackHitmarks[c.NormalCounter][i])
	}

	defer c.AdvanceNormalIndex()

	act := action.ActionInfo{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackHitmarks[c.NormalCounter][len(attackHitmarks[c.NormalCounter])-1],
		State:           action.NormalAttackState,
	}

	if c.NormalCounter == 1 {
		act.UseNormalizedTime = func(next action.Action) bool {
			return next == action.ActionCharge
		}
	}

	return act
}

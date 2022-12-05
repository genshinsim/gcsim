package ayato

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var (
	shunsuikenFrames      []int
	attackFrames          [][]int
	attackHitmarks        = [][]int{{12}, {18}, {20}, {22, 25}, {41}}
	attackHitlagHaltFrame = [][]float64{{0.03}, {0.03}, {0.06}, {0, 0}, {0.08}}
	attackDefHalt         = [][]bool{{true}, {true}, {true}, {false, false}, {true}}
	attackRadius          = []float64{1.7, 1.7, 1.61, 1.64, 3.16}
)

const normalHitNum = 5
const shunsuikenHitmark = 5

func init() {
	// NA cancels
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0][0], 21) // N1 -> CA
	attackFrames[0][action.ActionAttack] = 15                                // N1 -> N2

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1][0], 27) // N2 -> N3/CA

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2][0], 34) // N3 -> CA
	attackFrames[2][action.ActionAttack] = 30                                // N3 -> N4

	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3][1], 29) // N4 -> CA
	attackFrames[3][action.ActionAttack] = 27                                // N4 -> N5

	attackFrames[4] = frames.InitNormalCancelSlice(attackHitmarks[4][0], 63) // N5 -> N1
	attackFrames[4][action.ActionCharge] = 500                               // N5 -> CA, TODO: this action is illegal; need better way to handle it

	// NA (in skill) -> x
	shunsuikenFrames = frames.InitNormalCancelSlice(shunsuikenHitmark, 23)
}

func (c *char) Attack(p map[string]int) action.ActionInfo {

	if c.StatusIsActive(skillBuffKey) {
		return c.SoukaiKanka(p)
	}

	for i, mult := range attack[c.NormalCounter] {
		ai := combat.AttackInfo{
			Abil:               fmt.Sprintf("Normal %v", c.NormalCounter),
			ActorIndex:         c.Index,
			AttackTag:          combat.AttackTagNormal,
			ICDTag:             combat.ICDTagNormalAttack,
			ICDGroup:           combat.ICDGroupDefault,
			StrikeType:         combat.StrikeTypeSlash,
			Element:            attributes.Physical,
			Durability:         25,
			Mult:               mult[c.TalentLvlAttack()],
			HitlagFactor:       0.01,
			HitlagHaltFrames:   attackHitlagHaltFrame[c.NormalCounter][i] * 60,
			CanBeDefenseHalted: attackDefHalt[c.NormalCounter][i],
		}
		radius := attackRadius[c.NormalCounter]
		c.QueueCharTask(func() {
			c.Core.QueueAttack(
				ai,
				combat.NewCircleHit(c.Core.Combat.Player(), radius),
				0,
				0,
			)
		}, attackHitmarks[c.NormalCounter][i])
	}

	defer c.AdvanceNormalIndex()

	// normal state
	return action.ActionInfo{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackHitmarks[c.NormalCounter][len(attackHitmarks[c.NormalCounter])-1],
		State:           action.NormalAttackState,
	}
}

func (c *char) SoukaiKanka(p map[string]int) action.ActionInfo {

	ai := combat.AttackInfo{
		Abil:               fmt.Sprintf("Normal %v", c.NormalCounter),
		ActorIndex:         c.Index,
		AttackTag:          combat.AttackTagNormal,
		ICDTag:             combat.ICDTagNormalAttack,
		ICDGroup:           combat.ICDGroupDefault,
		StrikeType:         combat.StrikeTypeSlash,
		Element:            attributes.Physical,
		Durability:         25,
		Mult:               shunsuiken[c.NormalCounter][c.TalentLvlSkill()],
		HitlagFactor:       0.01,
		HitlagHaltFrames:   0.03 * 60,
		CanBeDefenseHalted: false,
	}
	c.Core.QueueAttack(ai, combat.NewCircleHit(c.Core.Combat.Player(), 5.32), 0, shunsuikenHitmark, c.generateParticles, c.skillStacks)

	defer c.AdvanceNormalIndex()

	return action.ActionInfo{
		Frames: func(next action.Action) int {
			return frames.AtkSpdAdjust(shunsuikenFrames[next], c.Stat(attributes.AtkSpd))
		},
		AnimationLength: shunsuikenFrames[action.InvalidAction],
		CanQueueAfter:   shunsuikenHitmark,
		State:           action.NormalAttackState,
	}
}

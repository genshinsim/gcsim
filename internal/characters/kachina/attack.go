package kachina

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

const normalHitNum = 4

var (
	attackFrames             [][]int
	chargeFrames             []int
	attackHitmark            = [][]int{{11}, {14, 18}, {21}, {16}}
	attackHitlagHaltFrame    = [][]float64{{0.06}, {0.06, 0.06}, {0}, {0.01}}
	attackCanBeDefenseHalted = [][]bool{{true}, {true, true}, {false}, {true}}
)

func init() {
	attackFrames = make([][]int, normalHitNum)
	attackFrames[0] = frames.InitNormalCancelSlice(11, 31)
	attackFrames[0][action.ActionAttack] = 18
	attackFrames[0][action.ActionCharge] = 20
	attackFrames[1] = frames.InitNormalCancelSlice(18, 52)
	attackFrames[1][action.ActionAttack] = 41
	attackFrames[1][action.ActionCharge] = 41
	attackFrames[2] = frames.InitNormalCancelSlice(21, 52)
	attackFrames[2][action.ActionAttack] = 44
	attackFrames[2][action.ActionCharge] = 45
	attackFrames[3] = frames.InitNormalCancelSlice(16, 51)
	attackFrames[3][action.ActionAttack] = 51
	attackFrames[3][action.ActionCharge] = 500

	chargeFrames = frames.InitAbilSlice(55)
	chargeFrames[action.ActionAttack] = 50
	chargeFrames[action.ActionCharge] = 500
	chargeFrames[action.ActionSkill] = 50
	chargeFrames[action.ActionBurst] = 50
	chargeFrames[action.ActionDash] = 50
	chargeFrames[action.ActionJump] = 50
	chargeFrames[action.ActionSwap] = 49

}

func (c *char) Attack(_ map[string]int) (action.Info, error) {
	if c.mounted && c.nightsoulState.HasBlessing() {
		return c.mountedAttack(), nil
	}

	var multipliers [][]float64
	switch c.NormalCounter {
	case 0:
		multipliers = [][]float64{attack[0]}
	case 1:
		multipliers = [][]float64{attack[1], attack[2]}
	case 2:
		multipliers = [][]float64{attack[3]}
	default:
		multipliers = [][]float64{attack[4]}
	}
	for i, mult := range multipliers {
		ai := info.AttackInfo{
			ActorIndex:         c.Index(),
			Abil:               fmt.Sprintf("Normal %v", c.NormalCounter+1),
			AttackTag:          attacks.AttackTagNormal,
			ICDTag:             attacks.ICDTagNormalAttack,
			ICDGroup:           attacks.ICDGroupDefault,
			StrikeType:         attacks.StrikeTypeSpear,
			Element:            attributes.Physical,
			Durability:         25,
			Mult:               mult[c.TalentLvlAttack()],
			HitlagFactor:       0.01,
			HitlagHaltFrames:   attackHitlagHaltFrame[c.NormalCounter][i] * 60,
			CanBeDefenseHalted: attackCanBeDefenseHalted[c.NormalCounter][i],
		}
		c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 2.5), attackHitmark[c.NormalCounter][i], attackHitmark[c.NormalCounter][i])
	}

	defer c.AdvanceNormalIndex()
	return action.Info{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackHitmark[c.NormalCounter][len(attackHitmark[c.NormalCounter])-1],
		State:           action.NormalAttackState,
	}, nil
}

func (c *char) mountedAttack() action.Info {
	ai := c.twirlyAttackInfo(fmt.Sprintf("Turbo Twirly Normal %v", c.NormalCounter+1), mounted[0][c.TalentLvlSkill()])
	c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, c.twirlyRadius()), mountedHitmark(c.NormalCounter), mountedHitmark(c.NormalCounter), c.twirlyParticleCB, func(info.AttackCB) {
		c.consumeTwirlyPoints(10)
	})
	defer c.AdvanceNormalIndex()
	return action.Info{
		Frames:          frames.NewAttackFunc(c.Character, mountedFrames),
		AnimationLength: mountedFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   mountedHitmark(c.NormalCounter),
		State:           action.NormalAttackState,
	}
}

var mountedFrames = [][]int{
	frames.InitNormalCancelSlice(38, 51),
	frames.InitNormalCancelSlice(36, 47),
	frames.InitNormalCancelSlice(35, 46),
	frames.InitNormalCancelSlice(34, 46),
}

func mountedHitmark(index int) int {
	return []int{38, 36, 35, 34}[index]
}

func (c *char) ChargeAttack(_ map[string]int) (action.Info, error) {
	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Charge",
		AttackTag:  attacks.AttackTagExtra,
		ICDTag:     attacks.ICDTagExtraAttack,
		ICDGroup:   attacks.ICDGroupPoleExtraAttack,
		StrikeType: attacks.StrikeTypeSpear,
		Element:    attributes.Physical,
		Durability: 25,
		Mult:       charge_attack[c.TalentLvlAttack()],
		PoiseDMG:   120,
	}
	c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 3.5), 20, 20)
	return action.Info{
		Frames:          frames.NewAbilFunc(chargeFrames),
		AnimationLength: chargeFrames[action.InvalidAction],
		CanQueueAfter:   20,
		State:           action.ChargeAttackState,
	}, nil
}

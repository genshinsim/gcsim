package barbara

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/core/targets"
)

var (
	attackFrames   [][]int
	attackHitmarks = []int{6, 11, 12, 32}
	attackRadius   = []float64{1, 1, 1, 2}
)

const normalHitNum = 4

func init() {
	attackFrames = make([][]int, normalHitNum)

	// N1 -> x
	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0], 23)
	attackFrames[0][action.ActionAttack] = 15
	attackFrames[0][action.ActionCharge] = 18

	// N2 -> x
	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1], 28)
	attackFrames[1][action.ActionAttack] = 21
	attackFrames[1][action.ActionCharge] = 24

	// N3 -> x
	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2], 30)
	attackFrames[2][action.ActionAttack] = 22
	attackFrames[2][action.ActionCharge] = 28

	// N4 -> x
	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3], 60)
	attackFrames[3][action.ActionCharge] = 500 //TODO: this action is illegal; need better way to handle it
	attackFrames[3][action.ActionDash] = 2
	attackFrames[3][action.ActionJump] = 3
	attackFrames[3][action.ActionSwap] = 2
	attackFrames[3][action.ActionWalk] = 57
}

// Standard attack function with seal handling
func (c *char) Attack(p map[string]int) (action.Info, error) {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       fmt.Sprintf("Normal %v", c.NormalCounter),
		AttackTag:  attacks.AttackTagNormal,
		ICDTag:     attacks.ICDTagNormalAttack,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Hydro,
		Durability: 25,
		Mult:       attack[c.NormalCounter][c.TalentLvlAttack()],
	}
	done := false
	cb := func(a combat.AttackCB) {
		if a.Target.Type() != targets.TargettableEnemy {
			return
		}
		if done {
			return
		}
		// check for healing
		if c.Core.Status.Duration(barbSkillKey) > 0 {
			c.Core.Player.Heal(player.HealInfo{
				Caller:  c.Index,
				Target:  -1,
				Message: "Melody Loop (Normal Attack)",
				Src:     prochpp[c.TalentLvlSkill()]*c.MaxHP() + prochp[c.TalentLvlSkill()],
				Bonus:   c.Stat(attributes.Heal),
			})
			done = true
		}
	}
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHit(
			c.Core.Combat.Player(),
			c.Core.Combat.PrimaryTarget(),
			nil,
			attackRadius[c.NormalCounter],
		),
		attackHitmarks[c.NormalCounter],
		attackHitmarks[c.NormalCounter],
		cb,
	)

	defer c.AdvanceNormalIndex()

	return action.Info{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackHitmarks[c.NormalCounter],
		State:           action.NormalAttackState,
	}, nil
}

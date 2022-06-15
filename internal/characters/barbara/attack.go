package barbara

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player"
)

var attackHitmarks = []int{7, 25 - 7, 45 - 25, 92 - 45}
var attackFrames [][]int

const normalHitNum = 4

func init() {
	attackFrames = make([][]int, normalHitNum)
	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0], attackHitmarks[0])
	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1], attackHitmarks[1])
	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2], attackHitmarks[2])
	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3], attackHitmarks[3])
}

// Standard attack function with seal handling
func (c *char) Attack(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       fmt.Sprintf("Normal %v", c.NormalCounter),
		AttackTag:  combat.AttackTagNormal,
		ICDTag:     combat.ICDTagNormalAttack,
		ICDGroup:   combat.ICDGroupDefault,
		Element:    attributes.Hydro,
		Durability: 25,
		Mult:       attack[c.NormalCounter][c.TalentLvlAttack()],
	}
	done := false
	cb := func(a combat.AttackCB) {
		if done {
			return
		}
		//check for healing
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
	c.Core.QueueAttack(ai,
		combat.NewDefCircHit(0.1, false, combat.TargettableEnemy),
		0,
		attackHitmarks[c.NormalCounter],
		cb,
	)
	c.AdvanceNormalIndex()

	// return animation cd
	return action.ActionInfo{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackHitmarks[c.NormalCounter],
		State:           action.NormalAttackState,
	}
}

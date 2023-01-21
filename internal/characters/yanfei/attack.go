package yanfei

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

var attackFrames [][]int
var attackHitmarks = []int{12, 16, 37}

const (
	normalHitNum = 3
	sealBuffKey  = "yanfei-seal"
)

func init() {
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0], 26) // N1 -> N2
	attackFrames[0][action.ActionCharge] = 21                             // N1 -> CA

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1], 28) // N2 -> N3
	attackFrames[1][action.ActionCharge] = 16                             // N2 -> CA

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2], 73) // N3 -> N1
	attackFrames[2][action.ActionCharge] = 42                             // N3 -> CA

}

// Standard attack function with seal handling
func (c *char) Attack(p map[string]int) action.ActionInfo {
	travel, ok := p["travel"]
	if !ok {
		travel = 10
	}

	done := false
	addSeal := func(a combat.AttackCB) {
		if a.Target.Type() != combat.TargettableEnemy {
			return
		}
		// doesn't gain seals off-field
		if c.Core.Player.Active() != c.Index {
			return
		}
		if done {
			return
		}
		if c.sealCount < c.maxTags {
			c.sealCount++
		}
		c.AddStatus(sealBuffKey, 600, true)
		c.Core.Log.NewEvent("yanfei gained a seal from normal attack", glog.LogCharacterEvent, c.Index).
			Write("current_seals", c.sealCount)
		done = true
	}

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       fmt.Sprintf("Normal %v", c.NormalCounter),
		AttackTag:  combat.AttackTagNormal,
		ICDTag:     combat.ICDTagNormalAttack,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeDefault,
		Element:    attributes.Pyro,
		Durability: 25,
		Mult:       attack[c.NormalCounter][c.TalentLvlAttack()],
	}
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHit(
			c.Core.Combat.Player(),
			c.Core.Combat.PrimaryTarget(),
			nil,
			0.75,
		),
		attackHitmarks[c.NormalCounter],
		attackHitmarks[c.NormalCounter]+travel,
		addSeal,
	)

	defer c.AdvanceNormalIndex()

	return action.ActionInfo{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackHitmarks[c.NormalCounter],
		State:           action.NormalAttackState,
	}
}

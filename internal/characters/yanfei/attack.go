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
var attackHitmarks = []int{13, 28, 49}

const normalHitNum = 3

func init() {
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0], 13)
	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1], 28)
	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2], 49)
}

// Standard attack function with seal handling
func (c *char) Attack(p map[string]int) action.ActionInfo {
	travel, ok := p["travel"]
	if !ok {
		travel = 10
	}

	done := false
	addSeal := func(a combat.AttackCB) {
		if done {
			return
		}
		if c.Tags["seal"] < c.maxTags {
			c.Tags["seal"]++
		}
		c.sealExpiry = c.Core.F + 600
		c.Core.Log.NewEvent("yanfei gained a seal from normal attack", glog.LogCharacterEvent, c.Index, "current_seals", c.Tags["seal"], "expiry", c.sealExpiry)
		done = true
	}

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       fmt.Sprintf("Normal %v", c.NormalCounter),
		AttackTag:  combat.AttackTagNormal,
		ICDTag:     combat.ICDTagNormalAttack,
		ICDGroup:   combat.ICDGroupDefault,
		Element:    attributes.Pyro,
		Durability: 25,
		Mult:       attack[c.NormalCounter][c.TalentLvlAttack()],
	}
	c.Core.QueueAttack(
		ai,
		combat.NewDefCircHit(0.1, false, combat.TargettableEnemy),
		attackHitmarks[c.NormalCounter],
		attackHitmarks[c.NormalCounter]+travel,
		addSeal,
	)

	defer c.AdvanceNormalIndex()

	return action.ActionInfo{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackHitmarks[c.NormalCounter],
		Post:            attackHitmarks[c.NormalCounter],
		State:           action.NormalAttackState,
	}
}

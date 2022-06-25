package ningguang

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

var attackFrames [][]int
var attackHitmarks = []int{10}

const normalHitNum = 1

func init() {
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0], 10)
}

func (c *char) Attack(p map[string]int) action.ActionInfo {
	travel, ok := p["travel"]
	if !ok {
		travel = 10
	}

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       fmt.Sprintf("Normal %v", c.NormalCounter),
		AttackTag:  combat.AttackTagNormal,
		ICDTag:     combat.ICDTagNormalAttack,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeBlunt,
		Element:    attributes.Geo,
		Durability: 25,
		Mult:       attack[c.TalentLvlAttack()],
	}

	done := false
	cb := func(a combat.AttackCB) {
		if done {
			return
		}
		count := c.Tags["jade"]
		//if we're at 7 dont increase but also dont reset back to 3
		if count != 7 {
			count++
			if count > 3 {
				count = 3
			} else {
				c.Core.Log.NewEvent("adding star jade", glog.LogCharacterEvent, c.Index, "count", count)
			}
			c.Tags["jade"] = count
		}
		done = true
	}

	r := 0.1
	if c.Base.Cons >= 1 {
		r = 2
	}
	c.Core.QueueAttack(ai, combat.NewDefCircHit(r, false, combat.TargettableEnemy), attackHitmarks[0], attackHitmarks[0]+travel, cb)
	c.Core.QueueAttack(ai, combat.NewDefCircHit(r, false, combat.TargettableEnemy), attackHitmarks[0], attackHitmarks[0]+travel, cb)

	//defer c.AdvanceNormalIndex()

	return action.ActionInfo{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[0][action.InvalidAction],
		CanQueueAfter:   attackHitmarks[0],
		State:           action.NormalAttackState,
	}
}

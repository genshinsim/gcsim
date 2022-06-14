package yoimiya

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var attackFrames [][]int
var attackHitmarks = [][]int{{15, 24}, {17}, {25}, {11, 26}, {17}}

const normalHitNum = 5

func init() {
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0][1], 35)
	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1][0], 26)
	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2][0], 39)
	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3][1], 44)
	attackFrames[4] = frames.InitNormalCancelSlice(attackHitmarks[4][0], 52)
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
		Element:    attributes.Physical,
		Durability: 25,
	}

	particleCB := func(combat.AttackCB) {
		if c.Core.Status.Duration("yoimiyaskill") <= 0 {
			return
		}
		if c.Core.F < c.lastPart {
			return
		}
		c.lastPart = c.Core.F + 300 //every 5 second

		var count float64 = 2
		if c.Core.Rand.Float64() < 0.5 {
			count = 3
		}
		c.Core.QueueParticle("yoimiya", count, attributes.Pyro, 100)
	}

	var totalMV float64
	for i, mult := range attack[c.NormalCounter] {
		ai.Mult = mult[c.TalentLvlAttack()]
		totalMV += mult[c.TalentLvlAttack()]
		c.Core.QueueAttack(
			ai,
			combat.NewDefCircHit(0.1, false, combat.TargettableEnemy),
			attackHitmarks[c.NormalCounter][i],
			attackHitmarks[c.NormalCounter][i]+travel,
			particleCB,
		)
	}

	if c.Base.Cons >= 6 && c.Core.Status.Duration("yoimiyaskill") > 0 && c.Core.Rand.Float64() < 0.5 {
		//trigger attack
		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       fmt.Sprintf("Kindling (C6) - N%v", c.NormalCounter),
			AttackTag:  combat.AttackTagNormal,
			ICDTag:     combat.ICDTagNormalAttack,
			ICDGroup:   combat.ICDGroupDefault,
			Element:    attributes.Pyro,
			Durability: 25,
			Mult:       totalMV * 0.6,
		}
		//TODO: frames?
		c.Core.QueueAttack(ai, combat.NewDefCircHit(0.1, false, combat.TargettableEnemy), 0, attackHitmarks[c.NormalCounter][0]+travel+5)
	}

	defer c.AdvanceNormalIndex()

	return action.ActionInfo{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackHitmarks[c.NormalCounter][len(attackHitmarks[c.NormalCounter])-1],
		Post:            attackHitmarks[c.NormalCounter][len(attackHitmarks[c.NormalCounter])-1],
		State:           action.NormalAttackState,
	}
}

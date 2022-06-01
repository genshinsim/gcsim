package ayato

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

var attackFrames [][]int
var attackHitmarks = [][]int{{12}, {18}, {20}, {22, 25}, {41}}
var shunsuikenFrames []int

const shunsuikenHitmark = 5

func (c *char) Attack(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		Abil:       fmt.Sprintf("Normal %v", c.NormalCounter),
		ActorIndex: c.Index,
		AttackTag:  combat.AttackTagNormal,
		ICDTag:     combat.ICDTagNormalAttack,
		ICDGroup:   combat.ICDGroupDefault,
		Element:    attributes.Physical,
		Durability: 25,
	}

	if c.Core.Status.Duration("soukaikanka") > 0 {
		ai.Mult = shunsuiken[c.NormalCounter][c.TalentLvlSkill()]
		c.Core.QueueAttack(ai, combat.NewDefCircHit(2, false, combat.TargettableEnemy), 0, shunsuikenHitmark, c.generateParticles, c.skillStacks)
	} else {
		for i, mult := range attack[c.NormalCounter] {
			ai.Mult = mult[c.TalentLvlAttack()]
			c.Core.QueueAttack(
				ai,
				combat.NewDefSingleTarget(1, combat.TargettableEnemy),
				attackHitmarks[c.NormalCounter][i],
				attackHitmarks[c.NormalCounter][i],
			)
		}
	}

	defer c.AdvanceNormalIndex()

	// while in skill
	if c.Core.Status.Duration("soukaikanka") > 0 {
		return action.ActionInfo{
			Frames: func(next action.Action) int {
				return frames.AtkSpdAdjust(shunsuikenFrames[next], c.Stat(attributes.AtkSpd))
			},
			AnimationLength: shunsuikenFrames[action.InvalidAction],
			CanQueueAfter:   shunsuikenHitmark,
			Post:            shunsuikenHitmark,
			State:           action.NormalAttackState,
		}
	}

	// normal state
	return action.ActionInfo{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackHitmarks[c.NormalCounter][len(attackHitmarks[c.NormalCounter])-1],
		Post:            attackHitmarks[c.NormalCounter][len(attackHitmarks[c.NormalCounter])-1],
		State:           action.NormalAttackState,
	}
}

func (c *char) generateParticles(ac combat.AttackCB) {
	if c.Core.F > c.particleICD {
		c.particleICD = c.Core.F + 114
		var count float64 = 1
		if c.Core.Rand.Float64() < 0.5 {
			count++
		}
		c.Core.QueueParticle("ayato", count, attributes.Hydro, 80)
	}
}

func (c *char) skillStacks(ac combat.AttackCB) {
	if c.stacks < c.stacksMax {
		c.stacks++
		c.Core.Log.NewEvent("gained namisen stack", glog.LogCharacterEvent, c.Index, "stacks", c.stacks)
	}
}

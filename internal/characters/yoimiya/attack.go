package yoimiya

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
)

var (
	attackFrames   [][]int
	attackHitmarks = [][]int{{15, 24}, {17}, {25}, {11, 26}, {17}}
)

const normalHitNum = 5

func init() {
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0][1], 35)
	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1][0], 26)
	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2][0], 39)
	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3][1], 44)
	attackFrames[4] = frames.InitNormalCancelSlice(attackHitmarks[4][0], 52)
}

func (c *char) Attack(p map[string]int) (action.Info, error) {
	travel, ok := p["travel"]
	if !ok {
		travel = 10
	}

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       fmt.Sprintf("Normal %v", c.NormalCounter),
		AttackTag:  attacks.AttackTagNormal,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypePierce,
		Element:    attributes.Physical,
		Durability: 25,
	}

	if c.StatusIsActive("yoimiyaskill") {
		ai.ICDTag = attacks.ICDTagNormalAttack
	}

	var particleCB combat.AttackCBFunc
	if c.StatusIsActive(skillKey) {
		particleCB = c.particleCB
	}
	a1CB := c.makeA1CB()
	c2CB := c.makeC2CB()

	var totalMV float64
	for i, mult := range attack[c.NormalCounter] {
		ai.Mult = mult[c.TalentLvlAttack()]
		totalMV += mult[c.TalentLvlAttack()]
		c.Core.QueueAttack(
			ai,
			combat.NewBoxHit(
				c.Core.Combat.Player(),
				c.Core.Combat.PrimaryTarget(),
				geometry.Point{Y: -0.5},
				0.1,
				1,
			),
			attackHitmarks[c.NormalCounter][i],
			attackHitmarks[c.NormalCounter][i]+travel,
			particleCB,
			a1CB,
			c2CB,
		)
	}

	if c.Base.Cons >= 6 && c.StatusIsActive(skillKey) && c.Core.Rand.Float64() < 0.5 {
		// trigger attack
		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       fmt.Sprintf("Kindling (C6) - N%v", c.NormalCounter),
			AttackTag:  attacks.AttackTagNormal,
			ICDTag:     attacks.ICDTagNormalAttack,
			ICDGroup:   attacks.ICDGroupDefault,
			StrikeType: attacks.StrikeTypePierce,
			Element:    attributes.Pyro,
			Durability: 25,
			Mult:       totalMV * 0.6,
		}
		// TODO: frames?
		c.Core.QueueAttack(
			ai,
			combat.NewBoxHit(
				c.Core.Combat.Player(),
				c.Core.Combat.PrimaryTarget(),
				geometry.Point{Y: -0.5},
				0.1,
				1,
			),
			0,
			attackHitmarks[c.NormalCounter][0]+travel+5,
			a1CB,
			c2CB,
		)
	}

	defer c.AdvanceNormalIndex()

	return action.Info{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackHitmarks[c.NormalCounter][len(attackHitmarks[c.NormalCounter])-1],
		State:           action.NormalAttackState,
	}, nil
}

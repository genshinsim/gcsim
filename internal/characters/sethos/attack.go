package sethos

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
)

const normalHitNum = 3

var (
	attackFrames   [][]int
	attackHitmarks = [][]int{{10}, {12, 12 + 3}, {39}}
)

func init() {
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0][0], 19) // N1 -> Walk
	attackFrames[0][action.ActionAttack] = 17
	attackFrames[0][action.ActionCharge] = 15

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1][1], 52) // N2 -> Walk
	attackFrames[1][action.ActionAttack] = 32
	attackFrames[1][action.ActionCharge] = 45

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2][0], 118) // N3 -> Walk
	attackFrames[2][action.ActionAttack] = 63
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

	ap := combat.NewBoxHit(
		c.Core.Combat.Player(),
		c.Core.Combat.PrimaryTarget(),
		geometry.Point{Y: -0.5},
		0.1,
		1,
	)

	for i, mult := range attack[c.NormalCounter] {
		c.QueueCharTask(func() {
			var c4cb combat.AttackCBFunc
			if c.StatusIsActive(burstBuffKey) {
				ai.Abil = fmt.Sprintf("Dusk Bolt %v", c.NormalCounter)
				ai.AttackTag = attacks.AttackTagExtra
				ai.ICDTag = attacks.ICDTagElementalBurst
				ai.Element = attributes.Electro
				ai.FlatDmg += burstEM[c.TalentLvlBurst()] * c.Stat(attributes.EM)

				deltaPos := c.Core.Combat.Player().Pos().Sub(c.Core.Combat.PrimaryTarget().Pos())
				dist := deltaPos.Magnitude()

				// simulate piercing. Extends 15 units from player
				ap = combat.NewBoxHit(
					c.Core.Combat.Player(),
					c.Core.Combat.PrimaryTarget(),
					geometry.Point{Y: -dist},
					0.1,
					15,
				)
				c4cb = c.makeC4cb()
			}
			ai.Mult = mult[c.TalentLvlAttack()]
			c.Core.QueueAttack(
				ai,
				ap,
				0,
				travel,
				c4cb,
			)
		}, attackHitmarks[c.NormalCounter][i])
	}

	defer c.AdvanceNormalIndex()

	return action.Info{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackHitmarks[c.NormalCounter][len(attackHitmarks[c.NormalCounter])-1],
		State:           action.NormalAttackState,
	}, nil
}

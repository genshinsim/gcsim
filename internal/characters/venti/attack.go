package venti

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var (
	attackFrames   [][]int
	attackHitmarks = [][]int{{17, 27}, {19}, {28}, {15, 28}, {17}, {49}}
)

const normalHitNum = 6

func init() {
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0][1], 30)
	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1][0], 38)
	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2][0], 33)
	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3][1], 31)
	attackFrames[4] = frames.InitNormalCancelSlice(attackHitmarks[4][0], 22)
	attackFrames[5] = frames.InitNormalCancelSlice(attackHitmarks[5][0], 98)
}

func (c *char) Attack(p map[string]int) (action.Info, error) {
	travel, ok := p["travel"]
	if !ok {
		travel = 10
	}

	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       fmt.Sprintf("Normal %v", c.NormalCounter),
		AttackTag:  attacks.AttackTagNormal,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypePierce,
		Element:    attributes.Physical,
		Durability: 25,
	}

	for i, mult := range attack[c.NormalCounter] {
		ai.Mult = mult[c.TalentLvlAttack()]

		if c.IsHexerei && c.Core.Player.GetHexereiCount() > 1 && c.StatusIsActive(burstKey) {
			deltaPos := c.Core.Combat.Player().Pos().Sub(c.Core.Combat.PrimaryTarget().Pos())
			dist := deltaPos.Magnitude()

			ai.Element = attributes.Anemo

			ap := combat.NewBoxHit(
				c.Core.Combat.Player(),
				c.Core.Combat.PrimaryTarget(),
				info.Point{Y: -dist},
				0.1,
				15,
			)

			c.Core.QueueAttack(
				ai,
				ap,
				attackHitmarks[c.NormalCounter][i],
				attackHitmarks[c.NormalCounter][i]+travel,
				c.hexAttackCB,
			)

			c.c1(ai, attackHitmarks[c.NormalCounter][i], travel)
		} else {
			c.Core.QueueAttack(
				ai,
				combat.NewBoxHit(
					c.Core.Combat.Player(),
					c.Core.Combat.PrimaryTarget(),
					info.Point{Y: -0.5},
					0.1,
					1,
				),
				attackHitmarks[c.NormalCounter][i],
				attackHitmarks[c.NormalCounter][i]+travel,
			)
		}
	}

	defer c.AdvanceNormalIndex()

	return action.Info{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackHitmarks[c.NormalCounter][len(attackHitmarks[c.NormalCounter])-1],
		State:           action.NormalAttackState,
	}, nil
}

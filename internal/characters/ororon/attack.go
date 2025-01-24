package ororon

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
)

var attackFrames [][]int
var attackHitmarks = []int{12, 10, 20}

const normalHitNum = 3

func init() {
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0], 32) // N1 -> Walk
	attackFrames[0][action.ActionAttack] = 24
	attackFrames[0][action.ActionAim] = 19

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1], 53) // N2 -> Walk
	attackFrames[1][action.ActionAttack] = 27
	attackFrames[1][action.ActionAim] = 19

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2], 70) // N3 -> N1
	attackFrames[2][action.ActionWalk] = 61
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
		ICDTag:     attacks.ICDTagNormalAttack,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypePierce,
		Element:    attributes.Physical,
		Durability: 25,
		Mult:       attack[c.NormalCounter][c.TalentLvlAttack()],
	}
	ap := combat.NewBoxHit(c.Core.Combat.Player(), c.Core.Combat.PrimaryTarget(), geometry.Point{Y: -0.5}, 0.1, 1)
	c.Core.QueueAttack(
		ai,
		ap,
		attackHitmarks[c.NormalCounter],
		attackHitmarks[c.NormalCounter]+travel,
	)

	defer c.AdvanceNormalIndex()

	return action.Info{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackHitmarks[c.NormalCounter],
		State:           action.NormalAttackState,
	}, nil
}

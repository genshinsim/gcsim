package lyney

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
	attackReleases = [][]int{{17}, {12}, {20, 29}, {29}}
)

const normalHitNum = 4

func init() {
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackReleases[0][0], 32) // N1 -> Walk
	attackFrames[0][action.ActionAttack] = 22
	attackFrames[0][action.ActionAim] = 22

	attackFrames[1] = frames.InitNormalCancelSlice(attackReleases[1][0], 34) // N2 -> CA
	attackFrames[1][action.ActionAttack] = 24
	attackFrames[1][action.ActionWalk] = 31

	attackFrames[2] = frames.InitNormalCancelSlice(attackReleases[2][1], 86) // N3 -> Walk
	attackFrames[2][action.ActionAttack] = 39
	attackFrames[2][action.ActionAim] = 81

	attackFrames[3] = frames.InitNormalCancelSlice(attackReleases[3][0], 66) // N4 -> Walk
	attackFrames[3][action.ActionAttack] = 59
	attackFrames[3][action.ActionAim] = 500 // TODO: this action is illegal; need better way to handle it
}

func (c *char) Attack(p map[string]int) action.ActionInfo {
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

	for i, mult := range attack[c.NormalCounter] {
		ai.Mult = mult[c.TalentLvlAttack()]
		c.Core.QueueAttack(
			ai,
			combat.NewBoxHit(
				c.Core.Combat.Player(),
				c.Core.Combat.PrimaryTarget(),
				geometry.Point{Y: -0.5},
				0.1,
				1,
			),
			attackReleases[c.NormalCounter][i],
			attackReleases[c.NormalCounter][i]+travel,
		)
	}

	defer c.AdvanceNormalIndex()

	return action.ActionInfo{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackReleases[c.NormalCounter][len(attackReleases[c.NormalCounter])-1],
		State:           action.NormalAttackState,
	}
}

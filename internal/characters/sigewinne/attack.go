package sigewinne

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
var attackHitmarks = []int{12, 14, 38}

const normalHitNum = 3

func init() {
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0], 33) // N1 -> Walk
	attackFrames[0][action.ActionAttack] = 20

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1], 40) // N2 -> Walk
	attackFrames[1][action.ActionAttack] = 36

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2], 67) // N3 -> Walk
	attackFrames[2][action.ActionAttack] = 82
}

func (c *char) Attack(p map[string]int) (action.Info, error) {
	if c.burstEarlyCancelled {
		return action.Info{}, fmt.Errorf("%v: Cannot early cancel Super Saturated Syringing with Normal Attack", c.Base.Key)
	}
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
		Mult:       attack[c.NormalCounter][c.TalentLvlAttack()],
	}
	var ap combat.AttackPattern
	if c.NormalCounter != 0 {
		ap = combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 1)
	} else {
		ap = combat.NewBoxHit(c.Core.Combat.Player(), c.Core.Combat.PrimaryTarget(), geometry.Point{Y: -0.5}, 0.1, 1)
	}
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

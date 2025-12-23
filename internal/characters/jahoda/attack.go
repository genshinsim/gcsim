package jahoda

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
	attackHitmarks = [][]int{{10}, {10, 10}, {10}} // Frames needed
)

const normalHitNum = 3

func init() {
	// NA cancels
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0][0], 10) // Frames needed
	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1][1], 10) // Frames needed
	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2][0], 10) // Frames needed

}

func (c *char) Attack(p map[string]int) (action.Info, error) {
	travel, ok := p["travel"]
	if !ok {
		travel = 10
	}

	for i, mult := range attack[c.NormalCounter] {
		ai := info.AttackInfo{
			ActorIndex:   c.Index(),
			Abil:         fmt.Sprintf("Normal %v", c.NormalCounter),
			AttackTag:    attacks.AttackTagNormal,
			ICDTag:       attacks.ICDTagNone,
			ICDGroup:     attacks.ICDGroupDefault,
			StrikeType:   attacks.StrikeTypePierce,
			Element:      attributes.Physical,
			Durability:   25,
			Mult:         mult[c.TalentLvlAttack()],
			HitlagFactor: 0.01,
		}

		ap := combat.NewBoxHit(c.Core.Combat.Player(), c.Core.Combat.PrimaryTarget(), info.Point{Y: -0.5}, 0.1, 1)
		c.Core.QueueAttack(
			ai,
			ap,
			attackHitmarks[c.NormalCounter][i],
			attackHitmarks[c.NormalCounter][i]+travel,
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

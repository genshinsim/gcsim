package sandrone

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

// TODO: frames estimated from pre-release footage (60 FPS, no hit markers), +-2-4 frames.
var (
	attackFrames   [][]int
	attackHitmarks = []int{18, 20, 30}
)

func init() {
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0], 60) // N1 -> W
	attackFrames[0][action.ActionAttack] = 52                             // N1 -> N2

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1], 80) // N2 -> W
	attackFrames[1][action.ActionAttack] = 75                             // N2 -> N3

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2], 90) // N3 -> W
	attackFrames[2][action.ActionAttack] = 90                             // N3 -> N1
}

// Generates threads that control her weapon using formulae, performing up to 3 consecutive strikes.
func (c *char) Attack(p map[string]int) (action.Info, error) {
	c.endHovermech()

	ai := info.AttackInfo{
		Abil:       fmt.Sprintf("Normal %v", c.NormalCounter),
		ActorIndex: c.Index(),
		AttackTag:  attacks.AttackTagNormal,
		ICDTag:     attacks.ICDTagNormalAttack,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeBlunt,
		Element:    attributes.Physical,
		Durability: 25,
		Mult:       attack[c.NormalCounter][c.TalentLvlAttack()],
	}

	// N1/N2 fire the weapon forward like a beam, N3 slams it onto the target
	ap := combat.NewBoxHitOnTarget(c.Core.Combat.Player(), info.Point{Y: -1}, 2.5, 6)
	if c.NormalCounter == 2 {
		ap = combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 3)
	}

	c.QueueCharTask(func() {
		c.Core.QueueAttack(ai, ap, 0, 0)
	}, attackHitmarks[c.NormalCounter])

	defer c.AdvanceNormalIndex()

	return action.Info{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackHitmarks[c.NormalCounter],
		State:           action.NormalAttackState,
	}, nil
}

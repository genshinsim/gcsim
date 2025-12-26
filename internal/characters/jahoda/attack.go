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
	attackHitmarks = [][]int{{14}, {15, 29}, {40}}
)

const normalHitNum = 3

func init() {
	// NA cancels
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0][0], 35) // N1 -> Walk
	attackFrames[0][action.ActionAttack] = 30
	attackFrames[0][action.ActionAim] = 30

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1][1], 52) // N2 -> Walk
	attackFrames[1][action.ActionAttack] = 48
	attackFrames[1][action.ActionAim] = 47

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2][0], 99) // N3 -> Walk
	attackFrames[2][action.ActionAttack] = 88
	attackFrames[2][action.ActionAim] = 89

}

func (c *char) Attack(p map[string]int) (action.Info, error) {
	if c.StatusIsActive(shadowPursuitKey) {
		c.Core.Tasks.Add(c.drainFlask(c.skillSrc), 0)
		return action.Info{
			Frames:          frames.NewAbilFunc(skillCancelFrames),
			AnimationLength: skillCancelFrames[action.InvalidAction],
			CanQueueAfter:   skillCancelFrames[action.ActionDash], // earliest cancel
			State:           action.SkillState,
		}, nil
	}

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

package wriothesley

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/targets"
)

// TODO: heizou based frames
var (
	attackFrames          [][]int
	attackHitmarks        = [][]int{{12}, {13}, {21}, {13, 19}, {31}}
	attackHitlagHaltFrame = [][]float64{{0.01}, {0.01}, {0.03}, {0, 0.01}, {0.06}}
	attackHitboxes        = [][]float64{{2, 3}, {2, 3}, {2.5, 3}, {2, 3}, {3, 3}}
	attackHitboxesSkill   = [][]float64{{2.4, 3.4}, {2.4, 3.4}, {2.8, 3.4}, {2.4, 3.4}, {3.4, 3.4}}
)

const normalHitNum = 5

func init() {
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0][0], 21)
	attackFrames[0][action.ActionAttack] = 20
	attackFrames[0][action.ActionCharge] = 21

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1][0], 21)
	attackFrames[1][action.ActionAttack] = 17
	attackFrames[1][action.ActionCharge] = 21

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2][0], 46)
	attackFrames[2][action.ActionAttack] = 45
	attackFrames[2][action.ActionCharge] = 46

	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3][1], 38)
	attackFrames[3][action.ActionAttack] = 36
	attackFrames[3][action.ActionCharge] = 38

	attackFrames[4] = frames.InitNormalCancelSlice(attackHitmarks[4][0], 66)
	attackFrames[4][action.ActionAttack] = 66
	attackFrames[4][action.ActionCharge] = 500
}

func (c *char) Attack(p map[string]int) action.ActionInfo {
	for i, mult := range attack[c.NormalCounter] {
		ai := combat.AttackInfo{
			ActorIndex:         c.Index,
			Abil:               fmt.Sprintf("Normal %v", c.NormalCounter),
			AttackTag:          attacks.AttackTagNormal,
			ICDTag:             attacks.ICDTagNormalAttack,
			ICDGroup:           attacks.ICDGroupDefault,
			StrikeType:         attacks.StrikeTypeDefault,
			Element:            attributes.Cryo,
			Durability:         25,
			Mult:               mult[c.TalentLvlAttack()],
			HitlagFactor:       0.01,
			HitlagHaltFrames:   attackHitlagHaltFrame[c.NormalCounter][i] * 60,
			CanBeDefenseHalted: true,
		}
		if c.NormalCounter == 3 && i == 0 { // N4-1
			ai.HitlagFactor = 0
		}

		hithoxes := attackHitboxes
		callbacks := make([]combat.AttackCBFunc, 0)
		if c.StatusIsActive(skillKey) {
			hithoxes = attackHitboxesSkill
			callbacks = append(callbacks, c.particleCB)

			if c.CurrentHPRatio() > 0.5 {
				ai.Mult *= skill[c.TalentLvlSkill()]
				callbacks = append(callbacks, c.chillingPenalty)
			}
			if c.Base.Cons >= 1 && c.NormalCounter == 4 {
				callbacks = append(callbacks, func(a combat.AttackCB) {
					if a.Target.Type() != targets.TargettableEnemy {
						return
					}
					c.a1Add()
				})
			}
		}

		ap := combat.NewBoxHitOnTarget(
			c.Core.Combat.Player(),
			geometry.Point{Y: -0.2},
			hithoxes[c.NormalCounter][0],
			hithoxes[c.NormalCounter][1],
		)
		c.Core.QueueAttack(ai, ap, attackHitmarks[c.NormalCounter][i], attackHitmarks[c.NormalCounter][i], callbacks...)
	}

	defer c.AdvanceNormalIndex()

	return action.ActionInfo{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackHitmarks[c.NormalCounter][len(attackHitmarks[c.NormalCounter])-1],
		State:           action.NormalAttackState,
	}

}

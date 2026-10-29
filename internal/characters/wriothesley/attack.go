package wriothesley

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
	attackFrames          [][]int
	attackHitmarks        = [][]int{{12}, {10}, {18}, {25, 35}, {39}}
	attackHitlagFactor    = [][]float64{{0}, {0}, {0.01}, {0, 0}, {0.01}}
	attackHitlagHaltFrame = [][]float64{{0}, {0}, {0.03}, {0, 0}, {0.06}}
	attackHitboxes        = [][][]float64{{{2, 3}, {2, 3}, {2.5, 3}, {2, 3}, {3, 3}}, {{2.4, 3.4}, {2.4, 3.4}, {2.8, 3.4}, {2.4, 3.4}, {3.4, 3.4}}}
	attackOffsets         = []float64{-0.2, 0}
)

const normalHitNum = 5

func init() {
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0][0], 27)
	attackFrames[0][action.ActionAttack] = 14
	attackFrames[0][action.ActionCharge] = 23

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1][0], 25)
	attackFrames[1][action.ActionAttack] = 13
	attackFrames[1][action.ActionCharge] = 20

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2][0], 41)
	attackFrames[2][action.ActionAttack] = 24
	attackFrames[2][action.ActionCharge] = 19

	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3][1], 56)
	attackFrames[3][action.ActionAttack] = 41
	attackFrames[3][action.ActionCharge] = 41

	attackFrames[4] = frames.InitNormalCancelSlice(attackHitmarks[4][0], 59)
	attackFrames[4][action.ActionCharge] = 39
	attackFrames[4][action.ActionWalk] = 57
}

func (c *char) Attack(p map[string]int) (action.Info, error) {
	// Apart from this, Normal Attack combo count will not reset for a short time after using Icefang Rush or sprinting.
	switch c.Core.Player.CurrentState() {
	case action.DashState, action.SkillState:
		c.NormalCounter = c.savedNormalCounter
	}
	n := c.NormalCounter
	for i, mult := range attack[c.NormalCounter] {
		ai := info.AttackInfo{
			ActorIndex:         c.Index(),
			Abil:               fmt.Sprintf("Normal %v", n+1),
			AttackTag:          attacks.AttackTagNormal,
			ICDTag:             attacks.ICDTagNormalAttack,
			ICDGroup:           attacks.ICDGroupDefault,
			StrikeType:         attacks.StrikeTypeDefault,
			Element:            attributes.Cryo,
			Durability:         25,
			Mult:               mult[c.TalentLvlAttack()],
			HitlagFactor:       attackHitlagFactor[c.NormalCounter][i],
			HitlagHaltFrames:   attackHitlagHaltFrame[c.NormalCounter][i] * 60,
			CanBeDefenseHalted: true,
		}
		c1N5CB := c.makeC1N5CB(n) // here so that the normalcounter check is correct

		c.QueueCharTask(func() {
			// TODO: when should this check happen?
			skillIndex := 0
			var particleCB info.AttackCBFunc
			var chillingPenalty info.AttackCBFunc
			var snap info.Snapshot
			if c.skillBuffActive() {
				skillIndex = 1
				particleCB = c.particleCB
				chillingPenalty = c.chillingPenalty
				ai.Abil = fmt.Sprintf("Normal %v (Enhanced)", n+1)
				ai.Mult *= skill[c.TalentLvlSkill()]
				ai.Mult *= 1 + c.c2RevelationNAMult()
			}

			if c.skillBuffActive() && c.a1StellarEnabled() && (n == 2 || n == 4) {
				ai.Abil += stellarConductText
				ai.AttackTag = attacks.AttackTagDirectStellarConduct
				ai.IgnoreDefPercent = 1
				ai.Durability = 0
				if n == 2 {
					ai.Mult *= 0.6
				} else {
					ai.AdditionalTags = []attacks.AttackTag{attacks.AttackTagWrioSSCN5}
					ai.Mult *= 0.8
				}
				// this won't overlap with the NA mult because the NA mult only get applied when the stellar field isn't enabled
				ai.Mult *= 1 + c.c2RevelationSSCMult()
			}

			ap := combat.NewBoxHitOnTarget(
				c.Core.Combat.Player(),
				info.Point{Y: attackOffsets[skillIndex]},
				attackHitboxes[skillIndex][c.NormalCounter][0],
				attackHitboxes[skillIndex][c.NormalCounter][1],
			)

			snap = c.Snapshot(&ai)

			if c.skillBuffActive() && c.a1StellarEnabled() {
				c.addC6Buff(&snap)
			}

			c.Core.QueueAttackWithSnap(
				ai,
				snap,
				ap,
				0,
				c1N5CB,
				particleCB,
				chillingPenalty,
			)

			if n == 4 {
				c.c6OnN5(&ai, &snap, &ap)
			}
		}, attackHitmarks[n][i])
	}

	defer func() {
		c.AdvanceNormalIndex()
		c.savedNormalCounter = c.NormalCounter
	}()

	return action.Info{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[n][action.InvalidAction],
		CanQueueAfter:   attackHitmarks[n][len(attackHitmarks[n])-1],
		State:           action.NormalAttackState,
	}, nil
}

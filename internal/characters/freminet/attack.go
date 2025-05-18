package freminet

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

var (
	attackFrames          [][]int
	attackHitmarks        = []int{26, 23, 31, 42}
	attackPoiseDMG        = []float64{112.6437, 107.8804, 136.267, 165.5529}
	attackHitlagHaltFrame = []float64{0.06, 0.06, 0.06, 0.09}
	attackHitboxes        = [][]float64{{2.2}, {2.2}, {2, 3}, {2.4}}
	attackOffsets         = []float64{0.5, 0.5, -0.5, 0.5}
	attackFrostDelay      = []int{10, 9, 10, 10} // delay from hitmark, approximation
)

const normalHitNum = 4

func init() {
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0], 47)
	attackFrames[0][action.ActionAttack] = 32

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1], 49)
	attackFrames[1][action.ActionAttack] = 33

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2], 65)
	attackFrames[2][action.ActionAttack] = 59

	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3], 86)
	attackFrames[3][action.ActionWalk] = 68
}

func (c *char) Attack(p map[string]int) (action.Info, error) {
	if c.skillStacks >= 4 {
		c.NormalCounter = 0
		return c.detonateSkill()
	}

	ai := combat.AttackInfo{
		Abil:               fmt.Sprintf("Normal %v", c.NormalCounter),
		ActorIndex:         c.Index,
		AttackTag:          attacks.AttackTagNormal,
		ICDTag:             attacks.ICDTagNormalAttack,
		ICDGroup:           attacks.ICDGroupDefault,
		StrikeType:         attacks.StrikeTypeBlunt,
		PoiseDMG:           attackPoiseDMG[c.NormalCounter],
		Element:            attributes.Physical,
		Durability:         25,
		Mult:               attack[c.NormalCounter][c.TalentLvlAttack()],
		HitlagFactor:       0.01,
		HitlagHaltFrames:   attackHitlagHaltFrame[c.NormalCounter] * 60,
		CanBeDefenseHalted: true,
	}

	ap := combat.NewCircleHitOnTarget(
		c.Core.Combat.Player(),
		geometry.Point{Y: attackOffsets[c.NormalCounter]},
		attackHitboxes[c.NormalCounter][0],
	)

	if c.NormalCounter == 2 {
		ap = combat.NewBoxHitOnTarget(
			c.Core.Combat.Player(),
			geometry.Point{Y: attackOffsets[c.NormalCounter]},
			attackHitboxes[c.NormalCounter][0],
			attackHitboxes[c.NormalCounter][1],
		)
	}

	c.Core.QueueAttack(ai, ap, attackHitmarks[c.NormalCounter], attackHitmarks[c.NormalCounter])

	if c.StatusIsActive(persTimeKey) {
		frostMod := skillAddNA[c.TalentLvlSkill()]
		if c.StatusIsActive(burstKey) {
			frostMod *= 2
		}

		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Pressurized Floe: Pers Time Frost",
			AttackTag:  attacks.AttackTagElementalArt,
			ICDTag:     attacks.ICDTagElementalArt,
			ICDGroup:   attacks.ICDGroupDefault,
			StrikeType: attacks.StrikeTypeBlunt,
			PoiseDMG:   25,
			Element:    attributes.Cryo,
			Durability: 25,
			Mult:       frostMod,
		}

		c.Core.QueueAttack(
			ai,
			ap,
			attackHitmarks[c.NormalCounter]+attackFrostDelay[c.NormalCounter],
			attackHitmarks[c.NormalCounter]+attackFrostDelay[c.NormalCounter],
		)

		amt := 1
		if c.StatusIsActive(burstKey) {
			amt = 2
		}

		c.skillStacks += amt
		if c.skillStacks > 4 {
			c.skillStacks = 4
		}
		c.Core.Log.NewEvent("freminet skill stacks gained", glog.LogCharacterEvent, c.Index).
			Write("stacks", c.skillStacks)
	}

	defer c.AdvanceNormalIndex()

	return action.Info{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackHitmarks[c.NormalCounter],
		State:           action.NormalAttackState,
	}, nil
}

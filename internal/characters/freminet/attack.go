package freminet

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
	attackFrames [][]int
	// TODO: Freminet; Copied from Chongyun
	attackHitmarks        = []int{26, 23, 31, 42}
	attackHitlagHaltFrame = []float64{.1, .09, .12, .12}
	attackHitboxes        = [][]float64{{2}, {2}, {2}, {2, 3}}
	attackOffsets         = []float64{1, 1, 1, -0.5}
	attackFanAngles       = []float64{360, 270, 360, 360}
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

func (c *char) Attack(p map[string]int) action.ActionInfo {

	if c.skillStacks >= 4 {
		return c.detonateSkill()
	}

	// TODO: Freminet; Copied from Chongyun
	ai := combat.AttackInfo{
		Abil:               fmt.Sprintf("Normal %v", c.NormalCounter),
		ActorIndex:         c.Index,
		AttackTag:          attacks.AttackTagNormal,
		ICDTag:             attacks.ICDTagNormalAttack,
		ICDGroup:           attacks.ICDGroupDefault,
		StrikeType:         attacks.StrikeTypeBlunt,
		Element:            attributes.Physical,
		Durability:         25,
		Mult:               attack[c.NormalCounter][c.TalentLvlAttack()],
		HitlagFactor:       0.01,
		HitlagHaltFrames:   attackHitlagHaltFrame[c.NormalCounter] * 60,
		CanBeDefenseHalted: true,
	}

	ap := combat.NewCircleHitOnTargetFanAngle(
		c.Core.Combat.Player(),
		geometry.Point{Y: attackOffsets[c.NormalCounter]},
		attackHitboxes[c.NormalCounter][0],
		attackFanAngles[c.NormalCounter],
	)

	if c.NormalCounter == 3 {
		ap = combat.NewBoxHitOnTarget(
			c.Core.Combat.Player(),
			geometry.Point{Y: attackOffsets[c.NormalCounter]},
			attackHitboxes[c.NormalCounter][0],
			attackHitboxes[c.NormalCounter][1],
		)
	}

	var cb combat.AttackCBFunc

	if c.StatusIsActive(persTimeKey) {
		cb = c.persTimeCB()
	}

	c.Core.QueueAttack(ai, ap, attackHitmarks[c.NormalCounter], attackHitmarks[c.NormalCounter], cb)

	defer c.AdvanceNormalIndex()

	return action.ActionInfo{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackHitmarks[c.NormalCounter],
		State:           action.NormalAttackState,
	}
}

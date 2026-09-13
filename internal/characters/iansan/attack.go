package iansan

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
	attackHitmarks        = []int{14, 7, 22}
	attackHitlagHaltFrame = []float64{0.01, 0.06, 0}
	attackHitboxes        = [][]float64{{1.75}, {2}, {2.2, 3}}
	attackOffsets         = []float64{0.4, 0.5, -0.5}
	attackFanAngles       = []float64{260, 360, 0}
	attackDefHalt         = []bool{true, true, false}
)

const normalHitNum = 3

func init() {
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0], 27) // N1 -> Walk
	attackFrames[0][action.ActionAttack] = 21
	attackFrames[0][action.ActionCharge] = 21

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1], 31) // N2 -> Walk
	attackFrames[1][action.ActionAttack] = 25
	attackFrames[1][action.ActionCharge] = 25

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2], 55) // N3 -> Walk
	attackFrames[2][action.ActionAttack] = 53
}

func (c *char) Attack(p map[string]int) (action.Info, error) {
	if c.StatusIsActive(fastSkill) {
		// TODO: or c.Core.Player.Exec(action.ActionCharge, c.Base.Key, nil)
		return c.ChargeAttack(p)
	}

	ai := info.AttackInfo{
		ActorIndex:         c.Index(),
		Abil:               fmt.Sprintf("Normal %v", c.NormalCounter),
		Mult:               attack[c.NormalCounter][c.TalentLvlAttack()],
		AttackTag:          attacks.AttackTagNormal,
		ICDTag:             attacks.ICDTagNormalAttack,
		ICDGroup:           attacks.ICDGroupDefault,
		StrikeType:         attacks.StrikeTypeSlash,
		Element:            attributes.Physical,
		Durability:         25,
		HitlagFactor:       0.01,
		HitlagHaltFrames:   attackHitlagHaltFrame[c.NormalCounter] * 60,
		CanBeDefenseHalted: attackDefHalt[c.NormalCounter],
	}
	ap := combat.NewCircleHitOnTargetFanAngle(
		c.Core.Combat.Player(),
		info.Point{Y: attackOffsets[c.NormalCounter]},
		attackHitboxes[c.NormalCounter][0],
		attackFanAngles[c.NormalCounter],
	)

	if c.NormalCounter == 2 {
		ai.PoiseDMG = 59.896
		ai.StrikeType = attacks.StrikeTypeBlunt
		ai.AdditionalTags = []attacks.AttackTag{attacks.AttackTagIansanBisonsaurus}
		ap = combat.NewBoxHitOnTarget(
			c.Core.Combat.Player(),
			info.Point{Y: attackOffsets[c.NormalCounter]},
			attackHitboxes[c.NormalCounter][0],
			attackHitboxes[c.NormalCounter][1],
		)
	}

	c.Core.QueueAttack(ai, ap, attackHitmarks[c.NormalCounter], attackHitmarks[c.NormalCounter])

	defer c.AdvanceNormalIndex()

	return action.Info{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackHitmarks[c.NormalCounter],
		State:           action.NormalAttackState,
	}, nil
}

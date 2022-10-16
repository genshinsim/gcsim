package nilou

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

const normalHitNum = 3

var attackFrames [][]int
var attackHitmarks = []int{12, 9, 17}

var attackHitlagHaltFrame = []float64{0.03, 0.03, 0.06}
var attackDefHalt = []bool{true, true, true}

func init() {
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0], 24) // N1 -> CA
	attackFrames[0][action.ActionAttack] = 20                             // N1 -> N2

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1], 21) // N2 -> CA
	attackFrames[1][action.ActionAttack] = 27                             // N2 -> N3

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2], 58) // N3 -> N1
}

func (c *char) Attack(p map[string]int) action.ActionInfo {
	if c.StatusIsActive(pirouetteStatus) {
		return c.Pirouette(p, NilouSkillTypeDance)
	}
	if c.StatusIsActive(lunarPrayerStatus) {
		return c.SwordDance(p)
	}

	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               fmt.Sprintf("Normal %v", c.NormalCounter),
		AttackTag:          combat.AttackTagNormal,
		ICDTag:             combat.ICDTagNormalAttack,
		ICDGroup:           combat.ICDGroupDefault,
		StrikeType:         combat.StrikeTypeSlash,
		Element:            attributes.Physical,
		Durability:         25,
		Mult:               auto[c.NormalCounter][c.TalentLvlAttack()],
		HitlagHaltFrames:   attackHitlagHaltFrame[c.NormalCounter] * 60,
		HitlagFactor:       0.01,
		CanBeDefenseHalted: attackDefHalt[c.NormalCounter],
	}
	// no multihits so no need for char queue here
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHit(c.Core.Combat.Player(), .3, false, combat.TargettableEnemy),
		attackHitmarks[c.NormalCounter],
		attackHitmarks[c.NormalCounter],
	)

	defer c.AdvanceNormalIndex()

	return action.ActionInfo{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackHitmarks[c.NormalCounter],
		State:           action.NormalAttackState,
	}
}

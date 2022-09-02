package cyno

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var attackFrames [][]int
var attackHitmarks = []int{25, 16, 13, 38}
var attackHitlagHaltFrame = []float64{0.1, 0.1, 0.1, 0.15}
var attackHitlagFactor = []float64{0.01, 0.01, 0.05, 0.01}

const normalHitNum = 4

func init() {
	attackFrames = make([][]int, normalHitNum) // should be 4

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0], 45)  // N1 -> N2
	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1], 33)  // N2 -> N3
	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2], 47)  // N3 -> N4
	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3], 116) // N4 -> N1
}

// TODO: Adjust the attack frame values (this ones are source: i made them the fuck up)
func (c *char) Attack(p map[string]int) action.ActionInfo {
	if c.StatusIsActive(burstKey) {
		return c.attackB(p) //go to burst mode attacks
	}

	for _, mult := range attack[c.NormalCounter] {
		ai := combat.AttackInfo{
			ActorIndex:         c.Index,
			Abil:               fmt.Sprintf("Normal %v", c.NormalCounter),
			Mult:               mult[c.TalentLvlAttack()],
			AttackTag:          combat.AttackTagNormal,
			ICDTag:             combat.ICDTagNormalAttack,
			ICDGroup:           combat.ICDGroupDefault,
			Element:            attributes.Physical,
			Durability:         25,
			HitlagFactor:       attackHitlagFactor[c.NormalCounter],
			HitlagHaltFrames:   attackHitlagHaltFrame[c.NormalCounter],
			CanBeDefenseHalted: true,
		}

		c.Core.QueueAttack(
			ai,
			combat.NewCircleHit(c.Core.Combat.Player(), 0.5, false, combat.TargettableEnemy),
			attackHitmarks[c.NormalCounter],
			attackHitmarks[c.NormalCounter],
		)
	}

	defer c.AdvanceNormalIndex()

	return action.ActionInfo{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackHitmarks[c.NormalCounter],
		State:           action.NormalAttackState,
	}
}

const BurstHitNum = 5 //Burst attack chains have 5

var AttackBFrames [][]int
var AttackBHitmarks = [][]int{{12}, {13}, {11}, {22, 33}, {33}}

func init() {
	// NA cancels (burst)
	AttackBFrames = make([][]int, BurstHitNum)

	AttackBFrames[0] = frames.InitNormalCancelSlice(AttackBHitmarks[0][0], 24)
	AttackBFrames[0][action.ActionAttack] = 19

	AttackBFrames[1] = frames.InitNormalCancelSlice(AttackBHitmarks[1][0], 26)
	AttackBFrames[1][action.ActionAttack] = 16

	AttackBFrames[2] = frames.InitNormalCancelSlice(AttackBHitmarks[2][0], 34)
	AttackBFrames[2][action.ActionAttack] = 16

	AttackBFrames[3] = frames.InitNormalCancelSlice(AttackBHitmarks[3][1], 67)
	AttackBFrames[3][action.ActionAttack] = 44

	AttackBFrames[4] = frames.InitNormalCancelSlice(AttackBHitmarks[4][0], 59)
	AttackBFrames[4][action.ActionCharge] = 500 //TODO: honestly idk what i am doing at dis point pls forgive me Koli

}

// TODO: Adjust the attack frame values (this ones are source: i made them the fuck up)
func (c *char) attackB(p map[string]int) action.ActionInfo {
	for i, mult := range attackB[c.NormalCounter] {
		ai := combat.AttackInfo{
			ActorIndex:         c.Index,
			Abil:               fmt.Sprintf("Pactsworn Pathclearer %v", c.NormalCounter),
			AttackTag:          combat.AttackTagNormal,
			ICDTag:             combat.ICDTagNormalAttack,
			ICDGroup:           combat.ICDGroupDefault,
			Element:            attributes.Electro,
			Durability:         25,
			HitlagFactor:       attackHitlagFactor[c.NormalCounter],
			HitlagHaltFrames:   attackHitlagHaltFrame[c.NormalCounter],
			CanBeDefenseHalted: true,
		}
		// i just copy pasted raiden code, why do we use ax:=ai i will never know
		ax := ai
		ax.Mult = mult[c.TalentLvlBurst()]
		ax.FlatDmg = c.Stat(attributes.EM) //this is A4

		c.QueueCharTask(func() {
			c.Core.QueueAttack(ax, combat.NewCircleHit(c.Core.Combat.Player(), 2, false, combat.TargettableEnemy), 0, 0)
		}, AttackBHitmarks[c.NormalCounter][i])
	}

	defer c.AdvanceNormalIndex()

	return action.ActionInfo{
		Frames:          frames.NewAttackFunc(c.Character, AttackBFrames),
		AnimationLength: AttackBFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   AttackBHitmarks[c.NormalCounter][len(AttackBHitmarks[c.NormalCounter])-1],
		State:           action.NormalAttackState,
	}
}

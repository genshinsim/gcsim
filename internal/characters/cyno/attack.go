package cyno

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var attackFrames [][]int
var attackEarliestCancel = []int{14, 16, 18, 12, 41}
var attackHitmarks = [][]int{{14}, {16}, {12, 21}, {41}}
var attackHitlagHaltFrame = []float64{0.1, 0.1, 0.1, 0.15, 0.15} //TODO:verify this with DM}
var attackHitlagFactor = []float64{0.01, 0.01, 0.05, 0.01, 0.01}

const normalHitNum = 4
const BurstHitNum = 5 //Burst attack chains have 5

func init() {
	attackFrames = make([][]int, normalHitNum) // should be 4

	attackFrames[0] = frames.InitNormalCancelSlice(attackEarliestCancel[0], 18) // N1 -> N2
	attackFrames[1] = frames.InitNormalCancelSlice(attackEarliestCancel[1], 37) // N2 -> N3
	attackFrames[2] = frames.InitNormalCancelSlice(attackEarliestCancel[2], 37) // N3 -> N4
	attackFrames[3] = frames.InitNormalCancelSlice(attackEarliestCancel[3], 34) // N4 -> N1
}

// TODO: Adjust the attack frame values (this ones are source: i made them the fuck up)
func (c *char) Attack(p map[string]int) action.ActionInfo {
	if c.StatusIsActive(burstKey) {
		c.NormalHitNum = BurstHitNum
		return c.attackB(p) //go to burst mode attacks
	}
	c.NormalHitNum = normalHitNum
	for i, mult := range attack[c.NormalCounter] {
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
			attackHitmarks[c.NormalCounter][i],
			attackHitmarks[c.NormalCounter][i],
		)
	}

	defer c.AdvanceNormalIndex()

	return action.ActionInfo{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackEarliestCancel[c.NormalCounter],
		State:           action.NormalAttackState,
	}
}

var attackBFrames [][]int
var attackBEarliestCancel = []int{14, 16, 18, 12, 41}
var attackBHitmarks = [][]int{{14}, {16}, {18}, {12, 21}, {41}}

func init() {
	// NA cancels (burst)
	attackBFrames = make([][]int, BurstHitNum) //Should be 5

	attackBFrames[0] = frames.InitNormalCancelSlice(attackEarliestCancel[0], 18) // N1 -> N2
	attackBFrames[1] = frames.InitNormalCancelSlice(attackEarliestCancel[1], 37) // N2 -> N3
	attackBFrames[2] = frames.InitNormalCancelSlice(attackEarliestCancel[2], 37) // N3 -> N4
	attackBFrames[3] = frames.InitNormalCancelSlice(attackEarliestCancel[3], 34) // N4 -> N5
	attackBFrames[4] = frames.InitNormalCancelSlice(attackEarliestCancel[4], 61) // N5 -> N1
	//TODO: honestly idk what i am doing at dis point pls forgive me Koli

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
		}, attackBHitmarks[c.NormalCounter][i])
	}

	defer c.AdvanceNormalIndex()

	return action.ActionInfo{
		Frames:          frames.NewAttackFunc(c.Character, attackBFrames),
		AnimationLength: attackBFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackBEarliestCancel[c.NormalCounter],
		State:           action.NormalAttackState,
	}
}

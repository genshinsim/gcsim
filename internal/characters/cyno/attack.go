package cyno

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var (
	attackFrames          [][]int
	attackHitmarks        = [][]int{{14}, {17}, {13, 22}, {27}}
	attackHitlagHaltFrame = []float64{0.1, 0.1, 0.1, 0.15, 0.15} // TODO:verify this with DM}
	attackHitlagFactor    = []float64{0.01, 0.01, 0.05, 0.01, 0.01}
)

const normalHitNum = 4

func init() {
	attackFrames = make([][]int, normalHitNum)                               // should be 4
	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0][0], 28) // N1 -> N2
	attackFrames[0][action.ActionAttack] = 15

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1][0], 23) // N2 -> N3
	attackFrames[1][action.ActionAttack] = 22

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2][1], 27) // N3 -> N4
	attackFrames[2][action.ActionCharge] = 26

	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3][0], 58) // N4 -> N1
	attackFrames[3][action.ActionCharge] = 500                               // impossible action
}

func (c *char) Attack(p map[string]int) action.ActionInfo {
	if c.StatusIsActive(burstKey) {
		c.NormalHitNum = BurstHitNum
		return c.attackB(p) // go to burst mode attacks
	}
	if c.NormalHitNum >= BurstHitNum { // this should avoid the panic error
		c.NormalHitNum = normalHitNum
		c.ResetNormalCounter() // TODO:verify is cyno resets his attack string if burst expires
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
		CanQueueAfter:   attackHitmarks[c.NormalCounter][len(attackHitmarks[c.NormalCounter])-1],
		State:           action.NormalAttackState,
	}
}

const BurstHitNum = 5 // Burst attack chains have 5

var (
	attackBFrames         [][]int
	attackBEarliestCancel = []int{14, 16, 18, 12, 41}
	attackBHitmarks       = [][]int{{14}, {16}, {18}, {12, 21}, {41}}
)

func init() {
	// NA cancels (burst)
	attackBFrames = make([][]int, BurstHitNum) // Should be 5

	attackBFrames[0] = frames.InitNormalCancelSlice(attackBEarliestCancel[0], 18) // N1 -> N2
	attackBFrames[1] = frames.InitNormalCancelSlice(attackBEarliestCancel[1], 37) // N2 -> N3
	attackBFrames[2] = frames.InitNormalCancelSlice(attackBEarliestCancel[2], 37) // N3 -> N4
	attackBFrames[3] = frames.InitNormalCancelSlice(attackBEarliestCancel[3], 34) // N4 -> N5
	attackBFrames[4] = frames.InitNormalCancelSlice(attackBEarliestCancel[4], 61) // N5 -> N1
	// TODO: honestly idk what i am doing at dis point pls forgive me Koli
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
		ax.FlatDmg = c.Stat(attributes.EM) * 1.5 // this is A4

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

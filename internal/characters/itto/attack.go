package itto

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

var attackHitmarks = []int{23, 25, 16, 48}
var attackHitlagHaltFrame = []float64{0.08, 0.08, 0.10, 0.10}

var attackFrames [][][]int

const normalHitNum = 4

type IttoAttackState int

const (
	InvalidAttackState IttoAttackState = iota - 1
	attack0Stacks
	attack1PlusStacks
	attackEndState
)

func init() {
	attackFrames = make([][][]int, attackEndState)
	attackFrames[attack0Stacks] = make([][]int, normalHitNum)
	attackFrames[attack0Stacks][0] = frames.InitNormalCancelSlice(attackHitmarks[0], 33) // N1 -> N2
	attackFrames[attack0Stacks][1] = frames.InitNormalCancelSlice(attackHitmarks[1], 36) // N2 -> N3
	attackFrames[attack0Stacks][2] = frames.InitNormalCancelSlice(attackHitmarks[2], 43) // N3 -> N4
	attackFrames[attack0Stacks][3] = frames.InitNormalCancelSlice(attackHitmarks[3], 83) // N4 -> N1

	attackFrames[attack0Stacks][0][action.ActionCharge] = 41  // N1 -> CA0
	attackFrames[attack0Stacks][1][action.ActionCharge] = 51  // N2 -> CA0
	attackFrames[attack0Stacks][2][action.ActionCharge] = 57  // N3 -> CA0
	attackFrames[attack0Stacks][3][action.ActionCharge] = 500 // N4 -> CA0, TODO: this action is illegal; need better way to handle it

	attackFrames[attack1PlusStacks] = make([][]int, normalHitNum)
	attackFrames[attack1PlusStacks][0] = frames.InitNormalCancelSlice(attackHitmarks[0], 33) // N1 -> N2
	attackFrames[attack1PlusStacks][1] = frames.InitNormalCancelSlice(attackHitmarks[1], 36) // N2 -> N3
	attackFrames[attack1PlusStacks][2] = frames.InitNormalCancelSlice(attackHitmarks[2], 43) // N3 -> N4
	attackFrames[attack1PlusStacks][3] = frames.InitNormalCancelSlice(attackHitmarks[3], 83) // N4 -> N1

	attackFrames[attack1PlusStacks][0][action.ActionCharge] = 23 // N1 -> CA1/CAF
	attackFrames[attack1PlusStacks][1][action.ActionCharge] = 27 // N2 -> CA1/CAF
	attackFrames[attack1PlusStacks][2][action.ActionCharge] = 21 // N3 -> CA1/CAF
	attackFrames[attack1PlusStacks][3][action.ActionCharge] = 52 // N4 -> CA1/CAF
}

func (c *char) attackState() IttoAttackState {
	if c.Tags[c.stackKey] == 0 {
		// 0 stacks: use NX -> CA0 frames
		return attack0Stacks
	}
	// 1+ stacks: use NX -> CA1/CAF frames (they are the same here)
	return attack1PlusStacks
}

// Normal Attack:
// Perform up to 4 consecutive strikes.
// When the 2nd and 4th strikes hit opponents, Itto will gain 1 and 2 stacks of Superlative Superstrength, respectively.
// Max 5 stacks. Triggering this effect will refresh the current duration of any existing stacks.
// Additionally, Itto's Normal Attack combo does not immediately reset after sprinting or using his Elemental Skill, "Masatsu Zetsugi: Akaushi Burst!"
func (c *char) Attack(p map[string]int) action.ActionInfo {
	// handle Dasshu
	lastWasItto := c.Core.Player.LastAction.Char == c.Index
	lastAction := c.Core.Player.LastAction.Type

	// don't reset attack string if previous action was NA/Dash/Skill
	if lastWasItto && (lastAction == action.ActionAttack || lastAction == action.ActionDash || lastAction == action.ActionSkill) {
		c.NormalCounter = c.savedNormalCounter
	} else {
		c.NormalCounter = 0
	}

	// Attack
	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               fmt.Sprintf("Normal %v", c.NormalCounter),
		Mult:               attack[c.NormalCounter][c.TalentLvlAttack()],
		AttackTag:          combat.AttackTagNormal,
		ICDTag:             combat.ICDTagNormalAttack,
		ICDGroup:           combat.ICDGroupDefault,
		StrikeType:         combat.StrikeTypeBlunt,
		Element:            attributes.Physical,
		Durability:         25,
		HitlagHaltFrames:   attackHitlagHaltFrame[c.NormalCounter] * 60,
		HitlagFactor:       0.01,
		CanBeDefenseHalted: true,
	}

	// check burst status for radius
	// TODO: proper hitbox
	radius := 1.0
	if c.StatModIsActive(c.burstBuffKey) {
		radius = 2
	}

	// TODO: hitmark is not getting adjusted for atk speed
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHit(c.Core.Combat.Player(), radius, false, combat.TargettableEnemy),
		attackHitmarks[c.NormalCounter],
		attackHitmarks[c.NormalCounter],
	)

	// Add superlative strength stacks on damage
	amount := 0
	switch c.NormalCounter {
	case 0:
		// N1
		if c.StatModIsActive(c.burstBuffKey) {
			amount = 1
		}
	case 1:
		// N2
		amount = 1
	case 2:
		// N3
		if c.StatModIsActive(c.burstBuffKey) {
			amount = 1
		}
	case 3:
		// N4
		amount = 2
	}

	if amount > 0 {
		c.changeStacks(amount)
		c.Core.Log.NewEvent("itto attack stack added", glog.LogCharacterEvent, c.Index).
			Write("stacks", c.Tags[c.stackKey])
	}

	// handle NX -> CA0/CA1/CAF frames
	state := c.attackState()

	defer c.AdvanceNormalIndex()

	// save the next NA in case of Dasshu
	c.savedNormalCounter = c.NormalCounter + 1
	if c.savedNormalCounter == c.NormalHitNum {
		c.savedNormalCounter = 0
	}

	return action.ActionInfo{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames[state]),
		AnimationLength: attackFrames[state][c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackHitmarks[c.NormalCounter],
		State:           action.NormalAttackState,
	}
}

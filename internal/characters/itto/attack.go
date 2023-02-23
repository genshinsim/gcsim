package itto

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var (
	attackFrames          [][][]int
	attackHitmarks        = []int{23, 25, 16, 48}
	attackHitlagHaltFrame = []float64{0.08, 0.08, 0.10, 0.10}
	attackHitboxes        = [][][]float64{{{2.5}, {2.5}, {2.5}, {3.2, 6}}, {{3.5}, {3.5}, {3.5}, {3.8, 8}}}
	attackOffsets         = [][]float64{{0.8, 0.8, 0.85, -1.5}, {0.8, 0.8, 0.8, -1.7}}
)

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
	attackFrames[attack0Stacks][0] = frames.InitNormalCancelSlice(attackHitmarks[0], 41) // N1 -> CA0
	attackFrames[attack0Stacks][1] = frames.InitNormalCancelSlice(attackHitmarks[1], 51) // N2 -> CA0
	attackFrames[attack0Stacks][2] = frames.InitNormalCancelSlice(attackHitmarks[2], 57) // N3 -> CA0
	attackFrames[attack0Stacks][3] = frames.InitNormalCancelSlice(attackHitmarks[3], 83) // N4 -> N1

	attackFrames[attack0Stacks][0][action.ActionAttack] = 33  // N1 -> N2
	attackFrames[attack0Stacks][1][action.ActionAttack] = 36  // N2 -> N3
	attackFrames[attack0Stacks][2][action.ActionAttack] = 43  // N3 -> N4
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
	if c.Tags[strStackKey] == 0 {
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

	// Additionally, Itto's Normal Attack combo does not immediately reset after sprinting or using his Elemental Skill
	switch c.Core.Player.CurrentState() {
	case action.DashState, action.SkillState:
		c.NormalCounter = c.savedNormalCounter
	}

	// Attack
	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               fmt.Sprintf("Normal %v", c.NormalCounter),
		Mult:               attack[c.NormalCounter][c.TalentLvlAttack()],
		AttackTag:          attacks.AttackTagNormal,
		ICDTag:             combat.ICDTagNormalAttack,
		ICDGroup:           combat.ICDGroupDefault,
		StrikeType:         attacks.StrikeTypeBlunt,
		Element:            attributes.Physical,
		Durability:         25,
		HitlagHaltFrames:   attackHitlagHaltFrame[c.NormalCounter] * 60,
		HitlagFactor:       0.01,
		CanBeDefenseHalted: true,
	}

	// check burst status for hitbox
	attackIndex := 0
	if c.StatModIsActive(burstBuffKey) {
		attackIndex = 1
	}
	ap := combat.NewCircleHitOnTarget(
		c.Core.Combat.Player(),
		combat.Point{Y: attackOffsets[attackIndex][c.NormalCounter]},
		attackHitboxes[attackIndex][c.NormalCounter][0],
	)
	if c.NormalCounter == 3 {
		ap = combat.NewBoxHitOnTarget(
			c.Core.Combat.Player(),
			combat.Point{Y: attackOffsets[attackIndex][c.NormalCounter]},
			attackHitboxes[attackIndex][c.NormalCounter][0],
			attackHitboxes[attackIndex][c.NormalCounter][1],
		)
	}
	// TODO: hitmark is not getting adjusted for atk speed
	c.Core.QueueAttack(ai, ap, attackHitmarks[c.NormalCounter], attackHitmarks[c.NormalCounter])

	// TODO: assume NAs always hit. since it is not possible to know if the next CA is CA0 or CA1/CAF when deciding what CA frames to return.
	// Add superlative strength stacks on damage
	n := c.NormalCounter
	if n == 1 {
		c.addStrStack("attack", 1)
	} else if n == 3 {
		c.addStrStack("attack", 2)
	}
	if c.StatModIsActive(burstBuffKey) && (n == 0 || n == 2) {
		c.addStrStack("q-attack", 1)
	}

	// handle NX -> CA0/CA1/CAF frames
	state := c.attackState()

	defer func() {
		c.AdvanceNormalIndex()
		c.savedNormalCounter = c.NormalCounter
	}()

	return action.ActionInfo{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames[state]),
		AnimationLength: attackFrames[state][c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackHitmarks[c.NormalCounter],
		State:           action.NormalAttackState,
	}
}

package itto

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

// 0 stacks
var chargeFramesDefaultCA0 []int
var chargeFramesN1CA0 []int
var chargeFramesN2CA0 []int

// 1 stack
var chargeFramesDefaultCAF []int
var chargeFramesCACAF []int
var chargeFramesECAF []int

// 2 stacks
var chargeFramesDefaultCA1CAF []int
var chargeFramesCA2CA1CAF []int
var chargeFramesECA1CAF []int

var chargeFramesDefaultCA2CAF []int

// 3+ stacks
var chargeFramesDefaultCA1CA2 []int
var chargeFramesCA2CA1CA2 []int
var chargeFramesECA1CA2 []int

var chargeFramesDefaultCA2CA1 []int

var chargeFrames [][]int

type IttoChargeState int

const (
	InvalidState IttoChargeState = iota - 1
	defaultToCA0
	n1CAFeToCA0
	n2n3ToCA0

	defaultToCAF
	CA1CA2ToCAF
	eToCAF

	defaultToCA1ToCAF
	ca2ToCA1ToCAF
	eToCA1ToCAF

	defaultToCA2ToCAF

	defaultToCA1ToCA2
	CA2ToCA1ToCA2
	eToCA1ToCA2

	defaultToCA2ToCA1

	endState
)

func init() {
	// 0 stacks -> do a CA0
	chargeFramesDefaultCA0 = frames.InitNormalCancelSlice(89, 131)
	chargeFramesN1CA0 = frames.InitNormalCancelSlice(89-14, 131-14) // previous action is N1/CAF/E
	chargeFramesN2CA0 = frames.InitNormalCancelSlice(89-21, 131-21) // previous action is N2/N3

	// 1 stack -> do a CAF
	chargeFramesDefaultCAF = frames.InitNormalCancelSlice(71, 110)
	chargeFramesDefaultCAF[action.ActionSkill] = 76
	chargeFramesDefaultCAF[action.ActionBurst] = 76
	chargeFramesDefaultCAF[action.ActionSwap] = 76

	chargeFramesCACAF = frames.InitNormalCancelSlice(71-25, 110-25) // previous action is CA1/CA2
	chargeFramesCACAF[action.ActionSkill] = 76 - 25
	chargeFramesCACAF[action.ActionBurst] = 76 - 25
	chargeFramesCACAF[action.ActionSwap] = 76 - 25

	chargeFramesECAF = frames.InitNormalCancelSlice(71-17, 110-17) // previous action is E
	chargeFramesECAF[action.ActionSkill] = 76 - 17
	chargeFramesECAF[action.ActionBurst] = 76 - 17
	chargeFramesECAF[action.ActionSwap] = 76 - 17

	// 2 stacks
	// we are doing a CA1, so the next CA has to be CAF
	chargeFramesDefaultCA1CAF = frames.InitNormalCancelSlice(51, 104)
	chargeFramesDefaultCA1CAF[action.ActionCharge] = 60
	chargeFramesCA2CA1CAF = frames.InitNormalCancelSlice(51-28, 104-28) // previous action is CA2
	chargeFramesCA2CA1CAF[action.ActionCharge] = 60 - 28
	chargeFramesECA1CAF = frames.InitNormalCancelSlice(51-17, 104-17) // previous action is E
	chargeFramesECA1CAF[action.ActionCharge] = 60 - 17
	// we are doing a CA2, so the next CA has to be CAF
	chargeFramesDefaultCA2CAF = frames.InitNormalCancelSlice(24, 77)
	chargeFramesDefaultCA2CAF[action.ActionCharge] = 32

	// 3+ stacks
	// we are doing a CA1, so the next CA has to be CA2
	chargeFramesDefaultCA1CA2 = frames.InitNormalCancelSlice(51, 104)
	chargeFramesDefaultCA1CA2[action.ActionCharge] = 57
	chargeFramesCA2CA1CA2 = frames.InitNormalCancelSlice(51-28, 104-28) // previous action is CA2
	chargeFramesCA2CA1CA2[action.ActionCharge] = 57 - 28
	chargeFramesECA1CA2 = frames.InitNormalCancelSlice(51-17, 104-17) // previous action is E
	chargeFramesECA1CA2[action.ActionCharge] = 57 - 17
	// we are doing a CA2, so the next CA has to be CA1
	chargeFramesDefaultCA2CA1 = frames.InitNormalCancelSlice(24, 77)
	chargeFramesDefaultCA2CA1[action.ActionCharge] = 29

	chargeFrames = make([][]int, endState)
	chargeFrames[defaultToCA0] = chargeFramesDefaultCA0
	chargeFrames[n1CAFeToCA0] = chargeFramesN1CA0
	chargeFrames[n2n3ToCA0] = chargeFramesN2CA0

	chargeFrames[defaultToCAF] = chargeFramesDefaultCAF
	chargeFrames[CA1CA2ToCAF] = chargeFramesCACAF
	chargeFrames[eToCAF] = chargeFramesECAF

	chargeFrames[defaultToCA1ToCAF] = chargeFramesDefaultCA1CAF
	chargeFrames[ca2ToCA1ToCAF] = chargeFramesCA2CA1CAF
	chargeFrames[eToCA1ToCAF] = chargeFramesECA1CAF

	chargeFrames[defaultToCA2ToCAF] = chargeFramesDefaultCA2CAF

	chargeFrames[defaultToCA1ToCA2] = chargeFramesDefaultCA1CA2
	chargeFrames[CA2ToCA1ToCA2] = chargeFramesCA2CA1CA2
	chargeFrames[eToCA1ToCA2] = chargeFramesECA1CA2

	chargeFrames[defaultToCA2ToCA1] = chargeFramesDefaultCA2CA1
}

func (c *char) chargeState() IttoChargeState {
	lastWasItto := c.Core.Player.LastAction.Char == c.Index
	lastAction := c.Core.Player.LastAction.Type
	if c.Tags[c.stackKey] == 0 {
		return c.determineChargeForCA0(lastWasItto, lastAction)
	} else if c.Tags[c.stackKey] == 1 {
		return c.determineChargeForCAF(lastWasItto, lastAction)
	} else {
		if c.chargedCount == -1 || c.chargedCount == 2 || c.chargedCount == 3 {
			return c.determineChargeForCA1(lastWasItto, lastAction)
		} else {
			return c.determineChargeForCA2(lastWasItto, lastAction)
		}
	}
}

func (c *char) determineChargeForCA0(lastWasItto bool, lastAction action.Action) IttoChargeState {
	state := InvalidState
	if c.NormalCounter == 1 ||
		(lastWasItto && lastAction == action.ActionCharge && c.chargedCount == 3) ||
		(lastWasItto && lastAction == action.ActionSkill) {
		// CA0 is 14 frames shorter if prior action was N1/CAF/E
		state = n1CAFeToCA0
	} else if c.NormalCounter == 3 || c.NormalCounter == 4 {
		// CA0 is 21 frames shorter if prior action was N2/N3
		state = n2n3ToCA0
	} else {
		// default
		state = defaultToCA0
	}
	c.chargedCount = 0 // CA0 was used
	return state
}

func (c *char) determineChargeForCAF(lastWasItto bool, lastAction action.Action) IttoChargeState {
	state := InvalidState
	// CAF -> X
	// CAF is 25 frames shorter if CA1/CA2 -> CAF
	if (lastWasItto && lastAction == action.ActionCharge) && (c.chargedCount == 1 || c.chargedCount == 2) {
		state = CA1CA2ToCAF
	} else if lastAction == action.ActionSkill {
		// CAF is 17 frames shorter if E -> CAF
		state = eToCAF
	} else {
		// default
		state = defaultToCAF
	}
	c.chargedCount = 3 // CAF was used
	return state
}

func (c *char) determineChargeForCA1(lastWasItto bool, lastAction action.Action) IttoChargeState {
	state := InvalidState
	if c.Tags[c.stackKey] == 2 {
		// CA1 -> CAF
		if (lastWasItto && lastAction == action.ActionCharge) && c.chargedCount == 2 {
			// CA1 is 28 frames shorter if CA2 -> CA1
			state = ca2ToCA1ToCAF
		} else if lastAction == action.ActionSkill {
			// CA1 is 17 frames shorter if E -> CA1
			state = eToCA1ToCAF
		} else {
			// default
			state = defaultToCA1ToCAF
		}
	} else {
		// CA1 -> CA2
		if (lastWasItto && lastAction == action.ActionCharge) && c.chargedCount == 2 {
			// CA1 is 28 frames shorter if CA2 -> CA1
			state = CA2ToCA1ToCA2
		} else if lastAction == action.ActionSkill {
			// CA1 is 17 frames shorter if E -> CA1
			state = eToCA1ToCA2
		} else {
			// default
			state = defaultToCA1ToCA2
		}
	}
	c.chargedCount = 1 // CA1 was used
	return state
}

func (c *char) determineChargeForCA2(lastWasItto bool, lastAction action.Action) IttoChargeState {
	state := InvalidState
	// CA2 -> X
	if c.Tags[c.stackKey] == 2 {
		// CA2 -> CAF
		state = defaultToCA2ToCAF
	} else {
		// CA2 -> CA1
		state = defaultToCA2ToCA1
	}
	c.chargedCount = 2 // CA2 was used
	return state
}

func (c *char) ChargeAttack(p map[string]int) action.ActionInfo {
	state := c.chargeState()

	// check burst status for radius
	// TODO: proper hitbox
	r := 1.0
	if c.StatModIsActive(c.burstBuffKey) {
		r = 3
	}

	// handle text to show in debug
	text := ""
	switch c.chargedCount {
	case 0:
		text = "Saichimonji Slash"
	case 1:
		text = "Arataki Kesagiri Combo Slash Left"
	case 2:
		text = "Arataki Kesagiri Combo Slash Right"
	case 3:
		text = "Arataki Kesagiri Final Slash"
	}

	// Attack
	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               fmt.Sprintf("%v, Stacks %v", text, c.Tags[c.stackKey]),
		Mult:               akCombo[c.TalentLvlAttack()],
		AttackTag:          combat.AttackTagExtra,
		ICDTag:             combat.ICDTagNormalAttack,
		ICDGroup:           combat.ICDGroupDefault,
		StrikeType:         combat.StrikeTypeBlunt,
		Element:            attributes.Physical,
		Durability:         25,
		HitlagFactor:       0.01,
		CanBeDefenseHalted: true,
	}

	// handle CA multiplier, CA hitlag and A4
	if c.Tags[c.stackKey] == 0 {
		// Saichimonji Slash
		ai.Mult = saichiSlash[c.TalentLvlAttack()]
		ai.FlatDmg = 0
		ai.HitlagHaltFrames = 0.10 * 60
	} else if c.Tags[c.stackKey] == 1 {
		// Arataki Kesagiri Final Slash
		ai.Mult = akFinal[c.TalentLvlAttack()]
		ai.HitlagHaltFrames = 0.10 * 60
		// apply A4
		c.a4(&ai)
	} else {
		// Arataki Kesagiri Combo Slash
		if c.Tags[c.stackKey] == 2 {
			ai.HitlagHaltFrames = 0.07 * 60
		} else if c.Tags[c.stackKey] == 3 {
			ai.HitlagHaltFrames = 0.05 * 60
		} else {
			ai.HitlagHaltFrames = 0.03 * 60
		}
		// apply A4
		c.a4(&ai)
	}

	// TODO: Does Itto CA snapshot at the start of CA?
	c.Core.QueueAttack(ai, combat.NewDefCircHit(r, false, combat.TargettableEnemy), 0, chargeFrames[state][action.ActionDash])

	// handle A1
	c.a1Update()

	// handle C6
	if c.Base.Cons >= 6 {
		c.c6StackHandler()
	} else {
		c.changeStacks(-1)
	}

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(chargeFrames[state]),
		AnimationLength: chargeFrames[state][action.InvalidAction],
		CanQueueAfter:   chargeFrames[state][action.ActionDash],
		State:           action.ChargeAttackState,
	}
}

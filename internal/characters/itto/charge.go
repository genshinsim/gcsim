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

func init() {
	// 0 stacks
	chargeFramesDefaultCA0 = frames.InitNormalCancelSlice(89, 131)
	chargeFramesN1CA0 = frames.InitNormalCancelSlice(89-14, 131-14)
	chargeFramesN2CA0 = frames.InitNormalCancelSlice(89-21, 131-21)

	// 1 stack
	chargeFramesDefaultCAF = frames.InitNormalCancelSlice(71, 110)
	chargeFramesDefaultCAF[action.ActionSkill] = 76
	chargeFramesDefaultCAF[action.ActionBurst] = 76
	chargeFramesDefaultCAF[action.ActionSwap] = 76

	chargeFramesCACAF = frames.InitNormalCancelSlice(71-25, 110-25)
	chargeFramesCACAF[action.ActionSkill] = 76 - 25
	chargeFramesCACAF[action.ActionBurst] = 76 - 25
	chargeFramesCACAF[action.ActionSwap] = 76 - 25

	chargeFramesECAF = frames.InitNormalCancelSlice(71-17, 110-17)
	chargeFramesECAF[action.ActionSkill] = 76 - 17
	chargeFramesECAF[action.ActionBurst] = 76 - 17
	chargeFramesECAF[action.ActionSwap] = 76 - 17

	// 2 stacks
	chargeFramesDefaultCA1CAF = frames.InitNormalCancelSlice(51, 104)
	chargeFramesDefaultCA1CAF[action.ActionCharge] = 60
	chargeFramesCA2CA1CAF = frames.InitNormalCancelSlice(51-28, 104-28)
	chargeFramesCA2CA1CAF[action.ActionCharge] = 60 - 28
	chargeFramesECA1CAF = frames.InitNormalCancelSlice(51-17, 104-17)
	chargeFramesECA1CAF[action.ActionCharge] = 60 - 17

	chargeFramesDefaultCA2CAF = frames.InitNormalCancelSlice(24, 77)
	chargeFramesDefaultCA2CAF[action.ActionCharge] = 32

	// 3+ stacks
	chargeFramesDefaultCA1CA2 = frames.InitNormalCancelSlice(51, 104)
	chargeFramesDefaultCA1CA2[action.ActionCharge] = 57
	chargeFramesCA2CA1CA2 = frames.InitNormalCancelSlice(51-28, 104-28)
	chargeFramesCA2CA1CA2[action.ActionCharge] = 57 - 28
	chargeFramesECA1CA2 = frames.InitNormalCancelSlice(51-17, 104-17)
	chargeFramesECA1CA2[action.ActionCharge] = 57 - 17

	chargeFramesDefaultCA2CA1 = frames.InitNormalCancelSlice(24, 77)
	chargeFramesDefaultCA2CA1[action.ActionCharge] = 29
}

func (c *char) ChargeAttack(p map[string]int) action.ActionInfo {
	lastWasItto := c.Core.Player.LastAction.Char == c.Index
	lastAction := c.Core.Player.LastAction.Type

	chargeFrames := make([]int, action.EndActionType)
	if c.Tags[c.stackKey] == 0 {
		if c.NormalCounter == 1 ||
			(lastWasItto && lastAction == action.ActionCharge && c.chargedCount == 3) ||
			(lastWasItto && lastAction == action.ActionSkill) {
			// CA0 is 14 frames shorter if prior action was N1/CAF/E
			copy(chargeFrames, chargeFramesN1CA0)
		} else if c.NormalCounter == 3 || c.NormalCounter == 4 {
			// CA0 is 21 frames shorter if prior action was N2/N3
			copy(chargeFrames, chargeFramesN2CA0)
		} else {
			// default
			copy(chargeFrames, chargeFramesDefaultCA0)
		}
		c.chargedCount = 0 // CA0 was used
	} else if c.Tags[c.stackKey] == 1 {
		// CAF -> X
		// CAF is 25 frames shorter if CA1/CA2 -> CAF
		if (lastWasItto && lastAction == action.ActionCharge) && (c.chargedCount == 1 || c.chargedCount == 2) {
			copy(chargeFrames, chargeFramesCACAF)
		} else if lastAction == action.ActionSkill {
			// CAF is 17 frames shorter if E -> CAF
			copy(chargeFrames, chargeFramesECAF)
		} else {
			// default
			copy(chargeFrames, chargeFramesDefaultCAF)
		}
		c.chargedCount = 3 // CAF was used
	} else {
		// CA1/CA2 -> X
		if c.chargedCount == 0 || c.chargedCount == 2 || c.chargedCount == 3 {
			if c.Tags[c.stackKey] == 2 {
				// CA1 -> CAF
				if (lastWasItto && lastAction == action.ActionCharge) && c.chargedCount == 2 {
					// CA1 is 28 frames shorter if CA2 -> CA1
					copy(chargeFrames, chargeFramesCA2CA1CAF)
				} else if lastAction == action.ActionSkill {
					// CA1 is 17 frames shorter if E -> CA1
					copy(chargeFrames, chargeFramesECA1CAF)
				} else {
					// default
					copy(chargeFrames, chargeFramesDefaultCA1CAF)
				}
			} else {
				// CA1 -> CA2
				if (lastWasItto && lastAction == action.ActionCharge) && c.chargedCount == 2 {
					// CA1 is 28 frames shorter if CA2 -> CA1
					copy(chargeFrames, chargeFramesCA2CA1CA2)
				} else if lastAction == action.ActionSkill {
					// CA1 is 17 frames shorter if E -> CA1
					copy(chargeFrames, chargeFramesECA1CA2)
				} else {
					// default
					copy(chargeFrames, chargeFramesDefaultCA1CA2)
				}
			}
			c.chargedCount = 1 // CA1 was used
		} else {
			// CA2 -> X
			if c.Tags[c.stackKey] == 2 {
				// CA2 -> CAF
				copy(chargeFrames, chargeFramesDefaultCA2CAF)
			} else {
				// CA2 -> CA1
				copy(chargeFrames, chargeFramesDefaultCA2CA1)
			}
			c.chargedCount = 2 // CA2 was used
		}
	}

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
	c.Core.QueueAttack(ai, combat.NewDefCircHit(r, false, combat.TargettableEnemy), 0, chargeFrames[action.ActionDash])

	// handle A1
	c.a1Update()

	// handle C6
	if c.Base.Cons >= 6 {
		c.c6StackHandler()
	} else {
		c.changeStacks(-1)
	}

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(chargeFrames),
		AnimationLength: chargeFrames[action.InvalidAction],
		CanQueueAfter:   chargeFrames[action.ActionDash],
		State:           action.ChargeAttackState,
	}
}

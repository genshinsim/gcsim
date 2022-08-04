package itto

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var chargeFrames []int
var chargeLongestCancel int
var chargeHitmark int

func init() {
}

func (c *char) ChargeAttack(p map[string]int) action.ActionInfo {
	lastWasItto := c.Core.Player.LastAction.Char == c.Index
	lastAction := c.Core.Player.LastAction.Type

	if c.Tags[c.stackKey] == 0 {
		// CA0 -> X

		// CA0 from idle or after burst (c.NormalCounter == 0)
		chargeLongestCancel = 131
		chargeHitmark = 89

		// CA0 is 14 frames shorter if prior action was N1/CAF/E
		if c.NormalCounter == 1 ||
			(lastWasItto && lastAction == action.ActionCharge && c.chargedCount == 3) ||
			(lastWasItto && lastAction == action.ActionSkill) {
			chargeLongestCancel -= 14
			chargeHitmark -= 14
		}

		// CA0 is 21 frames shorter if prior action was N2/N3
		if c.NormalCounter == 3 || c.NormalCounter == 4 {
			chargeLongestCancel -= 21
			chargeHitmark -= 21
		}

		chargeFrames = frames.InitAbilSlice(chargeLongestCancel) // CA0 -> N1/E/Q
		chargeFrames[action.ActionCharge] = 500                  // CA0 -> CA1/CAF, TODO: this action is illegal; need better way to handle it
		chargeFrames[action.ActionDash] = chargeHitmark          // CA0 -> D
		chargeFrames[action.ActionJump] = chargeHitmark          // CA0 -> J

		c.chargedCount = 0 // CA0 was used
	} else if c.Tags[c.stackKey] == 1 {
		// CAF -> X
		chargeLongestCancel = 110
		chargeHitmark = 71

		// CAF is 25 frames shorter if CA1/CA2 -> CAF
		if (lastWasItto && lastAction == action.ActionCharge) && (c.chargedCount == 1 || c.chargedCount == 2) {
			chargeLongestCancel -= 25
			chargeHitmark -= 25
		}
		// CAF is 17 frames shorter if E -> CAF
		if lastAction == action.ActionSkill {
			chargeLongestCancel -= 17
			chargeHitmark -= 17
		}

		chargeFrames = frames.InitAbilSlice(chargeLongestCancel)    // CAF -> N1/CA0 (potentially CA1 as well?)
		chargeFrames[action.ActionSkill] = chargeLongestCancel - 34 // CAF -> E (listed 76 = 110 - 34)
		chargeFrames[action.ActionBurst] = chargeLongestCancel - 34 // CAF -> Q (listed 76 = 110 - 34)
		chargeFrames[action.ActionDash] = chargeHitmark             // CAF -> D
		chargeFrames[action.ActionJump] = chargeHitmark             // CAF -> J
		chargeFrames[action.ActionSwap] = chargeLongestCancel - 34  // CAF -> Swap (listed 76 = 110 - 34)

		c.chargedCount = 3 // CAF was used
	} else {
		// CA1/CA2 -> X
		if c.chargedCount == 0 || c.chargedCount == 2 || c.chargedCount == 3 {
			// CA1 -> X
			chargeLongestCancel = 104
			chargeHitmark = 51

			// CA1 is 28 frames shorter if CA2 -> CA1
			if (lastWasItto && lastAction == action.ActionCharge) && c.chargedCount == 2 {
				chargeLongestCancel -= 28
				chargeHitmark -= 28
			}
			// CA1 is 17 frames shorter if E -> CA1
			if lastAction == action.ActionSkill {
				chargeLongestCancel -= 17
				chargeHitmark -= 17
			}

			chargeFrames = frames.InitNormalCancelSlice(chargeHitmark, chargeLongestCancel)
			if c.Tags[c.stackKey] == 2 {
				chargeFrames[action.ActionCharge] = chargeLongestCancel - 44 // CA1 -> CAF (listed 60 = 104 - 44)
			} else {
				// next CA has to be CA1
				chargeFrames[action.ActionCharge] = chargeLongestCancel - 47 // CA1 -> CA2 (listed 57 = 104 - 47)
			}
			c.chargedCount = 1 // CA1 was used
		} else {
			// CA2 -> X
			chargeLongestCancel = 77
			chargeHitmark = 24
			chargeFrames = frames.InitNormalCancelSlice(chargeHitmark, chargeLongestCancel)
			if c.Tags[c.stackKey] == 2 {
				// next CA has to be CAF
				chargeFrames[action.ActionCharge] = 32 // CA2 -> CAF
			} else {
				// next CA has to be CA1
				chargeFrames[action.ActionCharge] = 29 // CA2 -> CA1
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
	c.Core.QueueAttack(ai, combat.NewDefCircHit(r, false, combat.TargettableEnemy), 0, chargeHitmark)

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
		CanQueueAfter:   chargeHitmark,
		State:           action.ChargeAttackState,
	}
}

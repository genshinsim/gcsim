package raiden

import (
	"github.com/genshinsim/gcsim/pkg/core"
)

func (c *char) ActionFrames(a core.ActionType, p map[string]int) (int, int) {
	switch a {
	case core.ActionAttack:
		f := 0
		a := 0
		if c.Core.Status.Duration("raidenburst") == 0 {
			switch c.NormalCounter {
			//TODO: need to add atkspd mod
			case 0:
				f = 14
				a = 18
			case 1:
				f = 9
				a = 13
			case 2:
				f = 14
				a = 26
			case 3:
				f = 27
				a = 41
			case 4:
				f = 34
				a = 50
			}
		} else {
			switch c.NormalCounter {
			//TODO: need to add atkspd mod
			case 0:
				f = 12
				a = 19
			case 1:
				f = 13
				a = 16
			case 2:
				f = 11
				a = 16
			case 3:
				f = 33
				a = 44
			case 4:
				f = 33
				a = 59
			}
		}
		f = int(float64(f) / (1 + c.Stat(core.AtkSpd)))
		return f, a
	case core.ActionCharge:
		if c.Core.Status.Duration("raidenburst") == 0 {
			return 22, 37
		}
		return 24, 56
	case core.ActionSkill:
		return 17, 37
	case core.ActionBurst:
		return 98, 111
	default:
		c.Core.Log.NewEventBuildMsg(core.LogActionEvent, c.Index, "unknown action (invalid frames): ", a.String())
		return 0, 0
	}
}

func (c *char) InitCancelFrames() {

	//normal cancels
	c.SetNormalCancelFrames(0, core.ActionAttack, 18-14) //n1 -> next attack
	c.SetNormalCancelFrames(0, core.ActionCharge, 24-14) //n1 -> charge

	c.SetNormalCancelFrames(1, core.ActionAttack, 13-9) //n2 -> next attack
	c.SetNormalCancelFrames(1, core.ActionCharge, 26-9) //n2 -> charge

	c.SetNormalCancelFrames(2, core.ActionAttack, 26-14) //n3 -> next attack
	c.SetNormalCancelFrames(2, core.ActionCharge, 36-14) //n3 -> charge

	c.SetNormalCancelFrames(3, core.ActionAttack, 41-27) //n4 -> next attack
	c.SetNormalCancelFrames(3, core.ActionCharge, 57-27) //n4 -> charge

	c.SetNormalCancelFrames(4, core.ActionAttack, 50-34) //n5 -> n1

	c.SetAbilCancelFrames(core.ActionCharge, core.ActionAttack, 37-22) //charge -> n1
	c.SetAbilCancelFrames(core.ActionCharge, core.ActionSkill, 37-22)  //charge -> skill
	c.SetAbilCancelFrames(core.ActionCharge, core.ActionBurst, 37-22)  //charge -> burst
	c.SetAbilCancelFrames(core.ActionCharge, core.ActionSwap, 36-22)   //charge -> swap

	c.SetAbilCancelFrames(core.ActionBurst, core.ActionAttack, 111-98) //burst -> n1
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionSkill, 111-98)  //burst -> skill
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionDash, 110-98)   //burst -> dash
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionJump, 112-98)   //burst -> jump
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionSwap, 110-98)   //burst -> swap

	c.SetAbilCancelFrames(core.ActionSkill, core.ActionAttack, 37-17)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionBurst, 37-17)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionDash, 17-17)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionJump, 17-17)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionSwap, 36-17)
}

func (c *char) ActionInterruptableDelay(next core.ActionType, p map[string]int) int {
	// Provide a custom override for Raiden's Q Normals and Charges
	if (c.Core.LastAction.Typ == core.ActionAttack ||
		c.Core.LastAction.Typ == core.ActionCharge) &&
		c.Core.Status.Duration("raidenburst") > 0 {
		f := 0
		switch c.Core.LastAction.Typ {
		case core.ActionAttack:
			f = burstNormalCancels(next, c.NormalCounter)
		case core.ActionCharge:
			f = burstChargeCancels(next)
		}
		f = int(float64(f) / (1 + c.Stat(core.AtkSpd)))
		return f
	}
	//otherwise use default implementation
	return c.Tmpl.ActionInterruptableDelay(next, p)
}

func burstNormalCancels(next core.ActionType, NormalCounter int) int {
	switch next {
	case core.ActionAttack:
		switch NormalCounter { //counter for the next attack
		case 1:
			return 19 - 12
		case 2:
			return 16 - 13
		case 3:
			return 16 - 11
		case 4:
			return 44 - 33
		case 0:
			return 59 - 33

		}
	case core.ActionCharge:
		switch NormalCounter {
		case 1:
			return 24 - 12
		case 2:
			return 26 - 13
		case 3:
			return 34 - 11
		case 4:
			return 67 - 33
		case 0:
			return 83 - 33
		}
	}
	//all other actions can cancel immediately
	return 0
}

func burstChargeCancels(next core.ActionType) int {
	switch next {
	case core.ActionAttack:
		return 56 - 24
	case core.ActionSkill:
		return 56 - 24
	case core.ActionDash:
		return 35 - 24
	case core.ActionJump:
		return 35 - 24
	case core.ActionSwap:
		return 55 - 24
	default:
		return 0
	}
}

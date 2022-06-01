package xiao

import "github.com/genshinsim/gcsim/pkg/core"

func (c *char) ActionFrames(a core.ActionType, p map[string]int) (int, int) {
	switch a {
	case core.ActionAttack:
		f := 0
		a := 0
		switch c.NormalCounter {
		//TODO: need to add atkspd mod
		case 0:
			f = 17
			a = 25
		case 1:
			f = 15
			a = 22
		case 2:
			f = 15
			a = 26
		case 3:
			f = 31
			a = 39
		case 4:
			f = 16
			a = 24
		case 5:
			f = 39
			a = 79
		}
		atkspd := c.Stat(core.AtkSpd)
		f = int(float64(f) / (1 + atkspd))
		return f, a
	case core.ActionCharge:
		return 17, 46
	case core.ActionHighPlunge:
		c.Core.Log.NewEvent("plunge jump check", core.LogCharacterEvent, c.Index, "previous", c.Core.LastAction)
		if c.Core.LastAction.Target == core.Xiao && c.Core.LastAction.Typ == core.ActionJump {
			return 46, 61
		}
		c.Core.Log.NewEvent("invalid plunge (missing jump)", core.LogActionEvent, c.Index, "action", a)
		return 0, 0
	case core.ActionLowPlunge:
		c.Core.Log.NewEvent("plunge jump check", core.LogCharacterEvent, c.Index, "previous", c.Core.LastAction)
		if c.Core.LastAction.Target == core.Xiao && c.Core.LastAction.Typ == core.ActionJump {
			return 44, 60
		}
		c.Core.Log.NewEvent("invalid plunge (missing jump)", core.LogActionEvent, c.Index, "action", a)
		return 0, 0
	case core.ActionSkill:
		return 4, 24
	case core.ActionBurst:
		return 57, 82
	case core.ActionDash:
		return 21, 21
	case core.ActionJump:
		if c.Core.Status.Duration("xiaoburst") > 0 {
			return 5, 58
		}
		return 32, 32
	default:
		c.Core.Log.NewEventBuildMsg(core.LogActionEvent, c.Index, "unknown action (invalid frames): ", a.String())
		return 0, 0
	}
}

func (c *char) InitCancelFrames() {

	//normal cancels
	c.SetNormalCancelFrames(0, core.ActionAttack, 25-17) //n1 -> next attack
	c.SetNormalCancelFrames(0, core.ActionCharge, 26-17) //n1 -> charge

	c.SetNormalCancelFrames(1, core.ActionAttack, 22-15) //n2 -> next attack
	c.SetNormalCancelFrames(1, core.ActionCharge, 27-15) //n2 -> charge

	c.SetNormalCancelFrames(2, core.ActionAttack, 26-15) //n3 -> next attack
	c.SetNormalCancelFrames(2, core.ActionCharge, 38-15) //n3 -> charge

	c.SetNormalCancelFrames(3, core.ActionAttack, 39-31) //n4 -> next attack
	c.SetNormalCancelFrames(3, core.ActionCharge, 42-31) //n4 -> charge

	c.SetNormalCancelFrames(4, core.ActionAttack, 24-16) //n5 -> next attack
	c.SetNormalCancelFrames(4, core.ActionCharge, 30-16) //n5 -> charge

	c.SetNormalCancelFrames(5, core.ActionAttack, 79-39) //n5 -> n6

	c.SetAbilCancelFrames(core.ActionCharge, core.ActionAttack, 45-17) //charge -> n1
	c.SetAbilCancelFrames(core.ActionCharge, core.ActionSkill, 38-17)  //charge -> skill
	c.SetAbilCancelFrames(core.ActionCharge, core.ActionBurst, 37-17)  //charge -> burst
	c.SetAbilCancelFrames(core.ActionCharge, core.ActionSwap, 43-17)   //charge -> swap

	c.SetAbilCancelFrames(core.ActionBurst, core.ActionAttack, 82-57) //burst -> n1
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionSkill, 82-57)  //burst -> skill
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionDash, 57-57)   //burst -> dash
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionJump, 58-57)   //burst -> jump
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionSwap, 67-57)   //burst -> swap

	c.SetAbilCancelFrames(core.ActionSkill, core.ActionAttack, 24-4)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionSkill, 24-4)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionBurst, 24-4)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionDash, 35-4)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionJump, 37-4)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionSwap, 35-4)

	c.SetAbilCancelFrames(core.ActionHighPlunge, core.ActionAttack, 61-46)
	c.SetAbilCancelFrames(core.ActionHighPlunge, core.ActionSkill, 66-46)
	c.SetAbilCancelFrames(core.ActionHighPlunge, core.ActionDash, 66-46)
	c.SetAbilCancelFrames(core.ActionHighPlunge, core.ActionJump, 65-46)
	c.SetAbilCancelFrames(core.ActionHighPlunge, core.ActionSwap, 64-46)

	c.SetAbilCancelFrames(core.ActionLowPlunge, core.ActionAttack, 60-44)
	c.SetAbilCancelFrames(core.ActionLowPlunge, core.ActionSkill, 59-44)
	c.SetAbilCancelFrames(core.ActionLowPlunge, core.ActionDash, 60-44)
	c.SetAbilCancelFrames(core.ActionLowPlunge, core.ActionJump, 61-44)
	c.SetAbilCancelFrames(core.ActionLowPlunge, core.ActionSwap, 62-44)
}

func (c *char) ActionInterruptableDelay(next core.ActionType, p map[string]int) int {
	// Provide a custom override for Xiao's Jump
	if c.Core.LastAction.Typ == core.ActionJump &&
		c.Core.Status.Duration("xiaoburst") > 0 {
		return BurstJumpFrames(next)
	}
	//otherwise use default implementation
	return c.Tmpl.ActionInterruptableDelay(next, p)
}

func BurstJumpFrames(next core.ActionType) int {
	switch next {
	case core.ActionHighPlunge:
		return 6 - 5
	case core.ActionLowPlunge:
		return 5 - 5
	default:
		return 58 - 5 //if jump was not followed by plunge, he must wait to float down
	}
}

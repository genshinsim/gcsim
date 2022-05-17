package sucrose

import "github.com/genshinsim/gcsim/pkg/core"

func (c *char) ActionFrames(a core.ActionType, p map[string]int) (int, int) {
	switch a {
	case core.ActionAttack:
		f := 0
		a := 0
		switch c.NormalCounter {
		//TODO: need to add atkspd mod
		case 0:
			f = 2
			a = 20
		case 1:
			f = 4
			a = 26
		case 2:
			f = 16
			a = 33
		case 3:
			f = 28
			a = 51
		}
		f = int(float64(f) / (1 + c.Stat(core.AtkSpd)))

		return f, a
	case core.ActionCharge:
		return 53, 69
	case core.ActionSkill:
		return 11, 57
	case core.ActionBurst:
		return 46, 49
	case core.ActionDash:
		return 1, 24
	default:
		c.Core.Log.NewEventBuildMsg(core.LogActionEvent, c.Index, "unknown action (invalid frames): ", a.String())
		return 0, 0
	}
}

func (c *char) InitCancelFrames() {

	//normal cancels
	c.SetNormalCancelFrames(0, core.ActionAttack, 20-2) //n1 -> next attack
	c.SetNormalCancelFrames(0, core.ActionCharge, 2-2)  //n1 -> charge
	c.SetNormalCancelFrames(0, core.ActionSkill, 17-2)  //only CA can be done before hitmark
	c.SetNormalCancelFrames(0, core.ActionBurst, 17-2)
	c.SetNormalCancelFrames(0, core.ActionDash, 17-2)
	c.SetNormalCancelFrames(0, core.ActionJump, 17-2)
	c.SetNormalCancelFrames(0, core.ActionSwap, 17-2)

	c.SetNormalCancelFrames(1, core.ActionAttack, 26-4) //n2 -> next attack
	c.SetNormalCancelFrames(1, core.ActionCharge, 4-4)  //n2 -> charge
	c.SetNormalCancelFrames(1, core.ActionSkill, 18-4)  //only CA can be done before hitmark
	c.SetNormalCancelFrames(1, core.ActionBurst, 18-4)
	c.SetNormalCancelFrames(1, core.ActionDash, 18-4)
	c.SetNormalCancelFrames(1, core.ActionJump, 18-4)
	c.SetNormalCancelFrames(1, core.ActionSwap, 18-4)

	c.SetNormalCancelFrames(2, core.ActionAttack, 33-16) //n3 -> next attack
	c.SetNormalCancelFrames(2, core.ActionCharge, 16-16) //n3 -> charge
	c.SetNormalCancelFrames(2, core.ActionSkill, 28-16)  //only CA can be done before hitmark
	c.SetNormalCancelFrames(2, core.ActionBurst, 28-16)
	c.SetNormalCancelFrames(2, core.ActionDash, 28-16)
	c.SetNormalCancelFrames(2, core.ActionJump, 28-16)
	c.SetNormalCancelFrames(2, core.ActionSwap, 28-16)

	c.SetNormalCancelFrames(3, core.ActionAttack, 51-28) //n4 -> n1
	c.SetNormalCancelFrames(3, core.ActionCharge, 42-28) //n4 -> charge
	//CA is after hitmark this time, so no need for the rest

	c.SetAbilCancelFrames(core.ActionCharge, core.ActionAttack, 69-53) //charge -> n1
	c.SetAbilCancelFrames(core.ActionCharge, core.ActionCharge, 66-53) //charge -> charge
	c.SetAbilCancelFrames(core.ActionCharge, core.ActionSkill, 60-53)  //charge -> skill
	c.SetAbilCancelFrames(core.ActionCharge, core.ActionBurst, 61-53)  //charge -> burst
	c.SetAbilCancelFrames(core.ActionCharge, core.ActionSwap, 54-53)   //charge -> swap

	c.SetAbilCancelFrames(core.ActionBurst, core.ActionAttack, 49-46) //burst -> n1
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionCharge, 48-46) //burst -> charge
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionSkill, 48-46)  //burst -> skill
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionDash, 47-46)   //burst -> dash
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionJump, 47-46)   //burst -> jump
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionSwap, 47-46)   //burst -> swap

	c.SetAbilCancelFrames(core.ActionSkill, core.ActionAttack, 57-11)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionCharge, 56-11)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionSkill, 56-11)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionBurst, 57-11)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionDash, 11-11)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionJump, 11-11)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionSwap, 56-11)

	c.SetAbilCancelFrames(core.ActionDash, core.ActionAttack, 24-1)
	c.SetAbilCancelFrames(core.ActionDash, core.ActionCharge, 24-1)
	//sucrose cancel her dash with her E and Q
	c.SetAbilCancelFrames(core.ActionDash, core.ActionSkill, 1-1)
	c.SetAbilCancelFrames(core.ActionDash, core.ActionBurst, 1-1)
	c.SetAbilCancelFrames(core.ActionDash, core.ActionDash, 24-1)
	c.SetAbilCancelFrames(core.ActionDash, core.ActionJump, 24-1)
	c.SetAbilCancelFrames(core.ActionDash, core.ActionSwap, 24-1)
}

func (c *char) ActionInterruptableDelay(next core.ActionType, p map[string]int) int {
	return c.Tmpl.ActionInterruptableDelay(next, p)
}

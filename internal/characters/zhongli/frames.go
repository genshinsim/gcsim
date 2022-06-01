package zhongli

import "github.com/genshinsim/gcsim/pkg/core"

func (c *char) ActionFrames(a core.ActionType, p map[string]int) (int, int) {
	switch a {
	case core.ActionAttack:
		f := 0
		a := 0
		switch c.NormalCounter {
		//TODO: need to add atkspd mod
		case 0:
			f = 11
			a = 18
		case 1:
			f = 9
			a = 13
		case 2:
			f = 8
			a = 19
		case 3:
			f = 16
			a = 34
		case 4:
			f = 27
			a = 27
		case 5:
			f = 29
			a = 54
		}
		f = int(float64(f) / (1 + c.Stat(core.AtkSpd)))
		return f, a
	case core.ActionCharge:
		return 4, 47
	case core.ActionSkill:
		hold := p["hold"]
		holdNoStele := p["hold_nostele"]
		if hold == 0 && holdNoStele == 0 {
			//no hold
			return 23, 37
		}
		//yes hold
		return 48, 96
	case core.ActionBurst:
		return 101, 138
	case core.ActionDash:
		return 19, 19
	default:
		c.Core.Log.NewEventBuildMsg(core.LogActionEvent, c.Index, "unknown action (invalid frames): ", a.String())
		return 0, 0
	}
}

func (c *char) InitCancelFrames() {

	//normal cancels
	c.SetNormalCancelFrames(0, core.ActionAttack, 18-11) //n1 -> next attack
	c.SetNormalCancelFrames(0, core.ActionCharge, 30-11) //n1 -> charge

	c.SetNormalCancelFrames(1, core.ActionAttack, 13-9) //n2 -> next attack
	c.SetNormalCancelFrames(1, core.ActionCharge, 30-9) //n2 -> charge

	c.SetNormalCancelFrames(2, core.ActionAttack, 19-8) //n3 -> next attack
	c.SetNormalCancelFrames(2, core.ActionCharge, 28-8) //n3 -> charge

	c.SetNormalCancelFrames(3, core.ActionAttack, 34-16) //n4 -> next attack
	c.SetNormalCancelFrames(3, core.ActionCharge, 33-16) //n4 -> charge

	//n5 -> n6 is before hitmark, so we delay everything that must wait until hitmark
	c.SetNormalCancelFrames(4, core.ActionCharge, 31-27) //n5 -> charge
	c.SetNormalCancelFrames(4, core.ActionSkill, 29-27)  //n5 -> skill
	c.SetNormalCancelFrames(4, core.ActionBurst, 29-27)  //n5 -> burst
	c.SetNormalCancelFrames(4, core.ActionDash, 29-27)   //n5 -> dash
	c.SetNormalCancelFrames(4, core.ActionJump, 29-27)   //n5 -> jump
	c.SetNormalCancelFrames(4, core.ActionSwap, 29-27)   //n5 -> swap

	c.SetNormalCancelFrames(5, core.ActionAttack, 54-29) //n6 -> n1

	c.SetAbilCancelFrames(core.ActionCharge, core.ActionAttack, 47-4) //charge -> n1
	c.SetAbilCancelFrames(core.ActionCharge, core.ActionSkill, 33-4)  //charge -> skill
	c.SetAbilCancelFrames(core.ActionCharge, core.ActionBurst, 33-4)  //charge -> burst
	c.SetAbilCancelFrames(core.ActionCharge, core.ActionSwap, 31-4)   //charge -> swap

	c.SetAbilCancelFrames(core.ActionBurst, core.ActionAttack, 138-101) //burst -> n1
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionSkill, 139-101)  //burst -> skill
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionDash, 122-101)   //burst -> dash
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionJump, 122-101)   //burst -> jump
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionSwap, 138-101)   //burst -> swap

	c.SetAbilCancelFrames(core.ActionSkill, core.ActionAttack, 37-23)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionBurst, 38-23)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionDash, 23-23)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionJump, 23-23)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionSwap, 37-23)
}

func (c *char) ActionInterruptableDelay(next core.ActionType, p map[string]int) int {
	// Provide a custom override for Zhongli's Hold E
	if c.Core.LastAction.Typ == core.ActionSkill {
		h := c.Core.LastAction.Param["hold"]
		hNS := c.Core.LastAction.Param["hold_nostele"]
		if h > 0 || hNS > 0 {
			return SkillHoldFrames(next)
		}
	}
	//otherwise use default implementation
	return c.Tmpl.ActionInterruptableDelay(next, p)
}

func SkillHoldFrames(next core.ActionType) int {
	switch next {
	case core.ActionAttack:
		return 96 - 48
	case core.ActionBurst:
		return 96 - 48
	case core.ActionDash:
		return 55 - 48
	case core.ActionJump:
		return 55 - 48
	case core.ActionSwap:
		return 96 - 48
	default:
		return 96 - 48
	}
}

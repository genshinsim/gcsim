package lisa

import "github.com/genshinsim/gcsim/pkg/core"

func (c *char) ActionFrames(a core.ActionType, p map[string]int) (int, int) {
	switch a {
	case core.ActionAttack:
		f := 0
		a := 0
		switch c.NormalCounter {
		//TODO: need to add atkspd mod
		case 0:
			f = 15
			a = 30
		case 1:
			f = 12
			a = 20
		case 2:
			f = 17
			a = 34
		case 3:
			f = 31
			a = 57
		}
		f = int(float64(f) / (1 + c.Stat(core.AtkSpd)))
		return f, a
	case core.ActionCharge:
		return 70, 91
	case core.ActionSkill:
		hold := p["hold"]
		if hold == 0 {
			return 20, 38 //no hold
		}
		//yes hold
		return 114, 141
	case core.ActionBurst:
		return 56, 86
	case core.ActionDash:
		return 22, 22
	case core.ActionJump:
		return 33, 33
	default:
		c.Core.Log.NewEventBuildMsg(core.LogActionEvent, c.Index, "unknown action (invalid frames): ", a.String())
		return 0, 0
	}
}

func (c *char) InitCancelFrames() {

	//normal cancels
	c.SetNormalCancelFrames(0, core.ActionAttack, 30-15) //n1 -> next attack
	c.SetNormalCancelFrames(0, core.ActionCharge, 15-15) //n1 -> charge
	c.SetNormalCancelFrames(0, core.ActionSkill, 26-15)  //only CA can be done before hitmark
	c.SetNormalCancelFrames(0, core.ActionBurst, 26-15)
	c.SetNormalCancelFrames(0, core.ActionDash, 26-15)
	c.SetNormalCancelFrames(0, core.ActionJump, 26-15)
	c.SetNormalCancelFrames(0, core.ActionSwap, 26-15)

	c.SetNormalCancelFrames(1, core.ActionAttack, 20-12) //n2 -> next attack
	c.SetNormalCancelFrames(1, core.ActionCharge, 12-12) //n2 -> charge
	c.SetNormalCancelFrames(1, core.ActionSkill, 18-12)  //only CA can be done before hitmark
	c.SetNormalCancelFrames(1, core.ActionBurst, 18-12)
	c.SetNormalCancelFrames(1, core.ActionDash, 18-12)
	c.SetNormalCancelFrames(1, core.ActionJump, 18-12)
	c.SetNormalCancelFrames(1, core.ActionSwap, 18-12)

	c.SetNormalCancelFrames(2, core.ActionAttack, 34-17) //n3 -> n4
	c.SetNormalCancelFrames(2, core.ActionCharge, 26-17) //n3 -> charge
	//CA is after hitmark this time, so no need for the rest

	c.SetNormalCancelFrames(3, core.ActionAttack, 57-31) //n4 -> n1
	//Missing N4->CA - it's extremely long though.

	c.SetAbilCancelFrames(core.ActionCharge, core.ActionAttack, 86-70) //charge -> n1
	c.SetAbilCancelFrames(core.ActionCharge, core.ActionCharge, 90-70) //charge -> charge
	c.SetAbilCancelFrames(core.ActionCharge, core.ActionSkill, 94-70)  //charge -> skill
	c.SetAbilCancelFrames(core.ActionCharge, core.ActionBurst, 93-70)  //charge -> burst
	c.SetAbilCancelFrames(core.ActionCharge, core.ActionSwap, 90-70)   //charge -> swap

	c.SetAbilCancelFrames(core.ActionBurst, core.ActionAttack, 86-56) //burst -> n1
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionCharge, 86-56) //burst -> charge
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionSkill, 87-56)  //burst -> skill
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionDash, 88-56)   //burst -> dash
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionJump, 57-56)   //burst -> jump
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionSwap, 56-56)   //burst -> swap

	c.SetAbilCancelFrames(core.ActionSkill, core.ActionAttack, 37-20)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionCharge, 38-20)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionBurst, 40-20)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionDash, 35-20)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionJump, 20-20)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionSwap, 23-20)
}

func (c *char) ActionInterruptableDelay(next core.ActionType, p map[string]int) int {
	// Provide a custom override for Lisa's Hold E
	if c.Core.LastAction.Typ == core.ActionSkill &&
		c.Core.LastAction.Param["hold"] == 1 {
		return SkillHoldFrames(next)
	}
	//otherwise use default implementation
	return c.Tmpl.ActionInterruptableDelay(next, p)
}

func SkillHoldFrames(next core.ActionType) int {
	switch next {
	case core.ActionAttack:
		return 143 - 114
	case core.ActionCharge:
		return 125 - 114
	case core.ActionBurst:
		return 138 - 114
	case core.ActionDash:
		return 116 - 114
	case core.ActionJump:
		return 117 - 114
	default:
		return 114 - 114
	}
}

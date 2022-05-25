package keqing

import "github.com/genshinsim/gcsim/pkg/core"

func (c *char) ActionFrames(a core.ActionType, p map[string]int) (int, int) {
	switch a {
	case core.ActionAttack:
		f := 11
		a := 0
		switch c.NormalCounter {
		//TODO: need to add atkspd mod
		case 0:
			f = 11
			a = 15
		case 1:
			f = 11
			a = 16
		case 2:
			f = 15
			a = 27
		case 3:
			f = 22
			a = 31
		case 4:
			f = 26
			a = 66
		}
		return f, a
	case core.ActionCharge:
		return 24, 36
	case core.ActionSkill:
		if c.Core.Status.Duration(stilettoKey) > 0 {
			//2nd part
			return 15, 42
		}
		//first part
		return 21, 36
	case core.ActionBurst:
		return 56, 124
	case core.ActionDash:
		return 20, 20
	}
	c.Core.Log.NewEventBuildMsg(core.LogActionEvent, c.Index, "unknown action (invalid frames): ", a.String())
	return 0, 0
}

func (c *char) InitCancelFrames() {

	//normal cancels
	c.SetNormalCancelFrames(0, core.ActionAttack, 15-11) //n1 -> next attack
	c.SetNormalCancelFrames(0, core.ActionCharge, 21-11) //n1 -> charge

	c.SetNormalCancelFrames(1, core.ActionAttack, 16-11) //n2 -> next attack
	c.SetNormalCancelFrames(1, core.ActionCharge, 24-11) //n2 -> charge

	c.SetNormalCancelFrames(2, core.ActionAttack, 27-15) //n3 -> next attack
	c.SetNormalCancelFrames(2, core.ActionCharge, 36-15) //n3 -> charge

	c.SetNormalCancelFrames(3, core.ActionAttack, 31-22) //n4 -> next attack
	c.SetNormalCancelFrames(3, core.ActionCharge, 58-22) //n4 -> charge

	c.SetNormalCancelFrames(4, core.ActionAttack, 66-26) //n5 -> n1

	c.SetAbilCancelFrames(core.ActionCharge, core.ActionAttack, 36-24) //charge -> n1
	c.SetAbilCancelFrames(core.ActionCharge, core.ActionSkill, 35-24)  //charge -> skill
	c.SetAbilCancelFrames(core.ActionCharge, core.ActionBurst, 35-24)  //charge -> burst

	c.SetAbilCancelFrames(core.ActionBurst, core.ActionAttack, 124-56) //burst -> n1
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionSkill, 124-56)  //burst -> skill
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionDash, 122-56)   //burst -> dash
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionJump, 124-56)   //burst -> jump
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionSwap, 123-56)   //burst -> swap

	c.SetAbilCancelFrames(core.ActionSkill, core.ActionAttack, 36-21)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionSkill, 35-21)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionBurst, 37-21)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionDash, 21-21)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionJump, 21-21)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionSwap, 28-21)
}

func (c *char) ActionInterruptableDelay(next core.ActionType, p map[string]int) int {
	// Provide a custom override for Keqing's E recast
	if c.Core.LastAction.Typ == core.ActionSkill &&
		c.Core.Status.Duration(stilettoKey) == 0 {
		return SkillRecastFrames(next)
	}
	//otherwise use default implementation
	return c.Tmpl.ActionInterruptableDelay(next, p)
}

func SkillRecastFrames(next core.ActionType) int {
	switch next {
	case core.ActionAttack:
		return 42 - 15
	case core.ActionBurst:
		return 43 - 15
	case core.ActionDash:
		return 15 - 15
	case core.ActionJump:
		return 16 - 15
	case core.ActionSwap:
		return 42 - 15
	default:
		return 42 - 15
	}
}

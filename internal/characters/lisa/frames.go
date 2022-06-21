package lisa

import "github.com/genshinsim/gcsim/pkg/core"

func (c *char) ActionFrames(a action.Action, p map[string]int) (int, int) {
	switch a {
	case action.ActionAttack:
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
		f = int(float64(f) / (1 + c.Stat(attributes.AtkSpd)))
		return f, a
	case action.ActionCharge:
		return 70, 91
	case action.ActionSkill:
		hold := p["hold"]
		if hold == 0 {
			return 20, 38	//no hold
		}
		//yes hold
		return 114, 141
	case action.ActionBurst:
		return 56, 86
	case action.ActionDash:
		return 22, 22
	case action.ActionJump:
		return 33, 33
	default:
		c.Core.Log.NewEventBuildMsg(glog.LogActionEvent, c.Index, "unknown action (invalid frames): ", a.String())
		return 0, 0
	}
}

func (c *char) InitCancelFrames() {

	//normal cancels
	c.SetNormalCancelFrames(0, action.ActionAttack, 30-15)	//n1 -> next attack
	c.SetNormalCancelFrames(0, action.ActionCharge, 15-15)	//n1 -> charge
	c.SetNormalCancelFrames(0, action.ActionSkill, 26-15)	//only CA can be done before hitmark
	c.SetNormalCancelFrames(0, action.ActionBurst, 26-15)
	c.SetNormalCancelFrames(0, action.ActionDash, 26-15)
	c.SetNormalCancelFrames(0, action.ActionJump, 26-15)
	c.SetNormalCancelFrames(0, action.ActionSwap, 26-15)

	c.SetNormalCancelFrames(1, action.ActionAttack, 20-12)	//n2 -> next attack
	c.SetNormalCancelFrames(1, action.ActionCharge, 12-12)	//n2 -> charge
	c.SetNormalCancelFrames(1, action.ActionSkill, 18-12)	//only CA can be done before hitmark
	c.SetNormalCancelFrames(1, action.ActionBurst, 18-12)
	c.SetNormalCancelFrames(1, action.ActionDash, 18-12)
	c.SetNormalCancelFrames(1, action.ActionJump, 18-12)
	c.SetNormalCancelFrames(1, action.ActionSwap, 18-12)

	c.SetNormalCancelFrames(2, action.ActionAttack, 34-17)	//n3 -> n4
	c.SetNormalCancelFrames(2, action.ActionCharge, 26-17)	//n3 -> charge
	//CA is after hitmark this time, so no need for the rest

	c.SetNormalCancelFrames(3, action.ActionAttack, 57-31)	//n4 -> n1
	//Missing N4->CA - it's extremely long though.

	c.SetAbilCancelFrames(action.ActionCharge, action.ActionAttack, 86-70)	//charge -> n1
	c.SetAbilCancelFrames(action.ActionCharge, action.ActionCharge, 90-70)	//charge -> charge
	c.SetAbilCancelFrames(action.ActionCharge, action.ActionSkill, 94-70)	//charge -> skill
	c.SetAbilCancelFrames(action.ActionCharge, action.ActionBurst, 93-70)	//charge -> burst
	c.SetAbilCancelFrames(action.ActionCharge, action.ActionSwap, 90-70)	//charge -> swap

	c.SetAbilCancelFrames(action.ActionBurst, action.ActionAttack, 86-56)	//burst -> n1
	c.SetAbilCancelFrames(action.ActionBurst, action.ActionCharge, 86-56)	//burst -> charge
	c.SetAbilCancelFrames(action.ActionBurst, action.ActionSkill, 87-56)	//burst -> skill
	c.SetAbilCancelFrames(action.ActionBurst, action.ActionDash, 88-56)	//burst -> dash
	c.SetAbilCancelFrames(action.ActionBurst, action.ActionJump, 57-56)	//burst -> jump
	c.SetAbilCancelFrames(action.ActionBurst, action.ActionSwap, 56-56)	//burst -> swap

	c.SetAbilCancelFrames(action.ActionSkill, action.ActionAttack, 37-20)
	c.SetAbilCancelFrames(action.ActionSkill, action.ActionCharge, 38-20)
	c.SetAbilCancelFrames(action.ActionSkill, action.ActionBurst, 40-20)
	c.SetAbilCancelFrames(action.ActionSkill, action.ActionDash, 35-20)
	c.SetAbilCancelFrames(action.ActionSkill, action.ActionJump, 20-20)
	c.SetAbilCancelFrames(action.ActionSkill, action.ActionSwap, 23-20)
}

func (c *char) ActionInterruptableDelay(next action.Action, p map[string]int) int {
	// Provide a custom override for Lisa's Hold E
	if c.Core.LastAction.Typ == action.ActionSkill &&
		c.Core.LastAction.Param["hold"] == 1 {
		return SkillHoldFrames(next)
	}
	//otherwise use default implementation
	return c.Tmpl.ActionInterruptableDelay(next, p)
}

func SkillHoldFrames(next action.Action,) int {
	switch next {
	case action.ActionAttack:
		return 143 - 114
	case action.ActionCharge:
		return 125 - 114
	case action.ActionBurst:
		return 138 - 114
	case action.ActionDash:
		return 116 - 114
	case action.ActionJump:
		return 117 - 114
	default:
		return 114 - 114
	}
}

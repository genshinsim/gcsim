package ayato

import "github.com/genshinsim/gcsim/pkg/core"

func (c *char) ActionFrames(a core.ActionType, p map[string]int) (int, int) {
	switch a {
	case core.ActionAttack:
		f := 0
		a := 0
		if c.Core.Status.Duration("soukaikanka") > 0 {
			f = 5
			a = 10
			f = int(float64(f) / (1 + c.Stat(core.AtkSpd)))

			return f, a
		}

		switch c.NormalCounter {
		//TODO: need to add atkspd mod
		case 0:
			f = 12
			a = 15
		case 1:
			f = 18
			a = 27
		case 2:
			f = 20
			a = 30
		case 3:
			f = 25
			a = 27
		case 4:
			f = 41
			a = 63
		}
		f = int(float64(f) / (1 + c.Stat(core.AtkSpd)))

		return f, a
	case core.ActionCharge:
		return 24, 55
	case core.ActionSkill:
		return 21, 35
	case core.ActionBurst:
		return 101, 168
	default:
		c.Core.Log.NewEventBuildMsg(core.LogActionEvent, c.Index, "unknown action (invalid frames): ", a.String())
		return 0, 0
	}
}

func (c *char) InitCancelFrames() {

	//normal cancels
	c.SetNormalCancelFrames(0, core.ActionAttack, 15-12) //n1 -> next attack
	c.SetNormalCancelFrames(0, core.ActionCharge, 24-12) //n1 -> charge

	c.SetNormalCancelFrames(1, core.ActionAttack, 27-18) //n2 -> next attack
	//c.SetNormalCancelFrames(1, core.ActionCharge, 29-13) //n2 -> aim

	c.SetNormalCancelFrames(2, core.ActionAttack, 30-20) //n3 -> next attack
	//c.SetNormalCancelFrames(2, core.ActionCharge, 35-19) //n3 -> aim

	c.SetNormalCancelFrames(3, core.ActionAttack, 27-25) //n4 -> next attack
	//c.SetNormalCancelFrames(3, core.ActionCharge, 33-17) //n4 -> aim

	c.SetNormalCancelFrames(4, core.ActionAttack, 63-41) //n5 -> n1

	c.SetAbilCancelFrames(core.ActionCharge, core.ActionAttack, 55-24) //charge -> n1
	c.SetAbilCancelFrames(core.ActionCharge, core.ActionSkill, 55-24)  //charge -> skill
	c.SetAbilCancelFrames(core.ActionCharge, core.ActionBurst, 55-24)  //charge -> burst
	c.SetAbilCancelFrames(core.ActionCharge, core.ActionSwap, 53-24)   //charge -> swap

	c.SetAbilCancelFrames(core.ActionBurst, core.ActionAttack, 102-101) //burst -> n1
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionSkill, 102-101)  //burst -> skill
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionDash, 102-101)   //burst -> dash
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionJump, 102-101)   //burst -> jump
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionSwap, 101-101)   //burst -> swap

	c.SetAbilCancelFrames(core.ActionSkill, core.ActionAttack, 21-21)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionBurst, 22-21)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionDash, 21-21)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionJump, 21-21)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionSwap, 21-21)
}

func (c *char) ActionInterruptableDelay(next core.ActionType, p map[string]int) int {
	// Provide a custom override for differentiating Ayato's E stance attacks from his usual attacks
	// we only over ride if prev is attack and next is also attack
	if c.Core.LastAction.Typ == core.ActionAttack &&
		next == core.ActionAttack &&
		c.Core.Status.Duration("soukaikanka") > 0 {
		f := 23 - 5
		f = int(float64(f) / (1 + c.Stat(core.AtkSpd)))
		return f
	}
	//otherwise use default implementation
	return c.Tmpl.ActionInterruptableDelay(next, p)
}

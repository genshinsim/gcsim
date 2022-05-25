package beidou

import "github.com/genshinsim/gcsim/pkg/core"

func (c *char) ActionFrames(a core.ActionType, p map[string]int) (int, int) {
	switch a {
	case core.ActionAttack:
		f := 0
		a := 0
		switch c.NormalCounter {
		//TODO: need to add atkspd mod
		case 0:
			f = 23
			a = 31
		case 1:
			f = 22
			a = 36
		case 2:
			f = 45
			a = 54
		case 3:
			f = 25
			a = 36
		case 4:
			f = 43
			a = 96
		}
		atkspd := c.Stat(core.AtkSpd)
		if c.Core.Status.Duration("beidoua4") > 0 {
			atkspd += 0.15
		}
		f = int(float64(f) / (1 + atkspd))
		return f, a
	case core.ActionCharge:
		f := 35 //frames from keqing lib
		atkspd := c.Stat(core.AtkSpd)
		if c.Core.Status.Duration("beidoua4") > 0 {
			atkspd += 0.15
		}
		f = int(float64(f) / (1 + atkspd))
		return f, f
	case core.ActionSkill:
		return 24, 44
	case core.ActionBurst:
		return 28, 55
	case core.ActionDash:
		return 22, 22
	case core.ActionJump:
		return 32, 32
	default:
		c.Core.Log.NewEventBuildMsg(core.LogActionEvent, c.Index, "unknown action (invalid frames): ", a.String())
		return 0, 0
	}
}

func (c *char) InitCancelFrames() {

	//normal cancels
	c.SetNormalCancelFrames(0, core.ActionAttack, 31-23) //n1 -> next attack
	c.SetNormalCancelFrames(1, core.ActionAttack, 36-22) //n2 -> next attack
	c.SetNormalCancelFrames(2, core.ActionAttack, 54-45) //n3 -> next attack
	c.SetNormalCancelFrames(3, core.ActionAttack, 36-25) //n4 -> next attack
	c.SetNormalCancelFrames(4, core.ActionAttack, 96-43) //n5 -> n1

	c.SetAbilCancelFrames(core.ActionBurst, core.ActionAttack, 55-28) //burst -> n1
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionSkill, 58-28)  //burst -> skill
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionDash, 48-28)   //burst -> dash
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionJump, 48-28)   //burst -> jump
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionSwap, 46-28)   //burst -> swap

	c.SetAbilCancelFrames(core.ActionSkill, core.ActionAttack, 44-24)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionBurst, 45-24)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionDash, 24-24)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionJump, 24-24)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionSwap, 44-24)
}

func (c *char) ActionInterruptableDelay(next core.ActionType, p map[string]int) int {
	return c.Tmpl.ActionInterruptableDelay(next, p)
}

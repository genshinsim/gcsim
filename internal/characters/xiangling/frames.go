package xiangling

import "github.com/genshinsim/gcsim/pkg/core"

func (c *char) ActionFrames(a core.ActionType, p map[string]int) (int, int) {
	switch a {
	case core.ActionAttack:
		f := 0
		a := 0
		switch c.NormalCounter {
		case 0:
			f = 12
			a = 20
		case 1:
			f = 8
			a = 17
		case 2:
			f = 18
			a = 28
		case 3:
			f = 29
			a = 37
		case 4:
			f = 21
			a = 71
		}
		f = int(float64(f) / (1 + c.Stat(core.AtkSpd)))
		return f, a
	case core.ActionSkill:
		return 14, 29
	case core.ActionBurst:
		return 18, 56
	case core.ActionCharge:
		return 24, 67
	case core.ActionDash:
		return 20, 20
	case core.ActionJump:
		return 30, 30
	default:
		c.Core.Log.NewEventBuildMsg(core.LogActionEvent, c.Index, "unknown action (invalid frames): ", a.String())
		return 0, 0
	}
}

func (c *char) InitCancelFrames() {

	//normal cancels
	c.SetNormalCancelFrames(0, core.ActionAttack, 20-12) //n1 -> next attack
	c.SetNormalCancelFrames(0, core.ActionCharge, 20-12) //n1 -> charge

	c.SetNormalCancelFrames(1, core.ActionAttack, 17-8) //n2 -> next attack
	c.SetNormalCancelFrames(1, core.ActionCharge, 17-8) //n2 -> charge

	c.SetNormalCancelFrames(2, core.ActionAttack, 28-18) //n3 -> next attack
	c.SetNormalCancelFrames(2, core.ActionCharge, 24-18) //n3 -> charge

	c.SetNormalCancelFrames(3, core.ActionAttack, 37-29) //n4 -> next attack
	c.SetNormalCancelFrames(3, core.ActionCharge, 34-29) //n4 -> charge

	c.SetNormalCancelFrames(4, core.ActionAttack, 71-21) //n5 -> n1

	c.SetAbilCancelFrames(core.ActionCharge, core.ActionAttack, 67-24) //charge -> n1
	c.SetAbilCancelFrames(core.ActionCharge, core.ActionSkill, 69-24)  //charge -> skill
	c.SetAbilCancelFrames(core.ActionCharge, core.ActionBurst, 67-24)  //charge -> burst
	c.SetAbilCancelFrames(core.ActionCharge, core.ActionSwap, 66-24)   //charge -> swap

	c.SetAbilCancelFrames(core.ActionBurst, core.ActionAttack, 80-18) //burst -> n1
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionSkill, 80-18)  //burst -> skill
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionDash, 80-18)   //burst -> dash
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionJump, 80-18)   //burst -> jump
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionSwap, 79-18)   //burst -> swap

	c.SetAbilCancelFrames(core.ActionSkill, core.ActionAttack, 39-14)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionBurst, 39-14)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionDash, 14-14)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionJump, 14-14)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionSwap, 38-14)
}

func (c *char) ActionInterruptableDelay(next core.ActionType, p map[string]int) int {
	return c.Tmpl.ActionInterruptableDelay(next, p)
}

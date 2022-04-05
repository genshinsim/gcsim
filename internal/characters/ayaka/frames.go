package ayaka

import "github.com/genshinsim/gcsim/pkg/core"

func (c *char) ActionFrames(a core.ActionType, p map[string]int) (int, int) {
	switch a {
	case core.ActionAttack:
		f := 0
		a := 0
		switch c.NormalCounter {
		//TODO: need to add atkspd mod
		case 0:
			f = 8
			a = 9
		case 1:
			f = 10
			a = 19
		case 2:
			f = 16
			a = 31
		case 3:
			f = 22
			a = 22
		case 4:
			f = 27
			a = 66
		}
		f = int(float64(f) / (1 + c.Stat(core.AtkSpd)))
		return f, a
	case core.ActionCharge:
		return 31, 62
	case core.ActionSkill:
		return 30, 48
	case core.ActionBurst:
		return 104, 124
	case core.ActionDash:
		return 20, 35
	default:
		c.Core.Log.NewEventBuildMsg(core.LogActionEvent, c.Index, "unknown action (invalid frames): ", a.String())
		return 0, 0
	}
}

func (c *char) InitCancelFrames() {
	//normal cancels
	c.SetNormalCancelFrames(0, core.ActionAttack, 9-8)  //n1 -> next attack
	c.SetNormalCancelFrames(0, core.ActionCharge, 22-8) //n1 -> charge

	c.SetNormalCancelFrames(1, core.ActionAttack, 19-10) //n2 -> next attack
	c.SetNormalCancelFrames(1, core.ActionCharge, 20-10) //n2 -> charge

	c.SetNormalCancelFrames(2, core.ActionAttack, 32-16) //n3 -> next attack
	c.SetNormalCancelFrames(2, core.ActionCharge, 31-16) //n3 -> charge

	c.SetNormalCancelFrames(3, core.ActionAttack, 22-22) //n4 -> next attack
	c.SetNormalCancelFrames(3, core.ActionCharge, 23-22) //n4 -> charge

	c.SetNormalCancelFrames(4, core.ActionAttack, 66-27) //n5 -> next attack (n1)

	c.SetAbilCancelFrames(core.ActionCharge, core.ActionAttack, 71-31) //charge -> n1
	c.SetAbilCancelFrames(core.ActionCharge, core.ActionSkill, 62-31)  //charge -> skill
	c.SetAbilCancelFrames(core.ActionCharge, core.ActionBurst, 63-31)  //charge -> burst

	c.SetAbilCancelFrames(core.ActionBurst, core.ActionAttack, 124-104) //burst -> n1
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionSkill, 125-104)  //burst -> skill
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionDash, 124-104)   //burst -> dash
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionJump, 114-104)   //burst -> jump
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionSwap, 123-104)   //burst -> swap

	c.SetAbilCancelFrames(core.ActionSkill, core.ActionAttack, 49-30)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionBurst, 48-30)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionDash, 30-30)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionJump, 32-30)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionSwap, 48-30)

	c.SetAbilCancelFrames(core.ActionDash, core.ActionAttack, 35-20)
	c.SetAbilCancelFrames(core.ActionDash, core.ActionSkill, 35-20)
	c.SetAbilCancelFrames(core.ActionDash, core.ActionBurst, 35-20)
	c.SetAbilCancelFrames(core.ActionDash, core.ActionDash, 30-20)
	c.SetAbilCancelFrames(core.ActionDash, core.ActionSwap, 34-20)
}

func (c *char) ActionInterruptableDelay(next core.ActionType, p map[string]int) int {
	return c.Tmpl.ActionInterruptableDelay(next, p)
}

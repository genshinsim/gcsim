package rosaria

import "github.com/genshinsim/gcsim/pkg/core"

func (c *char) ActionFrames(a core.ActionType, p map[string]int) (int, int) {
	switch a {
	case core.ActionAttack:
		f := 0
		a := 0

		switch c.NormalCounter {
		//TODO: need to add atkspd mod
		case 0:
			f = 9
			a = 19
		case 1:
			f = 13
			a = 23
		case 2:
			f = 28
			a = 31
		case 3:
			f = 32
			a = 44
		case 4:
			f = 40
			a = 66
		}
		f = int(float64(f) / (1 + c.Stat(core.AtkSpd)))

		return f, a
	case core.ActionCharge:
		return 22, 69
	case core.ActionSkill:
		return 24, 51
	case core.ActionBurst:
		return 15, 70
	default:
		c.Core.Log.NewEventBuildMsg(core.LogActionEvent, c.Index, "unknown action (invalid frames): ", a.String())
		return 0, 0
	}
}

func (c *char) InitCancelFrames() {

	//normal cancels
	c.SetNormalCancelFrames(0, core.ActionAttack, 19-9) //n1 -> next attack
	c.SetNormalCancelFrames(0, core.ActionCharge, 24-9) //n1 -> charge

	c.SetNormalCancelFrames(1, core.ActionAttack, 23-13) //n2 -> next attack
	c.SetNormalCancelFrames(1, core.ActionCharge, 27-13) //n2 -> charge

	c.SetNormalCancelFrames(2, core.ActionAttack, 31-28) //n3 -> next attack
	c.SetNormalCancelFrames(2, core.ActionCharge, 34-28) //n3 -> charge

	c.SetNormalCancelFrames(3, core.ActionAttack, 44-32) //n4 -> next attack
	c.SetNormalCancelFrames(3, core.ActionCharge, 52-32) //n4 -> charge

	c.SetNormalCancelFrames(4, core.ActionAttack, 66-40) //n5 -> n1

	c.SetAbilCancelFrames(core.ActionCharge, core.ActionAttack, 69-22) //charge -> n1
	c.SetAbilCancelFrames(core.ActionCharge, core.ActionSkill, 69-22)  //charge -> skill
	c.SetAbilCancelFrames(core.ActionCharge, core.ActionBurst, 69-22)  //charge -> burst
	c.SetAbilCancelFrames(core.ActionCharge, core.ActionDash, 39-22)   //charge -> dash
	c.SetAbilCancelFrames(core.ActionCharge, core.ActionJump, 25-22)   //charge -> jump
	c.SetAbilCancelFrames(core.ActionCharge, core.ActionSwap, 66-22)   //charge -> swap

	c.SetAbilCancelFrames(core.ActionBurst, core.ActionAttack, 70-15) //burst -> n1
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionSkill, 70-15)  //burst -> skill
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionDash, 57-15)   //burst -> dash
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionJump, 59-15)   //burst -> jump
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionSwap, 69-15)   //burst -> swap

	c.SetAbilCancelFrames(core.ActionSkill, core.ActionAttack, 51-24)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionBurst, 51-24)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionDash, 38-24)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionJump, 40-24)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionSwap, 50-24)
}

func (c *char) ActionInterruptableDelay(next core.ActionType, p map[string]int) int {
	return c.Tmpl.ActionInterruptableDelay(next, p)
}

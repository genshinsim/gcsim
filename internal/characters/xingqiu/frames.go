package xingqiu

import (
	"github.com/genshinsim/gcsim/pkg/core"
)

func (c *char) ActionFrames(a core.ActionType, p map[string]int) (int, int) {
	switch a {
	case core.ActionAttack:
		f := 0
		a := 0
		switch c.NormalCounter {
		//TODO: need to add atkspd mod
		case 0:
			f = 10
			a = 18
		case 1:
			f = 13
			a = 24
		case 2:
			f = 19
			a = 26
		case 3:
			f = 17
			a = 28
		case 4:
			f = 39
			a = 66
		}
		f = int(float64(f) / (1 + c.Stat(core.AtkSpd)))
		return f, a
	case core.ActionCharge:
		return 20, 58
	case core.ActionSkill:
		return 30, 67
	case core.ActionBurst:
		return 18, 33
	default:
		c.Core.Log.NewEventBuildMsg(core.LogActionEvent, c.Index, "unknown action (invalid frames): ", a.String())
		return 0, 0
	}
}

func (c *char) InitCancelFrames() {
	//normal cancels
	c.SetNormalCancelFrames(0, core.ActionAttack, 18-10) //n1 -> next attack
	c.SetNormalCancelFrames(0, core.ActionCharge, 35-10) //n1 -> charge

	c.SetNormalCancelFrames(1, core.ActionAttack, 24-13) //n2 -> next attack
	c.SetNormalCancelFrames(1, core.ActionCharge, 29-13) //n2 -> charge

	c.SetNormalCancelFrames(2, core.ActionAttack, 26-19) //n3 -> next attack
	c.SetNormalCancelFrames(2, core.ActionCharge, 35-19) //n3 -> charge

	c.SetNormalCancelFrames(3, core.ActionAttack, 28-17) //n4 -> next attack
	c.SetNormalCancelFrames(3, core.ActionCharge, 33-17) //n4 -> charge

	c.SetNormalCancelFrames(4, core.ActionAttack, 66-39) //n5 -> next attack (n1)

	c.SetAbilCancelFrames(core.ActionCharge, core.ActionAttack, 58-20) //charge -> n1
	c.SetAbilCancelFrames(core.ActionCharge, core.ActionSkill, 32-20)  //charge -> skill
	c.SetAbilCancelFrames(core.ActionCharge, core.ActionBurst, 32-20)  //charge -> burst
	c.SetAbilCancelFrames(core.ActionCharge, core.ActionSwap, 31-20)   //charge -> swap

	c.SetAbilCancelFrames(core.ActionBurst, core.ActionAttack, 33-18) //burst -> n1
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionSkill, 33-18)  //burst -> skill
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionDash, 33-18)   //burst -> dash
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionJump, 33-18)   //burst -> jump
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionSwap, 40-18)   //burst -> swap

	c.SetAbilCancelFrames(core.ActionSkill, core.ActionAttack, 67-30)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionSkill, 65-30)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionBurst, 67-30)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionDash, 30-30)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionJump, 34-30)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionSwap, 67-30)

}

func (c *char) ActionInterruptableDelay(next core.ActionType, p map[string]int) int {
	return c.Tmpl.ActionInterruptableDelay(next, p)
}

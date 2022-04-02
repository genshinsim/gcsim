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
			a = 17
		case 1:
			f = 16
			a = 32
		case 2:
			f = 43 //here and below is not recounted
			a = 43
		case 3:
			f = 36
			a = 36
		case 4:
			f = 78
			a = 78
		}
		f = int(float64(f) / (1 + c.Stat(core.AtkSpd)))
		return f, a
	case core.ActionCharge:
		return 63, 63 // not recounted
	case core.ActionSkill:
		return 27, 77 // a not recounted
	case core.ActionBurst:
		return 28, 39 // a not recounted
	default:
		c.Core.Log.NewEventBuildMsg(core.LogActionEvent, c.Index, "unknown action (invalid frames): ", a.String())
		return 0, 0
	}
}

func (c *char) InitCancelFrames() {
	//normal cancels
	c.SetNormalCancelFrames(0, core.ActionAttack, 17-10) //n1 -> next attack
	//c.SetNormalCancelFrames(0, core.ActionCharge, 16-10) //n1 -> charge

	c.SetNormalCancelFrames(1, core.ActionAttack, 32-16) //n2 -> next attack
	//c.SetNormalCancelFrames(1, core.ActionCharge, 25-11) //n2 -> charge

	//c.SetNormalCancelFrames(2, core.ActionAttack, 30-25) //n3 -> next attack
	//c.SetNormalCancelFrames(2, core.ActionCharge, 35-25) //n3 -> charge

	//c.SetNormalCancelFrames(3, core.ActionAttack, 40-15) //n4 -> next attack
	//c.SetNormalCancelFrames(3, core.ActionCharge, 36-15) //n4 -> charge

	//c.SetNormalCancelFrames(4, core.ActionAttack, 71-31) //n5 -> next attack (n1)
	// c.SetNormalCancelFrames(4, core.ActionCharge, 36-15) //n5 -> charge, missing this one

	/*c.SetAbilCancelFrames(core.ActionCharge, core.ActionAttack, 55-21) //charge -> n1
	c.SetAbilCancelFrames(core.ActionCharge, core.ActionSkill, 34-21)  //charge -> skill
	c.SetAbilCancelFrames(core.ActionCharge, core.ActionBurst, 33-21)  //charge -> burst
	c.SetAbilCancelFrames(core.ActionCharge, core.ActionSwap, 32-21)   //charge -> swap*/

	c.SetAbilCancelFrames(core.ActionBurst, core.ActionAttack, 29-28) //burst -> n1
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionSkill, 28-28)  //burst -> skill
	//c.SetAbilCancelFrames(core.ActionBurst, core.ActionDash, 100-95) //burst -> dash
	//c.SetAbilCancelFrames(core.ActionBurst, core.ActionJump, 100-95) //burst -> jump
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionSwap, 29-28) //burst -> swap

	//skill frames, dmg at 27
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionAttack, 85-14)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionBurst, 64-27)
	//c.SetAbilCancelFrames(core.ActionSkill, core.ActionDash, 85-14)
	//c.SetAbilCancelFrames(core.ActionSkill, core.ActionJump, 85-14)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionSwap, 78-27)

}

func (c *char) ActionInterruptableDelay(next core.ActionType, p map[string]int) int {
	return c.Tmpl.ActionInterruptableDelay(next, p)
}

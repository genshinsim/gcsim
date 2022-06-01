package hutao

import "github.com/genshinsim/gcsim/pkg/core"

var hitmarks = [][]int{{12}, {10}, {17}, {23}, {16, 27}, {27}}
var paramitaChargeHitmark = 6

func (c *char) ActionFrames(a core.ActionType, p map[string]int) (int, int) {
	switch a {
	case core.ActionAttack:
		f := 0
		a := 0

		switch c.NormalCounter {
		//TODO: need to add atkspd mod
		case 0:
			f = 12
			a = 16
		case 1:
			f = 10
			a = 12
		case 2:
			f = 17
			a = 28
		case 3:
			f = 23
			a = 30
		case 4:
			f = 27
			a = 36
		case 5:
			f = 27
			a = 72
		}
		f = int(float64(f) / (1 + c.Stat(core.AtkSpd)))

		return f, a
	case core.ActionCharge:
		if c.Core.Status.Duration("paramita") > 0 {
			return 2, 44
		}
		return 19, 57
	case core.ActionSkill:
		return 14, 29
	case core.ActionBurst:
		return 66, 98
	default:
		c.Core.Log.NewEventBuildMsg(core.LogActionEvent, c.Index, "unknown action (invalid frames): ", a.String())
		return 0, 0
	}
}

func (c *char) InitCancelFrames() {

	//normal cancels
	c.SetNormalCancelFrames(0, core.ActionAttack, 16-12) //n1 -> next attack
	c.SetNormalCancelFrames(0, core.ActionCharge, 19-12) //n1 -> charge

	c.SetNormalCancelFrames(1, core.ActionAttack, 12-10) //n2 -> next attack
	c.SetNormalCancelFrames(1, core.ActionCharge, 18-10) //n2 -> charge

	c.SetNormalCancelFrames(2, core.ActionAttack, 28-17) //n3 -> next attack
	c.SetNormalCancelFrames(2, core.ActionCharge, 22-17) //n3 -> charge

	c.SetNormalCancelFrames(3, core.ActionAttack, 30-23) //n4 -> next attack
	c.SetNormalCancelFrames(3, core.ActionCharge, 32-23) //n4 -> charge

	c.SetNormalCancelFrames(4, core.ActionAttack, 36-27) //n5 -> n6
	c.SetNormalCancelFrames(3, core.ActionCharge, 48-27) //n5 -> charge

	c.SetNormalCancelFrames(4, core.ActionAttack, 72-27) //n6 -> n1

	c.SetAbilCancelFrames(core.ActionCharge, core.ActionAttack, 57-19) //charge -> n1
	c.SetAbilCancelFrames(core.ActionCharge, core.ActionSkill, 57-19)  //charge -> skill
	c.SetAbilCancelFrames(core.ActionCharge, core.ActionBurst, 60-19)  //charge -> burst
	c.SetAbilCancelFrames(core.ActionCharge, core.ActionSwap, 62-19)   //charge -> swap

	c.SetAbilCancelFrames(core.ActionBurst, core.ActionAttack, 98-66) //burst -> n1
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionSkill, 97-66)  //burst -> skill
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionDash, 98-66)   //burst -> dash
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionJump, 100-66)  //burst -> jump
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionSwap, 95-66)   //burst -> swap

	c.SetAbilCancelFrames(core.ActionSkill, core.ActionAttack, 29-14)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionBurst, 28-14)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionDash, 37-14)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionJump, 37-14)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionSwap, 52-14)
}

func (c *char) ActionInterruptableDelay(next core.ActionType, p map[string]int) int {
	// Provide a custom override for her CA cancels during paramita
	if c.Core.LastAction.Typ == core.ActionCharge &&
		c.Core.Status.Duration("paramita") > 0 {
		f := 0
		switch next {
		case core.ActionAttack:
			f = 44 - 2
		case core.ActionBurst:
			f = 35 - 2
		case core.ActionDash:
			f = 2 - 2
		case core.ActionJump:
			f = 2 - 2
		case core.ActionSwap:
			f = 42 - 2
		}
		return f
	}
	//otherwise use default implementation
	return c.Tmpl.ActionInterruptableDelay(next, p)
}

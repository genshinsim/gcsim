package heizou

import "github.com/genshinsim/gcsim/pkg/core"

func (c *char) ActionFrames(a core.ActionType, p map[string]int) (int, int) {
	switch a {
	case core.ActionAttack:
		f := 0
		a := 0
		switch c.NormalCounter {
		//TODO: need to add atkspd mod
		case 0:
			f = 20
			a = 12
		case 1:
			f = 17
			a = 13
		case 2:
			f = 45
			a = 21
		case 3:
			f = 36
			a = 27
		case 4:
			f = 66
			a = 31
		}
		f = int(float64(f) / (1 + c.Stat(core.AtkSpd)))

		return f, a
	case core.ActionCharge:
		return 38, 24
	case core.ActionSkill:
		//TODO: is this accurate? should it be 44 if 4 stack??
		return 37, 21
	case core.ActionBurst:
		return 71, 34
	case core.ActionDash:
		return 21, 21
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
	c.SetNormalCancelFrames(0, core.ActionCharge, 21-12) //n1 -> charge

	c.SetNormalCancelFrames(1, core.ActionAttack, 17-13) //n2 -> next attack
	c.SetNormalCancelFrames(1, core.ActionCharge, 21-13) //n2 -> charge

	c.SetNormalCancelFrames(2, core.ActionAttack, 45-21) //n3 -> next attack
	c.SetNormalCancelFrames(2, core.ActionCharge, 46-21) //n3 -> charge

	c.SetNormalCancelFrames(3, core.ActionAttack, 36-27) //n4 -> n5
	c.SetNormalCancelFrames(3, core.ActionCharge, 38-27) //n4 -> charge

	c.SetNormalCancelFrames(4, core.ActionAttack, 66-31) //n5 -> n1

	//CA is after hitmark this time, so no need for the rest

	c.SetAbilCancelFrames(core.ActionCharge, core.ActionAttack, 38-24) //charge -> n1
	c.SetAbilCancelFrames(core.ActionCharge, core.ActionSkill, 38-24)  //charge -> skill
	c.SetAbilCancelFrames(core.ActionCharge, core.ActionBurst, 38-24)  //charge -> burst
	c.SetAbilCancelFrames(core.ActionCharge, core.ActionSwap, 46-24)   //charge -> swap

	c.SetAbilCancelFrames(core.ActionBurst, core.ActionAttack, 71-34) //burst -> n1
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionSkill, 71-34)  //burst -> skill
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionDash, 72-34)   //burst -> dash
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionJump, 70-34)   //burst -> jump
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionSwap, 69-34)   //burst -> swap

	c.SetAbilCancelFrames(core.ActionSkill, core.ActionAttack, 37-21)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionBurst, 37-21)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionDash, 31-21)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionJump, 32-21)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionSwap, 30-21)

}

func (c *char) ActionInterruptableDelay(next core.ActionType, p map[string]int) int {
	return c.Tmpl.ActionInterruptableDelay(next, p)
}

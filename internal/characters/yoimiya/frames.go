package yoimiya

import "github.com/genshinsim/gcsim/pkg/core"

func (c *char) ActionFrames(a core.ActionType, p map[string]int) (int, int) {
	switch a {
	case core.ActionAttack:
		f := 0
		a := 0
		switch c.NormalCounter {
		//TODO: need to add atkspd mod
		case 0:
			f = 24
			a = 35
		case 1:
			f = 17
			a = 26
		case 2:
			f = 25
			a = 39
		case 3:
			f = 26
			a = 44
		case 4:
			f = 17
			a = 52
		}
		f = int(float64(f) / (1 + c.Stat(core.AtkSpd)))
		return f, a
	case core.ActionAim:
		return 86, 97
	case core.ActionSkill:
		return 11, 31
	case core.ActionBurst:
		return 75, 114
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
	//normal cancels (missing Nx -> Aim)
	c.SetNormalCancelFrames(0, core.ActionAttack, 35-24) //n1 -> next attack
	//c.SetNormalCancelFrames(0, core.Actionaim, 35-10) //n1 -> aim

	c.SetNormalCancelFrames(1, core.ActionAttack, 26-17) //n2 -> next attack
	//c.SetNormalCancelFrames(1, core.Actionaim, 29-13) //n2 -> aim

	c.SetNormalCancelFrames(2, core.ActionAttack, 39-25) //n3 -> next attack
	//c.SetNormalCancelFrames(2, core.Actionaim, 35-19) //n3 -> aim

	c.SetNormalCancelFrames(3, core.ActionAttack, 44-26) //n4 -> next attack
	//c.SetNormalCancelFrames(3, core.Actionaim, 33-17) //n4 -> aim

	c.SetNormalCancelFrames(4, core.ActionAttack, 52-17) //n5 -> n1
	//c.SetNormalCancelFrames(4, core.Actionaim, 33-17) //n5 -> aim

	//todo: get separate counts for each cancel, currently using generic frames for all of them
	c.SetAbilCancelFrames(core.ActionAim, core.ActionAttack, 97-86) //aim -> n1
	c.SetAbilCancelFrames(core.ActionAim, core.ActionAim, 97-86)    //aim -> aim
	c.SetAbilCancelFrames(core.ActionAim, core.ActionSkill, 97-86)  //aim -> skill
	c.SetAbilCancelFrames(core.ActionAim, core.ActionBurst, 97-86)  //aim -> burst
	c.SetAbilCancelFrames(core.ActionAim, core.ActionSwap, 97-86)   //aim -> swap

	c.SetAbilCancelFrames(core.ActionBurst, core.ActionAttack, 114-75) //burst -> n1
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionAim, 114-75)    //burst -> aim (assumed)
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionSkill, 110-75)  //burst -> skill
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionDash, 111-75)   //burst -> dash
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionJump, 113-75)   //burst -> jump
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionSwap, 109-75)   //burst -> swap

	c.SetAbilCancelFrames(core.ActionSkill, core.ActionAttack, 22-11)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionAim, 22-11) //assumed
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionBurst, 23-11)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionDash, 34-11)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionJump, 32-11)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionSwap, 31-11)

}

func (c *char) ActionInterruptableDelay(next core.ActionType, p map[string]int) int {
	return c.Tmpl.ActionInterruptableDelay(next, p)
}

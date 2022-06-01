package fischl

import "github.com/genshinsim/gcsim/pkg/core"

func (c *char) ActionFrames(a core.ActionType, p map[string]int) (int, int) {
	switch a {
	case core.ActionAttack:
		f := 0
		a := 0
		switch c.NormalCounter {
		//TODO: need to add atkspd mod
		case 0:
			f = 15
			a = 25
		case 1:
			f = 11
			a = 22
		case 2:
			f = 24
			a = 38
		case 3:
			f = 26
			a = 32
		case 4:
			f = 21
			a = 67
		}
		f = int(float64(f) / (1 + c.Stat(core.AtkSpd)))
		return f, a
	case core.ActionAim:
		return 86, 96
	case core.ActionSkill:
		return 14, 43
	case core.ActionBurst:
		return 18, 115
	default:
		c.Core.Log.NewEventBuildMsg(core.LogActionEvent, c.Index, "unknown action (invalid frames): ", a.String())
		return 0, 0
	}
}

func (c *char) InitCancelFrames() {

	//normal cancels (missing Nx -> Aim)
	c.SetNormalCancelFrames(0, core.ActionAttack, 25-15) //n1 -> next attack
	//c.SetNormalCancelFrames(0, core.ActionAim, 33-13) //n1 -> aim

	c.SetNormalCancelFrames(1, core.ActionAttack, 22-11) //n2 -> next attack
	//c.SetNormalCancelFrames(1, core.ActionAim, 27-9) //n2 -> aim

	c.SetNormalCancelFrames(2, core.ActionAttack, 38-24) //n3 -> next attack
	//c.SetNormalCancelFrames(2, core.ActionAim, 46-13) //n3 -> aim

	c.SetNormalCancelFrames(3, core.ActionAttack, 32-26) //n4 -> next attack
	//c.SetNormalCancelFrames(3, core.ActionAim, 48-25) //n4 -> aim

	c.SetNormalCancelFrames(4, core.ActionAttack, 67-21) //n5 -> n1

	//aim cancel frames are currently generic, should record specific cancels for each one at some point
	c.SetAbilCancelFrames(core.ActionAim, core.ActionAttack, 96-86) //aim -> n1
	c.SetAbilCancelFrames(core.ActionAim, core.ActionSkill, 96-86)  //aim -> skill
	c.SetAbilCancelFrames(core.ActionAim, core.ActionBurst, 96-86)  //aim -> burst
	c.SetAbilCancelFrames(core.ActionAim, core.ActionSwap, 96-86)   //aim -> swap
	c.SetAbilCancelFrames(core.ActionAim, core.ActionAim, 96-86)    //aim -> aim

	c.SetAbilCancelFrames(core.ActionBurst, core.ActionAttack, 148-18) //burst -> n1
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionAim, 148-18)    //burst -> aim (assumed)
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionDash, 111-18)   //burst -> dash
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionJump, 115-18)   //burst -> jump
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionSwap, 24-18)    //burst -> swap

	c.SetAbilCancelFrames(core.ActionSkill, core.ActionAttack, 43-14)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionAim, 43-14) //assumed
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionBurst, 43-14)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionDash, 14-14)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionJump, 16-14)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionSwap, 42-14)
}

func (c *char) ActionInterruptableDelay(next core.ActionType, p map[string]int) int {
	return c.Tmpl.ActionInterruptableDelay(next, p)
}

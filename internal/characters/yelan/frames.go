package yelan

import "github.com/genshinsim/gcsim/pkg/core"

func (c *char) ActionFrames(a core.ActionType, p map[string]int) (int, int) {
	switch a {
	case core.ActionAttack:
		f := 0
		a := 0
		switch c.NormalCounter {
		//TODO: need to add atkspd mod
		case 0:
			f = 13
			a = 15
		case 1:
			f = 13
			a = 21
		case 2:
			f = 18
			a = 38
		case 3:
			f = 29
			a = 67
		}
		f = int(float64(f) / (1 + c.Stat(core.AtkSpd)))
		return f, a
	case core.ActionAim:
		if c.Core.Status.Duration("breakthrough") > 0 { //Reduce required by 80% time if she has breakthrough
			reduced_frames := 34
			return int(reduced_frames), int(reduced_frames)
		}
		return 74, 74 //kqm
	case core.ActionBurst:
		return 76, 93
	case core.ActionSkill:
		return 35, 42
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
	//normal cancels - missing Nx -> Aim
	c.SetNormalCancelFrames(0, core.ActionAttack, 15-13) //n1 -> next attack
	//c.SetNormalCancelFrames(0, core.Actionaim, 35-10) //n1 -> aim

	c.SetNormalCancelFrames(1, core.ActionAttack, 21-13) //n2 -> next attack
	//c.SetNormalCancelFrames(1, core.Actionaim, 29-13) //n2 -> aim

	c.SetNormalCancelFrames(2, core.ActionAttack, 38-18) //n3 -> next attack
	//c.SetNormalCancelFrames(2, core.Actionaim, 35-19) //n3 -> aim

	c.SetNormalCancelFrames(3, core.ActionAttack, 67-29) //n4 -> next attack (n1)
	//c.SetNormalCancelFrames(3, core.Actionaim, 33-17) //n4 -> aim

	//todo: aim cancels

	c.SetAbilCancelFrames(core.ActionBurst, core.ActionAttack, 92-76) //burst -> n1
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionAim, 92-76)    //burst -> aim (assumed)
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionSkill, 93-76)  //burst -> skill
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionDash, 93-76)   //burst -> dash
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionJump, 92-76)   //burst -> jump
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionSwap, 91-76)   //burst -> swap

	c.SetAbilCancelFrames(core.ActionSkill, core.ActionAttack, 42-35)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionAim, 42-35) //assumed
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionBurst, 41-35)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionDash, 41-35)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionJump, 41-35)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionSwap, 40-35)
	//missing skill->skill

}

func (c *char) ActionInterruptableDelay(next core.ActionType, p map[string]int) int {
	return c.Tmpl.ActionInterruptableDelay(next, p)
}

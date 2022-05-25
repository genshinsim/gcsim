package ganyu

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
			a = 19
		case 1:
			f = 14
			a = 27
		case 2:
			f = 20
			a = 38
		case 3:
			f = 26
			a = 37
		case 4:
			f = 21
			a = 28
		case 5:
			f = 22
			a = 59
		}
		f = int(float64(f) / (1 + c.Stat(core.AtkSpd)))
		return f, a
	case core.ActionAim:
		//check for c6, if active then return 10, otherwise 115
		if c.Core.Status.Duration("ganyuc6") > 0 {
			c.Core.Log.NewEvent("ganyu c6 proc used", core.LogCharacterEvent, c.Index, "char", c.Index)
			c.Core.Status.DeleteStatus("ganyuc6")
			return 10, 10
		}
		return 103, 113
	case core.ActionSkill:
		return 13, 28 //ok
	case core.ActionBurst:
		return 122, 130 //ok
	default:
		c.Core.Log.NewEventBuildMsg(core.LogActionEvent, c.Index, "unknown action (invalid frames): ", a.String())
		return 0, 0
	}
}

func (c *char) InitCancelFrames() {
	//normal cancels
	c.SetNormalCancelFrames(0, core.ActionAttack, 19-13) //n1 -> next attack
	//c.SetNormalCancelFrames(0, core.Actionaim, 35-10) //n1 -> aim

	c.SetNormalCancelFrames(1, core.ActionAttack, 27-14) //n2 -> next attack
	//c.SetNormalCancelFrames(1, core.Actionaim, 29-13) //n2 -> aim

	c.SetNormalCancelFrames(2, core.ActionAttack, 38-20) //n3 -> next attack
	//c.SetNormalCancelFrames(2, core.Actionaim, 35-19) //n3 -> aim

	c.SetNormalCancelFrames(3, core.ActionAttack, 37-26) //n4 -> next attack
	//c.SetNormalCancelFrames(3, core.Actionaim, 33-17) //n4 -> aim

	c.SetNormalCancelFrames(4, core.ActionAttack, 28-21) //n5 -> next attack
	//c.SetNormalCancelFrames(4, core.Actionaim, 33-17) //n5 -> aim

	c.SetNormalCancelFrames(5, core.ActionAttack, 59-22) //n6 -> next attack (n1)
	//c.SetNormalCancelFrames(5, core.Actionaim, 33-17) //n6 -> aim

	//todo: get separate counts for each cancel, currently using generic frames for all of them
	c.SetAbilCancelFrames(core.ActionAim, core.ActionAttack, 113-103) //aim -> n1
	c.SetAbilCancelFrames(core.ActionAim, core.ActionAim, 113-103)    //aim -> aim
	c.SetAbilCancelFrames(core.ActionAim, core.ActionSkill, 113-103)  //aim -> skill
	c.SetAbilCancelFrames(core.ActionAim, core.ActionBurst, 113-103)  //aim -> burst
	c.SetAbilCancelFrames(core.ActionAim, core.ActionSwap, 113-103)   //aim -> swap

	c.SetAbilCancelFrames(core.ActionBurst, core.ActionAttack, 124-122) //burst -> n1
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionAim, 124-122)    //burst -> aim (assumed)
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionSkill, 124-122)  //burst -> skill
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionDash, 124-122)   //burst -> dash
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionJump, 124-122)   //burst -> jump
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionSwap, 122-122)   //burst -> swap

	c.SetAbilCancelFrames(core.ActionSkill, core.ActionAttack, 28-13)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionAim, 28-13) //assumed
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionBurst, 28-13)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionDash, 28-13)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionJump, 28-13)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionSwap, 27-13)

}

func (c *char) ActionInterruptableDelay(next core.ActionType, p map[string]int) int {
	return c.Tmpl.ActionInterruptableDelay(next, p)
}

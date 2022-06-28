package traveleranemo

import (
	"github.com/genshinsim/gcsim/pkg/core"
)

func (c *char) ActionFrames(a core.ActionType, p map[string]int) (int, int) {
	switch a {
	case core.ActionAttack:
		f := 0
		switch c.NormalCounter {
		//TODO: need to add atkspd mod
		case 0:
			f = 14 //frames from keqing lib
		case 1:
			f = 12
		case 2:
			f = 17
		case 3:
			f = 22
		case 4:
			f = 19
		}
		f = int(float64(f) / (1 + c.Stat(core.AtkSpd)))
		return f, f
	case core.ActionCharge:
		return 54 - 13, 53 - 13
	case core.ActionSkill:
		holdTicks := 0
		if p["hold"] == 1 {
			holdTicks = 6
		}
		if 0 < p["hold_ticks"] && p["hold_ticks"] <= 6 {
			holdTicks = p["hold_ticks"]
		}
		if holdTicks == 0 {
			return 28, 28
		}

		f := 31 + 15*(holdTicks-1) + 4
		if holdTicks >= 2 {
			f += 5
		}
		return f, f

	case core.ActionBurst:
		return 91, 91
	default:
		c.Core.Log.NewEventBuildMsg(core.LogActionEvent, c.Index, "unknown action (invalid frames): ", a.String())
		return 0, 0
	}
}

func (c *char) InitCancelFrames() {
	//normal cancels
	//normal cancels
	c.SetNormalCancelFrames(0, core.ActionAttack, 6) //n1 -> next attack
	//c.SetNormalCancelFrames(0, core.Actionaim, 35-10) //n1 -> aim

	c.SetNormalCancelFrames(1, core.ActionAttack, 11) //n2 -> next attack
	//c.SetNormalCancelFrames(1, core.Actionaim, 29-13) //n2 -> aim

	c.SetNormalCancelFrames(2, core.ActionAttack, 7) //n3 -> next attack
	//c.SetNormalCancelFrames(2, core.Actionaim, 35-19) //n3 -> aim

	c.SetNormalCancelFrames(3, core.ActionAttack, 13) //n4 -> next attack
	//c.SetNormalCancelFrames(3, core.Actionaim, 33-17) //n4 -> aim

	c.SetNormalCancelFrames(4, core.ActionAttack, 44) //n5 -> next attack
	//c.SetNormalCancelFrames(4, core.Actionaim, 33-17) //n5 -> aim

	c.SetAbilCancelFrames(core.ActionBurst, core.ActionAttack, 106-91) //burst -> n1
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionWalk, 111-91)   //burst -> walk
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionSkill, 104-91)  //burst -> skill
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionDash, 91-91)    //burst -> dash
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionJump, 91-91)    //burst -> jump
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionSwap, 96-91)    //burst -> swap

	c.SetAbilCancelFrames(core.ActionSkill, core.ActionAttack, 61-28)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionWalk, 81-28)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionBurst, 60-28)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionDash, 28-28)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionJump, 28-28)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionSwap, 60-28)

	c.SetAbilCancelFrames(core.ActionSkillHoldFramesOnly, core.ActionAttack, 82-55)
	c.SetAbilCancelFrames(core.ActionSkillHoldFramesOnly, core.ActionWalk, 103-55)
	// c.SetAbilCancelFrames(core.ActionSkillHoldFramesOnly, core.ActionSkill, ???-55)	// didn't test
	c.SetAbilCancelFrames(core.ActionSkillHoldFramesOnly, core.ActionBurst, 83-55)
	c.SetAbilCancelFrames(core.ActionSkillHoldFramesOnly, core.ActionDash, 55-55)
	c.SetAbilCancelFrames(core.ActionSkillHoldFramesOnly, core.ActionJump, 55-55)
	c.SetAbilCancelFrames(core.ActionSkillHoldFramesOnly, core.ActionSwap, 82-55)

}

func (c *char) ActionInterruptableDelay(next core.ActionType, p map[string]int) int {
	return c.Tmpl.ActionInterruptableDelay(next, p)
}

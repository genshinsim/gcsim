package venti

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
			f = 27
			a = 30
		case 1:
			f = 19
			a = 38
		case 2:
			f = 28
			a = 33
		case 3:
			f = 28
			a = 31
		case 4:
			f = 17
			a = 22
		case 5:
			f = 49
			a = 98
		}
		f = int(float64(f) / (1 + c.Stat(core.AtkSpd)))
		return f, a
	case core.ActionAim:
		return 86, 94
	case core.ActionHighPlunge:
		c.Core.Log.NewEvent("plunge skill check", core.LogCharacterEvent, c.Index, "previous", c.Core.LastAction)
		if c.Core.LastAction.Target == core.Venti && c.Core.LastAction.Typ == core.ActionSkill {
			h := c.Core.LastAction.Param["hold"]
			if h > 0 {
				return 58, 58
			}
		}
		c.Core.Log.NewEvent("invalid plunge (missing hold skill use)", core.LogActionEvent, c.Index, "action", a)
		return 0, 0
	case core.ActionSkill:
		if p["hold"] == 0 {
			return 22, 98
		}
		return 116, 174
	case core.ActionBurst:
		return 94, 108
	default:
		c.Core.Log.NewEventBuildMsg(core.LogActionEvent, c.Index, "unknown action (invalid frames): ", a.String())
		return 0, 0
	}
}

func (c *char) InitCancelFrames() {
	//normal cancels
	//normal cancels
	c.SetNormalCancelFrames(0, core.ActionAttack, 30-27) //n1 -> next attack
	//c.SetNormalCancelFrames(0, core.Actionaim, 35-10) //n1 -> aim

	c.SetNormalCancelFrames(1, core.ActionAttack, 38-19) //n2 -> next attack
	//c.SetNormalCancelFrames(1, core.Actionaim, 29-13) //n2 -> aim

	c.SetNormalCancelFrames(2, core.ActionAttack, 33-28) //n3 -> next attack
	//c.SetNormalCancelFrames(2, core.Actionaim, 35-19) //n3 -> aim

	c.SetNormalCancelFrames(3, core.ActionAttack, 31-28) //n4 -> next attack
	//c.SetNormalCancelFrames(3, core.Actionaim, 33-17) //n4 -> aim

	c.SetNormalCancelFrames(4, core.ActionAttack, 22-17) //n5 -> next attack
	//c.SetNormalCancelFrames(4, core.Actionaim, 33-17) //n5 -> aim

	c.SetNormalCancelFrames(5, core.ActionAttack, 98-49) //n6 -> next attack (n1)
	//c.SetNormalCancelFrames(5, core.Actionaim, 33-17) //n6 -> aim

	//todo: get separate counts for each cancel, currently using generic frames for all of them
	c.SetAbilCancelFrames(core.ActionAim, core.ActionAttack, 94-86) //aim -> n1
	c.SetAbilCancelFrames(core.ActionAim, core.ActionAim, 94-86)    //aim -> aim
	c.SetAbilCancelFrames(core.ActionAim, core.ActionSkill, 94-86)  //aim -> skill
	c.SetAbilCancelFrames(core.ActionAim, core.ActionBurst, 94-86)  //aim -> burst
	c.SetAbilCancelFrames(core.ActionAim, core.ActionSwap, 94-86)   //aim -> swap

	c.SetAbilCancelFrames(core.ActionBurst, core.ActionAttack, 95-94) //burst -> n1
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionAim, 95-94)    //burst -> aim (assumed)
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionSkill, 96-94)  //burst -> skill
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionDash, 95-94)   //burst -> dash
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionJump, 95-94)   //burst -> jump
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionSwap, 94-94)   //burst -> swap

	c.SetAbilCancelFrames(core.ActionSkill, core.ActionAttack, 22-22)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionAim, 22-22) //assumed
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionBurst, 22-22)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionDash, 22-22)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionJump, 22-22)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionSwap, 98-22)

	c.SetAbilCancelFrames(core.ActionSkillHoldFramesOnly, core.ActionHighPlunge, 116-116) //plunge should take 58 frames to hit
	c.SetAbilCancelFrames(core.ActionSkillHoldFramesOnly, core.ActionAttack, 289-116)     //float down if you didn't plunge
	c.SetAbilCancelFrames(core.ActionSkillHoldFramesOnly, core.ActionAim, 289-116)
	c.SetAbilCancelFrames(core.ActionSkillHoldFramesOnly, core.ActionSkill, 289-116)
	c.SetAbilCancelFrames(core.ActionSkillHoldFramesOnly, core.ActionBurst, 289-116)
	c.SetAbilCancelFrames(core.ActionSkillHoldFramesOnly, core.ActionDash, 289-116)
	c.SetAbilCancelFrames(core.ActionSkillHoldFramesOnly, core.ActionJump, 289-116)
	c.SetAbilCancelFrames(core.ActionSkillHoldFramesOnly, core.ActionSwap, 289-116)

}

func (c *char) ActionInterruptableDelay(next core.ActionType, p map[string]int) int {
	return c.Tmpl.ActionInterruptableDelay(next, p)
}

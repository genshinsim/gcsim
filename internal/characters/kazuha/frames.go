package kazuha

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
			f = 12
			a = 18
		case 1:
			f = 11
			a = 25
		case 2:
			f = 25 //hit mark 16, 25
			a = 35
		case 3:
			f = 15 //hit mark
			a = 40
		case 4:
			f = 31 //hit mark 15, 23, 31
			a = 71
		}
		f = int(float64(f) / (1 + c.Stat(core.AtkSpd)))
		return f, a
	case core.ActionCharge:
		return 21, 55 // kqm lib
	case core.ActionHighPlunge:
		c.Core.Log.NewEvent("plunge skill check", core.LogCharacterEvent, c.Index, "previous", c.Core.LastAction)
		if c.Core.LastAction.Target == core.Kazuha && c.Core.LastAction.Typ == core.ActionSkill {
			h := c.Core.LastAction.Param["hold"]
			if h > 0 {
				return 41, 60
			}
			return 36, 55
		}
		c.Core.Log.NewEvent("invalid plunge (missing skill use)", core.LogActionEvent, c.Index, "action", a)
		return 0, 0
	case core.ActionSkill:
		h := p["hold"]
		if h > 0 {
			return 34, 180
		}
		return 14, 88
	case core.ActionBurst:
		return 95, 100
	default:
		c.Core.Log.NewEventBuildMsg(core.LogActionEvent, c.Index, "unknown action (invalid frames): ", a.String())
		return 0, 0
	}
}

func (c *char) InitCancelFrames() {
	//normal cancels
	c.SetNormalCancelFrames(0, core.ActionAttack, 16-12) //n1 -> next attack
	c.SetNormalCancelFrames(0, core.ActionCharge, 16-12) //n1 -> charge

	c.SetNormalCancelFrames(1, core.ActionAttack, 20-11) //n2 -> next attack
	c.SetNormalCancelFrames(1, core.ActionCharge, 25-11) //n2 -> charge

	c.SetNormalCancelFrames(2, core.ActionAttack, 30-25) //n3 -> next attack
	c.SetNormalCancelFrames(2, core.ActionCharge, 35-25) //n3 -> charge

	c.SetNormalCancelFrames(3, core.ActionAttack, 40-15) //n4 -> next attack
	c.SetNormalCancelFrames(3, core.ActionCharge, 36-15) //n4 -> charge

	c.SetNormalCancelFrames(4, core.ActionAttack, 71-31) //n5 -> next attack (n1)
	// c.SetNormalCancelFrames(4, core.ActionCharge, 36-15) //n5 -> charge, missing this one

	c.SetAbilCancelFrames(core.ActionCharge, core.ActionAttack, 55-21) //charge -> n1
	c.SetAbilCancelFrames(core.ActionCharge, core.ActionSkill, 34-21)  //charge -> skill
	c.SetAbilCancelFrames(core.ActionCharge, core.ActionBurst, 33-21)  //charge -> burst
	c.SetAbilCancelFrames(core.ActionCharge, core.ActionSwap, 32-21)   //charge -> swap

	// c.SetAbilCancelFrames(core.ActionBurst, core.ActionAttack, 55-21) //burst -> n1
	// c.SetAbilCancelFrames(core.ActionBurst, core.ActionSkill, 34-21)  //burst -> skill
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionDash, 100-95) //burst -> dash
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionJump, 100-95) //burst -> jump
	// c.SetAbilCancelFrames(core.ActionBurst, core.ActionSwap, 32-21) //burst -> swap

	//skill press frames, dmg at 14
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionAttack, 85-14)     //85 frames to float down
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionBurst, 85-14)      //85 frames to float down
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionDash, 85-14)       //85 frames to float down
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionJump, 85-14)       //85 frames to float down
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionSwap, 80-14)       //85 frames to float down
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionHighPlunge, 27-14) //58 frames before you can start plunge

	//skill hold frames
	c.SetAbilCancelFrames(core.ActionSkillHoldFramesOnly, core.ActionAttack, 177-34)    //85 frames to float down
	c.SetAbilCancelFrames(core.ActionSkillHoldFramesOnly, core.ActionBurst, 177-34)     //85 frames to float down
	c.SetAbilCancelFrames(core.ActionSkillHoldFramesOnly, core.ActionDash, 177-34)      //85 frames to float down
	c.SetAbilCancelFrames(core.ActionSkillHoldFramesOnly, core.ActionJump, 177-34)      //85 frames to float down
	c.SetAbilCancelFrames(core.ActionSkillHoldFramesOnly, core.ActionSwap, 177-34)      //85 frames to float down
	c.SetAbilCancelFrames(core.ActionSkillHoldFramesOnly, core.ActionHighPlunge, 58-34) //58 frames before you can start plunge

}

func (c *char) ActionInterruptableDelay(next core.ActionType, p map[string]int) int {
	//provide override for skill and plunge cause it depends on hold vs press
	switch next {
	case core.ActionSkill:
		//have to check for press vs hold
		h := p["hold"]
		if h > 0 {
			return c.Tmpl.ActionInterruptableDelay(core.ActionSkillHoldFramesOnly, p)
		}
		return c.Tmpl.ActionInterruptableDelay(next, p)
	case core.ActionHighPlunge:
		//depends on if it was hold or press before so this has to be custom
		prev := c.Core.LastAction
		if prev.Typ != core.ActionSkill {
			//this should not happen
			c.Core.Log.NewEvent("ERROR: plunge used without skill use on kazuha!!", core.LogActionEvent, c.Index)
			return 0
		}
		//check if hold
		h := prev.Param["hold"]
		if h > 0 {
			return plungeHoldSkillFrames(next)
		}
		return plungePressSkillFrames(next)
	default:
		return c.Tmpl.ActionInterruptableDelay(next, p)
	}
}

func plungePressSkillFrames(next core.ActionType) int {
	switch next {
	case core.ActionAttack:
		return 55 - 36
	case core.ActionBurst:
		return 55 - 36
	default:
		return 47 - 36
	}
}

func plungeHoldSkillFrames(next core.ActionType) int {
	switch next {
	case core.ActionAttack:
		return 60 - 42
	case core.ActionBurst:
		return 60 - 42
	default:
		return 50 - 42
	}
}

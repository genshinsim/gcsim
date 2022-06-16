package kazuha

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

func (c *char) ActionFrames(a action.Action, p map[string]int) (int, int) {
	switch a {
	case action.ActionAttack:
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
		f = int(float64(f) / (1 + c.Stat(attributes.AtkSpd)))
		return f, a
	case action.ActionCharge:
		return 21, 55 // kqm lib
	case action.ActionHighPlunge:
		c.Core.Log.NewEvent("plunge skill check", glog.LogCharacterEvent, c.Index, "previous", c.Core.LastAction)
		if c.Core.LastAction.Target == core.Kazuha && c.Core.LastAction.Typ == action.ActionSkill {
			h := c.Core.LastAction.Param["hold"]
			if h > 0 {
				return 41, 60
			}
			return 36, 55
		}
		c.Core.Log.NewEvent("invalid plunge (missing skill use)", glog.LogActionEvent, c.Index, "action", a)
		return 0, 0
	case action.ActionSkill:
		h := p["hold"]
		if h > 0 {
			return 34, 180
		}
		return 14, 88
	case action.ActionBurst:
		return 95, 100
	default:
		c.Core.Log.NewEventBuildMsg(glog.LogActionEvent, c.Index, "unknown action (invalid frames): ", a.String())
		return 0, 0
	}
}

func (c *char) InitCancelFrames() {
	//normal cancels
	c.SetNormalCancelFrames(0, action.ActionAttack, 16-12) //n1 -> next attack
	c.SetNormalCancelFrames(0, action.ActionCharge, 16-12) //n1 -> charge

	c.SetNormalCancelFrames(1, action.ActionAttack, 20-11) //n2 -> next attack
	c.SetNormalCancelFrames(1, action.ActionCharge, 25-11) //n2 -> charge

	c.SetNormalCancelFrames(2, action.ActionAttack, 30-25) //n3 -> next attack
	c.SetNormalCancelFrames(2, action.ActionCharge, 35-25) //n3 -> charge

	c.SetNormalCancelFrames(3, action.ActionAttack, 40-15) //n4 -> next attack
	c.SetNormalCancelFrames(3, action.ActionCharge, 36-15) //n4 -> charge

	c.SetNormalCancelFrames(4, action.ActionAttack, 71-31) //n5 -> next attack (n1)
	// c.SetNormalCancelFrames(4, core.ActionCharge, 36-15) //n5 -> charge, missing this one

	c.SetAbilCancelFrames(action.ActionCharge, action.ActionAttack, 55-21) //charge -> n1
	c.SetAbilCancelFrames(action.ActionCharge, action.ActionSkill, 34-21)  //charge -> skill
	c.SetAbilCancelFrames(action.ActionCharge, action.ActionBurst, 33-21)  //charge -> burst
	c.SetAbilCancelFrames(action.ActionCharge, action.ActionSwap, 32-21)   //charge -> swap

	// c.SetAbilCancelFrames(core.ActionBurst, core.ActionAttack, 55-21) //burst -> n1
	// c.SetAbilCancelFrames(core.ActionBurst, core.ActionSkill, 34-21)  //burst -> skill
	c.SetAbilCancelFrames(action.ActionBurst, action.ActionDash, 100-95) //burst -> dash
	c.SetAbilCancelFrames(action.ActionBurst, action.ActionJump, 100-95) //burst -> jump
	// c.SetAbilCancelFrames(core.ActionBurst, core.ActionSwap, 32-21) //burst -> swap

	//skill press frames, dmg at 14
	c.SetAbilCancelFrames(action.ActionSkill, action.ActionAttack, 85-14)     //85 frames to float down
	c.SetAbilCancelFrames(action.ActionSkill, action.ActionBurst, 85-14)      //85 frames to float down
	c.SetAbilCancelFrames(action.ActionSkill, action.ActionDash, 85-14)       //85 frames to float down
	c.SetAbilCancelFrames(action.ActionSkill, action.ActionJump, 85-14)       //85 frames to float down
	c.SetAbilCancelFrames(action.ActionSkill, action.ActionSwap, 80-14)       //85 frames to float down
	c.SetAbilCancelFrames(action.ActionSkill, action.ActionHighPlunge, 27-14) //58 frames before you can start plunge

	//skill hold frames
	c.SetAbilCancelFrames(action.ActionSkillHoldFramesOnly, action.ActionAttack, 177-34)    //85 frames to float down
	c.SetAbilCancelFrames(action.ActionSkillHoldFramesOnly, action.ActionBurst, 177-34)     //85 frames to float down
	c.SetAbilCancelFrames(action.ActionSkillHoldFramesOnly, action.ActionDash, 177-34)      //85 frames to float down
	c.SetAbilCancelFrames(action.ActionSkillHoldFramesOnly, action.ActionJump, 177-34)      //85 frames to float down
	c.SetAbilCancelFrames(action.ActionSkillHoldFramesOnly, action.ActionSwap, 177-34)      //85 frames to float down
	c.SetAbilCancelFrames(action.ActionSkillHoldFramesOnly, action.ActionHighPlunge, 58-34) //58 frames before you can start plunge

}

func (c *char) ActionInterruptableDelay(next action.Action, p map[string]int) int {
	//provide override for skill and plunge cause it depends on hold vs press
	switch next {
	case action.ActionSkill:
		//have to check for press vs hold
		h := p["hold"]
		if h > 0 {
			return c.Tmpl.ActionInterruptableDelay(action.ActionSkillHoldFramesOnly, p)
		}
		return c.Tmpl.ActionInterruptableDelay(next, p)
	case action.ActionHighPlunge:
		//depends on if it was hold or press before so this has to be custom
		prev := c.Core.LastAction
		if prev.Typ != action.ActionSkill {
			//this should not happen
			c.Core.Log.NewEvent("ERROR: plunge used without skill use on kazuha!!", glog.LogActionEvent, c.Index)
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

func plungePressSkillFrames(next action.Action) int {
	switch next {
	case action.ActionAttack:
		return 55 - 36
	case action.ActionBurst:
		return 55 - 36
	default:
		return 47 - 36
	}
}

func plungeHoldSkillFrames(next action.Action) int {
	switch next {
	case action.ActionAttack:
		return 60 - 42
	case action.ActionBurst:
		return 60 - 42
	default:
		return 50 - 42
	}
}

package bennett

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
			a = 20
		case 1:
			f = 9
			a = 17
		case 2:
			f = 13
			a = 37
		case 3:
			f = 25
			a = 44
		case 4:
			f = 24
			a = 60
		}
		f = int(float64(f) / (1 + c.Stat(core.AtkSpd)))
		return f, a
	case core.ActionCharge:
		return 22, 55
	case core.ActionSkill:
		holdc4 := p["hold_c4"]
		if holdc4 == 1 {
			return 95, 107
		}

		hold := p["hold"]
		switch hold {
		case 1:
			return 57, 98
		case 2:
			return 166, 175
		default:
			return 16, 42
		}
	case core.ActionBurst:
		return 37, 53
	default:
		c.Core.Log.NewEventBuildMsg(core.LogActionEvent, c.Index, "unknown action (invalid frames): ", a.String())
		return 0, 0
	}
}

func (c *char) InitCancelFrames() {

	//normal cancels
	c.SetNormalCancelFrames(0, core.ActionAttack, 20-13) //n1 -> next attack
	c.SetNormalCancelFrames(0, core.ActionCharge, 33-13) //n1 -> charge

	c.SetNormalCancelFrames(1, core.ActionAttack, 17-9) //n2 -> next attack
	c.SetNormalCancelFrames(1, core.ActionCharge, 27-9) //n2 -> charge

	c.SetNormalCancelFrames(2, core.ActionAttack, 37-13) //n3 -> next attack
	c.SetNormalCancelFrames(2, core.ActionCharge, 46-13) //n3 -> charge

	c.SetNormalCancelFrames(3, core.ActionAttack, 44-25) //n4 -> next attack
	c.SetNormalCancelFrames(3, core.ActionCharge, 48-25) //n4 -> charge

	c.SetNormalCancelFrames(4, core.ActionAttack, 60-24) //n5 -> n1

	c.SetAbilCancelFrames(core.ActionCharge, core.ActionAttack, 55-22) //charge -> n1
	c.SetAbilCancelFrames(core.ActionCharge, core.ActionSkill, 41-22)  //charge -> skill
	c.SetAbilCancelFrames(core.ActionCharge, core.ActionBurst, 41-22)  //charge -> burst
	c.SetAbilCancelFrames(core.ActionCharge, core.ActionSwap, 44-22)   //charge -> swap

	c.SetAbilCancelFrames(core.ActionBurst, core.ActionAttack, 53-37) //burst -> n1
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionSkill, 53-37)  //burst -> skill
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionDash, 49-37)   //burst -> dash
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionJump, 50-37)   //burst -> jump
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionSwap, 51-37)   //burst -> swap

	c.SetAbilCancelFrames(core.ActionSkill, core.ActionAttack, 42-16)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionBurst, 42-16)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionDash, 22-16)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionJump, 23-16)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionSwap, 41-16)
}

func (c *char) ActionInterruptableDelay(next core.ActionType, p map[string]int) int {
	// Provide a custom override for Bennett's Hold E varieties
	if c.Core.LastAction.Typ == core.ActionSkill {
		//depends on which type of hold E was used
		hc4 := c.Core.LastAction.Param["hold_c4"]
		h := c.Core.LastAction.Param["hold"]

		//override for c4 hold tech
		if hc4 == 1 {
			return SkillHoldC4Frames(next)
		}

		switch h { //determine cancel frames based on which Hold E was used
		case 1:
			return SkillHoldLevel1Frames(next)
		case 2:
			if c.ModIsActive("bennett-field") { //Charge Level 2 has different cancel frames while in burst
				return SkillHoldLevel2FramesA4(next)
			}
			return SkillHoldLevel2Frames(next)
		default:
			return c.Tmpl.ActionInterruptableDelay(next, p)
		}
	}
	//otherwise use default implementation
	return c.Tmpl.ActionInterruptableDelay(next, p)
}

func SkillHoldC4Frames(next core.ActionType) int {
	switch next {
	case core.ActionAttack:
		return 107 - 95
	case core.ActionBurst:
		return 107 - 95
	case core.ActionSwap:
		return 106 - 95
	default:
		return 95 - 95
	}
}

func SkillHoldLevel1Frames(next core.ActionType) int {
	switch next {
	case core.ActionAttack:
		return 98 - 57
	case core.ActionBurst:
		return 97 - 57
	case core.ActionDash:
		return 65 - 57
	case core.ActionJump:
		return 66 - 57
	case core.ActionSwap:
		return 96 - 57
	default:
		return 57 - 57
	}
}

func SkillHoldLevel2Frames(next core.ActionType) int {
	switch next {
	case core.ActionAttack:
		return 343 - 166
	case core.ActionBurst:
		return 339 - 166
	case core.ActionDash:
		return 231 - 166
	case core.ActionJump:
		return 340 - 166
	case core.ActionSwap:
		return 337 - 166
	default:
		return 340 - 166
	}
}

func SkillHoldLevel2FramesA4(next core.ActionType) int {
	switch next {
	case core.ActionAttack:
		return 175 - 166
	case core.ActionDash:
		return 171 - 166
	case core.ActionJump:
		return 174 - 166
	case core.ActionSwap:
		return 175 - 166
	default:
		return 175 - 166
	}
}

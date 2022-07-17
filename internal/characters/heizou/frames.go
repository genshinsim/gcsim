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
		h := p["hold"]

		if h > 0 {
			stacks := c.skillHoldStacks(h)
			switch stacks {
			case 1:
				return 84, 65
			case 2:
				return 128, 108
			case 3:
				return 172, 152
			case 4:
				return 219, 198
			default:
				return 57, 37
			}

		}
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
	// Provide a custom override for Heizou's Hold E varieties
	if c.Core.LastAction.Typ == core.ActionSkill {
		//depends on how much stacks you are waiting for

		h := c.Core.LastAction.Param["hold"]
		if h < 1 {
			return c.Tmpl.ActionInterruptableDelay(next, p)
		}

		stacks := c.skillHoldStacks(h) //this should max out to 4 stacks

		switch stacks { //determine cancel frames based on which Hold E was used
		case 1:
			return SkillHoldStack1Frames(next)
		case 2:
			return SkillHoldStack2Frames(next)
		case 3:
			return SkillHoldStack3Frames(next)
		case 4:
			return SkillHoldStack4Frames(next)
		default:
			return SkillHoldStack0Frames(next)
		}
	}
	//otherwise use default implementation
	return c.Tmpl.ActionInterruptableDelay(next, p)
}

func SkillHoldStack1Frames(next core.ActionType) int { //skill hold(3->4 stacks)
	switch next {
	case core.ActionAttack:
		return 84 - 65
	case core.ActionBurst:
		return 85 - 65
	case core.ActionDash:
		return 78 - 65
	case core.ActionJump:
		return 78 - 65
	case core.ActionSwap:
		return 77 - 65
	default:
		return 65 - 65
	}
}

func SkillHoldStack2Frames(next core.ActionType) int { //skill hold(2->4 stacks)
	switch next {
	case core.ActionAttack:
		return 128 - 108
	case core.ActionBurst:
		return 127 - 108
	case core.ActionDash:
		return 122 - 108
	case core.ActionJump:
		return 123 - 108
	case core.ActionSwap:
		return 120 - 108
	default:
		return 108 - 108
	}
}

func SkillHoldStack3Frames(next core.ActionType) int { //skill hold(1->4 stacks)
	switch next {
	case core.ActionAttack:
		return 172 - 152
	case core.ActionBurst:
		return 172 - 152
	case core.ActionDash:
		return 167 - 152
	case core.ActionJump:
		return 167 - 152
	case core.ActionSwap:
		return 165 - 152
	default:
		return 152 - 152
	}
}

func SkillHoldStack4Frames(next core.ActionType) int { //skill hold(0->4 stacks)
	switch next {
	case core.ActionAttack:
		return 219 - 198
	case core.ActionBurst:
		return 218 - 198
	case core.ActionDash:
		return 212 - 198
	case core.ActionJump:
		return 212 - 198
	case core.ActionSwap:
		return 210 - 198
	default:
		return 198 - 198
	}
}

func SkillHoldStack0Frames(next core.ActionType) int { //skill hold(4->4 stacks)
	switch next {
	case core.ActionAttack:
		return 57 - 37
	case core.ActionBurst:
		return 57 - 37
	case core.ActionDash:
		return 53 - 37
	case core.ActionJump:
		return 53 - 37
	case core.ActionSwap:
		return 51 - 37
	default:
		return 37 - 37
	}
}

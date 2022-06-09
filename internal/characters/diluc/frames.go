package diluc

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
			a = 32
		case 1:
			f = 39
			a = 46
		case 2:
			f = 26
			a = 34
		case 3:
			f = 49
			a = 99
		}
		f = int(float64(f) / (1 + c.Stat(core.AtkSpd)))
		return f, a
	case core.ActionSkill:
		switch c.eCounter {
		case 1:
			return 29, 38
		case 2:
			return 47, 58
		default:
			return 25, 32
		}
	case core.ActionBurst:
		return 101, 140
	case core.ActionDash:
		return 19, 19
	default:
		c.Core.Log.NewEventBuildMsg(core.LogActionEvent, c.Index, "unknown action (invalid frames): ", a.String())
		return 0, 0
	}
}

func (c *char) InitCancelFrames() {

	//normal cancels, missing Nx->CA
	c.SetNormalCancelFrames(0, core.ActionAttack, 32-24) //n1 -> next attack
	//c.SetNormalCancelFrames(0, core.ActionCharge, 33-13) //n1 -> charge

	c.SetNormalCancelFrames(1, core.ActionAttack, 46-39) //n2 -> next attack
	//c.SetNormalCancelFrames(1, core.ActionCharge, 27-9) //n2 -> charge

	c.SetNormalCancelFrames(2, core.ActionAttack, 34-26) //n3 -> next attack
	//c.SetNormalCancelFrames(2, core.ActionCharge, 46-13) //n3 -> charge

	c.SetNormalCancelFrames(3, core.ActionAttack, 99-49) //n4 -> n1

	c.SetAbilCancelFrames(core.ActionBurst, core.ActionAttack, 140-101) //burst -> n1
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionSkill, 139-101)  //burst -> skill
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionDash, 139-101)   //burst -> dash
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionJump, 141-101)   //burst -> jump
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionSwap, 138-101)   //burst -> swap

	c.SetAbilCancelFrames(core.ActionSkill, core.ActionAttack, 32-25)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionBurst, 31-25)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionBurst, 32-25)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionSwap, 30-25)
}

func (c *char) ActionInterruptableDelay(next core.ActionType, p map[string]int) int {
	// Provide a custom override for Diluc's E
	if c.Core.LastAction.Typ == core.ActionSkill {
		switch c.eCounter { //depends on how many E's have been used
		case 1: //1 E - use default implementation
			return c.Tmpl.ActionInterruptableDelay(next, p)
		case 2: //2 Es
			return SecondSkillFrames(next)
		default: //3 Es
			return ThirdSkillFrames(next)
		}
	}
	//otherwise use default implementation
	return c.Tmpl.ActionInterruptableDelay(next, p)
}

func SecondSkillFrames(next core.ActionType) int {
	switch next {
	case core.ActionAttack:
		return 38 - 29
	case core.ActionSkill:
		return 37 - 29
	case core.ActionBurst:
		return 37 - 29
	case core.ActionDash:
		return 29 - 29
	case core.ActionJump:
		return 31 - 29
	case core.ActionSwap:
		return 36 - 29
	default:
		return 38 - 29
	}
}

func ThirdSkillFrames(next core.ActionType) int {
	switch next {
	case core.ActionAttack:
		return 58 - 47
	case core.ActionBurst:
		return 57 - 47
	case core.ActionDash:
		return 47 - 47
	case core.ActionJump:
		return 48 - 47
	case core.ActionSwap:
		return 66 - 47
	default:
		return 58 - 47
	}
}

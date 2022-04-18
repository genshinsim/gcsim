package tartaglia

import "github.com/genshinsim/gcsim/pkg/core"

func (c *char) ActionFrames(a core.ActionType, p map[string]int) (int, int) {
	switch a {
	case core.ActionAttack:
		f := 0
		a := 0
		if c.Core.Status.Duration("tartagliamelee") > 0 {
			switch c.NormalCounter {
			case 0:
				f = 7
				a = 9
			case 1:
				f = 8
				a = 11
			case 2:
				f = 17
				a = 32
			case 3:
				f = 7
				a = 27
			case 4:
				f = 14
				a = 27
			case 5:
				f = 21
				a = 66
			}
		} else {
			switch c.NormalCounter {
			//TODO: need to add atkspd mod
			case 0:
				f = 18
				a = 25
			case 1:
				f = 10
				a = 27
			case 2:
				f = 17
				a = 33
			case 3:
				f = 21
				a = 32
			case 4:
				f = 13
				a = 33
			case 5:
				f = 16
				a = 67
			}
		}
		atkspd := c.Stat(core.AtkSpd)
		f = int(float64(f) / (1 + atkspd))
		return f, a
	case core.Actionaim:
		if c.Core.Status.Duration("tartagliamelee") > 0 {
			return 9, 55
		}
		c.Core.Log.NewEvent("aim called when not in melee stance", core.LogActionEvent, c.Index, "action", a)
		return 0, 0
	case core.ActionAim:
		return 86, 94
	case core.ActionSkill:
		return 1, 18
	case core.ActionBurst:
		if c.Core.Status.Duration("tartagliamelee") > 0 {
			return 69, 102
		}
		return 51, 54
	case core.ActionDash:
		return 1, 24
	default:
		c.Core.Log.NewEventBuildMsg(core.LogActionEvent, c.Index, "unknown action (invalid frames): ", a.String())
		return 0, 0
	}
}

func (c *char) InitCancelFrames() {

	//normal cancels
	c.SetNormalCancelFrames(0, core.ActionAttack, 25-18) //n1 -> next attack
	//c.SetNormalCancelFrames(0, core.Actionaim, 25-18) //n1 -> aim (missing)
	//c.SetNormalCancelFrames(0, core.Actionaim, 25-18) //n1 -> walk (missing)

	c.SetNormalCancelFrames(1, core.ActionAttack, 27-10) //n2 -> next attack
	//c.SetNormalCancelFrames(1, core.Actionaim, 29-13) //n2 -> aim (missing)
	//c.SetNormalCancelFrames(0, core.Actionaim, 25-18) //n2 -> walk (missing)

	c.SetNormalCancelFrames(2, core.ActionAttack, 33-17) //n3 -> next attack
	//c.SetNormalCancelFrames(2, core.Actionaim, 35-19) //n3 -> aim (missing)
	//c.SetNormalCancelFrames(0, core.Actionaim, 25-18) //n3 -> walk (missing)

	c.SetNormalCancelFrames(3, core.ActionAttack, 32-21) //n4 -> next attack
	//c.SetNormalCancelFrames(3, core.Actionaim, 33-17) //n4 -> aim (missing)
	//c.SetNormalCancelFrames(0, core.Actionaim, 25-18) //n4 -> walk (missing)

	c.SetNormalCancelFrames(4, core.ActionAttack, 33-13) //n5 -> n6
	//c.SetNormalCancelFrames(3, core.Actionaim, 33-17) //n5 -> aim (missing)
	//c.SetNormalCancelFrames(0, core.Actionaim, 25-18) //n5 -> walk (missing)

	c.SetNormalCancelFrames(4, core.ActionAttack, 67-16) //n6 -> n1
	//c.SetNormalCancelFrames(0, core.Actionaim, 25-18) //n6 -> walk (missing)

	//TODO: Get frame counts for each specific cancel. All bow characters currently use generic recovery frames.
	c.SetAbilCancelFrames(core.ActionAim, core.ActionAttack, 94-86) //aim -> n1
	c.SetAbilCancelFrames(core.ActionAim, core.ActionAim, 94-86)    //aim -> aim
	c.SetAbilCancelFrames(core.ActionAim, core.ActionSkill, 94-86)  //aim -> skill
	c.SetAbilCancelFrames(core.ActionAim, core.ActionBurst, 94-86)  //aim -> burst
	c.SetAbilCancelFrames(core.ActionAim, core.ActionSwap, 94-86)   //aim -> swap
	c.SetAbilCancelFrames(core.ActionAim, core.ActionWalk, 94-86)   //aim -> walk

	//ranged burst
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionAttack, 54-51) //burst -> n1
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionAim, 54-51)    //aim -> aim (missing)
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionSkill, 55-51)  //burst -> skill
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionDash, 51-51)   //burst -> dash
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionJump, 51-51)   //burst -> jump
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionSwap, 53-51)   //burst -> swap
	//c.SetNormalCancelFrames(0, core.Actionaim, 25-18) //burst -> walk (missing)

	//skill (not proceeded by Walk or Dash) from ranged to melee
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionAttack, 19-1)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionBurst, 19-1)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionDash, 14-1)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionJump, 19-1)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionSwap, 39-1)
	c.SetAbilCancelFrames(core.AcitonSkill, core.ActionWalk, 19-1) //skill -> walk (missing, placeholder frames)

	//dash - can cancel it with E
	c.SetAbilCancelFrames(core.ActionDash, core.ActionAttack, 24-1)
	c.SetAbilCancelFrames(core.ActionDash, core.ActionAim, 24-1)
	c.SetAbilCancelFrames(core.ActionDash, core.ActionSkill, 3-1)
	c.SetAbilCancelFrames(core.ActionDash, core.ActionBurst, 24-1)
	c.SetAbilCancelFrames(core.ActionDash, core.ActionDash, 24-1)
	c.SetAbilCancelFrames(core.ActionDash, core.ActionJump, 24-1)
	c.SetAbilCancelFrames(core.ActionDash, core.ActionWalk, 24-1)
}

func (c *char) ActionInterruptableDelay(next core.ActionType, p map[string]int) int {
	// Provide a custom override for Childe Melee and his special E mechanics
	if c.Core.LastAction.Typ == core.ActionSkill {
		return skillFrames(next, c.skillSetup)
	} else if c.Core.Status.Duration("tartagliamelee") > 0 {
		return meleeFrames(next)
	}
	//otherwise use default implementation
	return c.Tmpl.ActionInterruptableDelay(next, p)
}

//Childe's E has different cancel frames if it was preceded by a walk or a dash.
func (c *char) skillFrames(next, prev core.ActionType) int {
	switch prev {
	case core.ActionDash: //preceded by dash
		return dashSkillFrames(next)
	case core.ActionWalk: //preceded by walk
		return walkSkillFrames(next)
	default:
		if c.Core.Status.Duration("tartagliamelee") > 0 {
			return meleeFrames(next) //proceeded by neither - redirect to melee frames
		}
		return c.Tmpl.ActionInterruptableDelay(next, p) //proceeded by neither in ranged. Use default implementation.
	}
}

func (c *char) meleeFrames(next core.ActionType) int {
	//melee frames here
}

func (c *char) dashSkillFrames(next core.ActionType) int {
	//dash e frames here
}

func (c *char) walkSkillFrames(next core.ActionType) int {
	//walk e frames here
}

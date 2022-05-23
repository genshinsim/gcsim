package jean

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
			a = 22
		case 1:
			f = 6
			a = 14
		case 2:
			f = 17
			a = 28
		case 3:
			f = 37
			a = 44
		case 4:
			f = 25
			a = 68
		}
		f = int(float64(f) / (1 + c.Stat(core.AtkSpd)))

		return f, a
	case core.ActionCharge:
		return 36, 57
	case core.ActionSkill:
		hold := p["hold"]
		//hold for p up to 5 seconds
		if hold > 300 {
			hold = 300
		}

		return 21 + hold, 46 + hold
	case core.ActionBurst:
		return 40, 83
	default:
		c.Core.Log.NewEventBuildMsg(core.LogActionEvent, c.Index, "unknown action (invalid frames): ", a.String())
		return 0, 0
	}
}

func (c *char) InitCancelFrames() {

	//normal cancels
	c.SetNormalCancelFrames(0, core.ActionAttack, 22-13) //n1 -> next attack
	c.SetNormalCancelFrames(0, core.ActionCharge, 25-13) //n1 -> charge

	c.SetNormalCancelFrames(1, core.ActionAttack, 14-6) //n2 -> next attack
	c.SetNormalCancelFrames(1, core.ActionCharge, 20-6) //n2 -> charge

	c.SetNormalCancelFrames(2, core.ActionAttack, 28-17) //n3 -> next attack
	c.SetNormalCancelFrames(2, core.ActionCharge, 31-17) //n3 -> charge

	c.SetNormalCancelFrames(3, core.ActionAttack, 44-37) //n4 -> next attack
	c.SetNormalCancelFrames(3, core.ActionCharge, 49-37) //n4 -> charge

	c.SetNormalCancelFrames(4, core.ActionAttack, 68-36) //n5 -> n1

	c.SetAbilCancelFrames(core.ActionCharge, core.ActionAttack, 57-36) //charge -> n1
	c.SetAbilCancelFrames(core.ActionCharge, core.ActionSkill, 57-36)  //charge -> skill
	c.SetAbilCancelFrames(core.ActionCharge, core.ActionBurst, 56-36)  //charge -> burst
	c.SetAbilCancelFrames(core.ActionCharge, core.ActionSwap, 39-36)   //charge -> swap

	c.SetAbilCancelFrames(core.ActionBurst, core.ActionAttack, 83-40) //burst -> n1
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionSkill, 83-40)  //burst -> skill
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionDash, 70-40)   //burst -> dash
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionJump, 70-40)   //burst -> jump
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionSwap, 84-40)   //burst -> swap

	c.SetAbilCancelFrames(core.ActionSkill, core.ActionAttack, 46-21)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionBurst, 46-21)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionDash, 28-21)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionJump, 28-21)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionSwap, 45-21)
}

func (c *char) ActionInterruptableDelay(next core.ActionType, p map[string]int) int {
	return c.Tmpl.ActionInterruptableDelay(next, p)
}

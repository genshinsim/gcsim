package kokomi

import "github.com/genshinsim/gcsim/pkg/core"

func (c *char) ActionFrames(a core.ActionType, p map[string]int) (int, int) {
	switch a {
	case core.ActionAttack:
		var f, a int
		switch c.NormalCounter {
		//TODO: need to add atkspd mod
		case 0:
			f = 4
			a = 30
		case 1:
			f = 12
			a = 34
		case 2:
			f = 28
			a = 65
		}
		f = int(float64(f) / (1 + c.Stat(core.AtkSpd)))
		return f, a
	case core.ActionCharge:
		return 48, 76
	case core.ActionSkill:
		return 30, 61
	case core.ActionBurst:
		return 76, 76
	default:
		c.Core.Log.NewEventBuildMsg(core.LogActionEvent, c.Index, "unknown action (invalid frames): ", a.String())
		return 0, 0
	}
}

func (c *char) InitCancelFrames() {
	// NA/CA cancels
	// we can't make the next CA faster by 14f so we reduce it here instead
	c.SetNormalCancelFrames(0, core.ActionAttack, 14-4-1)    // n1 -> next attack
	c.SetNormalCancelFrames(0, core.ActionCharge, 19-4-1-14) // n1 -> charge
	c.SetNormalCancelFrames(0, core.ActionWalk, 30-4-1)      // n1 -> walk

	c.SetNormalCancelFrames(1, core.ActionAttack, 30-12-1)    // n2 -> next attack
	c.SetNormalCancelFrames(1, core.ActionCharge, 34-12-1-14) // n2 -> charge
	c.SetNormalCancelFrames(1, core.ActionWalk, 34-12-1)      // n2 -> walk

	c.SetNormalCancelFrames(2, core.ActionAttack, 65-28-1)    // n3 -> next attack
	c.SetNormalCancelFrames(2, core.ActionCharge, 60-28-1-14) // n3 -> charge
	c.SetNormalCancelFrames(2, core.ActionWalk, 60-28-1)      // n3 -> walk

	// charge -> x
	c.SetAbilCancelFrames(core.ActionCharge, core.ActionAttack, 62-48)
	c.SetAbilCancelFrames(core.ActionCharge, core.ActionCharge, 62-48)
	c.SetAbilCancelFrames(core.ActionCharge, core.ActionSkill, 62-48)
	c.SetAbilCancelFrames(core.ActionCharge, core.ActionBurst, 62-48)
	c.SetAbilCancelFrames(core.ActionCharge, core.ActionSwap, 62-48)
	c.SetAbilCancelFrames(core.ActionCharge, core.ActionWalk, 76-48)

	// skill -> x
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionAttack, 61-30)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionCharge, 61-30)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionSkill, 61-30) // uses burst frames
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionBurst, 61-30)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionSwap, 61-30)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionWalk, 61-30)
}

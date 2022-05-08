package noelle

import "github.com/genshinsim/gcsim/pkg/core"

func (c *char) ActionFrames(a core.ActionType, p map[string]int) (int, int) {
	switch a {
	case core.ActionAttack:
		var f, a int
		switch c.NormalCounter {
		//TODO: need to add atkspd mod
		case 0:
			f = 28
			a = 38
		case 1:
			f = 25
			a = 46
		case 2:
			f = 20
			a = 31
		case 3:
			f = 42
			a = 107
		}
		f = int(float64(f) / (1 + c.Stat(core.AtkSpd)))
		return f, a
	case core.ActionSkill:
		return 12, 78
	case core.ActionBurst:
		return 82, 121
	default:
		c.Core.Log.NewEventBuildMsg(core.LogActionEvent, c.Index, "unknown action (invalid frames): ", a.String())
		return 0, 0
	}
}

func (c *char) InitCancelFrames() {
	// NA/CA cancels. TODO: CA cancels?
	c.SetNormalCancelFrames(0, core.ActionAttack, 38-28-1)  // n1 -> next attack
	c.SetNormalCancelFrames(1, core.ActionAttack, 46-25-1)  // n2 -> next attack
	c.SetNormalCancelFrames(2, core.ActionAttack, 31-20-1)  // n3 -> next attack
	c.SetNormalCancelFrames(3, core.ActionAttack, 107-42-1) // n4 -> next attack

	// skill -> x
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionSwap, 78-12)
	c.SetAbilCancelFrames(core.ActionSkill, core.ActionWalk, 43-12)

	// burst -> x
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionSwap, 121-82)
	c.SetAbilCancelFrames(core.ActionBurst, core.ActionWalk, 89-82)
}

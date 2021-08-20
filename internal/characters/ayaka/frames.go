package ayaka

import "github.com/genshinsim/gsim/pkg/core"

func (c *char) ActionFrames(a core.ActionType, p map[string]int) int {
	switch a {
	case core.ActionAttack:
		f := 0
		switch c.NormalCounter {
		//TODO: need to add atkspd mod
		case 0:
			f = 8
		case 1:
			f = 20
		case 2:
			f = 28
		case 3:
			f = 43
		case 4:
			f = 37
		}
		f = int(float64(f) / (1 + c.Stats[core.AtkSpd]))
		return f
	case core.ActionCharge:
		return 53
	case core.ActionSkill:
		return 56 //should be 82
	case core.ActionBurst:
		return 95 //ok
	default:
		c.Core.Log.Warnw("unknown action", "event", core.LogActionEvent, "frame", c.Core.F, "action", a)
		return 0
	}
}

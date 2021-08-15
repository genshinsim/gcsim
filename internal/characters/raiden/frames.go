package raiden

import "github.com/genshinsim/gsim/pkg/core"

func (c *char) ActionFrames(a core.ActionType, p map[string]int) int {
	switch a {
	case core.ActionAttack:
		f := 0
		if c.Core.Status.Duration("raidenburst") == 0 {
			switch c.NormalCounter {
			//TODO: need to add atkspd mod
			case 0:
				f = 13
			case 1:
				f = 19
			case 2:
				f = 22
			case 3:
				f = 44
			case 4:
				f = 41
			}
		} else {
			switch c.NormalCounter {
			//TODO: need to add atkspd mod
			case 0:
				f = 17
			case 1:
				f = 15
			case 2:
				f = 22
			case 3:
				f = 44
			case 4:
				f = 42
			}
		}
		f = int(float64(f) / (1 + c.Stats[core.AtkSpd]))
		return f
	case core.ActionSkill:
		return 40 //eye appears
	case core.ActionBurst:
		return 108
	default:
		c.Log.Warnw("unknown action", "event", core.LogActionEvent, "frame", c.Core.F, "action", a)
		return 0
	}
}

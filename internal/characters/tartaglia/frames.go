package tartaglia

import "github.com/genshinsim/gcsim/pkg/core"

func (c *char) ActionFrames(a core.ActionType, p map[string]int) (int, int) {
	switch a {
	case core.ActionAttack:
		f := 0
		if c.Core.Status.Duration("tartagliamelee") > 0 {
			switch c.NormalCounter {
			case 0:
				f = 7 //frames from keqing lib
			case 1:
				f = 20 - 7
			case 2:
				f = 48 - 20
			case 3:
				f = 80 - 48
			case 4:
				f = 116 - 80
			case 5:
				f = 165 - 116
			}
		} else {
			switch c.NormalCounter {
			//TODO: need to add atkspd mod
			case 0:
				f = 17 //frames from keqing lib
			case 1:
				f = 30 - 17
			case 2:
				f = 64 - 30
			case 3:
				f = 101 - 64
			case 4:
				f = 123 - 101
			case 5:
				f = 162 - 123
			}
		}
		atkspd := c.Stat(core.AtkSpd)
		f = int(float64(f) / (1 + atkspd))
		return f, f
	case core.ActionCharge:
		if c.Core.Status.Duration("tartagliamelee") > 0 {
			return 73, 73
		}
		c.Core.Log.Warnw("Charge called when not in melee stance", "event", core.LogActionEvent, "frame", c.Core.F, "action", a)
		return 0, 0
	case core.ActionAim:
		return 84, 84
	case core.ActionSkill:
		if c.Core.Status.Duration("tartagliamelee") > 0 {
			return 20, 20
		}
		//TODO: need exact frame
		return 28, 28
	case core.ActionBurst:
		if c.Core.Status.Duration("tartagliamelee") > 0 {
			return 97, 97
		}
		return 52, 52
	default:
		c.Core.Log.Warnw("unknown action", "event", core.LogActionEvent, "frame", c.Core.F, "action", a)
		return 0, 0
	}
}

package childe

import "github.com/genshinsim/gsim/pkg/core"

func (c *char) ActionFrames(a core.ActionType, p map[string]int) (int, int) {
	switch a {
	case core.ActionAttack:
		f := 0
		if c.Core.Status.Duration("childemelee") == 0 {
			switch c.NormalCounter {
			//TODO: need to add atkspd mod
			case 0:
				if c.Core.LastAction.Target == "tartaglia" && c.Core.LastAction.Typ == core.ActionAttack {
					f += 53
				}
				f = 17 //frames from keqing lib
				c.caFrame = 73
			case 1:
				f = 30 - 17
			case 2:
				f = 64 - 30
				c.caFrame = 77 // N3C(48+77) has different CA frame to N1C, N2C (7+73, 20+73)
			case 3:
				f = 101 - 64
			case 4:
				f = 123 - 101
			case 5:
				f = 162 - 123
			}
		} else {
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
		}
		atkspd := c.Stats[core.AtkSpd]
		f = int(float64(f) / (1 + atkspd))
		return f, f
	case core.ActionCharge:
		if c.Core.Status.Duration("childemelee") > 0 {
			return c.caFrame, c.caFrame
		}
		return 0, 0
	case core.ActionAim:
		return 84, 84
	case core.ActionSkill:
		return 28, 28
	case core.ActionBurst:
		if c.Core.Status.Duration("childemelee") == 0 {
			return 52, 52
		}
		return 97, 97
	default:
		c.Core.Log.Warnw("unknown action", "event", core.LogActionEvent, "frame", c.Core.F, "action", a)
		return 0, 0
	}
}

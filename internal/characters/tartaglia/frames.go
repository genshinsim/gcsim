package tartaglia

import "github.com/genshinsim/gsim/pkg/core"

func (c *char) ActionFrames(a core.ActionType, p map[string]int) (int, int) {
	switch a {
	case core.ActionAttack:
		f := 0
		if c.Core.Status.Duration("tartagliamelee") == 0 {
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
		} else {
			// Frames from KQM lib.
			// TODO: Hitlag issues? They say counts are done against ruin guard
			switch c.NormalCounter {
			//TODO: need to add atkspd mod
			case 0:
				f = 7
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
	case core.ActionAim:
		return 84, 84 //frames from keqing lib
	case core.ActionCharge:
		if c.Core.Status.Duration("tartagliamelee") != 0 {
			// Implied from N1C/N2C frame counts.
			// TODO: N3C yields different values (77 instead of 73)
			return 80 - 7, 80 - 7
		}
		c.Core.Log.Warnw("Charge called when not in melee stance", "event", core.LogActionEvent, "frame", c.Core.F, "action", a)
		return 0, 0
	case core.ActionSkill:
		return 28, 28 //ok
	// Has ranged and melee modes to his burst
	case core.ActionBurst:
		if c.Core.Status.Duration("tartagliamelee") == 0 {
			return 52, 52
		}
		return 97, 97
	default:
		c.Core.Log.Warnw("unknown action", "event", core.LogActionEvent, "frame", c.Core.F, "action", a)
		return 0, 0
	}
}

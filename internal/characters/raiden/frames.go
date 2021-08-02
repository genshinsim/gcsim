package raiden

import "github.com/genshinsim/gsim/pkg/def"

func (c *char) ActionFrames(a def.ActionType, p map[string]int) int {
	switch a {
	case def.ActionAttack:
		f := 0
		if c.Sim.Status("raidenburst") == 0 {
			switch c.NormalCounter {
			//TODO: need to add atkspd mod
			case 0:
				f = 12
			case 1:
				f = 18
			case 2:
				f = 18
			case 3:
				f = 54
			case 4:
				f = 35
			}
		} else {
			switch c.NormalCounter {
			//TODO: need to add atkspd mod
			case 0:
				f = 12
			case 1:
				f = 18
			case 2:
				f = 15
			case 3:
				f = 38
			case 4:
				f = 48
			}
		}
		f = int(float64(f) / (1 + c.Stats[def.AtkSpd]))
		return f
	case def.ActionSkill:
		return 68 //eye appears
	case def.ActionBurst:
		return 139
	default:
		c.Log.Warnw("unknown action", "event", def.LogActionEvent, "frame", c.Sim.Frame(), "action", a)
		return 0
	}
}

package ayaka

import "github.com/genshinsim/gsim/pkg/def"

func (c *char) ActionFrames(a def.ActionType, p map[string]int) int {
	switch a {
	case def.ActionAttack:
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
		f = int(float64(f) / (1 + c.Stats[def.AtkSpd]))
		return f
	case def.ActionCharge:
		return 53
	case def.ActionSkill:
		return 56 //should be 82
	case def.ActionBurst:
		return 95 //ok
	default:
		c.Log.Warnw("unknown action", "event", def.LogActionEvent, "frame", c.Sim.Frame(), "action", a)
		return 0
	}
}

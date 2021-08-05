package albedo

import "github.com/genshinsim/gsim/pkg/def"

func (c *char) ActionFrames(a def.ActionType, p map[string]int) int {
	switch a {
	case def.ActionAttack:
		f := 0
		switch c.NormalCounter {
		//TODO: need to add atkspd mod
		case 0:
			f = 12 //frames from keqing lib
		case 1:
			f = 30 - 12
		case 2:
			f = 59 - 30
		case 3:
			f = 98 - 59
		case 4:
			f = 152 - 98
		}
		f = int(float64(f) / (1 + c.Stats[def.AtkSpd]))
		return f
	case def.ActionCharge:
		return 54 //frames from keqing lib
	case def.ActionSkill:
		return 32
	case def.ActionBurst:
		return 96
	default:
		c.Log.Warnf("%v: unknown action (%v), frames invalid", c.Base.Name, a)
		return 0
	}
}

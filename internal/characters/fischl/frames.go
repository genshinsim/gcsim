package fischl

import "github.com/genshinsim/gsim/pkg/def"

func (c *char) ActionFrames(a def.ActionType, p map[string]int) int {
	switch a {
	case def.ActionAttack:
		f := 0
		switch c.NormalCounter {
		//TODO: need to add atkspd mod
		case 0:
			f = 10 //frames from keqing lib
		case 1:
			f = 28 - 10
		case 2:
			f = 61 - 28
		case 3:
			f = 102 - 61
		case 4:
			f = 131 - 102
		}
		f = int(float64(f) / (1 + c.Stats[def.AtkSpd]))
		return f
	case def.ActionAim:
		return 95
	case def.ActionSkill:
		return 40
	case def.ActionBurst:
		return 21 //TODO: this is swap cancelling
	default:
		c.Log.Warnf("%v: unknown action (%v), frames invalid", c.Base.Name, a)
		return 0
	}
}

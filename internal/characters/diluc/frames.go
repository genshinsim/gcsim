package diluc

import "github.com/genshinsim/gcsim/pkg/core"

func (c *char) ActionFrames(a core.ActionType, p map[string]int) (int, int) {
	switch a {
	case core.ActionAttack:
		f := 0
		switch c.NormalCounter {
		//TODO: need to add atkspd mod
		case 0:
			f = 24 //frames from keqing lib
		case 1:
			f = 77 - 24
		case 2:
			f = 115 - 77
		case 3:
			f = 181 - 115
		}
		f = int(float64(f) / (1 + c.Stats[core.AtkSpd]))
		return f, f
	case core.ActionSkill:
		switch c.eCounter {
		case 1:
			return 52, 52
		case 2:
			return 81, 81
		default:
			return 45, 45
		}
	case core.ActionBurst:
		return 145, 145
	default:
		c.Core.Log.Warnf("%v: unknown action (%v), frames invalid", c.Base.Key.String(), a)
		return 0, 0
	}
}

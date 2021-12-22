package albedo

import "github.com/genshinsim/gcsim/pkg/core"

func (c *char) ActionFrames(a core.ActionType, p map[string]int) (int, int) {
	switch a {
	case core.ActionAttack:
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
		f = int(float64(f) / (1 + c.Stats[core.AtkSpd]))
		return f, f
	case core.ActionCharge:
		return 54, 54 //frames from keqing lib
	case core.ActionSkill:
		return 32, 32
	case core.ActionBurst:
		return 96, 96
	default:
		c.Core.Log.Warnf("%v: unknown action (%v), frames invalid", c.Base.Key.String(), a)
		return 0, 0
	}
}

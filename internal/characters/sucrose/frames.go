package sucrose

import "github.com/genshinsim/gcsim/pkg/core"

func (c *char) ActionFrames(a core.ActionType, p map[string]int) (int, int) {
	switch a {
	case core.ActionAttack:
		f := 0
		switch c.NormalCounter {
		//TODO: need to add atkspd mod
		case 0:
			f = 19 //frames from keqing lib
		case 1:
			f = 38 - 19
		case 2:
			f = 70 - 38
		case 3:
			f = 101 - 70
		}
		f = int(float64(f) / (1 + c.Stat(core.AtkSpd)))
		return f, f
	case core.ActionCharge:
		return 53, 53 //frames from keqing lib
	case core.ActionSkill:
		return 55, 55 //ok
	case core.ActionBurst:
		return 46, 46 //ok
	default:
		c.Core.Log.Warnf("%v: unknown action (%v), frames invalid", c.Base.Key.String(), a)
		return 0, 0
	}
}

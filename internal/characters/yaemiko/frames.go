package yaemiko

import "github.com/genshinsim/gcsim/pkg/core"

func (c *char) ActionFrames(a core.ActionType, p map[string]int) (int, int) {
	switch a {
	case core.ActionAttack:
		f := 0
		switch c.NormalCounter {
		//TODO: need to add atkspd mod
		case 0:
			f = 22 //frames from keqing lib
		case 1:
			f = 35
		case 2:
			f = 71
		case 3:
			f = 101 - 70
		}
		f = int(float64(f) / (1 + c.Stat(core.AtkSpd)))
		return f, f
	case core.ActionCharge:
		return 114, 114 //frames from keqing lib
	case core.ActionSkill:
		return 40, 40 //ok
	case core.ActionBurst:
		return 124, 124 //ok
	default:
		c.Core.Log.Warnf("%v: unknown action (%v), frames invalid", c.Base.Key.String(), a)
		return 0, 0
	}
}

package diona

import "github.com/genshinsim/gcsim/pkg/core"

func (c *char) ActionFrames(a core.ActionType, p map[string]int) (int, int) {
	switch a {
	case core.ActionAttack:
		f := 0
		switch c.NormalCounter {
		//TODO: need to add atkspd mod
		case 0:
			f = 16 //frames from keqing lib
		case 1:
			f = 37 - 16
		case 2:
			f = 67 - 37
		case 3:
			f = 101 - 67
		case 4:
			f = 152 - 101
		}
		f = int(float64(f) / (1 + c.Stats[core.AtkSpd]))
		return f, f
	case core.ActionAim:
		if c.Base.Cons >= 4 && c.Core.Status.Duration("dionaburst") > 0 {
			return 34, 34 //reduced by 60%
		}
		return 84, 84 //kqm
	case core.ActionBurst:
		return 21, 21
	case core.ActionSkill:
		switch p["hold"] {
		case 1:
			return 24, 24
		default:
			return 15, 15
		}
	default:
		c.Core.Log.Warnf("%v: unknown action (%v), frames invalid", c.Base.Key.String(), a)
		return 0, 0
	}
}

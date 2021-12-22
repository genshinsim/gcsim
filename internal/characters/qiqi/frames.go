package qiqi

import "github.com/genshinsim/gcsim/pkg/core"

func (c *char) ActionFrames(a core.ActionType, p map[string]int) (int, int) {
	switch a {
	case core.ActionAttack:
		f := 0
		switch c.NormalCounter {
		//TODO: need to add atkspd mod
		case 0:
			f = 11 //frames from keqing lib
		case 1:
			f = 29 - 11
		case 2:
			f = 71 - 29
		case 3:
			f = 111 - 71
		case 4:
			f = 140 - 111
		}
		atkspd := c.Stats[core.AtkSpd]
		f = int(float64(f) / (1 + atkspd))
		return f, f
	case core.ActionCharge:
		return 100 - 11, 100 - 11 // N1C
	case core.ActionSkill:
		return 57, 57
	case core.ActionBurst:
		return 112, 112
	default:
		c.Core.Log.Warnf("%v: unknown action (%v), frames invalid", c.Base.Key.String(), a)
		return 0, 0
	}
}

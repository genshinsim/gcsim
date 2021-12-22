package mona

import "github.com/genshinsim/gcsim/pkg/core"

func (c *char) ActionFrames(a core.ActionType, p map[string]int) (int, int) {
	switch a {
	case core.ActionAttack:
		f := 0
		switch c.NormalCounter {
		//TODO: need to add atkspd mod
		case 0:
			f = 18
		case 1:
			f = 23
		case 2:
			f = 33
		case 3:
			f = 39
		}
		f = int(float64(f) / (1 + c.Stats[core.AtkSpd]))
		return f, f
	case core.ActionCharge:
		//TODO: need count
		return 50, 50
	case core.ActionSkill:
		//TODO: need count
		return 60, 60
	case core.ActionBurst:
		return 127, 127
	default:
		c.Core.Log.Warnf("%v: unknown action (%v), frames invalid", c.Base.Key.String(), a)
		return 0, 0
	}
}

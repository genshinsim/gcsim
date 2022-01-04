package itto

import "github.com/genshinsim/gcsim/pkg/core"

func (c *char) ActionFrames(a core.ActionType, p map[string]int) (int, int) {
	switch a {
	case core.ActionAttack:
		f := 0
		switch c.NormalCounter {
		//TODO: need to add atkspd mod TODO:these are not accurate
		case 0:
			f = 23
		case 1:
			f = 67 - 23
		case 2:
			f = 101 - 67 - 23
		case 3:
			f = 181 - 101 - 67 - 23
		}
		f = int(float64(f) / (1 + c.Stats[core.AtkSpd]))
		return f, f
	case core.ActionCharge:
		f := 35
		f = int(float64(f) / (1 + c.Stats[core.AtkSpd]))
		return f, f //TODO: no idea
	case core.ActionSkill:
		return 45, 45
	case core.ActionBurst:
		return 145, 145
	default:
		c.Core.Log.Warnf("%v: unknown action (%v), frames invalid", c.Base.Key.String(), a)
		return 0, 0
	}
}

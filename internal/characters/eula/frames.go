package eula

import "github.com/genshinsim/gcsim/pkg/core"

func (c *char) ActionFrames(a core.ActionType, p map[string]int) (int, int) {
	switch a {
	case core.ActionAttack:
		f := 0
		switch c.NormalCounter {
		//TODO: need to add atkspd mod
		case 0:
			f = 29
		case 1:
			f = 25
		case 2:
			f = 65
		case 3:
			f = 33
		case 4:
			f = 88
		}
		f = int(float64(f) / (1 + c.Stats[core.AtkSpd]))
		return f, f
	case core.ActionCharge:
		return 35, 35 //TODO: no idea
	case core.ActionSkill:
		if p["hold"] == 0 {
			return 34, 34
		}
		return 80, 80
	case core.ActionBurst:
		return 116, 116 //ok
	default:
		c.Core.Log.Warnf("%v: unknown action (%v), frames invalid", c.Base.Key.String(), a)
		return 0, 0
	}
}

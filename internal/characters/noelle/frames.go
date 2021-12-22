package noelle

import "github.com/genshinsim/gcsim/pkg/core"

func (c *char) ActionFrames(a core.ActionType, p map[string]int) (int, int) {
	switch a {
	case core.ActionAttack:
		f := 0
		switch c.NormalCounter {
		//TODO: need to add atkspd mod
		case 0:
			f = 28 //frames from keqing lib
		case 1:
			f = 70 - 28
		case 2:
			f = 116 - 70
		case 3:
			f = 174 - 116
		}
		f = int(float64(f) / (1 + c.Stats[core.AtkSpd]))
		return f, f
	case core.ActionSkill:
		return 41, 41 //TODO: not ok
	case core.ActionBurst:
		return 111, 111 //ok
	default:
		c.Core.Log.Warnf("%v: unknown action (%v), frames invalid", c.Base.Key.String(), a)
		return 0, 0
	}
}

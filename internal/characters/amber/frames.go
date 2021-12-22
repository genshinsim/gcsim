package amber

import "github.com/genshinsim/gcsim/pkg/core"

func (c *char) ActionFrames(a core.ActionType, p map[string]int) (int, int) {
	switch a {
	case core.ActionAttack:
		f := 0
		switch c.NormalCounter {
		//TODO: need to add atkspd mod
		case 0:
			f = 15 //frames from keqing lib
		case 1:
			f = 33 - 15
		case 2:
			f = 72 - 33
		case 3:
			f = 113 - 72
		case 4:
			f = 155 - 113
		}
		f = int(float64(f) / (1 + c.Stats[core.AtkSpd]))
		return f, f
	case core.ActionAim:
		return 94, 94 //kqm
	case core.ActionBurst:
		return 74, 74 //swap canceled
	case core.ActionSkill:
		return 35, 35 //no cancel
	default:
		c.Core.Log.Warnf("%v: unknown action (%v), frames invalid", c.Base.Key.String(), a)
		return 0, 0
	}
}

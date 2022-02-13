package thoma

import "github.com/genshinsim/gcsim/pkg/core"

func (c *char) ActionFrames(a core.ActionType, p map[string]int) (int, int) {
	switch a {
	case core.ActionAttack: //not final from KQM
		f := 0
		switch c.NormalCounter {
		case 0:
			f = 11
		case 1:
			f = 49 - 11
		case 2:
			f = 89 - 49 - 11
		case 3:
			f = 114 - 89 - 49 - 11
		}
		f = int(float64(f) / (1 + c.Stats[core.AtkSpd]))
		return f, f
	case core.ActionSkill:
		return 44, 44
	case core.ActionBurst:
		return 56, 56
	case core.ActionCharge:
		return 14 + 56, 14 + 56
	default:
		c.Core.Log.Warnf("%v: unknown action (%v), frames invalid", c.Base.Key.String(), a)
		return 0, 0
	}
}

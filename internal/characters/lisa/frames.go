package lisa

import "github.com/genshinsim/gcsim/pkg/core"

func (c *char) ActionFrames(a core.ActionType, p map[string]int) (int, int) {
	switch a {
	case core.ActionAttack:
		f := 0
		switch c.NormalCounter {
		//TODO: need to add atkspd mod
		case 0:
			f = 25
		case 1:
			f = 46 - 25
		case 2:
			f = 70 - 46
		case 3:
			f = 114 - 70
		}
		f = int(float64(f) / (1 + c.Stats[core.AtkSpd]))
		return f, f
	case core.ActionCharge:
		return 95, 95
	case core.ActionSkill:
		hold := p["hold"]
		if hold == 0 {
			return 21, 21 //no hold
		}
		//yes hold
		return 116, 116
	case core.ActionBurst:
		return 30, 30
	default:
		c.Core.Log.Warnf("%v: unknown action (%v), frames invalid", c.Base.Key.String(), a)
		return 0, 0
	}
}

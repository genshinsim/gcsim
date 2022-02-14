package zhongli

import "github.com/genshinsim/gcsim/pkg/core"

func (c *char) ActionFrames(a core.ActionType, p map[string]int) (int, int) {
	switch a {
	case core.ActionAttack:
		f := 0
		switch c.NormalCounter {
		//TODO: need to add atkspd mod
		case 0:
			f = 8
		case 1:
			f = 29 - 8
		case 2:
			f = 45 - 29
		case 3:
			f = 71 - 45
		case 4:
			f = 109 - 71
		case 5:
			f = 153 - 109
		}
		f = int(float64(f) / (1 + c.Stat(core.AtkSpd)))
		return f, f
	case core.ActionCharge:
		return 95, 95
	case core.ActionSkill:
		hold := p["hold"]
		if hold == 0 {
			//no hold
			return 39, 39
		}
		//yes hold
		return 97, 97
	case core.ActionBurst:
		return 140, 140
	default:
		c.Core.Log.Warnw("unknown action", "event", core.LogActionEvent, "frame", c.Core.F, "action", a)
		return 0, 0
	}
}

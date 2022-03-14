package shenhe

import "github.com/genshinsim/gcsim/pkg/core"

func (c *char) ActionFrames(a core.ActionType, p map[string]int) (int, int) {
	switch a {
	case core.ActionAttack:
		f := 0
		switch c.NormalCounter {
		case 0:
			f = 23
		case 1:
			f = 19
		case 2:
			f = 42
		case 3:
			f = 30
		case 4:
			f = 81
		}
		f = int(float64(f) / (1 + c.Stat(core.AtkSpd)))
		return f, f
	case core.ActionSkill:
		if h := p["hold"]; h == 1 {
			return 44, 44
		}
		return 31, 31
	case core.ActionBurst:
		return 99, 99
	case core.ActionCharge:
		return 49, 49
	default:
		c.coretype.Log.NewEventBuildMsg(core.LogActionEvent, c.Index, "unknown action (invalid frames): ", a.String())
		return 0, 0
	}
}

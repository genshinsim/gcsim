package xiao

import "github.com/genshinsim/gcsim/pkg/core"

func (c *char) ActionFrames(a core.ActionType, p map[string]int) (int, int) {
	switch a {
	case core.ActionAttack:
		f := 0
		a := 0
		switch c.NormalCounter {
		//TODO: need to add atkspd mod
		case 0:
			f = 17
			a = 25
		case 1:
			f = 15
			a = 22
		case 2:
			f = 15
			a = 26
		case 3:
			f = 31
			a = 39
		case 4:
			f = 16
			a = 24
		case 5:
			f = 39
			a = 79
		}
		atkspd := c.Stat(core.AtkSpd)
		f = int(float64(f) / (1 + atkspd))
		return f, a
	case core.ActionCharge:
		return 17, 46
	case core.ActionHighPlunge:
		return 54, 67
	case core.ActionLowPlunge:
		return 49, 66
	case core.ActionSkill:
		return 4, 24
	case core.ActionBurst:
		return 57, 82
	case core.ActionDash:
		return 21, 21
	default:
		c.Core.Log.NewEventBuildMsg(core.LogActionEvent, c.Index, "unknown action (invalid frames): ", a.String())
		return 0, 0
	}
}

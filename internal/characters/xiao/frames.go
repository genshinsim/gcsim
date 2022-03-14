package xiao

import "github.com/genshinsim/gcsim/pkg/core"

func (c *char) ActionFrames(a core.ActionType, p map[string]int) (int, int) {
	switch a {
	case core.ActionAttack:
		f := 0
		switch c.NormalCounter {
		//TODO: need to add atkspd mod
		case 0:
			f = 16 //frames from keqing lib
		case 1:
			f = 42 - 16
		case 2:
			f = 68 - 42
		case 3:
			f = 115 - 68
		case 4:
			f = 145 - 115
		case 5:
			f = 197 - 145
		}
		atkspd := c.Stat(core.AtkSpd)
		f = int(float64(f) / (1 + atkspd))
		return f, f
	case core.ActionCharge:
		return 80 - 16, 80 - 16 // Taken as N1C - N1. Not entirely right due to different recovery frames
	case core.ActionHighPlunge:
		return 75, 75
	case core.ActionLowPlunge:
		return 65, 65
	case core.ActionSkill:
		return 36, 36
	case core.ActionBurst:
		return 58, 58
	default:
		c.coretype.Log.NewEventBuildMsg(core.LogActionEvent, c.Index, "unknown action (invalid frames): ", a.String())
		return 0, 0
	}
}

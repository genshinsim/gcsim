package venti

import "github.com/genshinsim/gcsim/pkg/core"

func (c *char) ActionFrames(a core.ActionType, p map[string]int) (int, int) {
	switch a {
	case core.ActionAttack:
		f := 0
		switch c.NormalCounter {
		//TODO: need to add atkspd mod
		case 0:
			f = 21 //frames from keqing lib
		case 1:
			f = 44 - 21
		case 2:
			f = 90 - 44
		case 3:
			f = 123 - 90
		case 4:
			f = 140 - 123
		case 5:
			f = 191 - 140
		}
		f = int(float64(f) / (1 + c.Stats[core.AtkSpd]))
		return f, f
	case core.ActionAim:
		return 86, 86 //frames from keqing libf
	case core.ActionSkill:
		if p["hold"] == 0 {
			return 38, 38
		}
		return 65, 65
	case core.ActionBurst:
		return 94, 94
	default:
		c.Core.Log.NewEventBuildMsg(core.LogActionEvent, c.Index, "unknown action (invalid frames): ", a.String())
		return 0, 0
	}
}

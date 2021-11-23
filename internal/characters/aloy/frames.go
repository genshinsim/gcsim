package aloy

import "github.com/genshinsim/gcsim/pkg/core"

func (c *char) ActionFrames(a core.ActionType, p map[string]int) (int, int) {
	switch a {
	case core.ActionAttack:
		f := 0
		switch c.NormalCounter {
		//TODO: need to add atkspd mod
		case 0:
			f = 30 //frames from keqing lib
		case 1:
			f = 48 - 30
		case 2:
			f = 85 - 48
		case 3:
			f = 128 - 85
		}
		atkspd := c.Stats[core.AtkSpd]
		f = int(float64(f) / (1 + atkspd))
		return f, f
	case core.ActionAim:
		return 84, 84 //TODO: kqm doesn't have frames lol
	case core.ActionSkill:
		return 67, 67
	case core.ActionBurst:
		return 118, 118
	default:
		c.Core.Log.Warnw("unknown action", "event", core.LogActionEvent, "frame", c.Core.F, "action", a)
		return 0, 0
	}
}

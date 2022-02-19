package yanfei

import "github.com/genshinsim/gcsim/pkg/core"

func (c *char) ActionFrames(a core.ActionType, p map[string]int) (int, int) {
	switch a {
	case core.ActionAttack:
		f := 0
		switch c.NormalCounter {
		//TODO: need to add atkspd mod
		case 0:
			f = 13 //frames from keqing lib
		case 1:
			f = 41 - 13
		case 2:
			f = 90 - 41
		}
		atkspd := c.Stat(core.AtkSpd)
		f = int(float64(f) / (1 + atkspd))
		return f, f
	case core.ActionCharge:
		return 107 - 41, 107 - 41 // Use N2C for now
	case core.ActionSkill:
		return 46, 46
	case core.ActionBurst:
		return 65, 65
	default:
		c.Core.Log.NewEventBuildMsg(core.LogActionEvent, c.Index, "unknown action (invalid frames): ", a.String())
		return 0, 0
	}
}

package rosaria

import "github.com/genshinsim/gcsim/pkg/core"

func (c *char) ActionFrames(a core.ActionType, p map[string]int) (int, int) {
	switch a {
	case core.ActionAttack:
		f := 0
		switch c.NormalCounter {
		//TODO: need to add atkspd mod
		case 0:
			f = 10 //frames from keqing lib
		case 1:
			f = 36 - 10
		case 2:
			f = 81 - 36
		case 3:
			f = 115 - 81
		case 4:
			f = 175 - 115
		}
		atkspd := c.Stat(core.AtkSpd)
		f = int(float64(f) / (1 + atkspd))
		return f, f
	case core.ActionCharge:
		return 89, 89 //frames from keqing libf
	case core.ActionSkill:
		return 65, 65 //ok
	case core.ActionBurst:
		return 74, 74 //ok
	default:
		c.coretype.Log.NewEventBuildMsg(core.LogActionEvent, c.Index, "unknown action (invalid frames): ", a.String())
		return 0, 0
	}
}

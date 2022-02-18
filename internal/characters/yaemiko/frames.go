package yaemiko

import "github.com/genshinsim/gcsim/pkg/core"

func (c *char) ActionFrames(a core.ActionType, p map[string]int) (int, int) {
	switch a {
	case core.ActionAttack:
		f := 0
		switch c.NormalCounter {
		//TODO: need to add atkspd mod
		case 0:
			f = 22 //frames from keqing lib
		case 1:
			f = 46 - 22
		case 2:
			f = 90 - 46 - 22
		}
		f = int(float64(f) / (1 + c.Stat(core.AtkSpd)))
		return f, f
	case core.ActionCharge:
		return 114, 114 //frames from keqing lib
	case core.ActionSkill:
		return 21, 21 //ok
	case core.ActionBurst:
		return 111, 111 //ok
	default:
		c.Core.Log.NewEventBuildMsg(core.LogActionEvent, c.Index, "unknown action (invalid frames): ", a.String())
		return 0, 0
	}
}

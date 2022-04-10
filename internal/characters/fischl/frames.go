package fischl

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
			f = 28 - 10
		case 2:
			f = 61 - 28
		case 3:
			f = 102 - 61
		case 4:
			f = 131 - 102
		}
		f = int(float64(f) / (1 + c.Stat(core.AtkSpd)))
		return f, f
	case core.ActionAim:
		return 95, 95
	case core.ActionSkill:
		return 40, 40
	case core.ActionBurst:
		return 21, 21 //TODO: this is swap cancelling
	default:
		c.Core.Log.NewEventBuildMsg(core.LogActionEvent, c.Index, "unknown action (invalid frames): ", a.String())
		return 0, 0
	}
}

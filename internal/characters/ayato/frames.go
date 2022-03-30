package ayato

import "github.com/genshinsim/gcsim/pkg/core"

func (c *char) ActionFrames(a core.ActionType, p map[string]int) (int, int) {
	switch a {
	case core.ActionAttack:
		f := 0
		if c.Core.Status.Duration("soukaikanka") > 0 {
			f = 24
		} else {
			switch c.NormalCounter {
			//TODO: need to add atkspd mod
			case 0:
				f = 8
			case 1:
				f = 20
			case 2:
				f = 28
			case 3:
				f = 43
			case 4:
				f = 37
			}

			f = int(float64(f) / (1 + c.Stat(core.AtkSpd)))
		}
		return f, f
	case core.ActionCharge:
		return 53, 53
	case core.ActionSkill:
		return 56, 56 //should be 82
	case core.ActionBurst:
		return 95, 95 //ok
	default:
		c.Core.Log.NewEventBuildMsg(core.LogActionEvent, c.Index, "unknown action (invalid frames): ", a.String())
		return 0, 0
	}
}

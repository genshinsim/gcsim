package ayato

import "github.com/genshinsim/gcsim/pkg/core"

func (c *char) ActionFrames(a core.ActionType, p map[string]int) (int, int) {
	switch a {
	case core.ActionAttack:
		f := 0
		if c.Core.Status.Duration("soukaikanka") > 0 {
			switch c.NormalCounter {
			//TODO: need to add atkspd mod
			case 0:
				f = 9
			case 1:
				f = 33 - 9
			case 2:
				f = 57 - 33
			}

			f = int(float64(f) / (1 + c.Stat(core.AtkSpd)))
		} else {
			switch c.NormalCounter {
			//TODO: need to add atkspd mod
			case 0:
				f = 13
			case 1:
				f = 45 - 13
			case 2:
				f = 73 - 45
			case 3:
				f = 124 - 73
			case 4:
				f = 160 - 124
			}

			f = int(float64(f) / (1 + c.Stat(core.AtkSpd)))
		}
		return f, f
	case core.ActionCharge:
		return 53, 53
	case core.ActionSkill:
		return 20, 20 //should be 82
	case core.ActionBurst:
		return 121, 121 //ok
	default:
		c.Core.Log.NewEventBuildMsg(core.LogActionEvent, c.Index, "unknown action (invalid frames): ", a.String())
		return 0, 0
	}
}

package keqing

import "github.com/genshinsim/gcsim/pkg/core"

func (c *char) ActionFrames(a core.ActionType, p map[string]int) (int, int) {
	switch a {
	case core.ActionAttack:
		f := 11
		switch c.NormalCounter {
		//TODO: need to add atkspd mod
		case 1:
			f = 33 - 11
		case 2:
			f = 60 - 33
		case 3:
			f = 97 - 60
		case 4:
			f = 133 - 97
		}
		return f, f
	case core.ActionCharge:
		return 52, 52
	case core.ActionSkill:
		if c.Core.Status.Duration(stilettoKey) > 0 {
			//2nd part
			return 84, 84
		}
		//first part
		return 34, 34
	case core.ActionBurst:
		return 125, 125
	}
	c.Core.Log.NewEventBuildMsg(core.LogActionEvent, c.Index, "unknown action (invalid frames): ", a.String())
	return 0, 0
}

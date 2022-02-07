package kazuha

import (
	"github.com/genshinsim/gcsim/pkg/core"
)

func (c *char) ActionFrames(a core.ActionType, p map[string]int) (int, int) {
	switch a {
	case core.ActionAttack:
		f := 0
		switch c.NormalCounter {
		//TODO: need to add atkspd mod
		case 0:
			//add frames if last action is also attack
			if c.Core.LastAction.Target == core.Kazuha && c.Core.LastAction.Typ == core.ActionAttack {
				f += 60
			}
			f = 14
		case 1:
			f = 34 - 14
		case 2:
			f = 70 - 34 //hit at 60, 70
		case 3:
			f = 97 - 70
		case 4:
			f = 140 - 97
		}
		f = int(float64(f) / (1 + c.Stat(core.AtkSpd)))
		return f, f
	case core.ActionCharge:
		return 44, 44 // kqm lib
	case core.ActionHighPlunge:
		c.Core.Log.NewEvent("plunge skill check", core.LogCharacterEvent, c.Index, "previous", c.Core.LastAction)
		if c.Core.LastAction.Target == core.Kazuha && c.Core.LastAction.Typ == core.ActionSkill {
			_, ok := c.Core.LastAction.Param["hold"]
			if ok {
				return 41, 41
			}
			return 36, 36
		}
		c.Core.Log.NewEvent("invalid plunge (missing skill use)", core.LogActionEvent, c.Index, "action", a)
		return 0, 0
	case core.ActionSkill:
		_, ok := p["hold"]
		if ok {
			return 58, 58
		}
		return 27, 27
	case core.ActionBurst:
		return 95, 95
	default:
		c.Core.Log.NewEventBuildMsg(core.LogActionEvent, c.Index, "unknown action (invalid frames): ", a.String())
		return 0, 0
	}
}

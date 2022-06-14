package barbara

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

func (c *char) ActionFrames(a action.Action, p map[string]int) (int, int) {
	/*
		Source: https://library.keqingmains.com/characters/hydro/barbara
	*/
	switch a {
	case action.ActionAttack:
		f := 0
		switch c.NormalCounter {
		//TODO: need to add atkspd mod
		case 0:
			f = 7 //frames from keqing lib
		case 1:
			f = 25 - 7
		case 2:
			f = 45 - 25 - 7
		case 3:
			f = 92 - 45 - 25 - 7
		}
		atkspd := c.Stat(attributes.AtkSpd)
		f = int(float64(f) / (1 + atkspd))
		return f, f
	case action.ActionCharge:
		return 90, 90
	case action.ActionSkill:
		return 52, 52
	case action.ActionBurst:
		return 110, 110
	default:
		c.Core.Log.NewEventBuildMsg(glog.LogActionEvent, c.Index, "unknown action (invalid frames): ", a.String())
		return 0, 0
	}
}

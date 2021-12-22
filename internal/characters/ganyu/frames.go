package ganyu

import "github.com/genshinsim/gcsim/pkg/core"

func (c *char) ActionFrames(a core.ActionType, p map[string]int) (int, int) {
	switch a {
	case core.ActionAttack:
		f := 0
		switch c.NormalCounter {
		//TODO: need to add atkspd mod
		case 0:
			f = 18 //frames from keqing lib
		case 1:
			f = 43 - 18
		case 2:
			f = 73 - 43
		case 3:
			f = 117 - 73
		case 4:
			f = 153 - 117
		case 5:
			f = 190 - 153
		}
		f = int(float64(f) / (1 + c.Stats[core.AtkSpd]))
		return f, f
	case core.ActionAim:
		//check for c6, if active then return 10, otherwise 115
		if c.Core.Status.Duration("ganyuc6") > 0 {
			c.Core.Log.Debugw("ganyu c6 proc used", "frame", c.Core.F, "event", core.LogCharacterEvent, "char", c.Index)
			c.Core.Status.DeleteStatus("ganyuc6")
			return 10, 10
		}
		return 115, 115 //frames from keqing lib
	case core.ActionSkill:
		return 30, 30 //ok
	case core.ActionBurst:
		return 122, 122 //ok
	default:
		c.Core.Log.Warnf("%v: unknown action (%v), frames invalid", c.Base.Key.String(), a)
		return 0, 0
	}
}

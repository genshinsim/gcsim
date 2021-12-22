package kokomi

import "github.com/genshinsim/gcsim/pkg/core"

func (c *char) ActionFrames(a core.ActionType, p map[string]int) (int, int) {
	switch a {
	case core.ActionAttack:
		f := 0
		recovery := 0
		switch c.NormalCounter {
		//TODO: need to add atkspd mod
		case 0:
			f = 10 //frames from keqing lib
			recovery = f
		case 1:
			f = 36 - 10
			recovery = f
		case 2:
			f = 82 - 10
			recovery = 122 - 82
		}
		atkspd := c.Stats[core.AtkSpd]
		f = int(float64(f) / (1 + atkspd))
		return f, recovery
	case core.ActionCharge:
		return 45, 45
	case core.ActionSkill:
		// Took dash cancel frames for now
		return 51, 51
	case core.ActionBurst:
		return 75, 75
	default:
		c.Core.Log.Warnf("%v: unknown action (%v), frames invalid", c.Base.Key.String(), a)
		return 0, 0
	}
}

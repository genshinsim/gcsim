package travelerelectro

import "github.com/genshinsim/gcsim/pkg/core"

func (c *char) ActionFrames(a core.ActionType, p map[string]int) (int, int) {
	switch a {
	case core.ActionAttack:
		f := 0
		switch c.NormalCounter {
		//TODO: need to add atkspd mod
		case 0:
			f = 13 //frames from keqing lib
		case 1:
			f = 38 - 13
		case 2:
			f = 71 - 38
		case 3:
			f = 123 - 71
		case 4:
			f = 163 - 123
		}
		f = int(float64(f) / (1 + c.Stats[core.AtkSpd]))
		return f, f
	case core.ActionCharge:
		return 54 - 13, 53 - 13
	case core.ActionSkill:
		return 55, 55 //could be 52 if going into Q
	case core.ActionBurst:
		return 60, 60 //1573 start, 1610 cd starts, 1612 energy drained, 1633 first swapable
	default:
		c.Core.Log.Warnf("%v: unknown action, frames invalid", a)
		return 0, 0
	}
}

func (c *char) ActionStam(a core.ActionType, p map[string]int) float64 {
	switch a {
	case core.ActionDash:
		return 18
	case core.ActionCharge:
		return 25
	default:
		c.Core.Log.Warnf("%v ActionStam for %v not implemented; Character stam usage may be incorrect", c.Base.Key.String(), a.String())
		return 0
	}
}

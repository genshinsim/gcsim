package yunjin

import "github.com/genshinsim/gcsim/pkg/core"

// TODO: Currently uses beta frame counts. Need to update once we have final
func (c *char) ActionFrames(a core.ActionType, p map[string]int) (int, int) {
	switch a {
	case core.ActionAttack:
		f := 0
		switch c.NormalCounter {
		//TODO: need to add atkspd mod
		case 0:
			f = 24
		case 1:
			f = 22
		case 2:
			f = 28
		case 3:
			f = 33
		case 4:
			f = 39
		}
		atkspd := c.Stat(core.AtkSpd)
		f = int(float64(f) / (1 + atkspd))
		return f, f
	case core.ActionCharge:
		return 66, 66
	case core.ActionSkill:
		chargeLevel := p["hold"]
		// If perfect timing, assume press frames
		if p["perfect"] == 1 {
			chargeLevel = 0
		}
		switch chargeLevel {
		case 0:
			return 31, 31
		case 1:
			return 81, 81
		case 2:
			return 121, 121
		}
		return 0, 0
	case core.ActionBurst:
		return 53, 53
	default:
		c.Core.Log.Warnf("%v: unknown action (%v), frames invalid", c.Base.Key.String(), a)
		return 0, 0
	}
}

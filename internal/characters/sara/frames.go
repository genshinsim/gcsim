package sara

import "github.com/genshinsim/gcsim/pkg/core"

func (c *char) ActionFrames(a core.ActionType, p map[string]int) (int, int) {
	switch a {
	case core.ActionAttack:
		f := 0
		switch c.NormalCounter {
		//TODO: need to add atkspd mod
		case 0:
			f = 19 //frames from keqing lib
		case 1:
			f = 44 - 19
		case 2:
			f = 82 - 44
		case 3:
			f = 123 - 82
		case 4:
			f = 181 - 123
		}
		atkspd := c.Stats[core.AtkSpd]
		f = int(float64(f) / (1 + atkspd))
		return f, f
	case core.ActionAim:
		// A1 check:
		// While in the Crowfeather Cover state provided by Tengu Stormcall, Aimed Shot charge times are decreased by 60%.
		// TODO: Maybe not exactly right since some component of this is not the charge time
		// Nothing better in library yet though
		if c.Core.Status.Duration("saracover") > 0 {
			reduced_frames := 78 * 0.4
			return int(reduced_frames), int(reduced_frames)
		}
		return 78, 78
	case core.ActionSkill:
		return 50, 50
	case core.ActionBurst:
		// In line with most other cases in sim, assume that you swap cancel this instead of the full 80 frames
		return 60, 60
	default:
		c.Core.Log.Warnf("%v: unknown action (%v), frames invalid", c.Base.Key.String(), a)
		return 0, 0
	}
}

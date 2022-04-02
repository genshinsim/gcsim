package ayato

import "github.com/genshinsim/gcsim/pkg/core"

func (c *char) ActionFrames(a core.ActionType, p map[string]int) (int, int) {
	switch a {
	case core.ActionAttack:
		f := 0
		if c.Core.Status.Duration("soukaikanka") > 0 {
			switch c.NormalCounter {
			case 0:
				f = 5
			case 1:
				f = 5
			case 2:
				f = 5
			}
			f = int(float64(f) / (1 + c.Stat(core.AtkSpd)))

			// If the soukaikanka buff would expire during the normal attack, extend the buff so it expires immediately after instead.
			if c.Core.Status.Duration("soukaikanka") <= f {
				c.Core.Status.AddStatus("soukaikanka", f+1)
				c.Core.Log.NewEvent("Soukai Kanka extended", core.LogCharacterEvent, c.Index, "expiry", c.Core.F+f+1)
			}
			return f, f
		}

		switch c.NormalCounter {
		//TODO: need to add atkspd mod
		case 0:
			f = 13
		case 1:
			f = 43 - 13
		case 2:
			f = 73 - 43
		case 3:
			f = 111 - 73
		case 4:
			f = 159 - 111
		}
		f = int(float64(f) / (1 + c.Stat(core.AtkSpd)))

		return f, f
	case core.ActionCharge:
		return 53, 53
	case core.ActionSkill:
		return 20, 20
	case core.ActionBurst:
		return 125, 125
	default:
		c.Core.Log.NewEventBuildMsg(core.LogActionEvent, c.Index, "unknown action (invalid frames): ", a.String())
		return 0, 0
	}
}

func (c *char) InitCancelFrames() {

}

func (c *char) ActionInterruptableDelay(next core.ActionType, p map[string]int) int {
	// Provide a custom override for differentiating Ayato's E stance attacks from his usual attacks
	default_val := c.Tmpl.ActionInterruptableDelay(next, p)

	prev := c.Core.LastAction
	switch prev.Typ {
	case core.ActionAttack:
		if c.Core.Status.Duration("soukaikanka") > 0 {
			// In E stance
			switch next {
			case core.ActionAttack:
				f := 0
				switch c.NormalCounter {
				case 0:
					// The sim loses 1 frame for executing the attack. Because these attacks are buffered, we should compensate for that
					// The sim loses 1 frame for the animation delay. We should compensate for that here
					// Final results is that sim takes 26f from slash to slash instead of 24f
					f = 24 - 5
				case 1:
					f = 24 - 5
				case 2:
					f = 24 - 5
				}
				f = int(float64(f) / (1 + c.Stat(core.AtkSpd)))
				return f
			case core.ActionBurst:
				return 0
			case core.ActionDash:
				return 0
			case core.ActionJump:
				return 0
			case core.ActionSwap:
				return 0
			case core.ActionSkill: // I didn't actually test this so it's default for now
				return default_val
			}
		} else {
			// not in E stance
			switch next {
			case core.ActionAttack:
				f := 0
				switch c.NormalCounter {
				case 0:
					// N5 -> N1 frames
					// recovery frames from KQM discord
					f = 191 - 159
				case 1:
				case 2:
				case 3:
				case 4:
					f = 0
				}
				return f
			}
		}
	case core.ActionCharge:
		//  charge attack won't be in E stance
		switch next {
		case core.ActionAttack:
			// CA recovery frames from KQM discord
			return 84
		}
		// everything else is default because i didn't see frames for it
	}

	return default_val
}

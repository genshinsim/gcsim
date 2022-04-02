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
	// CA recovery frames from KQM discord
	c.SetAbilCancelFrames(core.ActionCharge, core.ActionAttack, 84) //charge -> n1
	c.SetNormalCancelFrames(4, core.ActionAttack, 191-159)          //n5 -> next attack (n1), recovery frames from KQM discord
}

func (c *char) ActionInterruptableDelay(next core.ActionType, p map[string]int) int {
	// Provide a custom override for differentiating Ayato's E stance attacks from his usual attacks
	// we only over ride if prev is attack and next is also attack
	if c.Core.LastAction.Typ == core.ActionAttack &&
		next == core.ActionAttack &&
		c.Core.Status.Duration("soukaikanka") > 0 {
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
	}
	//otherise use default implementation
	return c.Tmpl.ActionInterruptableDelay(next, p)
}

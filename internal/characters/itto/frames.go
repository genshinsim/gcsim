package itto

import (
	"github.com/genshinsim/gcsim/pkg/core"
)

func (c *char) ActionFrames(a core.ActionType, p map[string]int) (int, int) {
	switch a {
	case core.ActionAttack:
		f := 0
		hl := 0
		if c.Core.Status.Duration("ittoq") <= 0 {
			// Values outside of burst
			//TODO: need to add atkspd mod
			switch c.NormalCounter {
			case 0:
				f = 22 //frames from https://docs.google.com/spreadsheets/d/1nGl_oapCppBoCYnXnEvhkrUeQZ5VAKD2/edit#gid=1479101218
				hl = 9
			case 1:
				if c.dasshuUsed {
					f = 28
				} else {
					f = 67 - 22
				}
				hl = 9
			case 2:
				if c.dasshuUsed {
					f = 16
				} else {
					f = 102 - 67
				}
				hl = 10
			case 3:
				if c.dasshuUsed {
					f = 48
				} else {
					f = 187 - 102
				}
				hl = 10
			}
		} else {
			//Values inside of burst
			switch c.NormalCounter {
			case 0:
				f = 21
				hl = 9
			case 1:
				if c.dasshuUsed {
					f = 26
				} else {
					f = 63 - 21
				}
				hl = 9
			case 2:
				if c.dasshuUsed {
					f = 14
				} else {
					f = 97 - 63
				}
				hl = 10
			case 3:
				if c.dasshuUsed {
					f = 43
				} else {
					f = 175 - 97
				}
				hl = 10
			}
		}
		c.dasshuUsed = false
		f = int(float64(f) / (1 + c.Stats[core.AtkSpd]))
		return f, f + (hl * 0)
	case core.ActionCharge:
		f := 40
		hl := 0
		switch c.Tags["strStack"] {
		case 0:
			return 180, 180 // No idea, but it's forever
		case 1:
			// Last stack is a plunge, atk spd doesn't affect it much
			hl = 10
			if c.sCACount == 0 {
				// Add blink time, assumed used out of ushi
				f = 40 + 12
			}
			if c.sCACount < 3 {
				// Add endlag from previous swing
				f = 40 + 6 - c.sCACount
			} else {
				f = 36
			}
		default:
			switch c.sCACount {
			case 0:
				// Add blink time
				f = 17 + 24
				hl = 8
			case 1:
				// Endlag from previous swing
				f = 22 + 6
				hl = 7
			case 2:
				f = 19 + 5
				hl = 7
			default:
				// Odd swings are slightly faster
				f = 18 + 5 - (c.sCACount % 2)
				hl = 6
			}
		}
		f = int(float64(f) / (1 + c.Stats[core.AtkSpd]))
		return f, f + (hl * 0)
	case core.ActionSkill:
		// Ushi is 17, 32, but endlag isn't being handled?
		return 32, 32
	case core.ActionBurst:
		return 91, 91
	default:
		c.Core.Log.NewEventBuildMsg(core.LogActionEvent, c.Index, "unknown action (invalid frames): ", a.String())
		return 0, 0
	}
}

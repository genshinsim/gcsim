package yelan

import "github.com/genshinsim/gcsim/pkg/core"

func (c *char) ActionFrames(a core.ActionType, p map[string]int) (int, int) {
	switch a {
	case core.ActionAttack:
		f := 0
		switch c.NormalCounter {
		//TODO: need to add atkspd mod
		case 0:
			f = 17 //frames from keqing lib
		case 1:
			f = 20
		case 2:
			f = 39
		case 3:
			f = 74
		}
		f = int(float64(f) / (1 + c.Stat(core.AtkSpd)))
		return f, f
	case core.ActionAim:
		if c.Core.Status.Duration("breakthrough") > 0 {//Reduce required by 80% time if she has breakthrough
			reduced_frames := 34
			return int(reduced_frames), int(reduced_frames)
		}
		return 74, 74 //kqm
	case core.ActionBurst:
		return 105, 105 //??
	case core.ActionSkill:
		return 50, 50 //tap E
	default:
		c.Core.Log.NewEventBuildMsg(core.LogActionEvent, c.Index, "unknown action (invalid frames): ", a.String())
		return 0, 0
	}
}

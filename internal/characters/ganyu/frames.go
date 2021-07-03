package ganyu

import "github.com/genshinsim/gsim/pkg/def"

func (c *char) ActionFrames(a def.ActionType, p map[string]int) int {
	switch a {
	case def.ActionAttack:
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
		f = int(float64(f) / (1 + c.Stats[def.AtkSpd]))
		return f
	case def.ActionAim:
		//check for c6, if active then return 10, otherwise 115
		if c.Sim.Status("ganyuc6") > 0 {
			c.Log.Debugw("ganyu c6 proc used", "frame", c.Sim.Frame(), "event", def.LogCharacterEvent, "char", c.Index)
			c.Sim.DeleteStatus("ganyuc6")
			return 10
		}
		return 115 //frames from keqing lib
	case def.ActionSkill:
		return 30 //ok
	case def.ActionBurst:
		return 122 //ok
	default:
		c.Log.Warnf("%v: unknown action (%v), frames invalid", c.Base.Name, a)
		return 0
	}
}

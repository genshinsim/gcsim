package thunderingfury

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterSetFunc("viridescent venerer", New)
	core.RegisterSetFunc("viridescentvenerer", New)
}

func New(c core.Character, s *core.Core, count int) {
	if count >= 2 {
		m := make([]float64, core.EndStatType)
		m[core.AnemoP] = 0.15
		c.AddMod(core.CharStatMod{
			Key: "vv-2pc",
			Amount: func(a core.AttackTag) ([]float64, bool) {
				return m, true
			},
			Expiry: -1,
		})
	}
	if count >= 4 {
		s.Events.Subscribe(core.OnTransReaction, func(args ...interface{}) bool {
			ds := args[1].(*core.Snapshot)
			t := args[0].(core.Target)
			//ignore if source not current char
			if ds.ActorIndex != c.CharIndex() {
				return false
			}

			var ele core.EleType
			var key string
			switch ds.ReactionType {
			case core.SwirlCryo:
				ele = core.Cryo
				key = "vvcryo"
			case core.SwirlElectro:
				ele = core.Electro
				key = "vvelectro"
			case core.SwirlPyro:
				ele = core.Pyro
				key = "vvpyro"
			case core.SwirlHydro:
				ele = core.Hydro
				key = "vvhydro"
			default:
				return false
			}

			ds.ReactBonus += 0.6

			//ignore if character not on field
			if s.ActiveChar != c.CharIndex() {
				return false
			}

			t.AddResMod(key, core.ResistMod{
				Duration: 600, //10 seconds
				Ele:      ele,
				Value:    -0.4,
			})

			s.Log.Debugw("vv 4pc proc", "frame", s.F, "event", core.LogArtifactEvent, "reaction", ds.ReactionType, "char", c.CharIndex())
			return false
		}, fmt.Sprintf("vv4-%v", c.Name()))

	}
	//add flat stat to char
}

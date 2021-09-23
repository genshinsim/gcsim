package thunderingfury

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	core.RegisterSetFunc("viridescent venerer", New)
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
			if ds.ActorIndex != c.CharIndex() {
				return false
			}
			// s.Log.Println("ok")
			// s.Log.Println(ds.ReactionType)
			switch ds.ReactionType {
			case core.SwirlCryo:
				t.AddResMod("vvcryo", core.ResistMod{
					Duration: 600, //10 seconds
					Ele:      core.Cryo,
					Value:    -0.4,
				})
				// s.Log.Println(t.HasResMod("vvcryo"))
			case core.SwirlElectro:
				t.AddResMod("vvelectro", core.ResistMod{
					Duration: 600, //10 seconds
					Ele:      core.Electro,
					Value:    -0.4,
				})
			case core.SwirlPyro:
				t.AddResMod("vvpyro", core.ResistMod{
					Duration: 600, //10 seconds
					Ele:      core.Pyro,
					Value:    -0.4,
				})
			case core.SwirlHydro:
				t.AddResMod("vvhydro", core.ResistMod{
					Duration: 600, //10 seconds
					Ele:      core.Hydro,
					Value:    -0.4,
				})
			default:
				return false
			}
			ds.ReactBonus += 0.6
			s.Log.Debugw("vv 4pc proc", "frame", s.F, "event", core.LogArtifactEvent, "reaction", ds.ReactionType, "char", c.CharIndex())
			return false
		}, fmt.Sprintf("vv4-%v", c.Name()))

	}
	//add flat stat to char
}

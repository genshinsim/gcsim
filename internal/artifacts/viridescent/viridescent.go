package thunderingfury

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	combat.RegisterSetFunc("viridescent venerer", New)
}

func New(c core.Character, s core.Sim, logger core.Logger, count int) {
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
		s.AddOnTransReaction(func(t core.Target, ds *core.Snapshot) {
			// log.Println(ds)
			if ds.ActorIndex != c.CharIndex() {
				return
			}
			// log.Println("ok")
			// log.Println(ds.ReactionType)
			switch ds.ReactionType {
			case core.SwirlCryo:
				t.AddResMod("vvcryo", core.ResistMod{
					Duration: 600, //10 seconds
					Ele:      core.Cryo,
					Value:    -0.4,
				})
				// log.Println(t.HasResMod("vvcryo"))
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
				return
			}
			ds.ReactBonus += 0.6
			logger.Debugw("vv 4pc proc", "frame", s.Frame(), "event", core.LogArtifactEvent, "reaction", ds.ReactionType, "char", c.CharIndex())

		}, fmt.Sprintf("vv4-%v", c.Name()))

	}
	//add flat stat to char
}

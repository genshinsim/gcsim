package thunderingfury

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/def"
)

func init() {
	combat.RegisterSetFunc("viridescent venerer", New)
}

func New(c def.Character, s def.Sim, log def.Logger, count int) {
	if count >= 2 {
		m := make([]float64, def.EndStatType)
		m[def.AnemoP] = 0.15
		c.AddMod(def.CharStatMod{
			Key: "vv-2pc",
			Amount: func(a def.AttackTag) ([]float64, bool) {
				return m, true
			},
			Expiry: -1,
		})
	}
	if count >= 4 {
		s.AddOnTransReaction(func(t def.Target, ds *def.Snapshot) {
			if ds.ActorIndex != c.CharIndex() {
				return
			}
			switch ds.ReactionType {
			case def.SwirlCryo:
				t.AddResMod("vvcryo", def.ResistMod{
					Duration: 600, //10 seconds
					Ele:      def.Cryo,
					Value:    -0.4,
				})
			case def.SwirlElectro:
				t.AddResMod("vvelectro", def.ResistMod{
					Duration: 600, //10 seconds
					Ele:      def.Electro,
					Value:    -0.4,
				})
			case def.SwirlPyro:
				t.AddResMod("vvpyro", def.ResistMod{
					Duration: 600, //10 seconds
					Ele:      def.Pyro,
					Value:    -0.4,
				})
			case def.SwirlHydro:
				t.AddResMod("vvhydro", def.ResistMod{
					Duration: 600, //10 seconds
					Ele:      def.Hydro,
					Value:    -0.4,
				})
			default:
				return
			}
			ds.ReactBonus += 0.6
			log.Debugw("vv 4pc proc", "frame", s.Frame(), "event", def.LogArtifactEvent, "reaction", ds.ReactionType, "char", c.CharIndex())

		}, fmt.Sprintf("vv4-%v", c.Name()))

	}
	//add flat stat to char
}

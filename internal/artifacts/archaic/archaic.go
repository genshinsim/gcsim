package archaic

import (
	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	core.RegisterSetFunc("archaic petra", New)
	core.RegisterSetFunc("archaicpetra", New)
}

func New(c core.Character, s *core.Core, count int) {
	if count >= 2 {
		m := make([]float64, core.EndStatType)
		m[core.GeoP] = 0.15
		c.AddMod(core.CharStatMod{
			Key: "archaic-2pc",
			Amount: func(a core.AttackTag) ([]float64, bool) {
				return m, true
			},
			Expiry: -1,
		})
	}
	if count >= 4 {

		s.Events.Subscribe(core.OnShielded, func(args ...interface{}) bool {
			// Character that picks it up must be the petra set holder
			if s.ActiveChar != c.CharIndex() {
				return false
			}
			//check shield
			shd := s.Shields.Get(core.ShieldCrystallize)
			if shd != nil {
				//activate
				s.Status.AddStatus("archaic", 600)
				s.Log.Debugw("archaic petra proc'd", "frame", s.F, "event", core.LogArtifactEvent, "char", c.CharIndex(), "ele", shd.Element())

				// Apply mod to all characters
				for _, char := range s.Chars {
					char.AddMod(core.CharStatMod{
						Key: "archaic-4pc",
						Amount: func(ds core.AttackTag) ([]float64, bool) {
							if s.Status.Duration("archaic") == 0 {
								return nil, false
							}

							bonus := core.EleToDmgP(core.EleType(shd.Element()))
							m := make([]float64, core.EndStatType)
							m[bonus] = 0.35
							return m, true
						},
						Expiry: s.F + 600,
					})
				}
			}
			return false
		}, "archaic"+c.Name())
	}
}

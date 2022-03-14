package archaic

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

func init() {
	core.RegisterSetFunc("archaic petra", New)
	core.RegisterSetFunc("archaicpetra", New)
}

func New(c coretype.Character, s *core.Core, count int, params map[string]int) {
	if count >= 2 {
		m := make([]float64, core.EndStatType)
		m[core.GeoP] = 0.15
		c.AddMod(coretype.CharStatMod{
			Key: "archaic-2pc",
			Amount: func() ([]float64, bool) {
				return m, true
			},
			Expiry: -1,
		})
	}
	if count >= 4 {

		var ele coretype.StatType
		m := make([]float64, core.EndStatType)

		s.Subscribe(coretype.OnShielded, func(args ...interface{}) bool {
			// Character that picks it up must be the petra set holder
			if s.Player.ActiveChar != c.Index() {
				return false
			}
			//check shield
			shd := s.Player.GetShield(coretype.ShieldCrystallize)
			if shd != nil {
				//activate
				s.AddStatus("archaic", 600)
				s.Log.NewEvent("archaic petra proc'd", coretype.LogArtifactEvent, c.Index(), "ele", shd.Element())
				m[core.PyroP] = 0
				m[core.HydroP] = 0
				m[coretype.CryoP] = 0
				m[core.ElectroP] = 0
				m[core.AnemoP] = 0
				m[core.GeoP] = 0
				m[core.DendroP] = 0
				ele = coretype.EleToDmgP(coretype.EleType(shd.Element()))
				m[ele] = 0.35

				// Apply mod to all characters
				for _, char := range s.Player.Chars {
					char.AddMod(coretype.CharStatMod{
						Key: "archaic-4pc",
						Amount: func() ([]float64, bool) {
							if s.StatusDuration("archaic") == 0 {
								return nil, false
							}
							return m, true
						},
						Expiry: s.Frame + 600,
					})
				}
			}
			return false
		}, "archaic"+c.Key().String())
	}
}

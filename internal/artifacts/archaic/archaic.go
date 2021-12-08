package archaic

import (
	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterSetFunc("archaic petra", New)
	core.RegisterSetFunc("archaicpetra", New)
}

func New(c core.Character, s *core.Core, count int) {
	if count >= 2 {
		var m [core.EndStatType]float64
		m[core.GeoP] = 0.15
		c.AddMod(core.CharStatMod{
			Key: "archaic-2pc",
			Amount: func(a core.AttackTag) ([core.EndStatType]float64, bool) {
				return m, true
			},
			Expiry: -1,
		})
	}
	if count >= 4 {

		var ele core.StatType
		var m [core.EndStatType]float64

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
				m[core.PyroP] = 0
				m[core.HydroP] = 0
				m[core.CryoP] = 0
				m[core.ElectroP] = 0
				m[core.AnemoP] = 0
				m[core.GeoP] = 0
				m[core.DendroP] = 0
				ele = core.EleToDmgP(core.EleType(shd.Element()))
				m[ele] = 0.35

				// Apply mod to all characters
				for _, char := range s.Chars {
					char.AddMod(core.CharStatMod{
						Key: "archaic-4pc",
						Amount: func(ds core.AttackTag) ([core.EndStatType]float64, bool) {
							if s.Status.Duration("archaic") == 0 {
								return m, false
							}
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

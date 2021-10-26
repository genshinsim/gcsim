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
			if s.ActiveChar != c.CharIndex() {
				return false
			}
			//check shield
			shd := s.Shields.Get(core.ShieldCrystallize)
			if shd != nil {
				//activate
				s.Status.AddStatus("archaic", 600)
				s.SetCustomFlag("archaic", int(shd.Element()))
				s.Log.Debugw("archaic petra proc'd", "frame", s.F, "event", core.LogArtifactEvent, "char", c.CharIndex(), "ele", shd.Element())
			}

			return false
		}, "archaic")

		c.AddMod(core.CharStatMod{
			Key: "archaic-4pc",
			Amount: func(ds core.AttackTag) ([]float64, bool) {
				if s.Status.Duration("archaic") == 0 {
					return nil, false
				}
				ele, ok := s.GetCustomFlag("archaic")
				if !ok {
					return nil, false
				}

				bonus := core.EleToDmgP(core.EleType(ele))
				m := make([]float64, core.EndStatType)
				m[bonus] = 0.35
				s.Log.Debugw("archaic petra bonus", "frame", s.F, "event", core.LogSnapshotEvent, "char", c.CharIndex(), "ele", bonus)
				return m, true
			},
			Expiry: -1,
		})
	}
	//add flat stat to char
}

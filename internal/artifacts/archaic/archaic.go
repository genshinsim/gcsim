package archaic

import (
	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/def"
)

func init() {
	combat.RegisterSetFunc("archaic petra", New)
}

func New(c def.Character, s def.Sim, log def.Logger, count int) {
	if count >= 2 {
		m := make([]float64, def.EndStatType)
		m[def.GeoP] = 0.15
		c.AddMod(def.CharStatMod{
			Key: "archaic-2pc",
			Amount: func(a def.AttackTag) ([]float64, bool) {
				return m, true
			},
			Expiry: -1,
		})
	}
	if count >= 4 {

		s.AddEventHook(func(s def.Sim) bool {
			if s.ActiveCharIndex() != c.CharIndex() {
				return false
			}
			//check shield
			shd := s.GetShield(def.ShieldCrystallize)
			if shd != nil {
				//activate
				s.AddStatus("archaic", 600)
				s.SetCustomFlag("archaic", int(shd.Element()))
				log.Debugw("archaic petra proc'd", "frame", s.Frame(), "event", def.LogArtifactEvent, "char", c.CharIndex(), "ele", shd.Element())
			}

			return false
		}, "archaic", def.PostShieldHook) //ok to overwrite any other char's

		c.AddMod(def.CharStatMod{
			Key: "archaic-4pc",
			Amount: func(ds def.AttackTag) ([]float64, bool) {
				if s.Status("archaic") == 0 {
					return nil, false
				}
				ele, ok := s.GetCustomFlag("archaic")
				if !ok {
					return nil, false
				}

				bonus := def.EleToDmgP(def.EleType(ele))
				m := make([]float64, def.EndStatType)
				m[bonus] = 0.35
				log.Debugw("archaic petra bonus", "frame", s.Frame(), "event", def.LogSnapshotEvent, "char", c.CharIndex(), "ele", bonus)
				return m, true
			},
			Expiry: -1,
		})
	}
	//add flat stat to char
}

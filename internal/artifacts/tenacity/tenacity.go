package tenacity

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/def"
)

func init() {
	combat.RegisterSetFunc("tenacity of millelith", New)
}

func New(c def.Character, s def.Sim, log def.Logger, count int) {
	if count >= 2 {
		m := make([]float64, def.EndStatType)
		m[def.HPP] = 0.2
		c.AddMod(def.CharStatMod{
			Key: "tom-2pc",
			Amount: func(a def.AttackTag) ([]float64, bool) {
				return m, true
			},
			Expiry: -1,
		})
	}
	if count >= 4 {
		icd := 0
		m := make([]float64, def.EndStatType)
		m[def.ATKP] = 0.2

		s.AddOnAttackLanded(func(t def.Target, ds *def.Snapshot, dmg float64, crit bool) {
			if ds.ActorIndex != c.CharIndex() {
				return
			}
			if ds.AttackTag != def.AttackTagElementalArt {
				return
			}
			if icd > s.Frame() {
				return
			}
			s.AddStatus("tom-proc", 180)
			icd = s.Frame() + 30 //.5 second icd

			log.Debugw("tom 4pc proc", "frame", s.Frame(), "event", def.LogArtifactEvent, "expiry", s.Frame()+180, "icd", s.Frame()+30)
		}, fmt.Sprintf("pf4-%v", c.Name()))

		c.AddMod(def.CharStatMod{
			Key: "pf-4pc",
			Amount: func(a def.AttackTag) ([]float64, bool) {
				if s.Status("tom-proc") == 0 {
					return nil, false
				}
				return m, true
			},
			Expiry: -1,
		})

	}
	//add flat stat to char
}

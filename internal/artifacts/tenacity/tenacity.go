package tenacity

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	combat.RegisterSetFunc("tenacity of millelith", New)
}

func New(c core.Character, s core.Sim, log core.Logger, count int) {
	if count >= 2 {
		m := make([]float64, core.EndStatType)
		m[core.HPP] = 0.2
		c.AddMod(core.CharStatMod{
			Key: "tom-2pc",
			Amount: func(a core.AttackTag) ([]float64, bool) {
				return m, true
			},
			Expiry: -1,
		})
	}
	if count >= 4 {
		icd := 0
		m := make([]float64, core.EndStatType)
		m[core.ATKP] = 0.2

		s.AddOnAttackLanded(func(t core.Target, ds *core.Snapshot, dmg float64, crit bool) {
			if ds.ActorIndex != c.CharIndex() {
				return
			}
			if ds.AttackTag != core.AttackTagElementalArt {
				return
			}
			if icd > s.Frame() {
				return
			}
			s.AddStatus("tom-proc", 180)
			icd = s.Frame() + 30 //.5 second icd

			log.Debugw("tom 4pc proc", "frame", s.Frame(), "event", core.LogArtifactEvent, "expiry", s.Frame()+180, "icd", s.Frame()+30)
		}, fmt.Sprintf("pf4-%v", c.Name()))

		c.AddMod(core.CharStatMod{
			Key: "pf-4pc",
			Amount: func(a core.AttackTag) ([]float64, bool) {
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

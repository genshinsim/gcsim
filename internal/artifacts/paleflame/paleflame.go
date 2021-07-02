package paleflame

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/def"
)

func init() {
	combat.RegisterSetFunc("pale flame", New)
}

func New(c def.Character, s def.Sim, log def.Logger, count int) {
	if count >= 2 {
		m := make([]float64, def.EndStatType)
		m[def.PhyP] = 0.25
		c.AddMod(def.CharStatMod{
			Key: "maiden-2pc",
			Amount: func(a def.AttackTag) ([]float64, bool) {
				return m, true
			},
			Expiry: -1,
		})
	}
	if count >= 4 {
		stacks := 0
		icd := 0
		dur := 0
		m := make([]float64, def.EndStatType)

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
			stacks++
			if stacks > 2 {
				stacks = 2
				m[def.PhyP] = 0.25
			}
			m[def.ATKP] = 0.09 * float64(stacks)

			log.Debugw("pale flame 4pc proc", "frame", s.Frame(), "event", def.LogArtifactEvent, "stacks", stacks, "expiry", s.Frame()+420, "icd", s.Frame()+18)
			icd = s.Frame() + 18
			dur = s.Frame() + 420
		}, fmt.Sprintf("pf4-%v", c.Name()))

		c.AddMod(def.CharStatMod{
			Key: "pf-4pc",
			Amount: func(a def.AttackTag) ([]float64, bool) {
				if dur < s.Frame() {
					m[def.ATKP] = 0
					m[def.PhyP] = 0
					return nil, false
				}

				return m, true
			},
			Expiry: -1,
		})

	}
	//add flat stat to char
}

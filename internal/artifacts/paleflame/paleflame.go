package paleflame

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	combat.RegisterSetFunc("pale flame", New)
}

func New(c core.Character, s core.Sim, log core.Logger, count int) {
	if count >= 2 {
		m := make([]float64, core.EndStatType)
		m[core.PhyP] = 0.25
		c.AddMod(core.CharStatMod{
			Key: "maiden-2pc",
			Amount: func(a core.AttackTag) ([]float64, bool) {
				return m, true
			},
			Expiry: -1,
		})
	}
	if count >= 4 {
		stacks := 0
		icd := 0
		dur := 0
		m := make([]float64, core.EndStatType)

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
			stacks++
			if stacks > 2 {
				stacks = 2
				m[core.PhyP] = 0.25
			}
			m[core.ATKP] = 0.09 * float64(stacks)

			log.Debugw("pale flame 4pc proc", "frame", s.Frame(), "event", core.LogArtifactEvent, "stacks", stacks, "expiry", s.Frame()+420, "icd", s.Frame()+18)
			icd = s.Frame() + 18
			dur = s.Frame() + 420
		}, fmt.Sprintf("pf4-%v", c.Name()))

		c.AddMod(core.CharStatMod{
			Key: "pf-4pc",
			Amount: func(a core.AttackTag) ([]float64, bool) {
				if dur < s.Frame() {
					m[core.ATKP] = 0
					m[core.PhyP] = 0
					return nil, false
				}

				return m, true
			},
			Expiry: -1,
		})

	}
	//add flat stat to char
}

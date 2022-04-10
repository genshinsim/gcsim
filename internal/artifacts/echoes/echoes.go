package echoes

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterSetFunc("echoes of an offering", New)
	core.RegisterSetFunc("echoesofanoffering", New)
	core.RegisterSetFunc("echoes", New)
}

func New(c core.Character, s *core.Core, count int, params map[string]int) {
	prob := 0.36
	icd := 0

	if count >= 2 {
		m := make([]float64, core.EndStatType)
		m[core.ATKP] = 0.18
		c.AddMod(core.CharStatMod{
			Key: "echoes-2pc",
			Amount: func() ([]float64, bool) {
				return m, true
			},
			Expiry: -1,
		})
	}

	if count >= 4 {
		s.Events.Subscribe(core.OnAttackWillLand, func(args ...interface{}) bool {
			// if the active char is not the equipped char then ignore
			if s.ActiveChar != c.CharIndex() {
				return false
			}

			atk := args[1].(*core.AttackEvent)

			// If this is not a normal attack then ignore
			if atk.Info.AttackTag != core.AttackTagNormal {
				return false
			}

			// If Artifact set effect is still on CD then ignore
			if s.F < icd {
				return false
			}

			if s.Rand.Float64() < prob {
				dmgAdded := (atk.Snapshot.BaseAtk*(1+atk.Snapshot.Stats[core.ATKP]) + atk.Snapshot.Stats[core.ATK]) * 0.7
				atk.Info.FlatDmg += dmgAdded
				s.Log.NewEvent("echoes 4pc proc", core.LogArtifactEvent, c.CharIndex(),
					"probability", prob,
					"dmg_added", dmgAdded,
				)

				prob = 0.36
				icd = s.F + 0.2*60
			} else {
				prob += 0.2
				icd = s.F + 0.2*60

				if prob > 1 {
					prob = 1
				}
			}

			return false
		}, fmt.Sprintf("echoes-4pc-%v", c.Name()))
	}
}

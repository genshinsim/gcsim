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
	probicd := 0
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
			if s.ActiveChar != c.CharIndex() {
				return false
			}
			if args[1].(*core.AttackEvent).Info.AttackTag != core.AttackTagNormal {
				return false
			}
			atk := args[1].(*core.AttackEvent)
			if s.Rand.Float64() < prob {
				snap := c.Snapshot(&atk.Info)
				dmgAdded := (snap.BaseAtk*(1+snap.Stats[core.ATKP]) + snap.Stats[core.ATK]) * 0.6
				atk.Info.FlatDmg += dmgAdded
			} else {
				if s.F > probicd {
					prob += 0.2
					probicd = s.F + 0.3*60
				}
				if prob > 1 {
					prob = 1
				}
			}
			return false
		}, fmt.Sprintf("echoes-4pc-%v", c.Name()))

	}
}

package dodoco

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	combat.RegisterWeaponFunc("dodoco tales", weapon)
}

func weapon(c core.Character, s core.Sim, log core.Logger, r int, param map[string]int) {
	atkExpiry := 0
	dmgExpiry := 0

	m := make([]float64, core.EndStatType)
	m[core.DmgP] = .12 + float64(r)*.04
	c.AddMod(core.CharStatMod{
		Key: "dodoco ca",
		Amount: func(a core.AttackTag) ([]float64, bool) {
			if a != core.AttackTagExtra {
				return nil, false
			}
			return m, dmgExpiry > s.Frame()
		},
		Expiry: -1,
	})

	n := make([]float64, core.EndStatType)
	n[core.ATKP] = .06 + float64(r)*0.02
	c.AddMod(core.CharStatMod{
		Key: "dodoco atk",
		Amount: func(a core.AttackTag) ([]float64, bool) {
			return n, atkExpiry > s.Frame()
		},
		Expiry: -1,
	})

	s.AddOnAttackLanded(func(t core.Target, ds *core.Snapshot, dmg float64, crit bool) {
		if ds.ActorIndex != c.CharIndex() {
			return
		}
		switch ds.AttackTag {
		case core.AttackTagNormal:
			dmgExpiry = s.Frame() + 360
		case core.AttackTagExtra:
			atkExpiry = s.Frame() + 360
		}
	}, fmt.Sprintf("dodoco-%v", c.Name()))

}

package dodoco

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/def"
)

func init() {
	combat.RegisterWeaponFunc("dodoco tales", weapon)
}

func weapon(c def.Character, s def.Sim, log def.Logger, r int, param map[string]int) {
	atkExpiry := 0
	dmgExpiry := 0

	m := make([]float64, def.EndStatType)
	m[def.DmgP] = .12 + float64(r)*.04
	c.AddMod(def.CharStatMod{
		Key: "dodoco ca",
		Amount: func(a def.AttackTag) ([]float64, bool) {
			if a != def.AttackTagExtra {
				return nil, false
			}
			return m, dmgExpiry > s.Frame()
		},
		Expiry: -1,
	})

	n := make([]float64, def.EndStatType)
	n[def.ATKP] = .06 + float64(r)*0.02
	c.AddMod(def.CharStatMod{
		Key: "dodoco atk",
		Amount: func(a def.AttackTag) ([]float64, bool) {
			return n, atkExpiry > s.Frame()
		},
		Expiry: -1,
	})

	s.AddOnAttackLanded(func(t def.Target, ds *def.Snapshot, dmg float64, crit bool) {
		if ds.ActorIndex != c.CharIndex() {
			return
		}
		switch ds.AttackTag {
		case def.AttackTagNormal:
			dmgExpiry = s.Frame() + 360
		case def.AttackTagExtra:
			atkExpiry = s.Frame() + 360
		}
	}, fmt.Sprintf("dodoco-%v", c.Name()))

}

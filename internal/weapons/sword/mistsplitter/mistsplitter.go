package mistsplitter

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/def"
)

func init() {
	combat.RegisterWeaponFunc("mistsplitter reforged", weapon)
}

func weapon(c def.Character, s def.Sim, log def.Logger, r int, param map[string]int) {
	m := make([]float64, def.EndStatType)
	base := 0.09 + float64(r)*0.03
	m[def.PyroP] = base
	m[def.HydroP] = base
	m[def.CryoP] = base
	m[def.ElectroP] = base
	m[def.AnemoP] = base
	m[def.GeoP] = base
	m[def.EleP] = base
	m[def.PhyP] = base
	m[def.DendroP] = base
	stack := 0.06 + float64(r)*0.02
	max := 0.21 + float64(r)*0.07
	bonus := def.EleToDmgP(c.Ele())

	normal := 0
	skill := 0

	s.AddOnAttackLanded(func(t def.Target, ds *def.Snapshot, base float64, crit bool) {
		if ds.ActorIndex != c.CharIndex() {
			return
		}
		if ds.AttackTag != def.AttackTagNormal {
			return
		}
		if ds.Element == def.Physical {
			return
		}
		normal = s.Frame() + 300 // lasts 5 seconds

	}, fmt.Sprintf("mistsplitter-%v", c.Name()))

	s.AddEventHook(func(s def.Sim) bool {
		if s.ActiveCharIndex() != c.CharIndex() {
			return false
		}
		skill = s.Frame() + 600
		return false

	}, fmt.Sprintf("mistsplitter-%v", c.Name()), def.PostBurstHook)

	c.AddMod(def.CharStatMod{
		Key: "mistsplitter",
		Amount: func(a def.AttackTag) ([]float64, bool) {
			count := 0
			if c.CurrentEnergy() < c.MaxEnergy() {
				count++
			}
			if normal > s.Frame() {
				count++
			}
			if skill > s.Frame() {
				count++
			}
			dmg := float64(count) * stack
			if count > 3 {
				count = 3 // should never happen
				dmg = max
			}
			//bonus for current char
			m[bonus] = base + dmg
			return m, true
		},
		Expiry: -1,
	})
}

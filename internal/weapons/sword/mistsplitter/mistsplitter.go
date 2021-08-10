package mistsplitter

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	combat.RegisterWeaponFunc("mistsplitter reforged", weapon)
}

func weapon(c core.Character, s core.Sim, log core.Logger, r int, param map[string]int) {
	m := make([]float64, core.EndStatType)
	base := 0.09 + float64(r)*0.03
	m[core.PyroP] = base
	m[core.HydroP] = base
	m[core.CryoP] = base
	m[core.ElectroP] = base
	m[core.AnemoP] = base
	m[core.GeoP] = base
	m[core.EleP] = base
	m[core.PhyP] = base
	m[core.DendroP] = base
	stack := 0.06 + float64(r)*0.02
	max := 0.21 + float64(r)*0.07
	bonus := core.EleToDmgP(c.Ele())

	normal := 0
	skill := 0

	s.AddOnAttackLanded(func(t core.Target, ds *core.Snapshot, base float64, crit bool) {
		if ds.ActorIndex != c.CharIndex() {
			return
		}
		if ds.AttackTag != core.AttackTagNormal {
			return
		}
		if ds.Element == core.Physical {
			return
		}
		normal = s.Frame() + 300 // lasts 5 seconds

	}, fmt.Sprintf("mistsplitter-%v", c.Name()))

	s.AddEventHook(func(s core.Sim) bool {
		if s.ActiveCharIndex() != c.CharIndex() {
			return false
		}
		skill = s.Frame() + 600
		return false

	}, fmt.Sprintf("mistsplitter-%v", c.Name()), core.PostBurstHook)

	c.AddMod(core.CharStatMod{
		Key: "mistsplitter",
		Amount: func(a core.AttackTag) ([]float64, bool) {
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

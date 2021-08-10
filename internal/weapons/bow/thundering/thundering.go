package thundering

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	combat.RegisterWeaponFunc("thundering pulse", weapon)
}

func weapon(c core.Character, s core.Sim, log core.Logger, r int, param map[string]int) {
	m := make([]float64, core.EndStatType)
	m[core.ATKP] = 0.15 + float64(r)*0.05
	stack := 0.09 + float64(r)*0.03
	max := 0.3 + float64(r)*0.1

	normal := 0
	skill := 0

	s.AddOnAttackLanded(func(t core.Target, ds *core.Snapshot, dmg float64, crit bool) {
		if ds.ActorIndex != c.CharIndex() {
			return
		}
		if ds.AttackTag != core.AttackTagNormal {
			return
		}
		normal = s.Frame() + 300 // lasts 5 seconds

	}, fmt.Sprintf("thundering-pulse-%v", c.Name()))

	s.AddEventHook(func(s core.Sim) bool {
		if s.ActiveCharIndex() != c.CharIndex() {
			return false
		}
		skill = s.Frame() + 600
		return false

	}, fmt.Sprintf("thundering-pulse-%v", c.Name()), core.PostSkillHook)

	c.AddMod(core.CharStatMod{
		Key: "thundering",
		Amount: func(a core.AttackTag) ([]float64, bool) {
			m[core.DmgP] = 0
			if a != core.AttackTagNormal {
				return m, true
			}
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
			m[core.DmgP] = dmg
			return m, true
		},
		Expiry: -1,
	})
}

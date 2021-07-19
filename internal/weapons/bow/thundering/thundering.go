package thundering

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/def"
)

func init() {
	combat.RegisterWeaponFunc("thundering pulse", weapon)
}

func weapon(c def.Character, s def.Sim, log def.Logger, r int, param map[string]int) {
	m := make([]float64, def.EndStatType)
	m[def.ATKP] = 0.15 + float64(r)*0.05
	stack := 0.09 + float64(r)*0.03

	normal := 0
	skill := 0

	s.AddOnAttackLanded(func(t def.Target, ds *def.Snapshot, dmg float64, crit bool) {
		if ds.ActorIndex != c.CharIndex() {
			return
		}
		if ds.AttackTag != def.AttackTagNormal {
			return
		}
		normal = s.Frame() + 300 // lasts 5 seconds

	}, fmt.Sprintf("thundering-pulse-%v", c.Name()))

	s.AddEventHook(func(s def.Sim) bool {
		if s.ActiveCharIndex() != c.CharIndex() {
			return false
		}
		skill = s.Frame() + 600
		return false

	}, fmt.Sprintf("thundering-pulse-%v", c.Name()), def.PostSkillHook)

	c.AddMod(def.CharStatMod{
		Key: "thundering",
		Amount: func(a def.AttackTag) ([]float64, bool) {
			m[def.DmgP] = 0
			if a != def.AttackTagNormal {
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
			if count > 3 {
				count = 3 // should never happen
			}
			m[def.DmgP] = float64(count) * stack
			return m, true
		},
		Expiry: -1,
	})
}

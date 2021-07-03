package aquila

import (
	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/def"
)

func init() {
	combat.RegisterWeaponFunc("aquila favonia", weapon)
}

func weapon(c def.Character, s def.Sim, log def.Logger, r int, param map[string]int) {
	m := make([]float64, def.EndStatType)
	m[def.ATKP] = .15 + .05*float64(r)
	c.AddMod(def.CharStatMod{
		Key: "acquila favonia",
		Amount: func(a def.AttackTag) ([]float64, bool) {
			return m, true
		},
		Expiry: -1,
	})

	dmg := 1.7 + .3*float64(r)
	heal := .85 + .15*float64(r)

	last := -1

	s.AddOnHurt(func(s def.Sim) {
		if s.ActiveCharIndex() != c.CharIndex() {
			return
		}
		if s.Frame()-last < 900 && last != -1 {
			return
		}
		last = s.Frame()
		d := c.Snapshot(
			"Aquila Favonia",
			def.AttackTagWeaponSkill,
			def.ICDTagNone,
			def.ICDGroupDefault,
			def.StrikeTypeDefault,
			def.Physical,
			100,
			dmg,
		)
		c.QueueDmg(&d, 1)

		atk := d.BaseAtk*(1+d.Stats[def.ATKP]) + d.Stats[def.ATK]

		log.Debugw("acquila heal triggered", "frame", s.Frame(), "event", def.LogWeaponEvent, "atk", atk, "heal amount", atk*heal)
		s.HealActive(atk * heal)
	})
}

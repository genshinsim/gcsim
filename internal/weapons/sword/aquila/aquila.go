package aquila

import (
	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	combat.RegisterWeaponFunc("aquila favonia", weapon)
}

func weapon(c core.Character, s core.Sim, log core.Logger, r int, param map[string]int) {
	m := make([]float64, core.EndStatType)
	m[core.ATKP] = .15 + .05*float64(r)
	c.AddMod(core.CharStatMod{
		Key: "acquila favonia",
		Amount: func(a core.AttackTag) ([]float64, bool) {
			return m, true
		},
		Expiry: -1,
	})

	dmg := 1.7 + .3*float64(r)
	heal := .85 + .15*float64(r)

	last := -1

	s.AddOnHurt(func(s core.Sim) {
		if s.ActiveCharIndex() != c.CharIndex() {
			return
		}
		if s.Frame()-last < 900 && last != -1 {
			return
		}
		last = s.Frame()
		d := c.Snapshot(
			"Aquila Favonia",
			core.AttackTagWeaponSkill,
			core.ICDTagNone,
			core.ICDGroupDefault,
			core.StrikeTypeDefault,
			core.Physical,
			100,
			dmg,
		)
		d.Targets = core.TargetAll
		c.QueueDmg(&d, 1)

		atk := d.BaseAtk*(1+d.Stats[core.ATKP]) + d.Stats[core.ATK]

		log.Debugw("acquila heal triggered", "frame", s.Frame(), "event", core.LogWeaponEvent, "atk", atk, "heal amount", atk*heal)
		s.HealActive(atk * heal)
	})
}

package kaeya

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player/shield"
	"github.com/genshinsim/gcsim/pkg/enemy"
)

func (c *char) c1() {
	m := make([]float64, attributes.EndStatType)
	c.AddAttackMod("kaeya-c1", -1, func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
		e, ok := t.(*enemy.Enemy)
		if !ok {
			return nil, false
		}
		if atk.Info.AttackTag != combat.AttackTagNormal && atk.Info.AttackTag != combat.AttackTagExtra {
			return nil, false
		}
		if !e.AuraContains(attributes.Cryo, attributes.Frozen) {
			return nil, false
		}
		m[attributes.CR] = 0.15
		return m, true
	})
}

func (c *char) c4() {
	c.Core.Events.Subscribe(event.OnCharacterHurt, func(args ...interface{}) bool {
		if c.Core.F < c.c4icd && c.c4icd != 0 {
			return false
		}
		maxhp := c.MaxHP()
		if c.HPCurrent/maxhp < .2 {
			c.c4icd = c.Core.F + 3600
			c.Core.Player.Shields.Add(&shield.Tmpl{
				Src:        c.Core.F,
				ShieldType: shield.ShieldKaeyaC4,
				Name:       "Kaeya C4",
				HP:         .3 * maxhp,
				Ele:        attributes.Cryo,
				Expires:    c.Core.F + 1200,
			})
		}
		return false
	}, "kaeya-c4")
}

package thoma

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/shield"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func (c *char) genShield(src string, shieldamt float64) {
	if c.Core.F > c.a1icd && c.a1Stack < 5 {
		c.a1Stack++
		c.a1icd = c.Core.F + 0.3*60
		c.Core.Status.Add("thoma-a1", 6*60)
	}
	if c.Core.Player.Shields.Get(shield.ShieldThomaSkill) != nil {
		maxHP := c.maxShieldHP()
		if c.Core.Player.Shields.Get(shield.ShieldThomaSkill).CurrentHP()+shieldamt > maxHP {
			shieldamt = maxHP - c.Core.Player.Shields.Get(shield.ShieldThomaSkill).CurrentHP()
		}
	}
	//add shield
	c.Core.Tasks.Add(func() {
		c.Core.Player.Shields.Add(&shield.Tmpl{
			Src:        c.Core.F,
			ShieldType: shield.ShieldThomaSkill,
			Name:       src,
			HP:         shieldamt,
			Ele:        attributes.Pyro,
			Expires:    c.Core.F + 8*60, //8 sec
		})
	}, 1)

	if c.Base.Cons >= 6 {
		for _, char := range c.Core.Player.Chars() {
			char.AddAttackMod(character.AttackMod{
				Base: modifier.NewBaseWithHitlag("thoma-c6", 360),
				Amount: func(ae *combat.AttackEvent, t combat.Target) ([]float64, bool) {
					switch ae.Info.AttackTag {
					case combat.AttackTagNormal, combat.AttackTagExtra, combat.AttackTagPlunge:
						return c.c6buff, true
					}
					return nil, false
				},
			})
		}
	}
}

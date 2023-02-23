package thoma

import (
	"math"

	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/shield"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func (c *char) genShield(src string, shieldamt float64, shouldStack bool) {
	if !c.StatusIsActive("thoma-a1-icd") && c.a1Stack < 5 {
		c.a1Stack++
		c.AddStatus("thoma-a1-icd", 18, true) // 0.3s * 60
		c.AddStatus("thoma-a1", 360, true)    // 6s * 60
	}
	existingShield := c.Core.Player.Shields.Get(shield.ShieldThomaSkill)
	if existingShield != nil {
		if shouldStack {
			shieldamt += existingShield.CurrentHP()
		} else {
			shieldamt = math.Max(shieldamt, existingShield.CurrentHP())
		}
		shieldamt = math.Min(shieldamt, c.maxShieldHP())
	}
	// add shield
	c.Core.Tasks.Add(func() {
		c.Core.Player.Shields.Add(&shield.Tmpl{
			Src:        c.Core.F,
			ShieldType: shield.ShieldThomaSkill,
			Name:       src,
			HP:         shieldamt,
			Ele:        attributes.Pyro,
			Expires:    c.Core.F + 8*60, // 8 sec
		})
	}, 1)

	if c.Base.Cons >= 6 {
		for _, char := range c.Core.Player.Chars() {
			char.AddAttackMod(character.AttackMod{
				Base: modifier.NewBaseWithHitlag("thoma-c6", 360),
				Amount: func(ae *combat.AttackEvent, _ combat.Target) ([]float64, bool) {
					switch ae.Info.AttackTag {
					case attacks.AttackTagNormal, attacks.AttackTagExtra, attacks.AttackTagPlunge:
						return c.c6buff, true
					}
					return nil, false
				},
			})
		}
	}
}

package dahlia

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/player/shield"
)

type shd struct {
	*shield.Tmpl
	c *char
}

var dahliaShield *shd

func (c *char) genShield() {
	dahliaShield = &shd{
		Tmpl: &shield.Tmpl{
			ActorIndex: c.Index(),
			Target:     -1,
			Src:        c.Core.F,
			Name:       "Radiant Psalter (Shield)",
			ShieldType: shield.DahliaBurst,
			HP:         c.shieldHP(),
			Ele:        attributes.Hydro,
			Expires:    c.favonianFavorMaxExpiry, // TO-DO: Shields don't support hitlag so this field's value is wrong
		},
		c: c,
	}

	c.Core.Player.Shields.Add(dahliaShield)
}

func (c *char) removeShield() {
	if !c.hasShield() || dahliaShield == nil {
		return
	}

	dahliaShield.Expires = c.Core.F + 1 // +1f to be sure
	dahliaShield = nil
}

// If shield is broken and there are Benison stacks, create shield (after some delay)
func (s *shd) OnDamage(dmg float64, ele attributes.Element, bonus float64) (float64, bool) {
	taken, ok := s.Tmpl.OnDamage(dmg, ele, bonus)

	if !ok && s.Tmpl.ShieldType == shield.DahliaBurst && s.c.currentBenisonStacks > 0 && s.c.StatusExpiry(burstFavonianFavor) > s.c.Core.F+burstShieldRegenerated {
		s.c.Core.Tasks.Add(func() {
			s.c.currentBenisonStacks--
			s.c.genShield()
			s.c.c2()
		}, burstShieldRegenerated)
	}
	return taken, ok
}

func (c *char) hasShield() bool {
	return c.Core.Player.Shields.Get(shield.DahliaBurst) != nil
}

func (c *char) shieldHP() float64 {
	return burstShieldPP[c.TalentLvlBurst()]*c.MaxHP() + burstShieldFlat[c.TalentLvlBurst()]
}

package dahlia

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/player/shield"
)

type shd struct {
	*shield.Tmpl
	c *char
}

func (c *char) genShield() {
	c.shield = &shd{
		Tmpl: &shield.Tmpl{
			ActorIndex: c.Index(),
			Target:     -1,
			Src:        c.Core.F,
			Name:       "Radiant Psalter (Shield)",
			ShieldType: shield.DahliaBurst,
			HP:         c.shieldHP(),
			Ele:        attributes.Hydro,
			Expires:    -1, // QueueCharTask in Burst handles cuz hitlag affected
		},
		c: c,
	}

	c.Core.Player.Shields.Add(c.shield)
}

func (c *char) removeShield() {
	if !c.hasShield() || c.shield == nil {
		return
	}

	c.shield.Expires = c.Core.F + 1 // +1f to be sure
	c.shield = nil
}

// If shield is broken and there are Benison stacks, create shield (after some delay)
func (s *shd) OnDamage(dmg float64, ele attributes.Element, bonus float64) (float64, bool) {
	taken, ok := s.Tmpl.OnDamage(dmg, ele, bonus)

	if !ok {
		s.c.tryRegenShield()
	}
	return taken, ok
}

func (c *char) hasShield() bool {
	return c.Core.Player.Shields.Get(shield.DahliaBurst) != nil
}

func (c *char) shieldHP() float64 {
	return burstShieldPP[c.TalentLvlBurst()]*c.MaxHP() + burstShieldFlat[c.TalentLvlBurst()]
}

func (c *char) tryRegenShield() {
	if c.currentBenisonStacks <= 0 {
		return
	}
	if c.StatusDuration(burstFavonianFavor) < burstShieldRegenerated {
		return
	}
	c.Core.Tasks.Add(func() {
		c.currentBenisonStacks--
		c.genShield()
		c.c2()
	}, burstShieldRegenerated)
}

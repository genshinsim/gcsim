package dahlia

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/player/shield"
)

type shd struct {
	*shield.Tmpl
	c *char
}

func (c *char) newShield(base float64, expiry int) *shd {
	n := &shd{}
	n.Tmpl = &shield.Tmpl{}
	n.ActorIndex = c.Index()
	n.Target = -1
	n.Src = c.Core.F
	n.Name = "Radiant Psalter (Shield)"
	n.ShieldType = shield.DahliaBurst
	n.HP = base
	n.Ele = attributes.Hydro
	n.Expires = expiry
	n.c = c
	return n
}

func (c *char) genShield() {
	c.Core.Player.Shields.Add(c.newShield(c.shieldHP(), favonianFavorExpiry))
}

// If shield is broken and there are Benison stacks, create shield (after some delay)
func (s *shd) OnDamage(dmg float64, ele attributes.Element, bonus float64) (float64, bool) {
	taken, ok := s.Tmpl.OnDamage(dmg, ele, bonus)

	if !ok && currentBenisonStacks > 0 && favonianFavorExpiry > s.c.Core.F+burstShieldRegenerated {
		currentBenisonStacks--

		s.c.QueueCharTask(func() {
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
	return burstShieldPP[c.TalentLvlSkill()]*c.MaxHP() + burstShieldFlat[c.TalentLvlSkill()]
}

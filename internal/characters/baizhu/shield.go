package baizhu

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/player/shield"
)

func (c *char) removeShield() {
	c.summonSeamlessShieldHealing()
	c.summonSpiritvein()
}

func (c *char) newShield(base float64, dur int) *shd {
	n := &shd{}
	n.Tmpl = &shield.Tmpl{}
	n.Tmpl.ActorIndex = c.Index
	n.Tmpl.Src = c.Core.F
	n.Tmpl.ShieldType = shield.ShieldBaizhuBurst
	n.Tmpl.Ele = attributes.Dendro
	n.Tmpl.HP = base
	n.Tmpl.Name = "Baizhu Seamless shield"
	n.Tmpl.Expires = c.Core.F + dur
	n.c = c
	return n
}

type shd struct {
	*shield.Tmpl
	c *char
}

func (s *shd) OnExpire() {
	s.c.removeShield()
}

func (s *shd) OnDamage(dmg float64, ele attributes.Element, bonus float64) (float64, bool) {
	taken, ok := s.Tmpl.OnDamage(dmg, ele, bonus)
	if !ok {
		s.c.removeShield()
	}
	return taken, ok
}

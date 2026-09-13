package noelle

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/player/shield"
)

type noelleShield struct {
	*shield.Tmpl
	c *char
}

func (c *char) newShield(base float64, t shield.Type, dur int) *noelleShield {
	n := &noelleShield{}
	n.Tmpl = &shield.Tmpl{
		ActorIndex: c.Index(),
		Target:     -1,
		Name:       "Noelle Skill",
		Src:        c.Core.F,
		ShieldType: t,
		Ele:        attributes.Geo,
		HP:         base,
		Expires:    c.Core.F + dur,
	}
	n.c = c
	return n
}

func (n *noelleShield) OnExpire() {
	if n.c.Base.Cons >= 4 {
		n.c.explodeShield()
	}
}

func (n *noelleShield) OnOverwrite() {
	if n.c.Base.Cons >= 4 {
		n.c.explodeShield()
	}
}

func (n *noelleShield) OnDamage(dmg float64, ele attributes.Element, bonus float64) (float64, bool) {
	taken, ok := n.Tmpl.OnDamage(dmg, ele, bonus)
	if !ok && n.c.Base.Cons >= 4 {
		n.c.explodeShield()
	}
	return taken, ok
}

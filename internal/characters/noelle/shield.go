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
	n.Tmpl = &shield.Tmpl{}
	n.ActorIndex = c.Index()
	n.Target = -1
	n.Src = c.Core.F
	n.ShieldType = t
	n.Name = "Noelle Skill"
	n.HP = base
	n.Expires = c.Core.F + dur
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

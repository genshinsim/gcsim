package xinyan

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/player/shield"
)

type xinyanShield struct {
	*shield.Tmpl
	c *char
}

func (c *char) newShield(base float64, t shield.Type, dur int) *xinyanShield {
	n := &xinyanShield{}
	n.Tmpl = &shield.Tmpl{
		ActorIndex: c.Index(),
		Target:     -1,
		Name:       "Xinyan Skill",
		Src:        c.Core.F,
		ShieldType: t,
		Ele:        attributes.Pyro,
		HP:         base,
		Expires:    c.Core.F + dur,
	}
	n.c = c
	return n
}

func (n *xinyanShield) OnExpire() {
	n.c.shieldLevel = 1
}

func (n *xinyanShield) OnDamage(dmg float64, ele attributes.Element, bonus float64) (float64, bool) {
	taken, ok := n.Tmpl.OnDamage(dmg, ele, bonus)
	if !ok {
		n.c.shieldLevel = 1
	}
	return taken, ok
}

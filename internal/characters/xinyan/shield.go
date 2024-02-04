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
	n.Tmpl = &shield.Tmpl{}
	n.Tmpl.ActorIndex = c.Index
	n.Tmpl.Target = -1
	n.Tmpl.Src = c.Core.F
	n.Tmpl.ShieldType = t
	n.Tmpl.Name = "Xinyan Skill"
	n.Tmpl.HP = base
	n.Tmpl.Expires = c.Core.F + dur
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

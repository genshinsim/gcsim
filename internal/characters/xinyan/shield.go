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
	n.ActorIndex = c.Index()
	n.Target = -1
	n.Src = c.Core.F
	n.ShieldType = t
	n.Name = "Xinyan Skill"
	n.HP = base
	n.Expires = c.Core.F + dur
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

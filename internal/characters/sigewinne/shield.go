package sigewinne

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/player/shield"
)

// need to rewrite Expires
type shd struct {
	*shield.Tmpl
}

func (c *char) addC2Shield() {
	if c.Base.Cons < 2 {
		return
	}

	shieldHP := 0.3 * c.MaxHP()
	c.c2Shield = &shd{
		Tmpl: &shield.Tmpl{
			ActorIndex: c.Index,
			Target:     c.Index,
			Src:        c.Core.F,
			Name:       "Sigewinne C2 shield",
			ShieldType: shield.SigewinneC2,
			HP:         shieldHP,
			Ele:        attributes.Hydro,
			Expires:    c.Core.F + 15*60,
		},
	}
	c.Core.Player.Shields.Add(c.c2Shield)
}

func (c *char) removeC2Shield() {
	if c.Base.Cons < 2 || c.c2Shield == nil {
		return
	}
	// +1f to be sure
	c.c2Shield.Expires = c.Core.F + 1
	c.c2Shield = nil
}

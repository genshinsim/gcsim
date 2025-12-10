package ineffa

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/player/shield"
)

// need to rewrite Expires
type shd struct {
	*shield.Tmpl
}

func (c *char) addShield() {
	atk := c.TotalAtk()
	shieldHP := shieldAtk[c.TalentLvlSkill()]*atk + shieldFlat[c.TalentLvlSkill()]
	c.skillShield = &shd{
		Tmpl: &shield.Tmpl{
			ActorIndex: c.Index(),
			Target:     -1,
			Src:        c.Core.F,
			Name:       "Optical Flow Shield Barrier",
			ShieldType: shield.IneffaSkill,
			HP:         shieldHP,
			Ele:        attributes.Electro,
			Expires:    c.Core.F + 20*60,
		},
	}
	c.Core.Player.Shields.Add(c.skillShield)

	c.c1OnShield()
}

package citlali

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/player/shield"
)

// need to rewrite Expires
type shd struct {
	*shield.Tmpl
}

func (c *char) addShield() {
	em := c.Stat(attributes.EM)
	shieldHP := shieldEM[c.TalentLvlSkill()]*em + shieldFlat[c.TalentLvlSkill()]
	c.skillShield = &shd{
		Tmpl: &shield.Tmpl{
			ActorIndex: c.Index,
			Target:     -1,
			Src:        c.Core.F,
			Name:       "Citalali Skill Shield",
			ShieldType: shield.CitlaliSkill,
			HP:         shieldHP,
			Ele:        attributes.Cryo,
			Expires:    c.Core.F + 20*60,
		},
	}
	c.Core.Player.Shields.Add(c.skillShield)
}

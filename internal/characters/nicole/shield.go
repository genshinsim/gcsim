package nicole

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/player/shield"
)

// need to rewrite Expires
type shd struct {
	*shield.Tmpl
}

func (c *char) addShield() {
	shieldHP := skillShieldAtk[c.TalentLvlSkill()]*c.TotalAtk() + skillShieldFlat[c.TalentLvlSkill()]
	c.skillShield = &shd{
		Tmpl: &shield.Tmpl{
			ActorIndex: c.Index(),
			Target:     -1,
			Src:        c.Core.F,
			Name:       "Shield of Blazing Light (Shield)",
			ShieldType: shield.NicoleSkill,
			HP:         shieldHP,
			Ele:        attributes.Pyro,
			Expires:    c.Core.F + 20*60,
		},
	}
	c.Core.Player.Shields.Add(c.skillShield)
}

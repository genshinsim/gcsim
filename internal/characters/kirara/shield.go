package kirara

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/player/shield"
)

func (c *char) genShield(src string, shieldamt float64) {
	existingShield := c.Core.Player.Shields.Get(shield.KiraraSkill)
	if existingShield != nil {
		shieldamt += existingShield.CurrentHP()
	}
	shieldamt = min(shieldamt, c.maxShieldHP())

	// add shield
	c.Core.Tasks.Add(func() {
		c.Core.Player.Shields.Add(&shield.Tmpl{
			Src:        c.Core.F,
			ShieldType: shield.KiraraSkill,
			Name:       src,
			HP:         shieldamt,
			Ele:        attributes.Dendro,
			Expires:    c.Core.F + 12*60,
		})
	}, 1)
}

func (c *char) shieldHP() float64 {
	return shieldPP[c.TalentLvlSkill()]*c.MaxHP() + shieldFlat[c.TalentLvlSkill()]
}

func (c *char) maxShieldHP() float64 {
	return maxShieldPP[c.TalentLvlSkill()]*c.MaxHP() + maxShieldFlat[c.TalentLvlSkill()]
}

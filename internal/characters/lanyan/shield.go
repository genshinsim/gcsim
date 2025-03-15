package lanyan

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/player/shield"
)

func (c *char) genShield(ele attributes.Element) {
	c.Core.Player.Shields.Add(&shield.Tmpl{
		ActorIndex: c.Index,
		Target:     -1,
		Src:        c.Core.F,
		ShieldType: shield.LanyanShield,
		Name:       "Lanyan Skill",
		HP:         c.shieldHP(),
		Ele:        ele,
		Expires:    c.Core.F + 12.5*60,
	})
}

func (c *char) restoreShield(percent float64) {
	amt := c.shieldHP() * percent
	existingShield := c.Core.Player.Shields.Get(shield.LanyanShield)
	if existingShield == nil {
		return
	}
	amt = min(amt+existingShield.CurrentHP(), c.shieldHP())

	c.Core.Player.Shields.Add(&shield.Tmpl{
		ActorIndex: c.Index,
		Target:     -1,
		Src:        c.Core.F,
		ShieldType: shield.LanyanShield,
		Name:       existingShield.Desc(),
		HP:         amt,
		Ele:        c.absorbedElement,
		Expires:    existingShield.Expiry(),
	})
}

func (c *char) hasShield() bool {
	return c.Core.Player.Shields.Get(shield.LanyanShield) != nil
}

func (c *char) shieldHP() float64 {
	return shieldAmt[c.TalentLvlSkill()]*c.TotalAtk() + shieldFlat[c.TalentLvlSkill()]
}

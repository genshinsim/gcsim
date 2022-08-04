package kazuha

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

//After using Chihayaburu or Kazuha Slash, Kaedehara Kazuha gains an Anemo
//Infusion for 5s. Additionally, each point of Elemental Mastery will
//increase the DMG dealt by Kaedehara Kazuha's Normal, Charged, and Plunging
//Attacks by 0.2%.
func (c *char) c6() {
	c.AddStatus(c6BuffKey, 60*5, true)
	c.Core.Player.AddWeaponInfuse(
		c.Index,
		"kazuha-c6-infusion",
		attributes.Anemo,
		60*5,
		true,
		combat.AttackTagNormal, combat.AttackTagExtra, combat.AttackTagPlunge,
	)
}

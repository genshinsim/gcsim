package lynette

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

// TODO: C1 is not implemented, because vortex/pulling mechanics are not implemented

// Whenever the Bogglecat Box summoned by Magic Trick: Astonishing Shift fires a Vivid Shot, it will fire an extra Vivid Shot.
func (c *char) c2() {
	if c.Base.Cons < 2 {
		return
	}
	c.vividCount = 2
}

// Increases Enigmatic Feint's charges by 1.
func (c *char) c4() {
	if c.Base.Cons < 4 {
		return
	}
	c.SetNumCharges(action.ActionSkill, 2)
}

// When Lynette uses Enigmatic Feint's Enigma Thrust, she will gain an Anemo Infusion and 20% Anemo DMG Bonus for 6s.
func (c *char) c6() {
	duration := int((6 + 0.4) * 60)

	// add anemo infusion
	c.Core.Player.AddWeaponInfuse(
		c.Index,
		"lynette-c6-infusion",
		attributes.Anemo,
		duration,
		true,
		attacks.AttackTagNormal, attacks.AttackTagExtra, attacks.AttackTagPlunge,
	)

	// add anemo% buff
	m := make([]float64, attributes.EndStatType)
	m[attributes.AnemoP] = 0.2
	c.AddStatMod(character.StatMod{
		Base: modifier.NewBaseWithHitlag("lynette-c6-buff", duration),
		Amount: func() ([]float64, bool) {
			return m, true
		},
	})
}

package razor

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

// Decreases Claw and Thunder's CD by 18%.
func (c *char) a1CDReduction(cd int) int {
	if c.Base.Ascension < 1 {
		return cd
	}
	return int(float64(cd) * 0.82)
}

// Using Lightning Fang resets the CD of Claw and Thunder.
func (c *char) a1CDReset() {
	if c.Base.Ascension < 1 {
		return
	}
	c.ResetActionCooldown(action.ActionSkill)
}

// When Razor's Energy is below 50%, increases Energy Recharge by 30%.
func (c *char) a4() {
	if c.Base.Ascension < 4 {
		return
	}
	c.a4Bonus = make([]float64, attributes.EndStatType)
	c.a4Bonus[attributes.ER] = 0.3
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("razor-a4", -1),
		AffectedStat: attributes.ER,
		Amount: func() ([]float64, bool) {
			if c.Energy/c.EnergyMax >= 0.5 {
				return nil, false
			}
			return c.a4Bonus, true
		},
	})
}

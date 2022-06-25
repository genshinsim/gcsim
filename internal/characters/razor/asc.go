package razor

import "github.com/genshinsim/gcsim/pkg/core/attributes"

// When Razor's Energy is below 50%, increases Energy Recharge by 30%.
func (c *char) a4() {
	val := make([]float64, attributes.EndStatType)
	val[attributes.ER] = 0.3
	c.AddStatMod("er-sigil", -1, attributes.ER, func() ([]float64, bool) {
		if c.Energy/c.EnergyMax < 0.5 {
			return nil, false
		}

		return val, true
	})
}

package yaemiko

import "github.com/genshinsim/gcsim/pkg/core/attributes"

// When Sesshou Sakura lightning hits opponents, the Electro DMG Bonus of all nearby party members is increased by 20% for 5s.
func (c *char) c4() {
	m := make([]float64, attributes.EndStatType)
	m[attributes.ElectroP] = .20

	// TODO: does this trigger for yaemiko too? assuming it does
	for _, char := range c.Core.Player.Chars() {
		char.AddStatMod("yaemiko-c4", 5*60, attributes.ElectroP, func() ([]float64, bool) {
			return m, true
		})
	}
}

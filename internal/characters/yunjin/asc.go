package yunjin

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

// A1 is not implemented:
// TODO: Using Opening Flourish at the precise moment when Yun Jin is attacked will unleash its Level 2 Charged (Hold) form.

// count elemental types
func (c *char) a4Init() {
	if c.Base.Ascension < 4 {
		return
	}
	partyElementalTypes := make(map[attributes.Element]int)
	for _, char := range c.Core.Player.Chars() {
		partyElementalTypes[char.Base.Element]++
	}
	for range partyElementalTypes {
		c.partyElementalTypes += 1
	}
	c.Core.Log.NewEvent("Yun Jin Party Elemental Types (A4)", glog.LogCharacterEvent, c.Index).
		Write("party_elements", c.partyElementalTypes)
}

// The Normal Attack DMG Bonus granted by Flying Cloud Flag Formation is further increased by
// 2.5%/5%/7.5%/11.5% of Yun Jin's DEF when the party contains characters of 1/2/3/4 Elemental Types, respectively.
func (c *char) a4() float64 {
	if c.Base.Ascension < 4 {
		return 0
	}
	if c.partyElementalTypes == 4 {
		return 0.115
	}
	return 0.025 * float64(c.partyElementalTypes)
}

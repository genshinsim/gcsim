package zhongli

// When the Jade Shield takes DMG, it will Fortify:
//
// - Fortified characters have 5% increased Shield Strength.
//
// - Can stack up to 5 times, and lasts until the Jade Shield disappears.
func (c *char) a1() {
	if c.Base.Ascension < 1 {
		return
	}
	c.Core.Player.Shields.AddShieldBonusMod("zhongli-a1", -1, func() (float64, bool) {
		if c.Tags["shielded"] == 0 {
			return 0, false
		}
		return float64(c.Tags["a1"]) * 0.05, true
	})
}

// Zhongli deals bonus DMG based on his Max HP:
//
// - Normal Attack, Charged Attack, and Plunging Attack DMG is increased by 1.39% of Max HP.
func (c *char) a4Attacks() float64 {
	if c.Base.Ascension < 4 {
		return 0
	}
	return 0.0139 * c.MaxHP()
}

// Zhongli deals bonus DMG based on his Max HP:
//
// - Dominus Lapidis' Stone Stele, resonance, and hold DMG is increased by 1.9% of Max HP.
func (c *char) a4Skill() float64 {
	if c.Base.Ascension < 4 {
		return 0
	}
	return 0.019 * c.MaxHP()
}

// Zhongli deals bonus DMG based on his Max HP:
//
// - Planet Befall's DMG is increased by 33% of Max HP.
func (c *char) a4Burst() float64 {
	if c.Base.Ascension < 4 {
		return 0
	}
	return 0.33 * c.MaxHP()
}

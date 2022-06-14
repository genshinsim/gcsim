package zhongli

func (c *char) a1() {
	c.Core.Player.Shields.AddShieldBonusMod("zhongli-a1", -1, func() (float64, bool) {
		if c.Tags["shielded"] == 0 {
			return 0, false
		}
		return float64(c.Tags["a1"]) * 0.05, true
	})
}

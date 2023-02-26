package xiangling

// Increases the flame range of Guoba by 20%.
func (c *char) a1() {
	if c.Base.Ascension < 1 {
		return
	}
	c.guobaFlameRange *= 1.2
}

// A4 is not implemented:
// TODO: When Guoba Attack's effect ends, Guoba leaves a chili pepper on the spot where it disappeared. Picking up a chili pepper increases ATK by 10% for 10s.

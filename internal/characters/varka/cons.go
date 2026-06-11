package varka

func (c *char) c1OnSkill() {
	if c.Base.Cons < 1 {
		return
	}
	c.c1Extra = 1
	c.fourWindsChargesAva = 1
}

func (c *char) c1OnSpecialSkill() float64 {
	if c.Base.Cons < 1 {
		return 1.0
	}

	if c.c1Extra > 0 {
		c.c1Extra = 0
		return 2.0
	}
	return 1.0
}

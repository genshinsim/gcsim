package character

func (c *CharWrapper) ResetNormalCounter() {
	c.NormalCounter = 0
}

func (c *CharWrapper) AdvanceNormalIndex() {
	c.NormalCounter++
	if c.NormalCounter == c.NormalHitNum {
		c.NormalCounter = 0
	}
}

func (c *CharWrapper) NextNormalCounter() int {
	return c.NormalCounter + 1
}

package character

//advance normal index, return the current one
func (c *Tmpl) AdvanceNormalIndex() {
	c.NormalCounter++
	if c.NormalCounter == c.NormalHitNum {
		c.NormalCounter = 0
	}
}

func (c *Tmpl) NextNormalCounter() int {
	return c.NormalCounter + 1
}

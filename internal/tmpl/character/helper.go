package character

//advance normal index, return the current one
func (c *Tmpl) AdvanceNormalIndex() {
	c.NormalCounter++
	if c.NormalCounter == c.NormalHitNum {
		c.NormalCounter = 0
	}
}

func (c *Tmpl) CurrentNormalCounter() int {
	return c.NormalCounter
}

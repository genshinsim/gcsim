package character

import "github.com/genshinsim/gcsim/pkg/core"

//advance normal index, return the current one
func (c *Tmpl) AdvanceNormalIndex() {
	c.NormalCounter++
	if c.NormalCounter == c.NormalHitNum {
		c.NormalCounter = 0
	}
}

// c.NormalCounter returns the next normal index, and always resets
// when it reaches the max hit of the normal attack combo
func (c *Tmpl) NextNormalIndex() int {
	if c.NormalCounter == 0 && c.Core.LastAction.Typ == core.ActionAttack {
		return c.NormalHitNum
	}
	return c.NormalCounter
}

package character

import "github.com/genshinsim/gcsim/pkg/core/event"

func (c *Character) clampHPRatio() {
	if c.currentHPRatio > 1 {
		c.currentHPRatio = 1
	} else if c.currentHPRatio < 0 {
		c.currentHPRatio = 0
	}
}

func (c *Character) SetHPByAmount(amt float64) {
	c.currentHPRatio = amt / c.MaxHP()
	c.clampHPRatio()
}

func (c *Character) SetHPByRatio(r float64) {
	c.currentHPRatio = r
	c.clampHPRatio()
}

func (c *Character) ModifyHPByAmount(amt float64) {
	newHP := c.CurrentHP() + amt
	c.SetHPByAmount(newHP)
}

func (c *Character) ModifyHPByRatio(r float64) {
	newHPRatio := c.currentHPRatio + r
	c.SetHPByRatio(newHPRatio)
}

func (c *Character) clampHPDebt() {
	if c.currentHPDebt < 0 {
		c.currentHPDebt = 0
	}
}

func (c *Character) setHPDebtByAmount(amt float64) {
	c.currentHPDebt = amt
	c.clampHPDebt()
}

func (c *Character) ModifyHPDebtByAmount(amt float64) {
	if amt == 0 {
		return
	}
	newHPDebt := c.currentHPDebt + amt
	c.setHPDebtByAmount(newHPDebt)
	c.Core.Events.Emit(event.OnHPDebt, c.Index, amt)
}

func (c *Character) ModifyHPDebtByRatio(r float64) {
	amt := r * c.MaxHP()
	c.ModifyHPDebtByAmount(amt)
}

func (c *Character) CurrentHPRatio() float64 {
	return c.currentHPRatio
}

func (c *Character) CurrentHP() float64 {
	return c.currentHPRatio * c.MaxHP()
}

func (c *Character) CurrentHPDebt() float64 {
	return c.currentHPDebt
}

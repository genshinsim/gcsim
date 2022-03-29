package character

import (
	"github.com/genshinsim/gcsim/pkg/coretype"
)

func (c *Tmpl) Attack(p map[string]int) (int, int) {
	return 0, 0
}

func (c *Tmpl) Aimed(p map[string]int) (int, int) {
	return 0, 0
}

func (c *Tmpl) ChargeAttack(p map[string]int) (int, int) {
	return 0, 0
}

func (c *Tmpl) HighPlungeAttack(p map[string]int) (int, int) {
	return 0, 0
}

func (c *Tmpl) LowPlungeAttack(p map[string]int) (int, int) {
	return 0, 0
}

func (c *Tmpl) Skill(p map[string]int) (int, int) {
	return 0, 0
}

func (c *Tmpl) Burst(p map[string]int) (int, int) {
	return 0, 0
}

func (c *Tmpl) Dash(p map[string]int) (int, int) {
	return 24, 24
}

func (c *Tmpl) ActionStam(a coretype.ActionType, p map[string]int) float64 {
	switch a {
	case coretype.ActionDash:
		return 18
	default:
		c.Core.NewEvent("ActionStam not implemented", coretype.LogActionEvent, c.Index)
		return 0
	}
}

func (c *Tmpl) ResetNormalCounter() {
	c.NormalCounter = 0
}

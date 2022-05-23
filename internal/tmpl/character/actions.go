package character

import (
	"github.com/genshinsim/gcsim/pkg/core"
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
	return 24, 24 //likely too long
}

func (c *Tmpl) Jump(p map[string]int) (int, int) {
	return 33, 33 //likely too long
}

func (c *Tmpl) ActionStam(a core.ActionType, p map[string]int) float64 {
	switch a {
	case core.ActionDash:
		return 18
	default:
		c.Core.Log.NewEvent("ActionStam not implemented", core.LogActionEvent, c.Index)
		return 0
	}
}

func (c *Tmpl) ResetNormalCounter() {
	c.NormalCounter = 0
}

package actioncheck

import "github.com/genshinsim/gcsim/pkg/core/action"

type StamHandler struct {
	charge float64
	dash   float64
}

func New(chargeStam, dashStam float64) StamHandler {
	return StamHandler{
		charge: chargeStam,
		dash:   dashStam,
	}
}

func (c *StamHandler) ActionStam(a action.Action, p map[string]int) float64 {
	switch a {
	case action.ActionDash:
		return c.dash
	case action.ActionCharge:
		return c.charge
	default:
		return 0
	}
}

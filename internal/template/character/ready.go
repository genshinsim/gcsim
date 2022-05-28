package character

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

func (c *Character) ActionReady(a action.Action, p map[string]int) bool {
	//for dash and charge need to check for stam usage as well

	switch a {
	case action.ActionBurst:
		return c.Energy == c.EnergyMax
	case action.ActionSkill:
		return c.AvailableCDCharge[a] > 0
	case action.ActionCharge:
		req := c.Core.Player.AbilStamCost(c.Index, a, p)
		if c.Core.Player.Stam < req {
			c.Core.Log.NewEvent("insufficient stam: charge attack", glog.LogSimEvent, -1, "have", c.Core.Player.Stam)
			return false
		}
		return true
	case action.ActionDash:
		req := c.Core.Player.AbilStamCost(c.Index, a, p)
		if c.Core.Player.Stam < req {
			c.Core.Log.NewEvent("insufficient stam: dash", glog.LogSimEvent, -1, "have", c.Core.Player.Stam)
			return false
		}
		return true
	default:
		//all other actions should be always ready?
		return true
	}
}

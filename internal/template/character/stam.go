package character

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
)

//ActionStam provides default implementation for stam cost for charge and dash
//character should override this though
func (c *Character) ActionStam(a action.Action, p map[string]int) float64 {
	switch a {
	case action.ActionCharge:
		//20 sword (most)
		//25 polearm
		//40 per second claymore
		//50 catalyst
		switch c.Weapon.Class {
		case weapon.WeaponClassSword:
			return 20
		case weapon.WeaponClassSpear:
			return 25
		case weapon.WeaponClassCatalyst:
			return 50
		case weapon.WeaponClassClaymore:
			c.Core.Log.NewEvent("CLAYMORE CHARGE NOT IMPLEMENTED", glog.LogSimEvent, c.Index)
			return 0
		case weapon.WeaponClassBow:
			c.Core.Log.NewEvent("BOWS DONT HAVE CHARGE ATTACK; USE AIM", glog.LogSimEvent, c.Index)
			return 0
		default:
			return 0
		}
	case action.ActionDash:
		//18 per
		return 18
	default:
		return 0
	}
}

var defaultDash = action.ActionInfo{
	Frames:        func(action.Action) int { return 20 },
	CanQueueAfter: 20,
	Post:          20,
	State:         action.DashState,
}

func (c *Character) Dash(p map[string]int) action.ActionInfo {
	//consume stam at the end

	c.Core.Tasks.Add(func() {
		req := c.Core.Player.AbilStamCost(c.Index, action.ActionDash, p)
		c.Core.Player.Stam -= req
		//this really shouldn't happen??
		if c.Core.Player.Stam < 0 {
			c.Core.Player.Stam = 0
		}
		c.Core.Player.LastStamUse = c.Core.F
		c.Core.Events.Emit(event.OnStamUse, action.DashState)
	}, 19)
	return defaultDash
}

var defaultJump = action.ActionInfo{
	Frames:          func(action.Action) int { return 30 },
	AnimationLength: 30,
	CanQueueAfter:   30,
	Post:            30,
	State:           action.JumpState,
}

func (c *Character) Jump(p map[string]int) action.ActionInfo {
	return defaultJump
}

func (c *Character) Walk(p map[string]int) action.ActionInfo {
	f, ok := p["f"]
	if !ok {
		f = 1
	}
	return action.ActionInfo{
		Frames:          func(next action.Action) int { return f },
		AnimationLength: f,
		CanQueueAfter:   f,
		Post:            f,
		State:           action.WalkState,
	}
}

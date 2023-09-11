package character

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

// ActionStam provides default implementation for stam cost for charge and dash
// character should override this though
func (c *Character) ActionStam(a action.Action, p map[string]int) float64 {
	switch a {
	case action.ActionCharge:
		//20 sword (most)
		//25 polearm
		//40 per second claymore
		//50 catalyst
		switch c.Weapon.Class {
		case info.WeaponClassSword:
			return 20
		case info.WeaponClassSpear:
			return 25
		case info.WeaponClassCatalyst:
			return 50
		case info.WeaponClassClaymore:
			c.Core.Log.NewEvent("CLAYMORE CHARGE NOT IMPLEMENTED", glog.LogWarnings, c.Index)
			return 0
		case info.WeaponClassBow:
			c.Core.Log.NewEvent("BOWS DONT HAVE CHARGE ATTACK; USE AIM", glog.LogWarnings, c.Index)
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

func (c *Character) Dash(p map[string]int) action.ActionInfo {
	// Execute dash CD logic
	c.ApplyDashCD()

	// consume stamina at end of the dash
	c.QueueDashStaminaConsumption(p)

	length := c.DashLength()
	return action.ActionInfo{
		Frames:          func(action.Action) int { return length },
		AnimationLength: length,
		CanQueueAfter:   length,
		State:           action.DashState,
	}
}

// set the dash CD. If the dash was on CD when this dash executes, lockout dash
func (c *Character) ApplyDashCD() {
	var evt glog.Event

	if c.Core.Player.DashCDExpirationFrame > c.Core.F {
		c.Core.Player.DashLockout = true
		c.Core.Player.DashCDExpirationFrame = c.Core.F + 1.5*60
		evt = c.Core.Log.NewEvent("dash cooldown triggered", glog.LogCooldownEvent, c.Index)
	} else {
		c.Core.Player.DashLockout = false
		c.Core.Player.DashCDExpirationFrame = c.Core.F + 0.8*60
		evt = c.Core.Log.NewEvent("dash lockout evaluation started", glog.LogCooldownEvent, c.Index)
	}

	evt.Write("lockout", c.Core.Player.DashLockout).
		Write("expiry", c.Core.Player.DashCDExpirationFrame-c.Core.F).
		Write("expiry_frame", c.Core.Player.DashCDExpirationFrame)
}

func (c *Character) QueueDashStaminaConsumption(p map[string]int) {
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
	}, c.DashLength()-1)
}

func (c *Character) DashLength() int {
	switch c.CharBody {
	case info.BodyBoy, info.BodyLoli:
		return 21
	case info.BodyMale:
		return 19
	case info.BodyLady:
		return 22
	default:
		return 20
	}
}

func (c *Character) Jump(p map[string]int) action.ActionInfo {
	f := 30
	switch c.CharBody {
	case info.BodyBoy, info.BodyGirl:
		f = 31
	case info.BodyMale:
		f = 28
	case info.BodyLady:
		f = 32
	case info.BodyLoli:
		f = 29
	}
	return action.ActionInfo{
		Frames:          func(action.Action) int { return f },
		AnimationLength: f,
		CanQueueAfter:   f,
		State:           action.JumpState,
	}
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
		State:           action.WalkState,
	}
}

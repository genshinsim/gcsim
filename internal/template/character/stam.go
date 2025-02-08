package character

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player"
)

// ActionStam provides default implementation for stam cost for charge and dash
// character should override this though
func (c *Character) ActionStam(a action.Action, p map[string]int) float64 {
	switch a {
	case action.ActionCharge:
		// 20 sword (most)
		// 25 polearm
		// 40 per second claymore
		// 50 catalyst
		switch c.Weapon.Class {
		case info.WeaponClassSword:
			return 20
		case info.WeaponClassSpear:
			return 25
		case info.WeaponClassCatalyst:
			return 50
		case info.WeaponClassClaymore:
			return 0
		case info.WeaponClassBow:
			return 0
		default:
			return 0
		}
	case action.ActionDash:
		// 18 per
		return 18
	default:
		return 0
	}
}

func (c *Character) Dash(p map[string]int) (action.Info, error) {
	// Execute dash CD logic
	c.ApplyDashCD()

	// consume stamina at end of the dash
	c.QueueDashStaminaConsumption(p)

	length := c.DashLength()
	dashJumpLength := c.DashToJumpLength()
	return action.Info{
		Frames: func(a action.Action) int {
			switch a {
			case action.ActionJump:
				return dashJumpLength
			default:
				return length
			}
		},
		AnimationLength: length,
		CanQueueAfter:   dashJumpLength,
		State:           action.DashState,
	}, nil
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
	// consume stam at the end
	c.Core.Tasks.Add(func() {
		req := c.Core.Player.AbilStamCost(c.Index, action.ActionDash, p)
		c.Core.Player.Stam -= req
		// this really shouldn't happen??
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

func (c *Character) DashToJumpLength() int {
	switch c.CharBody {
	case info.BodyGirl, info.BodyLoli:
		return 4
	case info.BodyBoy:
		return 2
	default:
		return 3
	}
}

func (c *Character) Jump(p map[string]int) (action.Info, error) {
	if c.StatusIsActive(player.XianyunAirborneBuff) {
		c.Core.Player.SetAirborne(player.AirborneXianyun)
		// 4/8 for claymore/bow/catalyst and 5/9 for sword/polearm
		lowPlunge := 4
		highPlunge := 8
		switch c.Weapon.Class {
		case info.WeaponClassSword, info.WeaponClassSpear:
			lowPlunge = 5
			highPlunge = 9
		}

		animLength := 60 // Upperbound for jump for high/low plunge
		return action.Info{
			Frames: func(a action.Action) int {
				switch a {
				case action.ActionLowPlunge:
					return lowPlunge
				case action.ActionHighPlunge:
					return highPlunge
				default:
					return animLength // This is expected to later lead to action error because no other action besides plunges can be done while AirborneXianyun
				}
			},
			AnimationLength: animLength,
			CanQueueAfter:   lowPlunge, // earliest cancel
			State:           action.JumpState,
		}, nil
	}
	f := c.JumpLength()
	return action.Info{
		Frames:          func(action.Action) int { return f },
		AnimationLength: f,
		CanQueueAfter:   f,
		State:           action.JumpState,
	}, nil
}

func (c *Character) JumpLength() int {
	if c.Core.Player.LastAction.Type == action.ActionDash {
		switch c.CharBody {
		case info.BodyGirl, info.BodyBoy:
			return 34
		default:
			return 37
		}
	}
	switch c.CharBody {
	case info.BodyBoy, info.BodyGirl:
		return 31
	case info.BodyMale:
		return 28
	case info.BodyLady:
		return 32
	case info.BodyLoli:
		return 29
	default:
		return 30
	}
}

func (c *Character) Walk(p map[string]int) (action.Info, error) {
	f, ok := p["f"]
	if !ok {
		f = 1
	}
	return action.Info{
		Frames:          func(action.Action) int { return f },
		AnimationLength: f,
		CanQueueAfter:   f,
		State:           action.WalkState,
	}, nil
}

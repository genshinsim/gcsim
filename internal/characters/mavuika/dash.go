package mavuika

import (
	"errors"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

const (
	dashChargeCancelErr = "can only begin a charge cancelled dash while Mavuika is in the Flamestrider Charged Attack (Cyclic or Final) animation"
)

func (c *char) Dash(p map[string]int) (action.Info, error) {
	if c.armamentState != bike {
		return c.Character.Dash(p)
	}
	if !c.nightsoulState.HasBlessing() {
		return c.Character.Dash(p)
	}

	// Dash will do no damage, no matter how long it takes.
	// Frames for swap and jump will be set to 2f instead of travel frames.
	// Frames for charge will be set to 2f if allowed to early cancel
	collision, ok := p["collision"]
	if !ok {
		collision = 1
	}

	// Dash has a travelling hitbox that can hit as soon as 6f after start,
	// or as late as dash transitions to idle (exact end not tested)
	travel, ok := p["travel"]
	if !ok {
		travel = 6
	}
	if travel > 24 {
		travel = 24
	}
	if travel < 6 {
		travel = 6
	}

	// Travel is set to 2 to represent cancel frames instead of hitmark, if damage is skipped
	if collision == 0 {
		travel = 2
	}

	c.dashFrames[action.ActionCharge] = max(20, travel)
	c.dashFrames[action.ActionAttack] = max(18, travel)
	c.dashFrames[action.ActionSkill] = max(20, travel)
	c.dashFrames[action.ActionBurst] = max(20, travel)
	c.dashFrames[action.ActionSwap] = travel
	c.dashFrames[action.ActionJump] = travel

	// Only applies when preceded and followed by a charge, otherwise does nothing
	// If Mav is within the "clock lockout" of being able to cdc,
	// early cancel will be forced to false
	c.chargeCancel = false
	chargeCancel, ok := p["charge_cancel"]
	if !ok {
		chargeCancel = 0
	}
	if c.StatusIsActive(cdcLockoutStatus) {
		if chargeCancel != 0 {
			c.Core.Log.NewEvent(
				"Blocked from performing an early dash cancel- Mavuika within lockout frames.",
				glog.LogWarnings,
				c.Index(),
			)
		}
		chargeCancel = 0
	}
	if chargeCancel != 0 {
		if c.armamentState != bike {
			return action.Info{}, errors.New(dashChargeCancelErr)
		}
		if !c.nightsoulState.HasBlessing() {
			return action.Info{}, errors.New(dashChargeCancelErr)
		}

		// Used for n0 proc logic in charge.go
		c.chargeCancel = true
		c.dashFrames[action.ActionCharge] = travel
	}

	if collision != 0 {
		ai := info.AttackInfo{
			ActorIndex:     c.Index(),
			Abil:           "Flamestrider Sprint",
			AttackTag:      attacks.AttackTagNone,
			ICDTag:         attacks.ICDTagMavuikaFlamestrider,
			AdditionalTags: []attacks.AdditionalTag{attacks.AdditionalTagNightsoul},
			ICDGroup:       attacks.ICDGroupDefault,
			StrikeType:     attacks.StrikeTypeBlunt,
			PoiseDMG:       75.0,
			Element:        attributes.Pyro,
			Durability:     25,
			Mult:           skillDash[c.TalentLvlSkill()],
			IsDeployable:   true,
		}
		ap := combat.NewCircleHitOnTarget(
			c.Core.Combat.Player(),
			info.Point{Y: 1.0},
			1.2,
		)
		c.QueueCharTask(func() { c.Core.QueueAttack(ai, ap, 0, 0) }, travel)
	}

	c.reduceNightsoulPoints(10)
	if c.Core.Player.CurrentState() == action.NormalAttackState {
		// If dashing from NA while in bike, do not reset NA string
		c.savedNormalCounter = c.NormalCounter
	}

	// Execute dash CD logic
	c.ApplyDashCD()
	return action.Info{
		Frames:          frames.NewAbilFunc(c.dashFrames),
		AnimationLength: c.dashFrames[action.InvalidAction],
		CanQueueAfter:   travel,
		State:           action.DashState,
	}, nil
}

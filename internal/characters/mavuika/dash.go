package mavuika

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var dashFrames []int

func init() {
	dashFrames = frames.InitAbilSlice(24) // Dash -> Dash
}

func (c *char) Dash(p map[string]int) (action.Info, error) {
	// Dash will do no damage, no matter how long it takes.
	// Frames for swap and jump will be set to 2f instead of travel frames.
	// Frames for charge will be set to 2f if allowed to early cancel
	skipDmg, ok := p["skip_dmg"]
	if !ok {
		skipDmg = 0
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
	if skipDmg != 0 {
		travel = 2
	}

	// Only applies when preceded and followed by a charge, otherwise does nothing
	// If Mav is within the "clock lockout" of being able to cdc,
	//   early cancel will be forced to false
	earlyCancellable, ok := p["early_cancellable"]
	if !ok {
		earlyCancellable = 1
	}
	if c.StatusIsActive(cdcLockoutStatus) {
		if earlyCancellable != 0 {
			c.Core.Log.NewEvent(
				"Blocked from performing an early dash cancel- Mavuika within lockout frames.",
				glog.LogWarnings,
				c.Index(),
			)
		}
		earlyCancellable = 0
	}

	dashFrames[action.ActionCharge] = max(20, travel)
	dashFrames[action.ActionAttack] = max(18, travel)
	dashFrames[action.ActionSkill] = max(20, travel)
	dashFrames[action.ActionBurst] = max(20, travel)
	dashFrames[action.ActionSwap] = travel
	dashFrames[action.ActionJump] = travel

	if c.armamentState == bike && c.nightsoulState.HasBlessing() {
		if skipDmg != 0 {
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
				HitlagFactor:   0.05,
				IsDeployable:   true,
			}
			ap := combat.NewCircleHitOnTarget(
				c.Core.Combat.Player(),
				info.Point{Y: 1.0},
				1.2,
			)
			c.Core.QueueAttack(ai, ap, travel, travel)
		}

		c.reduceNightsoulPoints(10)
		x := c.Core.Player.CurrentState()
		c.isDashFromCA = false
		switch x {
		case action.NormalAttackState:
			// If dashing from NA while in bike, do not reset NA string
			c.savedNormalCounter = c.NormalCounter
		case action.ChargeAttackState:
			if earlyCancellable != 0 && c.armamentState == bike && c.nightsoulState.HasBlessing() {
				// Used for n0 proc logic in charge.go
				c.isDashFromCA = true
				dashFrames[action.ActionCharge] = travel
			}
		default:
		}

		// Execute dash CD logic
		c.ApplyDashCD()
		return action.Info{
			Frames:          frames.NewAbilFunc(dashFrames),
			AnimationLength: dashFrames[action.InvalidAction],
			CanQueueAfter:   travel,
			State:           action.DashState,
		}, nil
	}

	return c.Character.Dash(p)
}

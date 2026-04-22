package mavuika

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var dashFrames []int

func init() {
	dashFrames = frames.InitAbilSlice(24) // Dash -> Dash
	dashFrames[action.ActionAttack] = 18
	dashFrames[action.ActionSkill] = 20
	dashFrames[action.ActionBurst] = 20
	dashFrames[action.ActionSwap] = 0
	dashFrames[action.ActionJump] = 0
}

func (c *char) Dash(p map[string]int) (action.Info, error) {
	travel, ok := p["travel"]
	if !ok {
		travel = 6
	}
	if travel > 24 {
		travel = 24
	}
	if travel < 1 {
		travel = 1
	}

	// Only applies when preceded and followed by a charge, otherwise does nothing
	earlyCancellable, ok := p["early_cancellable"]
	if !ok {
		earlyCancellable = 1
	}

	dashFrames[action.ActionCharge] = 20

	if c.armamentState == bike && c.nightsoulState.HasBlessing() {
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
		c.reduceNightsoulPoints(10)
		x := c.Core.Player.CurrentState()
		c.isDashFromCA = false
		switch x {
		case action.NormalAttackState:
			// If dashing from NA while in bike, do not reset NA string
			c.savedNormalCounter = c.NormalCounter
		case action.ChargeAttackState:
			if earlyCancellable != 0 {
				// Used for n0 proc logic in charge.go
				c.isDashFromCA = true
				dashFrames[action.ActionCharge] = 0
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

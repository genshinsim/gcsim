package mavuika

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
)

var dashFrames []int

func init() {
	dashFrames = frames.InitAbilSlice(24) // Dash -> Dash
	dashFrames[action.ActionAttack] = 18
	dashFrames[action.ActionCharge] = 20
	dashFrames[action.ActionSkill] = 20
	dashFrames[action.ActionBurst] = 20
	dashFrames[action.ActionSwap] = 0
	dashFrames[action.ActionJump] = 0
}

func (c *char) Dash(p map[string]int) (action.Info, error) {
	if c.armamentState == bike && c.nightsoulState.HasBlessing() {
		ai := combat.AttackInfo{
			ActorIndex:     c.Index,
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
			geometry.Point{Y: 1.0},
			1.2,
		)
		c.Core.QueueAttack(ai, ap, 6, 6)
		c.reduceNightsoulPoints(10)
		// If dashing from NA while in bike, do not reset NA string
		if c.Core.Player.CurrentState() == action.NormalAttackState {
			c.savedNormalCounter = c.NormalCounter
		}

		// Execute dash CD logic
		c.ApplyDashCD()
		return action.Info{
			Frames:          frames.NewAbilFunc(dashFrames),
			AnimationLength: dashFrames[action.InvalidAction],
			CanQueueAfter:   dashFrames[action.ActionJump],
			State:           action.DashState,
		}, nil
	}

	return c.Character.Dash(p)
}

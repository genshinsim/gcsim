package alhaitham

import (
	"errors"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
)

var lowPlungeFrames []int

const lowPlungeHitmark = 38

func init() {
	lowPlungeFrames = frames.InitAbilSlice(70)
	lowPlungeFrames[action.ActionAttack] = 49
	lowPlungeFrames[action.ActionSkill] = 50
	lowPlungeFrames[action.ActionBurst] = 50
	lowPlungeFrames[action.ActionDash] = 40
	lowPlungeFrames[action.ActionSwap] = 58
}

func (c *char) LowPlungeAttack(p map[string]int) (action.Info, error) {
	// last action must be hold skill
	if c.Core.Player.LastAction.Type != action.ActionSkill ||
		c.Core.Player.LastAction.Param["hold"] != 1 {
		return action.Info{}, errors.New("only plunge after hold skill ends")
	}

	short := p["short"]
	skip := 0
	if short > 0 {
		skip = 20
	}

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Low Plunge Attack",
		AttackTag:  attacks.AttackTagPlunge,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeBlunt,
		PoiseDMG:   100,
		Element:    attributes.Dendro,
		Durability: 25,
		Mult:       lowPlunge[c.TalentLvlAttack()],
	}

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), geometry.Point{Y: 1}, 3),
		lowPlungeHitmark-skip,
		lowPlungeHitmark-skip,
		c.makeA1CB(), // A1 adds a stack before the mirror count for the Projection Attack is determined
		c.projectionAttack,
	)

	return action.Info{
		Frames:          func(next action.Action) int { return lowPlungeFrames[next] - skip },
		AnimationLength: lowPlungeFrames[action.InvalidAction] - skip,
		CanQueueAfter:   lowPlungeFrames[action.ActionDash] - skip,
		State:           action.PlungeAttackState,
	}, nil
}

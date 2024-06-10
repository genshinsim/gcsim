package chiori

import (
	"errors"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/player"
)

var lowPlungeFramesC []int

const lowPlungeHitmarkC = 41

const lowPlungePoiseDMG = 100.0
const lowPlungeRadius = 3.0

func init() {
	lowPlungeFramesC = frames.InitAbilSlice(69) // Low Plunge -> J/Walk
	lowPlungeFramesC[action.ActionAttack] = 54
	lowPlungeFramesC[action.ActionSkill] = 53
	lowPlungeFramesC[action.ActionBurst] = 54
	lowPlungeFramesC[action.ActionDash] = 44
	lowPlungeFramesC[action.ActionSwap] = 56
}

func (c *char) LowPlungeAttack(p map[string]int) (action.Info, error) {
	defer c.Core.Player.SetAirborne(player.Grounded)
	// last action hold skill
	if c.Core.Player.LastAction.Type == action.ActionSkill &&
		c.Core.Player.LastAction.Param["hold"] >= 1 {
		return c.lowPlungeC(), nil
	}

	return action.Info{}, errors.New("low_plunge can only be used after hold skill")
}

func (c *char) lowPlungeC() action.Info {
	c.tryTriggerA1TailoringNA()

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Low Plunge Attack",
		AttackTag:  attacks.AttackTagPlunge,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeBlunt,
		PoiseDMG:   lowPlungePoiseDMG,
		Element:    attributes.Physical,
		Durability: 25,
		Mult:       lowPlunge[c.TalentLvlAttack()],
	}

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), geometry.Point{Y: 1}, lowPlungeRadius),
		lowPlungeHitmarkC,
		lowPlungeHitmarkC,
	)

	return action.Info{
		Frames:          frames.NewAbilFunc(lowPlungeFramesC),
		AnimationLength: lowPlungeFramesC[action.InvalidAction],
		CanQueueAfter:   lowPlungeFramesC[action.ActionDash],
		State:           action.PlungeAttackState,
	}
}

package ayaka

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var dashFrames []int

const dashHitmark = 20

func init() {
	dashFrames = frames.InitAbilSlice(35)
	dashFrames[action.ActionDash] = 30
	dashFrames[action.ActionSwap] = 34
}

// TODO: move this into PostDash event instead
func (c *char) Dash(p map[string]int) action.ActionInfo {
	f, ok := p["f"]
	if !ok {
		f = 0
	}

	//no dmg attack at end of dash
	ai := combat.AttackInfo{
		Abil:       "Dash",
		ActorIndex: c.Index,
		AttackTag:  attacks.AttackTagNone,
		ICDTag:     attacks.ICDTagDash,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Cryo,
		Durability: 25,
	}

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), combat.Point{Y: 0.1}, 2),
		dashHitmark+f,
		dashHitmark+f,
		c.makeA4CB(),
	)

	//add cryo infuse
	//TODO: check weapon infuse timing; this SHOULD be ok?
	c.Core.Tasks.Add(func() {
		c.Core.Player.AddWeaponInfuse(
			c.Index,
			"ayaka-dash",
			attributes.Cryo,
			300,
			true,
			attacks.AttackTagNormal, attacks.AttackTagExtra, attacks.AttackTagPlunge,
		)
	}, dashHitmark+f)

	// handle stamina usage, avoid default dash implementation since dont want CD
	c.QueueDashStaminaConsumption(p)

	return action.ActionInfo{
		Frames:          func(next action.Action) int { return dashFrames[next] + f },
		AnimationLength: dashFrames[action.InvalidAction] + f,
		CanQueueAfter:   dashFrames[action.ActionDash] + f, // earliest cancel
		State:           action.DashState,
	}
}

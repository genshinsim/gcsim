package mona

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var dashFrames []int

const dashHitmark = 20

func init() {
	dashFrames = frames.InitAbilSlice(42) // D -> N1
	dashFrames[action.ActionCharge] = 36  // D -> CA
	dashFrames[action.ActionSkill] = 35   // D -> E
	dashFrames[action.ActionBurst] = 21   // D -> Q
	dashFrames[action.ActionDash] = 30    // D -> D
	dashFrames[action.ActionJump] = 500   // D -> J, TODO: this action is illegal; need better way to handle it
	dashFrames[action.ActionSwap] = 34    // D -> Swap
}

func (c *char) Dash(p map[string]int) (action.Info, error) {
	f, ok := p["f"]
	if !ok {
		f = 0
	}
	// no dmg attack at end of dash
	ai := info.AttackInfo{
		Abil:       "Dash",
		ActorIndex: c.Index(),
		AttackTag:  attacks.AttackTagNone,
		ICDTag:     attacks.ICDTagDash,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Hydro,
		Durability: 25,
	}
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), info.Point{Y: 0.1}, 2),
		dashHitmark+f,
		dashHitmark+f,
	)

	// A1
	c.Core.Tasks.Add(c.a1, 120)

	c.c6OnDash()

	// handle stamina usage, avoid default dash implementation since dont want CD
	c.QueueDashStaminaConsumption(p)

	return action.Info{
		Frames:          func(next action.Action) int { return dashFrames[next] + f },
		AnimationLength: dashFrames[action.InvalidAction] + f,
		CanQueueAfter:   dashHitmark + f,
		State:           action.DashState,
		OnRemoved:       c.c6OnDashEnd,
	}, nil
}

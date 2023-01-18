package mona

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
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

func (c *char) Dash(p map[string]int) action.ActionInfo {
	f, ok := p["f"]
	if !ok {
		f = 0
	}
	// no dmg attack at end of dash
	ai := combat.AttackInfo{
		Abil:       "Dash",
		ActorIndex: c.Index,
		AttackTag:  combat.AttackTagNone,
		ICDTag:     combat.ICDTagDash,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeDefault,
		Element:    attributes.Hydro,
		Durability: 25,
	}
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), combat.Point{Y: 0.1}, 2),
		dashHitmark+f,
		dashHitmark+f,
	)

	// A1
	c.Core.Tasks.Add(c.a1(), 120)
	// C6
	if c.Base.Cons >= 6 {
		// reset c6 stacks in case we dash again before using a CA
		c.c6Stacks = 0
		// need to keep track of src in case of Mona Dash Dash, where the second dash starts between two c6 ticks
		// without a src check the second Dash would gain a stack before 1s is up and a second one at 1s
		c.c6Src = c.Core.F
		c.Core.Tasks.Add(c.c6(c.Core.F), 60)
	}

	// call default implementation to handle stamina
	c.Character.Dash(p)

	return action.ActionInfo{
		Frames:          func(next action.Action) int { return dashFrames[next] + f },
		AnimationLength: dashFrames[action.InvalidAction] + f,
		CanQueueAfter:   dashHitmark + f,
		State:           action.DashState,
	}
}

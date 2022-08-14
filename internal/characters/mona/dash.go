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
		Element:    attributes.Hydro,
		Durability: 25,
	}
	c.Core.QueueAttack(ai, combat.NewCircleHit(c.Core.Combat.Player(), 2, false, combat.TargettableEnemy), dashHitmark+f, dashHitmark+f)

	// A1
	c.Core.Tasks.Add(c.a1(), 120)

	// call default implementation to handle stamina
	c.Character.Dash(p)

	return action.ActionInfo{
		Frames:          func(next action.Action) int { return dashFrames[next] + f },
		AnimationLength: dashFrames[action.InvalidAction] + f,
		CanQueueAfter:   dashHitmark + f,
		State:           action.DashState,
	}
}

package mona

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var dashFrames []int

const dashHitmark = 36

func (c *char) Dash(p map[string]int) action.ActionInfo {
	f, ok := p["f"]
	if !ok {
		f = 0
	}
	//no dmg attack at end of dash
	ai := combat.AttackInfo{
		Abil:       "Dash",
		ActorIndex: c.Index,
		AttackTag:  combat.AttackTagNone,
		ICDTag:     combat.ICDTagDash,
		ICDGroup:   combat.ICDGroupDefault,
		Element:    attributes.Hydro,
		Durability: 25,
	}
	c.Core.QueueAttack(ai, combat.NewDefCircHit(2, false, combat.TargettableEnemy), dashHitmark+f, dashHitmark+f)

	//After she has used Illusory Torrent for 2s, if there are any opponents nearby,
	//Mona will automatically create a Phantom.
	//A Phantom created in this manner lasts for 2s, and its explosion DMG is equal to 50% of Mirror Reflection of Doom.

	//TODO: a4 not implemented. needs to know if this can be created while already on one field
	//and if it overrides

	// call default implementation to handle stamina
	c.Character.Dash(p)

	return action.ActionInfo{
		Frames:          func(next action.Action) int { return dashFrames[next] + f },
		AnimationLength: dashFrames[action.InvalidAction] + f,
		CanQueueAfter:   dashHitmark + f,
		State:           action.DashState,
	}
}

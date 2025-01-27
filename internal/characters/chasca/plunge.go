package chasca

import (
	"errors"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var lowPlungeFrames []int

const lowPlungeHitmark = 57
const collisionHitmark = lowPlungeHitmark - 6

func init() {
	lowPlungeFrames = frames.InitAbilSlice(94)
	lowPlungeFrames[action.ActionAttack] = 71
	lowPlungeFrames[action.ActionCharge] = 70
	lowPlungeFrames[action.ActionBurst] = 69
	lowPlungeFrames[action.ActionDash] = lowPlungeHitmark
	lowPlungeFrames[action.ActionJump] = 84
	lowPlungeFrames[action.ActionWalk] = 85
}

func (c *char) LowPlungeAttack(p map[string]int) (action.Info, error) {
	// Not in falling state
	if !c.StatusIsActive(plungeAvailableKey) {
		return action.Info{}, errors.New("only plunge after skill ends")
	}
	c.DeleteStatus(plungeAvailableKey)

	collision, ok := p["collision"]
	if !ok {
		collision = 0 // Whether or not Wanderer does a collision hit
	}

	if collision > 0 {
		c.plungeCollision(collisionHitmark)
	}

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Low Plunge Attack",
		AttackTag:  attacks.AttackTagPlunge,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Physical,
		Durability: 25,
		Mult:       lowPlunge[c.TalentLvlAttack()],
	}

	// TODO: check snapshot delay
	c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 3),
		lowPlungeHitmark, lowPlungeHitmark)

	return action.Info{
		Frames:          func(next action.Action) int { return lowPlungeFrames[next] },
		AnimationLength: lowPlungeFrames[action.InvalidAction],
		CanQueueAfter:   lowPlungeHitmark,
		State:           action.PlungeAttackState,
	}, nil
}

func (c *char) plungeCollision(fullDelay int) {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Plunge Collision",
		AttackTag:  attacks.AttackTagPlunge,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Physical,
		Durability: 0,
		Mult:       plunge[c.TalentLvlAttack()],
	}
	c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 1.5), fullDelay, fullDelay)
}

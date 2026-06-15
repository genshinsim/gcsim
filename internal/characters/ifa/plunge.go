package ifa

import (
	"errors"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var (
	lowPlungeFrames   []int
	lowPlungeFramesNS []int
)

const (
	lowPlungeHitmark = 38
	collisionHitmark = lowPlungeHitmark - 6
)

func init() {
	lowPlungeFrames = frames.InitAbilSlice(71)
	lowPlungeFrames[action.ActionAttack] = 50
	lowPlungeFrames[action.ActionCharge] = 50
	lowPlungeFrames[action.ActionSkill] = lowPlungeHitmark // assuming it's the same as burst
	lowPlungeFrames[action.ActionBurst] = 50
	lowPlungeFrames[action.ActionDash] = 61 - 19
	lowPlungeFrames[action.ActionWalk] = 69
	lowPlungeFrames[action.ActionSwap] = 57

	lowPlungeFramesNS = frames.InitAbilSlice(69)
	lowPlungeFramesNS[action.ActionAttack] = 50
	lowPlungeFramesNS[action.ActionCharge] = 50
	lowPlungeFramesNS[action.ActionSkill] = lowPlungeHitmark // assuming it's the same as burst
	lowPlungeFramesNS[action.ActionBurst] = 50
	lowPlungeFramesNS[action.ActionDash] = 62 - 19
	lowPlungeFramesNS[action.ActionWalk] = 69
	lowPlungeFramesNS[action.ActionSwap] = 55
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

	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Low Plunge Attack",
		AttackTag:  attacks.AttackTagPlunge,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Anemo,
		Durability: 25,
		Mult:       plunge_low[c.TalentLvlAttack()],
	}

	// TODO: check snapshot delay
	c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 3.5),
		lowPlungeHitmark, lowPlungeHitmark)

	c.Core.Tasks.Add(func() {
		c.exitNightsoul()
	}, lowPlungeHitmark)

	if c.nightsoulState.HasBlessing() {
		ai.AdditionalTags = []attacks.AdditionalTag{attacks.AdditionalTagNightsoul}
	}

	return action.Info{
		Frames: func(next action.Action) int {
			if c.nightsoulState.HasBlessing() {
				return (lowPlungeFramesNS[next])
			}
			return (lowPlungeFrames[next])
		},
		AnimationLength: lowPlungeFrames[action.InvalidAction],
		CanQueueAfter:   lowPlungeHitmark,
		State:           action.PlungeAttackState,
	}, nil
}

func (c *char) plungeCollision(fullDelay int) {
	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Plunge Collision",
		AttackTag:  attacks.AttackTagPlunge,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Anemo,
		Durability: 0,
		Mult:       plunge_collision[c.TalentLvlAttack()],
	}

	if c.nightsoulState.HasBlessing() {
		ai.AdditionalTags = []attacks.AdditionalTag{attacks.AdditionalTagNightsoul}
	}

	c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 1.5), fullDelay, fullDelay)
}

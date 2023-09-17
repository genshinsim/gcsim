package xiao

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

var highPlungeFrames []int
var lowPlungeFrames []int

const collisionHitmark = 38
const highPlungeHitmark = 46
const lowPlungeHitmark = 44

func init() {
	// high_plunge -> x
	highPlungeFrames = frames.InitAbilSlice(66)
	highPlungeFrames[action.ActionAttack] = 61
	highPlungeFrames[action.ActionJump] = 65
	highPlungeFrames[action.ActionSwap] = 64

	// low_plunge -> x
	lowPlungeFrames = frames.InitAbilSlice(62)
	lowPlungeFrames[action.ActionAttack] = 60
	lowPlungeFrames[action.ActionSkill] = 59
	lowPlungeFrames[action.ActionDash] = 60
	lowPlungeFrames[action.ActionJump] = 61
}

// High Plunge attack damage queue generator
// Use the "collision" optional argument if you want to do a falling hit on the way down
// Default = 0
func (c *char) HighPlungeAttack(p map[string]int) action.Info {
	if c.Core.Player.CurrentState() != action.JumpState {
		c.Core.Log.NewEvent("only plunge after using jump", glog.LogActionEvent, c.Index).
			Write("action", action.ActionHighPlunge)
		return action.Info{
			Frames:          func(action.Action) int { return 1200 },
			AnimationLength: 1200,
			CanQueueAfter:   1200,
			State:           action.Idle,
		}
	}

	collision, ok := p["collision"]
	if !ok {
		collision = 0 // Whether or not Xiao does a collision hit
	}

	if collision > 0 {
		c.plungeCollision(collisionHitmark)
	}

	highPlungeRadius := 5.0
	if c.StatusIsActive(burstBuffKey) {
		highPlungeRadius = 6
	}

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "High Plunge",
		AttackTag:  attacks.AttackTagPlunge,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeBlunt,
		Element:    attributes.Physical,
		Durability: 25,
		Mult:       highplunge[c.TalentLvlAttack()],
	}
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, highPlungeRadius),
		highPlungeHitmark,
		highPlungeHitmark,
		c.c6cb(),
	)

	return action.Info{
		Frames:          frames.NewAbilFunc(highPlungeFrames),
		AnimationLength: highPlungeFrames[action.InvalidAction],
		CanQueueAfter:   highPlungeFrames[action.ActionAttack],
		State:           action.PlungeAttackState,
	}
}

// Low Plunge attack damage queue generator
// Use the "collision" optional argument if you want to do a falling hit on the way down
// Default = 0
func (c *char) LowPlungeAttack(p map[string]int) action.Info {
	if c.Core.Player.CurrentState() != action.JumpState {
		c.Core.Log.NewEvent("only plunge after using jump", glog.LogActionEvent, c.Index).
			Write("action", action.ActionLowPlunge)
		return action.Info{
			Frames:          func(action.Action) int { return 1200 },
			AnimationLength: 1200,
			CanQueueAfter:   1200,
			State:           action.Idle,
		}
	}

	collision, ok := p["collision"]
	if !ok {
		collision = 0 // Whether or not Xiao does a collision hit
	}

	if collision > 0 {
		c.plungeCollision(collisionHitmark)
	}

	lowPlungeRadius := 3.0
	if c.StatusIsActive(burstBuffKey) {
		lowPlungeRadius = 4
	}

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Low Plunge",
		AttackTag:  attacks.AttackTagPlunge,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeBlunt,
		Element:    attributes.Physical,
		Durability: 25,
		Mult:       lowplunge[c.TalentLvlAttack()],
	}
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, lowPlungeRadius),
		lowPlungeHitmark,
		lowPlungeHitmark,
		c.c6cb(),
	)

	return action.Info{
		Frames:          frames.NewAbilFunc(lowPlungeFrames),
		AnimationLength: lowPlungeFrames[action.InvalidAction],
		CanQueueAfter:   lowPlungeFrames[action.ActionSkill],
		State:           action.PlungeAttackState,
	}
}

// Plunge normal falling attack damage queue generator
// Standard - Always part of high/low plunge attacks
func (c *char) plungeCollision(delay int) {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Plunge Collision",
		AttackTag:  attacks.AttackTagPlunge,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeSlash,
		Element:    attributes.Physical,
		Durability: 0,
		Mult:       plunge[c.TalentLvlAttack()],
	}
	c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 1), delay, delay)
}

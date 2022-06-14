package xiao

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
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
func (c *char) HighPlungeAttack(p map[string]int) action.ActionInfo {
	collision, ok := p["collision"]
	if !ok {
		collision = 0 // Whether or not Xiao does a collision hit
	}

	if collision > 0 {
		c.plungeCollision(collisionHitmark)
	}

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "High Plunge",
		AttackTag:  combat.AttackTagPlunge,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeBlunt,
		Element:    attributes.Physical,
		Durability: 25,
		Mult:       highplunge[c.TalentLvlAttack()],
	}
	c.Core.QueueAttack(ai, combat.NewDefCircHit(2, false, combat.TargettableEnemy), highPlungeHitmark, highPlungeHitmark)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(highPlungeFrames),
		AnimationLength: highPlungeFrames[action.InvalidAction],
		CanQueueAfter:   highPlungeHitmark,
		Post:            highPlungeHitmark,
		State:           action.PlungeAttackState,
	}
}

// Low Plunge attack damage queue generator
// Use the "collision" optional argument if you want to do a falling hit on the way down
// Default = 0
func (c *char) LowPlungeAttack(p map[string]int) action.ActionInfo {
	collision, ok := p["collision"]
	if !ok {
		collision = 0 // Whether or not Xiao does a collision hit
	}

	if collision > 0 {
		c.plungeCollision(collisionHitmark)
	}

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Low Plunge",
		AttackTag:  combat.AttackTagPlunge,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeBlunt,
		Element:    attributes.Physical,
		Durability: 25,
		Mult:       lowplunge[c.TalentLvlAttack()],
	}
	c.Core.QueueAttack(ai, combat.NewDefCircHit(2, false, combat.TargettableEnemy), lowPlungeHitmark, lowPlungeHitmark)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(lowPlungeFrames),
		AnimationLength: lowPlungeFrames[action.InvalidAction],
		CanQueueAfter:   lowPlungeHitmark,
		Post:            lowPlungeHitmark,
		State:           action.PlungeAttackState,
	}
}

// Plunge normal falling attack damage queue generator
// Standard - Always part of high/low plunge attacks
func (c *char) plungeCollision(delay int) {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Plunge Collision",
		AttackTag:  combat.AttackTagPlunge,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		Element:    attributes.Physical,
		Durability: 0,
		Mult:       plunge[c.TalentLvlAttack()],
	}
	c.Core.QueueAttack(ai, combat.NewDefCircHit(0.1, false, combat.TargettableEnemy), delay, delay)
}

package xiao

import (
	"errors"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player"
)

var lowPlungeFramesX []int
var highPlungeFramesX []int

var lowPlungeFramesXY []int
var highPlungeFramesXY []int

var lowPlungeFramesXYX []int
var highPlungeFramesXYX []int

const lowPlungeHitmarkX = 44
const highPlungeHitmarkX = 46
const collisionHitmarkX = lowPlungeHitmarkX - 6

const lowPlungeHitmarkXY = 42 + 3
const highPlungeHitmarkXY = 43 + 3
const collisionHitmarkXY = lowPlungeHitmarkXY - 6

const lowPlungeHitmarkXYX = 43 + 3
const highPlungeHitmarkXYX = 44 + 3
const collisionHitmarkXYX = lowPlungeHitmarkXYX - 6

func init() {
	// from xiao
	lowPlungeFramesX = frames.InitAbilSlice(62)
	lowPlungeFramesX[action.ActionAttack] = 60
	lowPlungeFramesX[action.ActionSkill] = 59
	lowPlungeFramesX[action.ActionDash] = 60
	lowPlungeFramesX[action.ActionJump] = 61

	highPlungeFramesX = frames.InitAbilSlice(66)
	highPlungeFramesX[action.ActionAttack] = 61
	highPlungeFramesX[action.ActionJump] = 65
	highPlungeFramesX[action.ActionSwap] = 64

	// from xianyun
	lowPlungeFramesXY = frames.InitAbilSlice(73)
	lowPlungeFramesXY[action.ActionAttack] = 59
	lowPlungeFramesXY[action.ActionSkill] = 58
	lowPlungeFramesXY[action.ActionBurst] = 58
	lowPlungeFramesXY[action.ActionDash] = 58
	lowPlungeFramesXY[action.ActionJump] = 58
	lowPlungeFramesXY[action.ActionSwap] = 62

	highPlungeFramesXY = frames.InitAbilSlice(75)
	highPlungeFramesXY[action.ActionAttack] = 60
	highPlungeFramesXY[action.ActionSkill] = 60
	highPlungeFramesXY[action.ActionBurst] = 61
	highPlungeFramesXY[action.ActionDash] = 62
	highPlungeFramesXY[action.ActionJump] = 60
	highPlungeFramesXY[action.ActionSwap] = 63

	// from xiao + xianyun
	lowPlungeFramesXYX = frames.InitAbilSlice(73)
	lowPlungeFramesXYX[action.ActionAttack] = 59
	lowPlungeFramesXYX[action.ActionSkill] = 58
	lowPlungeFramesXYX[action.ActionBurst] = 58 // assumed to be same as skill
	lowPlungeFramesXYX[action.ActionDash] = 59
	lowPlungeFramesXYX[action.ActionJump] = 60
	lowPlungeFramesXYX[action.ActionSwap] = 62

	highPlungeFramesXYX = frames.InitAbilSlice(75)
	highPlungeFramesXYX[action.ActionAttack] = 61
	highPlungeFramesXYX[action.ActionSkill] = 62
	highPlungeFramesXYX[action.ActionBurst] = 62 // assumed to be same as skill
	highPlungeFramesXYX[action.ActionDash] = 62
	highPlungeFramesXYX[action.ActionJump] = 62
	highPlungeFramesXYX[action.ActionSwap] = 64
}

// Low Plunge attack damage queue generator
// Use the "collision" optional argument if you want to do a falling hit on the way down
// Default = 0
func (c *char) LowPlungeAttack(p map[string]int) (action.Info, error) {
	defer c.Core.Player.SetAirborne(player.Grounded)
	if c.Core.Player.CurrentState() != action.JumpState {
		return action.Info{}, errors.New("only plunge after using jump")
	}

	collision, ok := p["collision"]
	if !ok {
		collision = 0 // Whether or not Xiao does a collision hit
	}

	var a action.Info
	var lowPlungeRadius float64
	var lowPlungePoiseDMG float64
	var lowPlungeHitmark int
	var collisionHitmark int
	switch {
	case c.StatusIsActive(player.XianyunAirborneBuff) && c.StatusIsActive(burstBuffKey):
		a = action.Info{
			Frames:          frames.NewAbilFunc(lowPlungeFramesXYX),
			AnimationLength: lowPlungeFramesXYX[action.InvalidAction],
			CanQueueAfter:   lowPlungeFramesXYX[action.ActionSkill],
			State:           action.PlungeAttackState,
		}
		lowPlungePoiseDMG = 150
		lowPlungeRadius = 4
		lowPlungeHitmark = lowPlungeHitmarkXYX
		collisionHitmark = collisionHitmarkXYX
	case c.StatusIsActive(player.XianyunAirborneBuff):
		a = action.Info{
			Frames:          frames.NewAbilFunc(lowPlungeFramesXY),
			AnimationLength: lowPlungeFramesXY[action.InvalidAction],
			CanQueueAfter:   lowPlungeFramesXY[action.ActionSkill],
			State:           action.PlungeAttackState,
		}
		lowPlungePoiseDMG = 100
		lowPlungeRadius = 3
		lowPlungeHitmark = lowPlungeHitmarkXY
		collisionHitmark = collisionHitmarkXY
	default:
		// assumed to be Xiao burst active
		a = action.Info{
			Frames:          frames.NewAbilFunc(lowPlungeFramesX),
			AnimationLength: lowPlungeFramesX[action.InvalidAction],
			CanQueueAfter:   lowPlungeFramesX[action.ActionSkill],
			State:           action.PlungeAttackState,
		}
		lowPlungePoiseDMG = 150
		lowPlungeRadius = 4
		lowPlungeHitmark = highPlungeHitmarkX
		collisionHitmark = collisionHitmarkX
	}

	if collision > 0 {
		c.plungeCollision(collisionHitmark)
	}

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Low Plunge",
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
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, lowPlungeRadius),
		lowPlungeHitmark,
		lowPlungeHitmark,
		c.c6cb(),
	)

	return a, nil
}

// High Plunge attack damage queue generator
// Use the "collision" optional argument if you want to do a falling hit on the way down
// Default = 0
func (c *char) HighPlungeAttack(p map[string]int) (action.Info, error) {
	defer c.Core.Player.SetAirborne(player.Grounded)
	if c.Core.Player.CurrentState() != action.JumpState {
		return action.Info{}, errors.New("only plunge after using jump")
	}

	collision, ok := p["collision"]
	if !ok {
		collision = 0 // Whether or not Xiao does a collision hit
	}

	var a action.Info
	var highPlungeRadius float64
	var highPlungePoiseDMG float64
	var highPlungeHitmark int
	var collisionHitmark int
	switch {
	case c.StatusIsActive(player.XianyunAirborneBuff) && c.StatusIsActive(burstBuffKey):
		a = action.Info{
			Frames:          frames.NewAbilFunc(highPlungeFramesXYX),
			AnimationLength: highPlungeFramesXYX[action.InvalidAction],
			CanQueueAfter:   highPlungeFramesXYX[action.ActionAttack],
			State:           action.PlungeAttackState,
		}
		highPlungePoiseDMG = 225
		highPlungeRadius = 6
		highPlungeHitmark = highPlungeHitmarkXYX
		collisionHitmark = collisionHitmarkXYX
	case c.StatusIsActive(player.XianyunAirborneBuff):
		a = action.Info{
			Frames:          frames.NewAbilFunc(highPlungeFramesXY),
			AnimationLength: highPlungeFramesXY[action.InvalidAction],
			CanQueueAfter:   highPlungeFramesXY[action.ActionAttack],
			State:           action.PlungeAttackState,
		}
		highPlungePoiseDMG = 150
		highPlungeRadius = 5
		highPlungeHitmark = highPlungeHitmarkXY
		collisionHitmark = collisionHitmarkXY
	default:
		// assumed to be Xiao burst active
		a = action.Info{
			Frames:          frames.NewAbilFunc(highPlungeFramesX),
			AnimationLength: highPlungeFramesX[action.InvalidAction],
			CanQueueAfter:   highPlungeFramesX[action.ActionAttack],
			State:           action.PlungeAttackState,
		}
		highPlungePoiseDMG = 225
		highPlungeRadius = 6
		highPlungeHitmark = highPlungeHitmarkX
		collisionHitmark = collisionHitmarkX
	}

	if collision > 0 {
		c.plungeCollision(collisionHitmark)
	}

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "High Plunge",
		AttackTag:  attacks.AttackTagPlunge,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeBlunt,
		PoiseDMG:   highPlungePoiseDMG,
		Element:    attributes.Physical,
		Durability: 25,
		Mult:       highPlunge[c.TalentLvlAttack()],
	}
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, highPlungeRadius),
		highPlungeHitmark,
		highPlungeHitmark,
		c.c6cb(),
	)

	return a, nil
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

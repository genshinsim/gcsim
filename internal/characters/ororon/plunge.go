package ororon

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player"
)

const highPlungeHitmark = 52
const collisionHitmark = 51

// const highPlungePoiseDMG = 100.0 // Not needed since dmg type is pierce?
// const collisionPoiseDMG = 10.0

const highPlungeRadius = 3.5
const collisionRadius = 1.0

var plungeFrames []int

func init() {
	// Plunge -> X
	plungeFrames = frames.InitAbilSlice(66) // Default is From plunge animation start to swap icon un-gray
	plungeFrames[action.ActionAttack] = 68
	plungeFrames[action.ActionAim] = 66
	plungeFrames[action.ActionSkill] = 65
	plungeFrames[action.ActionBurst] = 65
	plungeFrames[action.ActionDash] = 53
	plungeFrames[action.ActionJump] = 82
	plungeFrames[action.ActionWalk] = 80
	plungeFrames[action.ActionSwap] = 66
}

func (c *char) fall() action.Info {
	// Fall cancel can't happen until after high_plunge can happen. Delay all side effects if try to fall cancel too early.
	delay := fallCancelFrames - (c.Core.F - c.jmpSrc)

	// Cleanup high jump.
	if delay <= 0 {
		delay = 0
		c.exitJumpBlessing()
	} else {
		c.Core.Log.NewEvent(
			fmt.Sprintf("Fall cancel cannot begin until %d frames after jump start; delaying fall by %d frames", fallCancelFrames, delay),
			glog.LogCooldownEvent,
			c.Index)

		c.QueueCharTask(c.exitJumpBlessing, delay)
	}
	// Allow stam to start regen when landing
	c.Core.Player.LastStamUse = c.Core.F + jumpHoldFrames[1][action.ActionSwap] + delay

	return action.Info{
		Frames: func(next action.Action) int {
			return frames.NewAbilFunc(jumpHoldFrames[1])(next) + delay
		},
		// Is this supposed to be whatever the max over Frames is?
		AnimationLength: jumpHoldFrames[1][action.ActionWalk] + delay,
		CanQueueAfter:   jumpHoldFrames[1][action.ActionSwap] + delay,
		State:           action.JumpState,
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
		StrikeType: attacks.StrikeTypePierce,
		Element:    attributes.Physical,
		Durability: 0,
		Mult:       collision[c.TalentLvlAttack()],
	}
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(
			c.Core.Combat.Player(),
			geometry.Point{},
			collisionRadius),
		delay,
		delay)
}

// High Plunge attack damage queue generator
// Use the "collision" optional argument if you want to do a falling hit on the way down
// Default = 0
func (c *char) HighPlungeAirborneOroron(p map[string]int) (action.Info, error) {
	// Cleanup high jump.
	c.exitJumpBlessing()
	// Allow player to resume stam as soon as plunge is initiated
	c.Core.Player.LastStamUse = c.Core.F

	collision := p["collision"]

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "High Plunge",
		AttackTag:  attacks.AttackTagPlunge,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypePierce,
		Element:    attributes.Physical,
		Durability: 25,
		Mult:       highPlunge[c.TalentLvlAttack()],
		UseDef:     true,
	}

	if collision > 0 {
		c.plungeCollision(collisionHitmark)
	}

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(
			c.Core.Combat.Player(),
			geometry.Point{},
			highPlungeRadius),
		highPlungeHitmark,
		highPlungeHitmark,
	)

	return action.Info{
		Frames:          frames.NewAbilFunc(plungeFrames),
		AnimationLength: plungeFrames[action.ActionJump],
		CanQueueAfter:   plungeFrames[action.ActionDash],
		State:           action.PlungeAttackState,
	}, nil
}

func (c *char) HighPlungeAttack(p map[string]int) (action.Info, error) {
	defer func() {
		c.Core.Player.SetAirborne(player.Grounded)
		c.jmpSrc = 0
	}()

	if c.Core.Player.Airborne() != player.AirborneOroron {
		return c.Character.HighPlungeAttack(p)
	}

	if p["fall"] != 0 {
		return c.fall(), nil
	}

	return c.HighPlungeAirborneOroron(p)
}

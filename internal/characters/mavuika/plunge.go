package mavuika

import (
	"errors"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/player"
)

const collisionHitmark = 5

var highPlungeFrames []int
var lowPlungeFrames []int
var bikePlungeFrames []int

func init() {

}

// Special bike Plunge will be here
func (c *char) LowPlungeAttack(p map[string]int) (action.Info, error) {
	if c.nightsoulState.HasBlessing() && c.flamestriderModeActive {
		return c.bikePlunge(p), nil
	}

	defer c.Core.Player.SetAirborne(player.Grounded)
	switch c.Core.Player.Airborne() {
	case player.AirborneXianyun:
		return c.lowPlungeXY(p), nil
	default:
		return action.Info{}, errors.New("low_plunge can only be used while airborne")
	}
}

func (c *char) lowPlungeXY(p map[string]int) action.Info {
	collision, ok := p["collision"]
	if !ok {
		collision = 0 // Whether or not collision hit
	}

	if collision > 0 {
		c.plungeCollision(collisionHitmark)
	}

	return action.Info{
		Frames:          frames.NewAbilFunc(lowPlungeFrames),
		AnimationLength: lowPlungeFrames[action.InvalidAction],
		CanQueueAfter:   lowPlungeFrames[action.ActionDash],
		State:           action.PlungeAttackState,
	}
}

// High Plunge attack damage queue generator
// Use the "collision" optional argument if you want to do a falling hit on the way down
// Default = 0
func (c *char) HighPlungeAttack(p map[string]int) (action.Info, error) {
	if c.nightsoulState.HasBlessing() && c.flamestriderModeActive {
		return c.bikePlunge(p), nil
	}

	defer c.Core.Player.SetAirborne(player.Grounded)
	switch c.Core.Player.Airborne() {
	case player.AirborneXianyun:
		return c.highPlungeXY(p), nil
	default:
		return action.Info{}, errors.New("high_plunge can only be used while airborne")
	}
}

func (c *char) highPlungeXY(p map[string]int) action.Info {
	collision, ok := p["collision"]
	if !ok {
		collision = 0 // Whether or not collision hit
	}

	if collision > 0 {
		c.plungeCollision(collisionHitmark)
	}

	return action.Info{
		Frames:          frames.NewAbilFunc(highPlungeFrames),
		AnimationLength: highPlungeFrames[action.InvalidAction],
		CanQueueAfter:   highPlungeFrames[action.ActionDash],
		State:           action.PlungeAttackState,
	}
}

func (c *char) plungeCollision(delay int) {
}

func (c *char) bikePlunge(_ map[string]int) action.Info {
	return action.Info{
		Frames:          frames.NewAbilFunc(bikePlungeFrames),
		AnimationLength: bikePlungeFrames[action.InvalidAction],
		CanQueueAfter:   bikePlungeFrames[action.ActionDash],
		State:           action.PlungeAttackState,
	}
}

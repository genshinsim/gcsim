package xianyun

import (
	"errors"
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player"
)

var leapFrames []int
var plungeHitmarks = []int{20, 30, 40}
var plungeRadius = []float64{4, 5, 6.5}

var highPlungeFrames []int
var lowPlungeFramesXY []int

const collisionHitmark = 38
const highPlungeHitmark = 46
const lowPlungeHitmark = 44

// TODO: missing plunge -> skill
func init() {
	// skill (press) -> high plunge -> x
	leapFrames = frames.InitAbilSlice(55) // max
	leapFrames[action.ActionDash] = 43
	leapFrames[action.ActionJump] = 50
	leapFrames[action.ActionSwap] = 50

	// high_plunge -> x
	highPlungeFrames = frames.InitAbilSlice(66)
	highPlungeFrames[action.ActionAttack] = 61
	highPlungeFrames[action.ActionJump] = 65
	highPlungeFrames[action.ActionSwap] = 64

	// low_plunge -> x
	lowPlungeFramesXY = frames.InitAbilSlice(62)
	lowPlungeFramesXY[action.ActionAttack] = 60
	lowPlungeFramesXY[action.ActionSkill] = 59
	lowPlungeFramesXY[action.ActionDash] = 60
	lowPlungeFramesXY[action.ActionJump] = 61
}

func (c *char) HighPlungeAttack(p map[string]int) (action.Info, error) {
	// last action must be skill (for leap)
	// dont need to check airborne for this because she can plunge if she's on the ground anyways
	if c.StatusIsActive(skillStateKey) {
		return c.driftcloudWave(p)
	}

	switch c.Core.Player.Airborne() {
	case player.AirborneVenti:
		return action.Info{}, fmt.Errorf("xiangyun plunge while airborne due to venti is unimplemented due to lack of frame data. Please see https://docs.gcsim.app/mechanics/frames for how to contribute")
	case player.AirborneXianyun:
		return c.highPlungeXY(p)
	default:
		return action.Info{}, fmt.Errorf("xiangyun high_plunge cannot be used")
	}
}

func (c *char) driftcloudWave(_ map[string]int) (action.Info, error) {
	skillArea := combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, plungeRadius[c.skillCounter-1])
	skillHitmark := plungeHitmarks[c.skillCounter-1]
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       fmt.Sprintf("Chasing Crane %v", c.skillCounter),
		AttackTag:  attacks.AttackTagPlunge,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Anemo,
		Durability: 25,
		Mult:       leap[c.skillCounter-1][c.TalentLvlSkill()],
	}

	c.Core.QueueAttack(
		ai,
		skillArea,
		skillHitmark,
		skillHitmark,
		c.particleCB,
		c.a1cb(),
	)
	// reset window after leap
	c.DeleteStatus(skillStateKey)
	c.skillCounter = 0
	c.skillEnemiesHit = nil
	c.skillSrc = noSrcVal

	return action.Info{
		Frames:          frames.NewAbilFunc(leapFrames),
		State:           action.PlungeAttackState,
		AnimationLength: leapFrames[action.InvalidAction],
		CanQueueAfter:   leapFrames[action.ActionSkill],
	}, nil
}

// Low Plunge attack damage queue generator
// Use the "collision" optional argument if you want to do a falling hit on the way down
// Default = 0
func (c *char) LowPlungeAttack(p map[string]int) (action.Info, error) {
	// last action must be skill (for leap)
	// dont need to check airborne for this because she can plunge if she's on the ground anyways
	if c.StatusIsActive(skillStateKey) {
		return c.driftcloudWave(p)
	}

	switch c.Core.Player.Airborne() {
	case player.AirborneVenti:
		return action.Info{}, fmt.Errorf("xiangyun plunge while airborne due to venti hold E is unimplemented due to lack of frame data. Please see https://docs.gcsim.app/mechanics/frames for how to contribute")
	case player.AirborneXianyun:
		return c.lowPlungeXY(p)
	default:
		return action.Info{}, fmt.Errorf("xiangyun low_plunge cannot be used")
	}
}

func (c *char) lowPlungeXY(p map[string]int) (action.Info, error) {
	if c.Core.Player.CurrentState() != action.JumpState {
		return action.Info{}, errors.New("only plunge after using jump")
	}

	collision, ok := p["collision"]
	if !ok {
		collision = 0 // Whether or not Xianyun does a collision hit
	}

	if collision > 0 {
		c.plungeCollision(collisionHitmark)
	}

	poiseDMG := 100.0
	lowPlungeRadius := 3.0

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Low Plunge",
		AttackTag:  attacks.AttackTagPlunge,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		PoiseDMG:   poiseDMG,
		Element:    attributes.Anemo,
		Durability: 25,
		Mult:       lowPlunge[c.TalentLvlAttack()],
	}
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, lowPlungeRadius),
		lowPlungeHitmark,
		lowPlungeHitmark,
	)

	return action.Info{
		Frames:          frames.NewAbilFunc(lowPlungeFramesXY),
		AnimationLength: lowPlungeFramesXY[action.InvalidAction],
		CanQueueAfter:   lowPlungeFramesXY[action.ActionSkill],
		State:           action.PlungeAttackState,
	}, nil
}

func (c *char) highPlungeXY(p map[string]int) (action.Info, error) {
	if c.Core.Player.CurrentState() != action.JumpState {
		return action.Info{}, errors.New("only plunge after using jump")
	}

	collision, ok := p["collision"]
	if !ok {
		collision = 0 // Whether or not Xianyun does a collision hit
	}

	if collision > 0 {
		c.plungeCollision(collisionHitmark)
	}

	poiseDMG := 150.0
	highPlungeRadius := 5.0

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "High Plunge",
		AttackTag:  attacks.AttackTagPlunge,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		PoiseDMG:   poiseDMG,
		Element:    attributes.Anemo,
		Durability: 25,
		Mult:       highPlunge[c.TalentLvlAttack()],
	}
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, highPlungeRadius),
		highPlungeHitmark,
		highPlungeHitmark,
	)

	return action.Info{
		Frames:          frames.NewAbilFunc(highPlungeFrames),
		AnimationLength: highPlungeFrames[action.InvalidAction],
		CanQueueAfter:   highPlungeFrames[action.ActionAttack],
		State:           action.PlungeAttackState,
	}, nil
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
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Anemo,
		Durability: 0,
		Mult:       collision[c.TalentLvlAttack()],
	}
	c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 1.5), delay, delay)
}

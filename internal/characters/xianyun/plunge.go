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

var driftcloudFrames [][]int
var plungeHitmarks = []int{35, 40, 46}
var plungeRadius = []float64{4, 5, 6.5}

var highPlungeFramesXY []int
var lowPlungeFramesXY []int

const collisionHitmark = 38
const highPlungeHitmark = 46
const lowPlungeHitmark = 44

// TODO: missing plunge -> skill
func init() {
	driftcloudFrames = make([][]int, 3)
	// skill (press) -> high plunge -> x
	driftcloudFrames[0] = frames.InitAbilSlice(65) // max
	driftcloudFrames[0][action.ActionAttack] = 57
	driftcloudFrames[0][action.ActionCharge] = 56
	driftcloudFrames[0][action.ActionSkill] = 56
	driftcloudFrames[0][action.ActionBurst] = 54
	driftcloudFrames[0][action.ActionDash] = 51
	driftcloudFrames[0][action.ActionJump] = 56
	driftcloudFrames[0][action.ActionSwap] = 49

	driftcloudFrames[1] = frames.InitAbilSlice(70) // max
	driftcloudFrames[1][action.ActionAttack] = 60
	driftcloudFrames[1][action.ActionCharge] = 61
	driftcloudFrames[1][action.ActionSkill] = 55
	driftcloudFrames[1][action.ActionBurst] = 61
	driftcloudFrames[1][action.ActionDash] = 55
	driftcloudFrames[1][action.ActionJump] = 62
	driftcloudFrames[1][action.ActionSwap] = 53

	driftcloudFrames[2] = frames.InitAbilSlice(76) // max
	driftcloudFrames[2][action.ActionAttack] = 66
	driftcloudFrames[2][action.ActionCharge] = 67
	driftcloudFrames[2][action.ActionSkill] = 64
	driftcloudFrames[2][action.ActionBurst] = 67
	driftcloudFrames[2][action.ActionDash] = 63
	driftcloudFrames[2][action.ActionJump] = 68
	driftcloudFrames[2][action.ActionSwap] = 62

	// high_plunge -> x
	highPlungeFramesXY = frames.InitAbilSlice(66)
	highPlungeFramesXY[action.ActionAttack] = 61
	highPlungeFramesXY[action.ActionJump] = 65
	highPlungeFramesXY[action.ActionSwap] = 64

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
		return action.Info{}, fmt.Errorf("%s high_plunge while airborne due to venti is unimplemented due to lack of frame data. Please see https://docs.gcsim.app/mechanics/frames for how to contribute", c.Base.Key.String())	case player.AirborneXianyun:
	case player.AirborneXianyun:
		return c.highPlungeXY(p)
	default:
        return action.Info{}, fmt.Errorf("%s high_plunge can only be used while airborne", c.Base.Key.String())
	}
}

func (c *char) driftcloudWave(_ map[string]int) (action.Info, error) {
	skillInd := c.skillCounter - 1
	skillArea := combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, plungeRadius[skillInd])
	skillHitmark := plungeHitmarks[skillInd]
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       fmt.Sprintf("Chasing Crane %v", c.skillCounter),
		AttackTag:  attacks.AttackTagPlunge,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Anemo,
		Durability: 25,
		Mult:       leap[skillInd][c.TalentLvlSkill()],
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
		Frames:          frames.NewAbilFunc(driftcloudFrames[skillInd]),
		State:           action.PlungeAttackState,
		AnimationLength: driftcloudFrames[skillInd][action.InvalidAction],
		CanQueueAfter:   driftcloudFrames[skillInd][action.ActionSwap],
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
		return action.Info{}, fmt.Errorf("%s low_plunge while airborne due to venti hold E is unimplemented due to lack of frame data. Please see https://docs.gcsim.app/mechanics/frames for how to contribute", c.Base.Key.String())
	case player.AirborneXianyun:
		return c.lowPlungeXY(p)
	default:
		return action.Info{}, fmt.Errorf("%s low_plunge can only be used while airborne", c.Base.Key.String())
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
		Frames:          frames.NewAbilFunc(highPlungeFramesXY),
		AnimationLength: highPlungeFramesXY[action.InvalidAction],
		CanQueueAfter:   highPlungeFramesXY[action.ActionAttack],
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

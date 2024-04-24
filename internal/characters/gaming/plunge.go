package gaming

import (
	"errors"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/core/targets"
)

var highPlungeFrames []int
var lowPlungeFrames []int
var specialPlungeFrames []int

const lowPlungeHitmark = 43
const highPlungeHitmark = 46
const collisionHitmark = lowPlungeHitmark - 6
const specialPlungeHitmark = 32

const lowPlungePoiseDMG = 150.0
const lowPlungeRadius = 3.0

const highPlungePoiseDMG = 200.0
const highPlungeRadius = 5.0

const hpDrainThreshold = 0.1
const specialPlungeKey = "Charmed Cloudstrider"
const particleICD = 3 * 60
const particleICDKey = "gaming-particle-icd"

func init() {
	// low_plunge -> x
	lowPlungeFrames = frames.InitAbilSlice(84)
	lowPlungeFrames[action.ActionAttack] = 56
	lowPlungeFrames[action.ActionSkill] = 56
	lowPlungeFrames[action.ActionBurst] = 56
	lowPlungeFrames[action.ActionDash] = lowPlungeHitmark
	lowPlungeFrames[action.ActionSwap] = 67

	// high_plunge -> x
	highPlungeFrames = frames.InitAbilSlice(87)
	highPlungeFrames[action.ActionAttack] = 58
	highPlungeFrames[action.ActionSkill] = 57
	highPlungeFrames[action.ActionBurst] = 57
	highPlungeFrames[action.ActionDash] = highPlungeHitmark
	highPlungeFrames[action.ActionWalk] = 86
	highPlungeFrames[action.ActionSwap] = 69

	// special plunge
	specialPlungeFrames = frames.InitAbilSlice(99)
	specialPlungeFrames[action.ActionAttack] = 52
	specialPlungeFrames[action.ActionSkill] = 52
	specialPlungeFrames[action.ActionBurst] = 52
	specialPlungeFrames[action.ActionDash] = specialPlungeHitmark // was 30
	specialPlungeFrames[action.ActionWalk] = 74
	specialPlungeFrames[action.ActionSwap] = 69
}

// Low Plunge attack damage queue generator
// Use the "collision" optional argument if you want to do a falling hit on the way down
// Default = 0
func (c *char) LowPlungeAttack(p map[string]int) (action.Info, error) {
	if c.Core.Player.LastAction.Type == action.ActionSkill {
		return c.specialPlunge(p)
	}

	defer c.Core.Player.SetAirborne(player.Grounded)
	switch c.Core.Player.Airborne() {
	case player.AirborneXianyun:
		return c.lowPlungeXY(p)
	default:
		return action.Info{}, errors.New("low_plunge can only be used while airborne")
	}
}

func (c *char) lowPlungeXY(p map[string]int) (action.Info, error) {
	collision, ok := p["collision"]
	if !ok {
		collision = 0 // Whether or not collision hit
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
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), geometry.Point{Y: 1}, lowPlungeRadius),
		lowPlungeHitmark,
		lowPlungeHitmark,
	)

	return action.Info{
		Frames:          frames.NewAbilFunc(lowPlungeFrames),
		AnimationLength: lowPlungeFrames[action.InvalidAction],
		CanQueueAfter:   lowPlungeFrames[action.ActionDash],
		State:           action.PlungeAttackState,
	}, nil
}

// High Plunge attack damage queue generator
// Use the "collision" optional argument if you want to do a falling hit on the way down
// Default = 0
func (c *char) HighPlungeAttack(p map[string]int) (action.Info, error) {
	if c.Core.Player.LastAction.Type == action.ActionSkill {
		return c.specialPlunge(p)
	}

	defer c.Core.Player.SetAirborne(player.Grounded)
	switch c.Core.Player.Airborne() {
	case player.AirborneXianyun:
		return c.highPlungeXY(p)
	default:
		return action.Info{}, errors.New("high_plunge can only be used while airborne")
	}
}

func (c *char) highPlungeXY(p map[string]int) (action.Info, error) {
	collision, ok := p["collision"]
	if !ok {
		collision = 0 // Whether or not collision hit
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
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), geometry.Point{Y: 1}, highPlungeRadius),
		highPlungeHitmark,
		highPlungeHitmark,
	)

	return action.Info{
		Frames:          frames.NewAbilFunc(highPlungeFrames),
		AnimationLength: highPlungeFrames[action.InvalidAction],
		CanQueueAfter:   highPlungeFrames[action.ActionDash],
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
		StrikeType: attacks.StrikeTypeSlash,
		Element:    attributes.Physical,
		Durability: 0,
		Mult:       collision[c.TalentLvlAttack()],
	}
	c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 1), delay, delay)
}

func (c *char) specialPlunge(p map[string]int) (action.Info, error) {
	if p[manChaiParam] > 0 {
		c.manChaiWalkBack = p[manChaiParam]
	} else {
		c.manChaiWalkBack = 92
	}

	ai := combat.AttackInfo{
		ActorIndex:     c.Index,
		Abil:           specialPlungeKey,
		AttackTag:      attacks.AttackTagPlunge,
		ICDTag:         attacks.ICDTagNone,
		ICDGroup:       attacks.ICDGroupDefault,
		StrikeType:     attacks.StrikeTypeBlunt,
		PoiseDMG:       lowPlungePoiseDMG,
		Element:        attributes.Pyro,
		Durability:     25,
		Mult:           skill[c.TalentLvlSkill()],
		IgnoreInfusion: true,
	}

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, c.specialPlungeRadius),
		specialPlungeHitmark,
		specialPlungeHitmark,
		c.particleCB,
		c.makeA1CB(),
		c.makeC4CB(),
	)

	// queue man chai and drain hp 1f after hitmark
	c.Core.Tasks.Add(func() {
		if c.StatusIsActive(burstKey) && c.CurrentHPRatio() > 0.5 {
			c.queueManChai()
		}
		// only drain HP when above 10% HP
		if c.CurrentHPRatio() > hpDrainThreshold {
			currentHP := c.CurrentHP()
			maxHP := c.MaxHP()
			hpdrain := 0.15 * currentHP
			// The HP consumption from using this skill can only bring him to 10% HP.
			if (currentHP-hpdrain)/maxHP <= hpDrainThreshold {
				hpdrain = currentHP - hpDrainThreshold*maxHP
			}
			c.Core.Player.Drain(info.DrainInfo{
				ActorIndex: c.Index,
				Abil:       specialPlungeKey,
				Amount:     hpdrain,
			})
		}
	}, specialPlungeHitmark+1)

	return action.Info{
		Frames:          frames.NewAbilFunc(specialPlungeFrames),
		AnimationLength: specialPlungeFrames[action.InvalidAction],
		CanQueueAfter:   specialPlungeFrames[action.ActionDash],
		State:           action.PlungeAttackState,
	}, nil
}

func (c *char) particleCB(a combat.AttackCB) {
	if a.Target.Type() != targets.TargettableEnemy {
		return
	}
	if c.StatusIsActive(particleICDKey) {
		return
	}
	c.AddStatus(particleICDKey, particleICD, true)

	c.Core.QueueParticle(c.Base.Key.String(), 2, attributes.Pyro, c.ParticleDelay)
}

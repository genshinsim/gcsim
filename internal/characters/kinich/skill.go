package kinich

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/targets"
)

const (
	skillCD                  = 18 * 60
	skillStart               = 9
	scalespikerHitmark       = 48
	pointsConsumptionsDelay  = 36
	generateNSPointDelay     = 30
	nightSoulEnterDelay      = 11
	scalespikerHoldFrameDiff = 18

	scalespikerAbil = "Scalespiker Cannon"
)

var (
	skillFrames       []int
	scalespikerFrames []int
)

var blindSpotAppearanceDelays = []int{30, 40}

func init() {
	skillFrames = frames.InitAbilSlice(42) // E -> D/J
	skillFrames[action.ActionAttack] = 30
	skillFrames[action.ActionBurst] = 27
	skillFrames[action.ActionWalk] = 41

	scalespikerFrames = frames.InitAbilSlice(100) // E -> Swap
	scalespikerFrames[action.ActionAttack] = 59
	scalespikerFrames[action.ActionBurst] = 59
	scalespikerFrames[action.ActionDash] = 67
	scalespikerFrames[action.ActionJump] = 67
	scalespikerFrames[action.ActionWalk] = 71
	scalespikerFrames[action.ActionSwap] = 100
}

func (c *char) Skill(p map[string]int) (action.Info, error) {
	if c.nightsoulState.HasBlessing() {
		if c.nightsoulState.Points() == c.nightsoulState.MaxPoints {
			return c.ScalespikerCannon(p)
		}
		return action.Info{}, fmt.Errorf("%v: Cannot use Scalespiker Cannon with %v Nightsoul points, should be %v",
			c.Base.Key, c.nightsoulState.Points(), c.nightsoulState.MaxPoints)
	}

	hold, ok := p["hold"]
	if !ok {
		hold = 0
	}
	if hold < 0 {
		hold = 0
	} else if hold > 301 {
		hold = 301
	}

	c.QueueCharTask(func() {
		src := c.Core.F
		c.nightsoulSrc = src
		c.nightsoulState.EnterBlessing(0.)
		c.setNightsoulExitTimer(10 * 60)
		c.c2AoeIncreased = false
		c.particlesGenerated = false
		c.SetCD(action.ActionSkill, skillCD)
		c.QueueCharTask(c.timePassGenerateNSPoints(src), generateNSPointDelay)
		c.QueueCharTask(c.createBlindSpot, blindSpotAppearanceDelays[0])
	}, skillStart+hold)

	return action.Info{
		Frames: func(next action.Action) int {
			return skillFrames[next] + hold
		},
		AnimationLength: skillFrames[action.InvalidAction] + hold,
		CanQueueAfter:   skillFrames[action.ActionDash] + hold,
		State:           action.SkillState,
	}, nil
}

func (c *char) ScalespikerCannon(p map[string]int) (action.Info, error) {
	hold, ok := p["hold"]
	if !ok {
		hold = 0
	}
	if hold < 0 {
		hold = 0
	} else if hold > 301 {
		hold = 301
	}
	c6Travel, ok := p["c6_travel"]
	if !ok {
		c6Travel = 0 // TODO: find exact frame
	}

	diff := 0
	if hold > 0 {
		diff = hold - (scalespikerHoldFrameDiff + 1)
	}

	ai := combat.AttackInfo{
		ActorIndex:     c.Index,
		Abil:           scalespikerAbil,
		AttackTag:      attacks.AttackTagElementalArt,
		AdditionalTags: []attacks.AdditionalTag{attacks.AdditionalTagNightsoul},
		ICDTag:         attacks.ICDTagKinichScalespikerCannon,
		ICDGroup:       attacks.ICDGroupKinichScalespikerCannon,
		StrikeType:     attacks.StrikeTypeDefault,
		Element:        attributes.Dendro,
		Durability:     25,
		Mult:           scalespikerCannon[c.TalentLvlSkill()],
		FlatDmg:        c.a4Amount(),
	}
	s, radius := c.c2Bonus(&ai)
	target := c.Core.Combat.PrimaryTarget()
	ap := combat.NewCircleHitOnTarget(target, nil, radius)

	c.QueueCharTask(func() {
		c.Core.QueueAttackWithSnap(ai, s, ap, 0, c.particleCB, c.a1CB, c.c2ResShredCB)
		c.c4()
		c.c6(ai, &s, radius, target, c6Travel)
	}, scalespikerHitmark+diff)

	c.QueueCharTask(func() {
		c.nightsoulState.ConsumePoints(c.nightsoulState.MaxPoints)
	}, pointsConsumptionsDelay+min(hold, 1)*(hold-scalespikerHoldFrameDiff-1))
	c.QueueCharTask(c.createBlindSpot, blindSpotAppearanceDelays[1])

	return action.Info{
		Frames: func(next action.Action) int {
			return scalespikerFrames[next] + diff
		},
		AnimationLength: scalespikerFrames[action.InvalidAction] + diff,
		CanQueueAfter:   scalespikerFrames[action.ActionAttack] + diff,
		State:           action.SkillState,
	}, nil
}

func (c *char) loopShotGenerateNSPoints() {
	if !c.nightsoulState.HasBlessing() {
		return
	}
	c.nightsoulState.GeneratePoints(3.)
}

func (c *char) timePassGenerateNSPoints(src int) func() {
	return func() {
		if c.nightsoulSrc != src {
			return
		}
		c.nightsoulState.GeneratePoints(1.)
		c.QueueCharTask(c.timePassGenerateNSPoints(src), generateNSPointDelay)
	}
}

func (c *char) createBlindSpot() {
	newBlindSpotAngularPosition := c.characterAngularPosition + float64(c.Core.Rand.Intn(2)*2-1)*90.
	newBlindSpotAngularPosition = NormalizeAngle360(newBlindSpotAngularPosition)
	c.blindSpotAngularPosition = newBlindSpotAngularPosition
}

func (c *char) cancelNightsoul() {
	c.nightsoulState.ExitBlessing()
	c.DeleteStatus(desolationKey)
	c.nightsoulSrc = -1
	c.blindSpotAngularPosition = -1
	c.exitStateF = -1
}

func (c *char) setNightsoulExitTimer(duration int) {
	src := c.Core.F + duration
	c.exitStateF = src
	c.QueueCharTask(func() {
		if c.exitStateF != src {
			return
		}
		c.cancelNightsoul()
	}, duration)
}

func (c *char) particleCB(a combat.AttackCB) {
	if a.Target.Type() != targets.TargettableEnemy {
		return
	}
	if c.particlesGenerated {
		return
	}
	c.particlesGenerated = true

	c.Core.QueueParticle(c.Base.Key.String(), 5, attributes.Dendro, c.ParticleDelay)
}

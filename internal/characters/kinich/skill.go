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
	timePassNSGenDelay       = 38
	nightSoulEnterDelay      = 11
	scalespikerHoldFrameDiff = 18

	scalespikerAbil   = "Scalespiker Cannon"
	scalespikerC6Abil = "Scalespiker Cannon (C6)"
)

var (
	skillFrames       []int
	scalespikerFrames []int
)

var blindSpotAppearanceDelays = []int{30, 40}

func init() {
	skillFrames = frames.InitAbilSlice(42) // E -> D/J
	skillFrames[action.ActionAttack] = 30
	skillFrames[action.ActionWalk] = 41
	skillFrames[action.ActionBurst] = 27

	scalespikerFrames = frames.InitAbilSlice(68) // E -> D/J
	scalespikerFrames[action.ActionAttack] = 60
	scalespikerFrames[action.ActionWalk] = 71
	scalespikerFrames[action.ActionBurst] = 59
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
	c.c2AoeIncreased = false
	c.nightsoulSrc = c.Core.F
	c.particlesGenerated = false
	c.QueueCharTask(func() { c.nightsoulState.EnterBlessing(0.) }, nightSoulEnterDelay)
	c.SetCDWithDelay(action.ActionSkill, skillCD-10*60, skillStart)
	c.QueueCharTask(c.timePassGenerateNSPoints, timePassNSGenDelay)
	c.QueueCharTask(c.createBlindSpot, blindSpotAppearanceDelays[0]-1) // just in case, since attack can be executed at the same frame(?)
	c.setNightsoulExitTimer(10 * 60)

	return action.Info{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionDash],
		State:           action.SkillState,
	}, nil
}

func (c *char) ScalespikerCannon(p map[string]int) (action.Info, error) {
	hold, ok := p["hold"]
	if !ok {
		hold = 0
	} else {
		if hold < 0 {
			hold = 0
		} else if hold > 301 {
			hold = 301
		}
	}
	ai := c.getScalespikerAi()
	radius := 3
	dmgBonus := 0
	if c.Base.Cons >= 2 {
		if !c.c2AoeIncreased {
			c.c2AoeIncreased = true
			radius = 5
		}
		dmgBonus = 1
	}
	ap := combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, float64(radius))
	c.QueueCharTask(func() {
		a4FlatDmg := c.a4Amount()
		ai.FlatDmg += a4FlatDmg
		s := c.C2Snapshot(ai, dmgBonus)
		c.Core.QueueAttackWithSnap(ai, s, ap, 0, c.particleCB, c.desolationCB, c.c2ResShredCB)
		if c.Base.Cons >= 6 {
			c6Travel, ok := p["c6_travel"]
			if !ok {
				c6Travel = 10
			}
			c.QueueCharTask(func() {
				c6Ai := c.getScalespikerAi()
				c6Ai.FlatDmg += a4FlatDmg
				c6Ai.Abil = scalespikerC6Abil
				s := c.C2Snapshot(c6Ai, dmgBonus)
				enemies := c.Core.Combat.Enemies()
				var target combat.Target
				for _, enemy := range enemies {
					if enemy.Key() != c.Core.Combat.PrimaryTarget().Key() {
						target = enemy
					}
				}
				if target == nil {
					return
				}
				apC6 := combat.NewCircleHitOnTarget(target, nil, float64(radius))
				c.Core.QueueAttackWithSnap(c6Ai, s, apC6, 0, c.particleCB, c.desolationCB, c.c2ResShredCB)
			}, scalespikerHitmark+min(hold, 1)*(hold-scalespikerHoldFrameDiff-1)+c6Travel)
		}
	}, scalespikerHitmark+min(hold, 1)*(hold-scalespikerHoldFrameDiff-1))
	c.QueueCharTask(func() {
		c.nightsoulState.ConsumePoints(c.nightsoulState.MaxPoints)
	}, pointsConsumptionsDelay+min(hold, 1)*(hold-scalespikerHoldFrameDiff-1))
	c.QueueCharTask(c.createBlindSpot, blindSpotAppearanceDelays[1]-1) // just in case, since attack can be executed at the same frame(?)
	return action.Info{
		Frames: func(next action.Action) int {
			return scalespikerFrames[next] + min(hold, 1)*(hold-scalespikerHoldFrameDiff-1)
		},
		AnimationLength: scalespikerFrames[action.InvalidAction],
		CanQueueAfter:   scalespikerFrames[action.ActionAttack],
		State:           action.SkillState,
	}, nil
}

func (c *char) getScalespikerAi() combat.AttackInfo {
	return combat.AttackInfo{
		ActorIndex:     c.Index,
		Abil:           scalespikerAbil,
		AttackTag:      attacks.AttackTagElementalArt,
		ICDTag:         attacks.ICDTagKinichScalespikerCannon,
		ICDGroup:       attacks.ICDGroupKinichScalespikerCannon,
		StrikeType:     attacks.StrikeTypeDefault,
		Element:        attributes.Dendro,
		Durability:     25,
		Mult:           scalespikerCannon[c.TalentLvlSkill()],
		HitlagFactor:   0.01,
		IgnoreInfusion: true,
	}
}

func (c *char) loopShotGenerateNSPoints() {
	if !c.nightsoulState.HasBlessing() {
		return
	}
	c.nightsoulState.GeneratePoints(3.)
}

func (c *char) timePassGenerateNSPoints() {
	if !c.nightsoulState.HasBlessing() {
		return
	}
	c.nightsoulState.GeneratePoints(1.)
	c.QueueCharTask(c.timePassGenerateNSPoints, 30)
}

func (c *char) createBlindSpot() {
	newBlindSpotAngularPosition := c.characterAngularPosition + float64(c.Core.Rand.Intn(2)*2-1)*90.
	newBlindSpotAngularPosition = NormalizeAngle360(newBlindSpotAngularPosition)
	c.blindSpotAngularPosition = newBlindSpotAngularPosition
}

func (c *char) cancelNightsoul() {
	c.nightsoulState.ClearPoints()
	c.nightsoulState.ExitBlessing()
	c.DeleteStatus(desolationKey)
	c.nightsoulSrc = -1
	c.blindSpotAngularPosition = -1
}

func (c *char) setNightsoulExitTimer(duration int) {
	src := c.Core.F
	c.QueueCharTask(func() {
		if c.nightsoulSrc != src {
			return
		}
		if c.skillDurationExtended {
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

func (c *char) desolationCB(a combat.AttackCB) {
	if c.Base.Ascension < 1 {
		return
	}
	c.AddStatus(desolationKey, -1, false)
}

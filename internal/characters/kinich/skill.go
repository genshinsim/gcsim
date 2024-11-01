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
	skillStart               = 1
	scalespikerHitmark       = 10
	pointsConsumptionsDelay  = 2
	timePassNSGenDelay       = 1
	blindSpotAppearanceDelay = 1
)

var (
	skillFrames       []int
	scalespikerFrames []int
)

func init() {
	skillFrames = frames.InitAbilSlice(33) // E -> Q
	skillFrames[action.ActionAttack] = 31
	skillFrames[action.ActionSkill] = 32
	skillFrames[action.ActionDash] = skillStart // ability doesn't start if dash is done before CD
	skillFrames[action.ActionJump] = 25
	skillFrames[action.ActionSwap] = 25
	skillFrames[action.ActionWalk] = 32

	scalespikerFrames = frames.InitAbilSlice(43) // E -> Walk
	scalespikerFrames[action.ActionAttack] = 24
	scalespikerFrames[action.ActionSkill] = 24
	scalespikerFrames[action.ActionBurst] = 24
	scalespikerFrames[action.ActionDash] = 25
	scalespikerFrames[action.ActionJump] = 25
	scalespikerFrames[action.ActionSwap] = 42
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
	c.nightsoulState.EnterBlessing(0.)
	c.SetCDWithDelay(action.ActionSkill, skillCD, skillStart)
	c.QueueCharTask(c.timePassGenerateNSPoints, timePassNSGenDelay)
	c.QueueCharTask(c.createBlindSpot, blindSpotAppearanceDelay)
	c.setNightsoulExitTimer(10 * 60)

	return action.Info{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionDash],
		State:           action.SkillState,
	}, nil
}

func (c *char) ScalespikerCannon(p map[string]int) (action.Info, error) {
	travel, ok := p["travel"]
	if !ok {
		travel = 10
	}
	c.scaleskiperAttackInfo = c.getScalespikerAi()
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
		c.scaleskiperAttackInfo.FlatDmg += a4FlatDmg
		s := c.C2Snapshot(dmgBonus)
		c.Core.QueueAttackWithSnap(c.scaleskiperAttackInfo, s, ap, 0, c.particleCB, c.desolationCB, c.c2ResShredCB)
		if c.Base.Cons >= 6 {
			c6Travel, ok := p["c6_travel"]
			if !ok {
				c6Travel = 10
			}
			c.QueueCharTask(func() {
				c6Ai := c.getScalespikerAi()
				c6Ai.FlatDmg += a4FlatDmg
				s := c.C2Snapshot(dmgBonus)
				enemies := c.Core.Combat.Enemies()
				var target combat.Target
				for _, enemy := range enemies {
					if enemy.Key() != c.Core.Combat.PrimaryTarget().Key() {
						target = enemy
					}
				}
				apC6 := combat.NewCircleHitOnTarget(target, nil, float64(radius))
				c.Core.QueueAttackWithSnap(c6Ai, s, apC6, 0, c.particleCB, c.desolationCB, c.c2ResShredCB)
			}, scalespikerHitmark+c6Travel)
		}
	}, scalespikerHitmark+travel)
	c.QueueCharTask(func() { c.nightsoulState.ConsumePoints(c.nightsoulState.MaxPoints) }, pointsConsumptionsDelay)
	c.QueueCharTask(c.createBlindSpot, blindSpotAppearanceDelay)
	return action.Info{
		Frames:          frames.NewAbilFunc(scalespikerFrames),
		AnimationLength: scalespikerFrames[action.InvalidAction],
		CanQueueAfter:   scalespikerFrames[action.ActionAttack],
		State:           action.SkillState,
	}, nil
}

func (c *char) getScalespikerAi() combat.AttackInfo {
	return combat.AttackInfo{
		ActorIndex:     c.Index,
		Abil:           "Scalespiker Cannon",
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
	c.nightsoulState.GeneratePoints(3.)
}

func (c *char) timePassGenerateNSPoints() {
	if !c.nightsoulState.HasBlessing() {
		return
	}
	c.nightsoulState.GeneratePoints(2.)
	c.QueueCharTask(c.timePassGenerateNSPoints, 60)
}

func (c *char) createBlindSpot() {
	newBlindSpotAngularPosition := c.characterAngularPosition + float64(c.Core.Rand.Intn(2)*2-1)*90.
	fmt.Println("Generated new blind spot", c.characterAngularPosition, newBlindSpotAngularPosition)
	newBlindSpotAngularPosition = NormalizeAngle(newBlindSpotAngularPosition)
	c.blindSpotAngularPosition = newBlindSpotAngularPosition
}

func (c *char) cancelNightsoul() {
	c.nightsoulState.ClearPoints()
	c.nightsoulState.ExitBlessing()
	c.DeleteStatus(desolationKey)
	c.nightsoulSrc = -1
}

func (c *char) setNightsoulExitTimer(duration int) {
	src := c.Core.F
	c.QueueCharTask(func() {
		if c.nightsoulSrc != src {
			return
		}
		if c.durationExtended {
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

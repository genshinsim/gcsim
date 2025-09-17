package kinich

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/enemy"
)

const (
	skillCD                  = 18 * 60
	skillStart               = 9
	scalespikerDefaultTravel = 12
	pointsConsumptionsDelay  = 1
	nightSoulEnterDelay      = 11
	scalespikerHoldFrameDiff = 18
	generateNSPointDelay     = 30
)

var (
	skillFrames       []int
	scalespikerFrames []int
)

var (
	blindSpotAppearanceDelays = []int{5, 30}  // tap, hold (both tap and hold for entering nightsoul)
	scalespikerReleases       = []int{35, 17} // tap, hold
)

func init() {
	skillFrames = frames.InitAbilSlice(42) // E -> D/J
	skillFrames[action.ActionAttack] = 29
	skillFrames[action.ActionBurst] = 27
	skillFrames[action.ActionWalk] = 41

	scalespikerFrames = frames.InitAbilSlice(100) // E -> Swap
	scalespikerFrames[action.ActionAttack] = 59 - scalespikerReleases[0]
	scalespikerFrames[action.ActionBurst] = 59 - scalespikerReleases[0]
	scalespikerFrames[action.ActionDash] = 67 - scalespikerReleases[0]
	scalespikerFrames[action.ActionJump] = 67 - scalespikerReleases[0]
	scalespikerFrames[action.ActionWalk] = 71 - scalespikerReleases[0]
	scalespikerFrames[action.ActionSwap] = 100 - scalespikerReleases[0]
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
	switch {
	case hold < 0:
		hold = 0
	case hold > 301:
		hold = 301
	}
	if hold > 0 {
		hold--
	}

	c.Core.Tasks.Add(func() {
		src := c.Core.F
		c.nightsoulSrc = src
		c.nightsoulState.EnterTimedBlessing(0., 10*60+10, c.cancelNightsoul)
		c.c2AoeIncreased = false
		c.particlesGenerated = false
		c.SetCD(action.ActionSkill, skillCD)
		c.Core.Tasks.Add(c.timePassGenerateNSPoints(src), generateNSPointDelay)
		c.Core.Tasks.Add(c.createBlindSpot, blindSpotAppearanceDelays[1]-skillStart)
	}, skillStart+hold)

	return action.Info{
		Frames: func(next action.Action) int {
			return skillFrames[next] + hold
		},
		AnimationLength: skillFrames[action.InvalidAction] + hold,
		CanQueueAfter:   skillFrames[action.ActionBurst] + hold,
		State:           action.SkillState,
	}, nil
}

func (c *char) ScalespikerCannon(p map[string]int) (action.Info, error) {
	hold, ok := p["hold"]
	if !ok {
		hold = 0
	}
	switch {
	case hold < 0:
		hold = 0
	case hold > 181:
		hold = 181
	}

	travel, ok := p["travel"]
	if !ok {
		travel = scalespikerDefaultTravel
	}
	c6Travel, ok := p["c6_travel"]
	if !ok {
		c6Travel = 50 // TODO: find exact frame
	}

	releaseFrame := scalespikerReleases[0]
	blindSpotDelay := blindSpotAppearanceDelays[0]
	if hold > 0 {
		hold--
		releaseFrame = scalespikerReleases[1]
		blindSpotDelay = blindSpotAppearanceDelays[1]
	}

	ai := info.AttackInfo{
		ActorIndex:     c.Index(),
		Abil:           "Scalespiker Cannon",
		AttackTag:      attacks.AttackTagElementalArt,
		AdditionalTags: []attacks.AdditionalTag{attacks.AdditionalTagNightsoul, attacks.AdditionalTagKinichCannon},
		ICDTag:         attacks.ICDTagKinichScalespikerCannon,
		ICDGroup:       attacks.ICDGroupKinichScalespikerCannon,
		StrikeType:     attacks.StrikeTypeDefault,
		Element:        attributes.Dendro,
		Durability:     25,
		Mult:           scalespikerCannon[c.TalentLvlSkill()],
		FlatDmg:        c.a4Amount(),
	}

	target := c.Core.Combat.PrimaryTarget()
	radius := 3.0

	var snap info.Snapshot
	c.QueueCharTask(func() {
		// Nightsoul points are drained before snapshot
		c.nightsoulState.ClearPoints()
		snap = c.Snapshot(&ai)
		c.Core.Tasks.Add(c.createBlindSpot, blindSpotDelay)
		c.Core.Tasks.Add(func() {
			if c.Base.Cons >= 2 && !c.c2AoeIncreased {
				c.c2AoeIncreased = true
				radius = 5.0
				snap.Stats[attributes.DmgP] += 1.0
				c.Core.Log.NewEvent("C2 bonus dmg% applied", glog.LogCharacterEvent, c.Index()).
					Write("final", snap.Stats[attributes.DmgP])
			}
			ap := combat.NewCircleHitOnTarget(target, nil, radius)
			c.Core.QueueAttackWithSnap(ai, snap, ap, 0, c.particleCB, c.a1CB, c.c2ResShredCB)
			c.c4()
			c.c6(ai, &snap, radius, target, c6Travel)
		}, travel)
	}, releaseFrame+hold+pointsConsumptionsDelay)

	return action.Info{
		Frames: func(next action.Action) int {
			return releaseFrame + hold + scalespikerFrames[next]
		},
		AnimationLength: releaseFrame + hold + scalespikerFrames[action.InvalidAction],
		CanQueueAfter:   releaseFrame + hold + scalespikerFrames[action.ActionAttack],
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
		c.Core.Tasks.Add(c.timePassGenerateNSPoints(src), generateNSPointDelay)
	}
}

func (c *char) createBlindSpot() {
	newBlindSpotAngularPosition := c.characterAngularPosition + float64(c.Core.Rand.Intn(2)*2-1)*90.
	newBlindSpotAngularPosition = NormalizeAngle360(newBlindSpotAngularPosition)
	c.blindSpotAngularPosition = newBlindSpotAngularPosition
}

func (c *char) cancelNightsoul() {
	c.nightsoulState.ClearPoints()
	c.nightsoulState.ExitBlessing()
	c.nightsoulSrc = -1
	c.blindSpotAngularPosition = -1

	// Clear desolation status from all enemies
	for _, t := range c.Core.Combat.Enemies() {
		if e, ok := t.(*enemy.Enemy); ok {
			e.DeleteStatus(desolationKey)
		}
	}
}

func (c *char) particleCB(a info.AttackCB) {
	if a.Target.Type() != info.TargettableEnemy {
		return
	}
	if c.particlesGenerated {
		return
	}
	c.particlesGenerated = true

	c.Core.QueueParticle(c.Base.Key.String(), 5, attributes.Dendro, c.ParticleDelay)
}

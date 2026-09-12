package jahoda

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var (
	skillFrames       []int
	skillCancelFrames []int
)

const (
	skillWindup              = 30
	shadowPursuitMaxDuration = 334

	firstFlaskFillDelay = 7
	flaskFillInterval   = 30
	flaskFillValue      = 20
	flaskFillWeakValue  = flaskFillValue / 2
	flaskGaugeMax       = 100

	drainFlaskHitmark = 4
	unfillHitmark     = 4
	fillHitmark       = 2

	firstMeowballFirstHitmark = 129
	meowballHitmarkInterval   = 116

	skillCD = 15 * 60

	c1BounceHitmark = 32

	shadowPursuitKey         = "jahoda-shadow-pursuit"
	meowballKey              = "jahoda-meowball"
	meowballFlatEnergyKey    = "jahoda-meowball-flat-energy"
	meowballFlatEnergyICDKey = "jahoda-meowball-flat-energy-icd"
	particleICDKey           = "jahoda-particle-icd"
)

func init() {
	skillFrames = frames.InitAbilSlice(12 + skillWindup) // E -> E

	skillCancelFrames = frames.InitAbilSlice(48) // E -> E -> J
	skillCancelFrames[action.ActionAttack] = 43  // E -> E -> N1
	skillCancelFrames[action.ActionAim] = 42     // E -> E -> Aim
	skillCancelFrames[action.ActionBurst] = 42   // E -> E -> Q
	skillCancelFrames[action.ActionDash] = 41    // E -> E -> D
	skillCancelFrames[action.ActionWalk] = 44    // E -> E -> W
	skillCancelFrames[action.ActionSwap] = 41    // E -> E -> Swap
}

func (c *char) Skill(p map[string]int) (action.Info, error) {
	if c.StatusIsActive(shadowPursuitKey) {
		c.Core.Tasks.Add(c.drainFlask(c.skillSrc), 0)
		return action.Info{
			Frames:          frames.NewAbilFunc(skillCancelFrames),
			AnimationLength: skillCancelFrames[action.InvalidAction],
			CanQueueAfter:   skillCancelFrames[action.ActionDash], // earliest cancel
			State:           action.SkillState,
		}, nil
	}

	travel, ok := p["travel"]
	if !ok {
		travel = 13
	}
	c.meowballTravel = travel

	c.skillSrc = c.Core.F
	src := c.skillSrc
	c.meowballSrc = -1
	c.DeleteStatus(meowballKey)

	c.flaskAbsorb = attributes.NoElement
	c.flaskGauge = 0
	c.flaskAbsorbCheckLocation = combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 5)

	// enter shadow pursuit
	c.pursuitDuration = shadowPursuitMaxDuration
	c.Core.Tasks.Add(func() {
		c.AddStatus(shadowPursuitKey, shadowPursuitMaxDuration, false)

		c.QueueCharTask(c.fillFlask(c.skillSrc), firstFlaskFillDelay)

		// schedule drain flask after max duration to avoid skill end abruptly
		// without doing damage when there are no elemental aura on enemy
		c.Core.Tasks.Add(
			c.drainFlask(c.skillSrc),
			shadowPursuitMaxDuration,
		)
	}, skillWindup)
	return action.Info{
		Frames: func(next action.Action) int {
			// during pursuit state, allow the frames to increase from current src
			if src == c.skillSrc {
				// except for E -> E, as it uses the cancel frames for the skill
				if next == action.ActionSkill {
					return skillFrames[action.ActionSkill]
				}
				return (c.Core.F - src) + skillCancelFrames[next]
			}

			// after pursuit state ended, pursuit duration freeze and the frames becomes fix
			return c.pursuitDuration + skillCancelFrames[next]
		},
		AnimationLength: skillWindup + shadowPursuitMaxDuration + skillCancelFrames[action.ActionJump],
		CanQueueAfter:   skillCancelFrames[action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}, nil
}

func (c *char) ParticleCB(a info.AttackCB) {
	if a.Target.Type() != info.TargettableEnemy {
		return
	}
	if c.StatusIsActive(particleICDKey) {
		return
	}
	c.AddStatus(particleICDKey, 0.5*60, false)
	c.Core.QueueParticle(c.Base.Key.String(), 4, attributes.Anemo, c.ParticleDelay)
}

func (c *char) meowballEnergyCB(a info.AttackCB) {
	if a.Target.Type() != info.TargettableEnemy {
		return
	}

	if c.StatusIsActive(meowballFlatEnergyICDKey) {
		return
	}

	c.AddEnergy(meowballFlatEnergyKey, 2)
	c.AddStatus(meowballFlatEnergyICDKey, 3.5*60, true)
}

func (c *char) fillFlask(src int) func() {
	return func() {
		if src != c.skillSrc {
			return
		}

		// determine which element should be absorbed
		priority := make([]attributes.Element, 0, len(elePriority))

		if c.flaskAbsorb != attributes.NoElement {
			priority = append(priority, c.flaskAbsorb)
		}

		for _, ele := range elePriority {
			if ele != c.flaskAbsorb {
				priority = append(priority, ele)
			}
		}

		// check the prioritized elemental aura on the enemy
		objectElem := c.enemyAuraInArea(c.flaskAbsorbCheckLocation, priority)

		c.pursuitDuration = c.Core.F - c.skillSrc

		if objectElem != attributes.NoElement {
			amount := flaskFillWeakValue

			switch {
			case c.flaskAbsorb == attributes.NoElement:
				// if there is an element to absorb AND the flask has not absorbed any element,
				// the flask fill up its gauge by full value
				c.flaskAbsorb = objectElem
				amount = flaskFillValue

			case objectElem == c.flaskAbsorb:
				// if there is an element to absorb AND the flask has absorbed an element,
				// the flask fill up its gauge by full value
				amount = flaskFillValue

			default:
				// otherwise the flask fill up its gauge by half of its full value
			}

			c.fillFlaskGauge(src, amount, c.flaskAbsorb)
		}

		c.QueueCharTask(c.fillFlask(src), flaskFillInterval)
	}
}

func (c *char) fillFlaskGauge(src, amount int, ele attributes.Element) {
	if src != c.skillSrc {
		return
	}

	prevFlaskGauge := c.flaskGauge
	c.flaskGauge = min(c.flaskGauge+amount, flaskGaugeMax)
	c.Core.Log.NewEvent("jahoda flask gauge increase", glog.LogCharacterEvent, c.Index()).
		Write("element", ele.String()).
		Write("amount", amount).
		Write("previous flask gauge", prevFlaskGauge).
		Write("current flask gauge", c.flaskGauge)

	// if the flask is full drain the flask
	if c.flaskGauge >= flaskGaugeMax {
		c.Core.Tasks.Add(c.drainFlask(c.skillSrc), drainFlaskHitmark)
		return
	}
}

func (c *char) consumeFlaskGauge(src, amount int) func() {
	return func() {
		if src != c.meowballSrc {
			return
		}

		prevFlaskGauge := c.flaskGauge
		c.flaskGauge -= amount
		c.Core.Log.NewEvent("jahoda flask gauge decrease", glog.LogCharacterEvent, c.Index()).
			Write("previous flask gauge", prevFlaskGauge).
			Write("current flask gauge", c.flaskGauge)
		if c.flaskGauge < 0 {
			c.flaskGauge = 0
		}
	}
}

func (c *char) drainFlask(src int) func() {
	return func() {
		if src != c.skillSrc {
			return
		}

		c.pursuitDuration = c.Core.F - src
		c.cancelPursuit() // exit state

		if c.flaskGauge >= flaskGaugeMax {
			c.flaskGauge = flaskGaugeMax
			c.meowballSrc = c.Core.F

			// if the flask is full, do filled flask damage
			ai := info.AttackInfo{
				ActorIndex: c.Index(),
				Abil:       "Filled Treasure Flask",
				AttackTag:  attacks.AttackTagElementalArt,
				ICDTag:     attacks.ICDTagElementalArt,
				ICDGroup:   attacks.ICDGroupDefault,
				StrikeType: attacks.StrikeTypeDefault,
				Element:    attributes.Anemo,
				Durability: 25,
				Mult:       filledFlask[c.TalentLvlSkill()],
			}

			c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), info.Point{Y: 2.5}, 5), 0, fillHitmark, c.ParticleCB)

			// if in ascendent gleam, do meowball damage
			if c.Core.Player.GetMoonsignLevel() >= 2 {
				c.c6()

				c.Core.Tasks.Add(
					c.meowballTick(c.meowballSrc),
					firstMeowballFirstHitmark,
				)

				ticks := flaskGaugeMax / 10
				c.AddStatus(
					meowballKey,
					firstMeowballFirstHitmark+(ticks-1)*meowballHitmarkInterval,
					false,
				)
			}
		} else {
			// if the flask is not full (early cancel or the duration expired), do unfill damage
			ai := info.AttackInfo{
				ActorIndex: c.Index(),
				Abil:       "Unfilled Treasure Flask",
				AttackTag:  attacks.AttackTagElementalArt,
				ICDTag:     attacks.ICDTagElementalArt,
				ICDGroup:   attacks.ICDGroupDefault,
				StrikeType: attacks.StrikeTypeDefault,
				Element:    attributes.Anemo,
				Durability: 25,
				Mult:       unfilledFlask[c.TalentLvlSkill()],
			}

			c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), info.Point{Y: 2.5}, 5), 0, unfillHitmark, c.ParticleCB)
		}
	}
}

func (c *char) cancelPursuit() {
	if !c.StatusIsActive(shadowPursuitKey) {
		return
	}
	c.SetCD(action.ActionSkill, skillCD+skillWindup)
	c.DeleteStatus(shadowPursuitKey)
	c.skillSrc = -1
}

func (c *char) meowballTick(src int) func() {
	return func() {
		if src != c.meowballSrc {
			return
		}

		if !c.StatusIsActive(meowballKey) {
			return
		}

		if c.flaskGauge < 0 {
			c.flaskGauge = 0
		}

		if c.flaskGauge <= 0 {
			return
		}

		c.Core.Tasks.Add(c.consumeFlaskGauge(src, 10), 0)

		ai := info.AttackInfo{
			ActorIndex: c.Index(),
			Abil:       "Meowball",
			AttackTag:  attacks.AttackTagElementalArt,
			ICDTag:     attacks.ICDTagNone,
			ICDGroup:   attacks.ICDGroupDefault,
			StrikeType: attacks.StrikeTypeDefault,
			Element:    c.flaskAbsorb,
			Durability: 25,
			Mult:       meowball[c.TalentLvlSkill()],
		}

		c.Core.QueueAttack(
			ai,
			combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 4),
			0,
			c.meowballTravel,
			c.meowballEnergyCB,
			c.makeC1CB,
		)

		c.Core.Tasks.Add(
			c.meowballTick(c.meowballSrc),
			meowballHitmarkInterval,
		)
	}
}

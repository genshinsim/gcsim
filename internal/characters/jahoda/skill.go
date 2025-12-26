package jahoda

import (
	"math"

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
	skillWindup               = 30
	shadowPursuitMaxDuration  = 334
	firstFillFlaskDelay       = 19
	fillFlaskInterval         = 29
	drainFlask                = 22
	unfillHitmark             = 4
	fillHitmark               = 2
	firstMeowballFirstHitmark = 129
	meowballHitmarkInterval   = 116
	skillCD                   = 15 * 60
	c1BounceHitmark           = 32
	shadowPursuitKey          = "jahoda-shadow-pursuit"
	meowballKey               = "jahoda-meowball"
	meowballFlatEnergyKey     = "jahoda-meowball-flat-energy"
	meowballFlatEnergyICDKey  = "jahoda-meowball-flat-energy-icd"
	particleICDKey            = "jahoda-particle-icd"
)

func init() {
	skillFrames = frames.InitAbilSlice(12) // E -> E

	skillCancelFrames = frames.InitAbilSlice(43) // E -> E -> N1
	skillCancelFrames[action.ActionAim] = 42     // E -> E -> Aim
	skillCancelFrames[action.ActionBurst] = 42   // E -> E -> Q
	skillCancelFrames[action.ActionDash] = 41    // E -> E -> D
	skillCancelFrames[action.ActionJump] = 48    // E -> E -> J
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

	c.Core.Player.SwapCD = math.MaxInt16
	travel, ok := p["travel"]
	if !ok {
		travel = 13
	}

	c.skillTravel = travel

	c.skillSrc = c.Core.F
	c.flaskAbsorb = attributes.NoElement
	c.flaskGauge = 0
	c.flaskGaugeMax = 100
	c.flaskAbsorbCheckLocation = combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 4)

	// Enter shadow pursuit
	c.pursuitDuration = shadowPursuitMaxDuration
	c.Core.Tasks.Add(func() {
		c.AddStatus(shadowPursuitKey, shadowPursuitMaxDuration, false)
		c.Core.Player.SwapCD = math.MaxInt16
		c.Core.Tasks.Add(
			c.fillFlask(c.skillSrc),
			firstFillFlaskDelay,
		)
	}, skillWindup)

	return action.Info{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionAttack],
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
	c.AddStatus(particleICDKey, 0.5*60, false) // Couldn't find anywhere in dm, assume to be the same as Sayu
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
	c.AddStatus(meowballFlatEnergyICDKey, int(3.5*60), true)
}

func (c *char) fillFlask(src int) func() {
	return func() {
		if src != c.skillSrc {
			return
		}

		c.pursuitDuration = c.Core.F - c.skillSrc

		// If the flask is full OR the max duration of the state is reached, drain the flask
		if c.flaskGauge >= c.flaskGaugeMax || !c.StatusIsActive(shadowPursuitKey) || c.Core.F >= shadowPursuitMaxDuration+c.skillSrc {
			c.Core.Tasks.Add(c.drainFlask(c.skillSrc), 0)
			return
		}

		// Check elemental aura in the area
		objectElem := c.Core.Combat.AbsorbCheck(c.Index(), c.flaskAbsorbCheckLocation, attributes.Pyro, attributes.Hydro, attributes.Electro, attributes.Cryo)
		if objectElem != attributes.NoElement {
			if c.flaskAbsorb == attributes.NoElement {
				// If there are an element to absorb AND the flask has not absorbed any elements, the flask absorb that element and increase its gauge
				c.flaskAbsorb = objectElem
				c.Core.Tasks.Add(c.changeFlaskGauge(c.flaskGaugeMax/4), 0)
			} else if objectElem == c.flaskAbsorb {
				// Else if there are an element to absorb AND if that element is the same as that the flask has absorbed the flask absorb that element
				// and increase its gauge
				c.Core.Tasks.Add(c.changeFlaskGauge(c.flaskGaugeMax/4), 0)
			}
			// Otherwise nothing happend since the flask will not change its elemental absorption mid way

			c.Core.Log.NewEventBuildMsg(glog.LogCharacterEvent, c.Index(),
				"jahoda flask absorbed ", c.flaskAbsorb.String(),
			)
		}
		c.Core.Tasks.Add(c.fillFlask(c.skillSrc), fillFlaskInterval)
	}
}

func (c *char) changeFlaskGauge(amount int) func() {
	return func() {
		prevFlaskGauge := c.flaskGauge
		c.flaskGauge = c.flaskGauge + amount
		c.Core.Log.NewEvent("Flask Gauge Change", glog.LogCharacterEvent, c.Index()).
			Write("previous flask gauge", prevFlaskGauge).
			Write("current flask gauge", c.flaskGauge)
	}
}

func (c *char) drainFlask(src int) func() {
	return func() {
		if src != c.skillSrc {
			return
		}

		c.cancelPursuit() // Exit state

		if c.flaskGauge >= c.flaskGaugeMax {
			c.flaskGauge = c.flaskGaugeMax

			// If the flask is full, do filled flask damage
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

			c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), info.Point{Y: 2.5}, 5), 0, drainFlask+fillHitmark, c.ParticleCB)

			// If in Ascendent Gleam, do meowball damage
			if c.Core.Player.GetMoonsignLevel() >= 2 {
				c.c6() // Apply buff from C6

				ticks := c.flaskGaugeMax / 10

				for i := range ticks {
					c.Core.Tasks.Add(
						c.meowballTick(c.skillSrc),
						firstMeowballFirstHitmark+i*meowballHitmarkInterval,
					)
					c.Core.Tasks.Add(c.changeFlaskGauge(-10), firstMeowballFirstHitmark+i*meowballHitmarkInterval)
				}

				c.AddStatus(
					meowballKey,
					firstMeowballFirstHitmark+(ticks-1)*meowballHitmarkInterval,
					false,
				)
			}

		} else {
			// If the flask is not full (early cancel or no furation expired), do unfill damage
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

			c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), info.Point{Y: 2.5}, 5), 0, unfillHitmark, c.ParticleCB)
		}
	}
}

func (c *char) cancelPursuit() {
	if !c.StatusIsActive(shadowPursuitKey) {
		return
	}
	c.SetCD(action.ActionSkill, skillCD+skillWindup)
	c.DeleteStatus(shadowPursuitKey)
	c.Core.Player.SwapCD = skillCancelFrames[action.ActionSwap]
	c.skillSrc = -1
}

func (c *char) meowballTick(src int) func() {
	return func() {
		if src != c.skillSrc {
			return
		}

		if c.flaskGauge <= 0 {
			return
		}

		if c.flaskGauge < 0 {
			c.flaskGauge = 0
		}

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
			c.skillTravel,
			c.meowballEnergyCB,
		)

		// 50% hit twice
		if c.Base.Cons >= 1 {
			if c.Core.Rand.Float64() < 0.5 {
				aiC1 := info.AttackInfo{
					ActorIndex: c.Index(),
					Abil:       "Meowball (C1)",
					AttackTag:  attacks.AttackTagElementalArt,
					ICDTag:     attacks.ICDTagJahodaCons,
					ICDGroup:   attacks.ICDGroupJahodaCons,
					StrikeType: attacks.StrikeTypeDefault,
					Element:    c.flaskAbsorb,
					Durability: 25,
					Mult:       meowball[c.TalentLvlSkill()],
				}

				c.Core.QueueAttack(
					aiC1,
					combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 4),
					0,
					c.skillTravel+c1BounceHitmark,
					nil,
				)
			}
		}

	}
}

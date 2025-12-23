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
	flaskAbsorbMaxDuration = 7 * 60   // Eyeball to be around 7 seconds. Frames needed
	fillFlaskInterval      = 0.5 * 60 // From wiki
	drainFlaskDuration     = 0.5 * 60 // From wiki
	meowballFirstTickDelay = 10       // Frames needed
	meowballInterval       = 2 * 60
	skillCD                = 15 * 60
	c1BounceDelay          = 10 // Frames needed
	meowballKey            = "jahoda-meowball"
	shadowPursuitKey       = "jahoda-shadow-pursuit"
	flatEnergyKey          = "jahoda-skill-flat-energy"
	particleICDKey         = "jahoda-particle-icd"
	flatEnergyICDKey       = "jahoda-flat-energy-icd"
)

func init() {
	skillFrames = frames.InitAbilSlice(10)                   // E -> E. Frames needed
	skillFrames[action.ActionBurst] = flaskAbsorbMaxDuration // E -> Q. cannot use Q during shadow pursuit state

	skillCancelFrames = frames.InitAbilSlice(10) // E -> E -> N1. Frames needed
	skillCancelFrames[action.ActionAim] = 10     // E -> E -> Aim. Frames needed
	skillCancelFrames[action.ActionBurst] = 10   // E -> E -> Q. Frames needed
	skillCancelFrames[action.ActionDash] = 10    // E -> E -> D. Frames needed
	skillCancelFrames[action.ActionJump] = 10    // E -> E -> J. Frames needed
	skillCancelFrames[action.ActionWalk] = 10    // E -> E -> W. Frames needed
	skillCancelFrames[action.ActionSwap] = 10    // E -> E -> Swap. Frames needed
}

func ceil(x float64) int {
	return int(math.Ceil(x))
}

func (c *char) Skill(p map[string]int) (action.Info, error) {
	if c.StatusIsActive(shadowPursuitKey) {
		c.Core.Tasks.Add(c.drainFlask(c.skillSrc), 0)
		return action.Info{
			Frames:          frames.NewAbilFunc(skillCancelFrames),
			AnimationLength: skillCancelFrames[action.InvalidAction],
			CanQueueAfter:   skillCancelFrames[action.ActionBurst], // earliest cancel, need checking
			State:           action.SkillState,
		}, nil
	}

	travel, ok := p["travel"]
	if !ok {
		travel = 10 // Frames needed
	}

	c.skillSrc = c.Core.F
	c.skillTravel = travel
	c.flaskAbsorb = attributes.NoElement
	c.flaskAbsorbDuration = c.Core.F + flaskAbsorbMaxDuration
	c.flaskGauge = 0
	c.flaskGaugeMax = 100
	c.flaskAbsorbCheckLocation = combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 4)

	// Enter shadow pursuit
	c.AddStatus(shadowPursuitKey, flaskAbsorbMaxDuration, false)
	c.Core.Tasks.Add(c.fillFlask(c.skillSrc), 10+fillFlaskInterval) // First absorbtion delay. Frames needed

	return action.Info{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionSwap], // earliest cancel, need checking
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
	c.AddStatus(particleICDKey, 0.5*60, false) // Couldn't find anywhere in dm, assume top be the same as Sayu
	c.Core.QueueParticle(c.Base.Key.String(), 4, attributes.Anemo, c.ParticleDelay)
}

func (c *char) meowballEnergyCB(a info.AttackCB) {
	if a.Target.Type() != info.TargettableEnemy {
		return
	}

	if c.StatusIsActive(flatEnergyICDKey) {
		return
	}

	c.AddEnergy(flatEnergyKey, 2)
	c.AddStatus(flatEnergyICDKey, int(3.5*60), true)
}

func (c *char) cancelPursuit() {
	if !c.StatusIsActive(shadowPursuitKey) {
		return
	}
	c.DeleteStatus(shadowPursuitKey)
	c.flaskAbsorbDuration = -1
	c.SetCD(action.ActionSkill, skillCD)
}

func (c *char) fillFlask(src int) func() {
	return func() {
		if src != c.skillSrc {
			return
		}

		if c.flaskGauge >= c.flaskGaugeMax || !c.StatusIsActive(shadowPursuitKey) || c.Core.F > c.flaskAbsorbDuration {
			c.Core.Tasks.Add(c.drainFlask(c.skillSrc), 0)
			return
		}

		objectElem := c.Core.Combat.AbsorbCheck(c.Index(), c.flaskAbsorbCheckLocation, attributes.Pyro, attributes.Hydro, attributes.Electro, attributes.Cryo)
		if objectElem != attributes.NoElement {
			if c.flaskAbsorb == attributes.NoElement {
				c.flaskAbsorb = objectElem
				c.flaskGauge = c.flaskGauge + c.flaskGaugeMax/4 // fill a quarter (out of 100) every iteration
			} else if objectElem == c.flaskAbsorb {
				c.flaskGauge = c.flaskGauge + c.flaskGaugeMax/4
			}

			c.Core.Log.NewEventBuildMsg(glog.LogCharacterEvent, c.Index(),
				"jahoda absorbed ", c.flaskAbsorb.String(),
			)
		}
		c.Core.Tasks.Add(c.fillFlask(c.skillSrc), fillFlaskInterval)
	}
}

// Copy from Escoffier skill
func (c *char) meowballTick(src int) func() {
	return func() {
		if src != c.skillSrc {
			return
		}

		if c.flaskGauge <= 0 {
			return
		}

		c.flaskGauge -= 10
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
					c.skillTravel+c1BounceDelay,
					c.meowballEnergyCB, // Does C1 trigger refund?
				)
			}
		}

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

			c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), info.Point{Y: 2.5}, 5), 0, drainFlaskDuration, c.ParticleCB)

			if c.Core.Player.GetMoonsignLevel() >= 2 {
				c.c6() // Apply buff from C6

				ticks := c.flaskGauge / 10

				for i := 0; i < ticks; i++ {
					c.Core.Tasks.Add(
						c.meowballTick(c.skillSrc),
						meowballFirstTickDelay+i*meowballInterval,
					)
				}

				c.AddStatus(
					meowballKey,
					meowballFirstTickDelay+(ticks-1)*meowballInterval,
					false,
				)
			}

		} else {
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

			c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), info.Point{Y: 2.5}, 5), 0, drainFlaskDuration, c.ParticleCB)
		}
	}
}

package mavuika

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
)

const (
	skillHitmark                     = 16
	ringsOfSearchingRadianceInterval = 120
)

var (
	skillFrames       []int
	skillSwitchFrames []int
)

func init() {
	skillFrames = frames.InitAbilSlice(20)       // E -> Swap
	skillSwitchFrames = frames.InitAbilSlice(18) // E -> N1
	// on one footage the E icon is switched 1 frame before the Q animation
	// on others the icon is not switched, still she's on a bike. no idea how that works
	// for now I assume E should be pressed to perform instant switch
	skillSwitchFrames[action.ActionBurst] = 1
}

func (c *char) Skill(p map[string]int) (action.Info, error) {
	if c.nightsoulState.HasBlessing() {
		c.allFireArmamnetsActive = !c.allFireArmamnetsActive
		if c.allFireArmamnetsActive {
			c.c2DeleteDefMod()
		} else {
			c.c2AddDefMod()
		}

		return action.Info{
			Frames:          frames.NewAbilFunc(skillSwitchFrames),
			AnimationLength: skillSwitchFrames[action.InvalidAction],
			CanQueueAfter:   skillSwitchFrames[action.ActionAttack], // change to earliest
			State:           action.SkillState,
		}, nil
	}

	c.nightsoulState.EnterBlessing(c.nightsoulState.MaxPoints)
	c.c2BaseIncrease(true)
	c.nightsoulSrc = c.Core.F
	c.QueueCharTask(c.nightsoulPointReduceFunc(c.Core.F), 12)
	hold, ok := p["hold"]
	if !ok {
		hold = 0
	}
	switch {
	case hold < 0:
		hold = 0
	case hold > 1:
		hold = 1
	}
	ai := combat.AttackInfo{
		ActorIndex:     c.Index,
		Abil:           "Skill DMG",
		AttackTag:      attacks.AttackTagElementalArt,
		AdditionalTags: []attacks.AdditionalTag{attacks.AdditionalTagNightsoul},
		ICDTag:         attacks.ICDTagNone,
		ICDGroup:       attacks.ICDGroupDefault,
		StrikeType:     attacks.StrikeTypeDefault,
		Element:        attributes.Pyro,
		Durability:     25,
		Mult:           1.339,
	}
	c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player().Pos(), geometry.Point{Y: 1}, 3.5), skillHitmark, skillHitmark, c.particleCB)
	c.SetCD(action.ActionSkill, 15*60)

	c.QueueCharTask(c.ringsOfSearchingRadianceHit(c.Core.F), ringsOfSearchingRadianceInterval)
	if hold == 1 {
		c.allFireArmamnetsActive = true
	} else {
		c.allFireArmamnetsActive = false
		c.c2AddDefMod()
	}

	return action.Info{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionBurst], // change to earliest
		State:           action.SkillState,
	}, nil
}

func (c *char) nightsoulPointReduceFunc(src int) func() {
	return func() {
		if c.nightsoulSrc != src {
			return
		}

		if !c.nightsoulState.HasBlessing() {
			return
		}

		if !c.StatusIsActive(crucibleOfDeathAndLifeStatus) {
			c.reduceNightsoulPoints(1)
		}
		c.QueueCharTask(c.nightsoulPointReduceFunc(src), 12)
	}
}

func (c *char) reduceNightsoulPoints(val float64) {
	c.nightsoulState.ConsumePoints(val)
	if c.nightsoulState.Points() <= 0.00001 {
		if !c.allFireArmamnetsActive {
			c.c2DeleteDefMod()
		}
		c.c2BaseIncrease(false)
		c.nightsoulState.ExitBlessing()
	}
}

func (c *char) ringsOfSearchingRadianceHit(src int) func() {
	return func() {
		if src != c.nightsoulSrc {
			return
		}
		if !c.nightsoulState.HasBlessing() {
			return
		}
		if !c.allFireArmamnetsActive {
			ai := combat.AttackInfo{
				ActorIndex:     c.Index,
				Abil:           "Rings of Searing Radiance DMG",
				AttackTag:      attacks.AttackTagElementalArt,
				AdditionalTags: []attacks.AdditionalTag{attacks.AdditionalTagNightsoul},
				ICDTag:         attacks.ICDTagNone,
				ICDGroup:       attacks.ICDGroupDefault,
				StrikeType:     attacks.StrikeTypeDefault,
				Element:        attributes.Pyro,
				Durability:     25,
				Mult:           2.304,
			}
			c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 3.5), 0, 0)
			// a hit of E comsumes 3 NS points
			c.nightsoulState.ConsumePoints(3)
		}
		c.QueueCharTask(c.ringsOfSearchingRadianceHit(src), ringsOfSearchingRadianceInterval)
	}
}

func (c *char) particleCB(a combat.AttackCB) {
	c.Core.QueueParticle(c.Base.Key.String(), 5, attributes.Pyro, c.ParticleDelay)
}

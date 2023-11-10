package furina

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/core/targets"
)

var skillFrames []int

const (
	particleICDKey   = "furina-skill-particle-icd"
	skillKey         = "furina-skill"
	skillMaxDuration = 1800
)

func init() {
	skillFrames = frames.InitAbilSlice(30)
}

func (c *char) Skill(p map[string]int) (action.Info, error) {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Salon Solitaire",
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Hydro,
		Durability: 25,
		FlatDmg:    0.134 * c.MaxHP(),
	}

	c.Core.QueueAttack(ai, combat.NewSingleTargetHit(c.Core.Combat.PrimaryTarget().Key()), 0, 0, func(ac combat.AttackCB) {
		currentFrame := c.Core.F
		c.lastSkillUseFrame = currentFrame

		c.surintendanteChevalmarin(currentFrame)()
		c.gentilhommeUsher(currentFrame)()
		c.mademoiselleCrabaletta(currentFrame)()
	})

	c.SetCDWithDelay(action.ActionSkill, 1200, 0)

	return action.Info{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionSwap],
		State:           action.SkillState,
	}, nil
}

func (c *char) surintendanteChevalmarin(src int) func() {
	return func() {
		if src != c.lastSkillUseFrame || c.Core.F-src > skillMaxDuration {
			return
		}

		alliesWithDrainedHPCounter := c.consumeAlliesHealth(0.016)
		damageMultiplier := 1 + 0.1*float64(alliesWithDrainedHPCounter)

		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Salon Solitaire: Surintendante Chevalmarin",
			AttackTag:  attacks.AttackTagElementalArt,
			ICDTag:     attacks.ICDTagFurinaSurintendanteChevalmarin,
			ICDGroup:   attacks.ICDGroupAlhaithamProjectionAttack,
			StrikeType: attacks.StrikeTypeDefault,
			Element:    attributes.Hydro,
			Durability: 25,
			FlatDmg:    0.0549 * c.MaxHP() * damageMultiplier,
		}

		c.Core.QueueAttack(ai, combat.NewSingleTargetHit(c.Core.Combat.PrimaryTarget().Key()), 0, 0, c.particleCB)

		c.Core.Tasks.Add(c.surintendanteChevalmarin(src), 1.5*60) // 1.5s interval
	}
}

func (c *char) gentilhommeUsher(src int) func() {
	return func() {
		if src != c.lastSkillUseFrame || c.Core.F-src > skillMaxDuration {
			return
		}

		alliesWithDrainedHPCounter := c.consumeAlliesHealth(0.024)
		damageMultiplier := 1 + 0.1*float64(alliesWithDrainedHPCounter)

		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Salon Solitaire: Gentilhomme Usher",
			AttackTag:  attacks.AttackTagElementalArt,
			ICDTag:     attacks.ICDTagFurinaGentilhommeUsher,
			ICDGroup:   attacks.ICDGroupAlhaithamProjectionAttack,
			StrikeType: attacks.StrikeTypeDefault,
			Element:    attributes.Hydro,
			Durability: 25,
			FlatDmg:    0.1013 * c.MaxHP() * damageMultiplier,
		}

		c.Core.QueueAttack(ai, combat.NewSingleTargetHit(c.Core.Combat.PrimaryTarget().Key()), 0, 0, c.particleCB)

		c.Core.Tasks.Add(c.gentilhommeUsher(src), 3.75*60) // 3.75s interval
	}
}

func (c *char) mademoiselleCrabaletta(src int) func() {
	return func() {
		if src != c.lastSkillUseFrame || c.Core.F-src > skillMaxDuration {
			return
		}

		alliesWithDrainedHPCounter := c.consumeAlliesHealth(0.036)
		damageMultiplier := 1 + 0.1*float64(alliesWithDrainedHPCounter)

		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Salon Solitaire: Mademoiselle Crabaletta",
			AttackTag:  attacks.AttackTagElementalArt,
			ICDTag:     attacks.ICDTagNone,
			ICDGroup:   attacks.ICDGroupDefault,
			StrikeType: attacks.StrikeTypeDefault,
			Element:    attributes.Hydro,
			Durability: 25,
			FlatDmg:    0.1409 * c.MaxHP() * damageMultiplier,
		}

		c.Core.QueueAttack(ai, combat.NewSingleTargetHit(c.Core.Combat.PrimaryTarget().Key()), 0, 0, c.particleCB)

		c.Core.Tasks.Add(c.mademoiselleCrabaletta(src), 256) // 4.28s interval
	}
}

func (c *char) particleCB(ac combat.AttackCB) {
	if ac.Target.Type() != targets.TargettableEnemy {
		return
	}

	if c.StatusIsActive(particleICDKey) {
		return
	}

	c.AddStatus(particleICDKey, 2.5*60, false)
	c.Core.QueueParticle(c.Base.Key.String(), 1, attributes.Hydro, c.ParticleDelay)
}

func (c *char) consumeAlliesHealth(hpDrainRatio float64) int {
	var alliesWithDrainedHPCounter = 0

	for _, char := range c.Core.Player.Chars() {
		currentHPRatio := char.CurrentHPRatio()

		if currentHPRatio <= 0.5 {
			continue
		}

		alliesWithDrainedHPCounter++

		hpDrain := char.MaxHP() * hpDrainRatio

		c.Core.Player.Drain(player.DrainInfo{
			ActorIndex: char.Index,
			Abil:       "Salon Solitaire",
			Amount:     hpDrain,
			External:   true,
		})
	}

	return alliesWithDrainedHPCounter
}

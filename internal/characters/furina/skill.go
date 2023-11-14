package furina

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/core/targets"
)

var skillFrames []int

const (
	skillHitmark     = 30
	salonInitialTick = 60
	summonDelay      = 30
	particleICDKey   = "furina-skill-particle-icd"
	skillKey         = "furina-skill"
	skillMaxDuration = 1800

	salonMemberKey = "Salon Member"

	chevalmarinIntervalMean   = 90
	chevalmarinIntervalStddev = 5
	chevalmarinTravelMean     = 10
	chevalmarinTravelStddev   = 3

	usherIntervalMean   = 225
	usherIntervalStddev = 5
	usherTravelMean     = 10
	usherTravelStddev   = 3

	crabalettaIntervalMean   = 256
	crabalettaIntervalStddev = 5
	crabalettaTravelMean     = 5
	crabalettaTravelStddev   = 2

	singerInterval = 120
)

func init() {
	skillFrames = frames.InitAbilSlice(30)
}

func (c *char) Skill(p map[string]int) (action.Info, error) {
	c.SetCDWithDelay(action.ActionSkill, 1200, 0)
	switch c.arkhe {
	case ousia:
		return c.skillOusia(p)
	case pneuma:
		return c.skillPneuma(p)
	default:
		return action.Info{}, fmt.Errorf("%v: character is in unknown arkhe: %v", c.CharWrapper.Base.Key, c.arkhe)
	}
}

func (c *char) skillPneuma(_ map[string]int) (action.Info, error) {
	c.AddStatus(skillKey, 1800+summonDelay, false)
	c.summonSinger(c.Core.F, summonDelay)

	return action.Info{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionSwap],
		State:           action.SkillState,
	}, nil
}
func (c *char) skillOusia(_ map[string]int) (action.Info, error) {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Salon Solitaire: Ousia Bubble",
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Hydro,
		Durability: 25,
		FlatDmg:    skillOusiaBubble[c.TalentLvlSkill()] * c.MaxHP(),
	}

	c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), geometry.Point{}, 5), skillHitmark, skillHitmark)

	c.AddStatus(skillKey, 1800+summonDelay, false)
	c.summonSalonMembers(c.Core.F, summonDelay)

	return action.Info{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionSwap],
		State:           action.SkillState,
	}, nil
}

func (c *char) calcRandNorm(mean, std int) int {
	val := int(c.Core.Rand.NormFloat64()*float64(std) + 0.5)
	if val < -2*std {
		val = -2 * std
	}
	if val > std*2 {
		val = std * 2
	}
	return mean + val
}

func (c *char) summonSalonMembers(src, delay int) {
	// TODO: figure out first action time
	c.lastSummonSrc = src
	c.Core.Tasks.Add(c.surintendanteChevalmarin(src), delay+salonInitialTick)
	c.Core.Tasks.Add(c.gentilhommeUsher(src), delay+salonInitialTick)
	c.Core.Tasks.Add(c.mademoiselleCrabaletta(src), delay+salonInitialTick)
}

func (c *char) summonSinger(src, delay int) {
	c.Core.Tasks.Add(c.singerOfManyWaters(src), delay+salonInitialTick)
}

func (c *char) surintendanteChevalmarin(src int) func() {
	return func() {
		if c.arkhe != ousia {
			return
		}

		if src != c.lastSummonSrc {
			return
		}

		if !c.StatusIsActive(skillKey) {
			return
		}

		alliesWithDrainedHPCounter := c.consumeAlliesHealth(0.016)
		damageMultiplier := 1 + 0.1*float64(alliesWithDrainedHPCounter)

		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       fmt.Sprintf("%v: Surintendante Chevalmarin", salonMemberKey),
			AttackTag:  attacks.AttackTagElementalArt,
			ICDTag:     attacks.ICDTagFurinaChevalmarin,
			ICDGroup:   attacks.ICDGroupFurinaSalonSolitaire,
			StrikeType: attacks.StrikeTypeDefault,
			Element:    attributes.Hydro,
			Durability: 25,
			FlatDmg:    skillChevalmarin[c.TalentLvlSkill()] * c.MaxHP() * damageMultiplier,
		}
		travel := c.calcRandNorm(chevalmarinTravelMean, chevalmarinTravelStddev)
		if c.Base.Cons >= 4 {
			c.Core.QueueAttack(ai, combat.NewSingleTargetHit(c.Core.Combat.PrimaryTarget().Key()), travel, travel, c.particleCB, c.c4cb)
		} else {
			c.Core.QueueAttack(ai, combat.NewSingleTargetHit(c.Core.Combat.PrimaryTarget().Key()), travel, travel, c.particleCB)
		}

		interval := c.calcRandNorm(chevalmarinIntervalMean, chevalmarinIntervalStddev)
		c.Core.Tasks.Add(c.surintendanteChevalmarin(src), interval)
	}
}

func (c *char) gentilhommeUsher(src int) func() {
	return func() {
		if c.arkhe != ousia {
			return
		}

		if src != c.lastSummonSrc {
			return
		}

		if !c.StatusIsActive(skillKey) {
			return
		}

		alliesWithDrainedHPCounter := c.consumeAlliesHealth(0.024)
		damageMultiplier := 1 + 0.1*float64(alliesWithDrainedHPCounter)

		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       fmt.Sprintf("%v: Gentilhomme Usher", salonMemberKey),
			AttackTag:  attacks.AttackTagElementalArt,
			ICDTag:     attacks.ICDTagFurinaUsher,
			ICDGroup:   attacks.ICDGroupFurinaSalonSolitaire,
			StrikeType: attacks.StrikeTypeDefault,
			Element:    attributes.Hydro,
			Durability: 25,
			FlatDmg:    skillUsher[c.TalentLvlSkill()] * c.MaxHP() * damageMultiplier,
		}

		travel := c.calcRandNorm(usherTravelMean, usherTravelStddev)
		if c.Base.Cons >= 4 {
			c.Core.QueueAttack(ai, combat.NewSingleTargetHit(c.Core.Combat.PrimaryTarget().Key()), travel, travel, c.particleCB, c.c4cb)
		} else {
			c.Core.QueueAttack(ai, combat.NewSingleTargetHit(c.Core.Combat.PrimaryTarget().Key()), travel, travel, c.particleCB)
		}

		interval := c.calcRandNorm(usherIntervalMean, usherIntervalStddev)
		c.Core.Tasks.Add(c.gentilhommeUsher(src), interval)
	}
}

func (c *char) mademoiselleCrabaletta(src int) func() {
	return func() {
		if c.arkhe != ousia {
			return
		}

		if src != c.lastSummonSrc {
			return
		}

		if !c.StatusIsActive(skillKey) {
			return
		}

		alliesWithDrainedHPCounter := c.consumeAlliesHealth(0.036)
		damageMultiplier := 1 + 0.1*float64(alliesWithDrainedHPCounter)

		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       fmt.Sprintf("%v: Mademoiselle Crabaletta", salonMemberKey),
			AttackTag:  attacks.AttackTagElementalArt,
			ICDTag:     attacks.ICDTagNone,
			ICDGroup:   attacks.ICDGroupDefault,
			StrikeType: attacks.StrikeTypeDefault,
			Element:    attributes.Hydro,
			Durability: 25,
			FlatDmg:    skillCrabaletta[c.TalentLvlSkill()] * c.MaxHP() * damageMultiplier,
		}

		travel := c.calcRandNorm(crabalettaTravelMean, crabalettaTravelStddev)

		if c.Base.Cons >= 4 {
			c.Core.QueueAttack(ai, combat.NewSingleTargetHit(c.Core.Combat.PrimaryTarget().Key()), travel, travel, c.particleCB, c.c4cb)
		} else {
			c.Core.QueueAttack(ai, combat.NewSingleTargetHit(c.Core.Combat.PrimaryTarget().Key()), travel, travel, c.particleCB)
		}

		interval := c.calcRandNorm(crabalettaIntervalMean, crabalettaIntervalStddev)
		c.Core.Tasks.Add(c.mademoiselleCrabaletta(src), interval)
	}
}

func (c *char) singerOfManyWaters(src int) func() {
	return func() {
		if c.arkhe != pneuma {
			return
		}

		if src != c.lastSummonSrc {
			return
		}

		if !c.StatusIsActive(skillKey) {
			return
		}
		// heal
		c.Core.Player.Heal(player.HealInfo{
			Caller:  c.Index,
			Target:  c.Core.Player.Active(),
			Message: "Singer of Many Waters",
			Src:     skillSingerHealFlat[c.TalentLvlSkill()] + skillSingerHealScale[c.TalentLvlSkill()]*c.MaxHP(),
			Bonus:   c.Stat(attributes.Heal),
		})
		intervalDelta := c.MaxHP() / 1000.0 * 0.004
		if intervalDelta > 0.16 {
			intervalDelta = 0.16
		}
		interval := int(singerInterval*(1-intervalDelta) + 0.5)
		c.Core.Tasks.Add(c.singerOfManyWaters(src), interval)
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

	for i, char := range c.Core.Player.Chars() {
		currentHPRatio := char.CurrentHPRatio()

		if currentHPRatio <= 0.5 {
			continue
		}

		alliesWithDrainedHPCounter++

		if c.Core.Player.Active() == i && (c.Core.Player.CurrentState() == action.BurstState || c.Core.Player.CurrentState() == action.DashState) {
			// her skill does not drain the HP of active characters that are in iframes (burst or dash)
			continue
		}
		hpDrain := char.MaxHP() * hpDrainRatio

		c.Core.Player.Drain(player.DrainInfo{
			ActorIndex: char.Index,
			Abil:       "Salon Solitaire",
			Amount:     hpDrain,
			External:   false,
		})
	}

	return alliesWithDrainedHPCounter
}

package furina

import (
	"fmt"
	"math"

	"github.com/genshinsim/gcsim/internal/common"
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/core/targets"
)

var skillFrames [][]int

const (
	skillHitmark     = 18 // TODO:
	salonInitialTick = 60 // TODO:
	particleICDKey   = "furina-skill-particle-icd"
	skillKey         = "furina-skill"
	skillMaxDuration = 30 * 60

	salonMemberKey = "Salon Member"

	chevalmarinInitialTick  = 72.3333
	chevalmarinIntervalMean = 97.5858
	chevalmarinTravel       = 20
	chevalmarinAoE          = 0.5

	usherInitialTick    = 74.5000
	usherIntervalMean   = 202.138
	usherIntervalStddev = 5
	usherTravel         = 40
	usherAoE            = 2.5

	crabalettaInitialTick  = 71.5926
	crabalettaIntervalMean = 313.859
	crabalettaTravel       = 41
	crabalettaAoE          = 3.5

	singerInterval = 120
)

func (c *char) calcSalonTick(tickNum int, initialTick, interval float64) int {
	// the distribution is left skewed. We approxiamated with boxcox with lambda 0.728
	// then used the transformation to convert from norm dist to the experimental distribution
	randOffset := math.Pow(common.Max(c.Core.Rand.NormFloat64()*1.0403+4.073023273, 0.0), (1/0.728)) - 7

	// this limits the offset to [-7, 7]
	randOffset = common.Min(randOffset, 7)
	return int(math.Round(initialTick + float64(tickNum)*interval + randOffset))
}

func init() {
	skillFrames = make([][]int, 2)
	skillFrames[ousia] = frames.InitAbilSlice(54) // E -> Q
	skillFrames[ousia][action.ActionAttack] = 53  // E -> N1
	skillFrames[ousia][action.ActionCharge] = 53  // E -> CA
	skillFrames[ousia][action.ActionBurst] = 54   // E -> Q
	skillFrames[ousia][action.ActionDash] = 18    // E -> D
	skillFrames[ousia][action.ActionJump] = 18    // E -> J
	skillFrames[ousia][action.ActionWalk] = 42    // E -> W
	skillFrames[ousia][action.ActionSwap] = 52    // TODO: E -> Swap

	skillFrames[pneuma] = frames.InitAbilSlice(57)
	skillFrames[pneuma][action.ActionAttack] = 56 // E -> N1
	skillFrames[pneuma][action.ActionCharge] = 56 // E -> CA
	skillFrames[pneuma][action.ActionBurst] = 57  // E -> Q
	skillFrames[pneuma][action.ActionDash] = 15   // E -> D
	skillFrames[pneuma][action.ActionJump] = 18   // E -> J
	skillFrames[pneuma][action.ActionWalk] = 41   // E -> W
	skillFrames[pneuma][action.ActionSwap] = 55   // TODO: E -> Swap
}
func (c *char) Skill(p map[string]int) (action.Info, error) {
	if c.Base.Cons >= 6 {
		c.c6Count = 0
		c.AddStatus(c6Key, 10*60, true)
	}
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
	c.AddStatus(skillKey, 1800+skillHitmark, false)
	c.summonSinger(c.Core.F, skillHitmark)
	c.SetCDWithDelay(action.ActionSkill, 1200, 10)
	return action.Info{
		Frames:          frames.NewAbilFunc(skillFrames[pneuma]),
		AnimationLength: skillFrames[pneuma][action.InvalidAction],
		CanQueueAfter:   skillFrames[pneuma][action.ActionDash],
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

	c.AddStatus(skillKey, 1736+skillHitmark, false)
	c.summonSalonMembers(skillHitmark)
	c.SetCDWithDelay(action.ActionSkill, 1200, 0)

	return action.Info{
		Frames:          frames.NewAbilFunc(skillFrames[ousia]),
		AnimationLength: skillFrames[ousia][action.InvalidAction],
		CanQueueAfter:   skillFrames[ousia][action.ActionDash],
		State:           action.SkillState,
	}, nil
}

func (c *char) summonSalonMembers(delay int) {
	c.Core.Tasks.Add(func() {
		src := c.Core.F
		c.lastSummonSrc = src

		c.Core.Tasks.Add(
			c.surintendanteChevalmarin(src, 0),
			c.calcSalonTick(0, chevalmarinInitialTick, chevalmarinIntervalMean),
		)
		c.Core.Tasks.Add(
			c.gentilhommeUsher(src, 0),
			c.calcSalonTick(0, usherInitialTick, usherIntervalMean),
		)
		c.Core.Tasks.Add(
			c.mademoiselleCrabaletta(src, 0),
			c.calcSalonTick(0, crabalettaInitialTick, crabalettaIntervalMean),
		)
	}, delay)
}

func (c *char) summonSinger(src, delay int) {
	c.Core.Tasks.Add(c.singerOfManyWaters(src), delay)
}

func (c *char) queueSalonAttack(src int, ai combat.AttackInfo, ap combat.AttackPattern, delay int, callbacks ...combat.AttackCBFunc) {
	c.QueueCharTask(func() {
		if src != c.lastSummonSrc {
			return
		}

		if !c.StatusIsActive(skillKey) {
			return
		}

		c.Core.QueueAttack(ai, ap, 0, 0, callbacks...)
	}, delay)
}

func (c *char) surintendanteChevalmarin(src, tick int) func() {
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
		ap := combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), geometry.Point{}, chevalmarinAoE)
		if c.Base.Cons >= 4 {
			c.queueSalonAttack(src, ai, ap, chevalmarinTravel, c.particleCB, c.c4cb)
		} else {
			c.queueSalonAttack(src, ai, ap, chevalmarinTravel, c.particleCB)
		}

		interval := c.calcSalonTick(tick+1, chevalmarinInitialTick, chevalmarinIntervalMean) - (c.Core.F - src)
		c.Core.Tasks.Add(c.surintendanteChevalmarin(src, tick+1), interval)
	}
}

func (c *char) gentilhommeUsher(src, tick int) func() {
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

		ap := combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), geometry.Point{}, usherAoE)
		if c.Base.Cons >= 4 {
			c.queueSalonAttack(src, ai, ap, usherTravel, c.particleCB, c.c4cb)
		} else {
			c.queueSalonAttack(src, ai, ap, usherTravel, c.particleCB)
		}

		interval := c.calcSalonTick(tick+1, usherInitialTick, usherIntervalMean) - (c.Core.F - src)
		c.Core.Tasks.Add(c.gentilhommeUsher(src, tick+1), interval)
	}
}

func (c *char) mademoiselleCrabaletta(src, tick int) func() {
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

		ap := combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), geometry.Point{}, crabalettaAoE)

		if c.Base.Cons >= 4 {
			c.queueSalonAttack(src, ai, ap, crabalettaTravel, c.particleCB, c.c4cb)
		} else {
			c.queueSalonAttack(src, ai, ap, crabalettaTravel, c.particleCB)
		}

		interval := c.calcSalonTick(tick+1, crabalettaInitialTick, crabalettaIntervalMean) - (c.Core.F - src)
		c.Core.Tasks.Add(c.mademoiselleCrabaletta(src, tick+1), interval)
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
		intervalDelta := common.Min(c.MaxHP()/1000.0*0.004, 0.16)

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

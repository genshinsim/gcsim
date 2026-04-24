package columbina

import (
	"cmp"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/enemy"
)

var (
	skillFrames      []int
	skillHitmarksLCr = 40
	skillHitmarksLC  = 40
	skillHitmarksLB  = []int{40, 45, 50, 55, 60}
	skillHitmarks    [3][]int
)

const (
	skillHitmark   = 24
	particleICDKey = "columbina-particle-icd"
	skillKey       = "columbina-skill"
	gravityKey     = "columbina-gravity"
	gravityMax     = 60
)

type lunarReaction int

const (
	LunarCharge lunarReaction = iota
	LunarBloom
	LunarCrystallize
)

func init() {
	skillFrames = frames.InitAbilSlice(26)

	skillHitmarks[LunarCharge] = []int{skillHitmarksLC}
	skillHitmarks[LunarBloom] = skillHitmarksLB
	skillHitmarks[LunarCrystallize] = []int{skillHitmarksLCr}
}

func (c *char) skillInit() {
	makeHook := func(reaction info.ReactionType) func(args ...any) {
		return func(args ...any) {
			if _, ok := args[0].(*enemy.Enemy); !ok {
				return
			}
			if !c.StatusIsActive(skillKey) {
				return
			}
			c.gravityLastReaction = reaction
			c.AddStatus(gravityKey, 2*60, false)
			if !c.gravityTask {
				c.gravityAccum()
			}
		}
	}

	c.Core.Events.Subscribe(event.OnLunarCharged, makeHook(info.ReactionTypeLunarCharged), "columbina-gravity-lc")
	c.Core.Events.Subscribe(event.OnLunarBloom, makeHook(info.ReactionTypeLunarBloom), "columbina-gravity-lb")
	c.Core.Events.Subscribe(event.OnLunarCrystallize, makeHook(info.ReactionTypeLunarCrystallize), "columbina-gravity-lcr")

	c.Core.Events.Subscribe(event.OnEnemyDamage, func(args ...any) {
		atk := args[1].(*info.AttackEvent)
		if !attacks.AttackTagIsLunar(atk.Info.AttackTag) {
			return
		}
		if !c.StatusIsActive(skillKey) {
			return
		}
		c.AddStatus(gravityKey, 2*60, false)
		if !c.gravityTask {
			c.gravityAccum()
		}
		switch atk.Info.AttackTag {
		case attacks.AttackTagDirectLunarCharged, attacks.AttackTagReactionLunarCharge:
			c.gravityLastReaction = info.ReactionTypeLunarCharged
		case attacks.AttackTagDirectLunarBloom:
			c.gravityLastReaction = info.ReactionTypeLunarBloom
		case attacks.AttackTagDirectLunarCrystallize, attacks.AttackTagReactionLunarCrystallize:
			c.gravityLastReaction = info.ReactionTypeLunarCrystallize
		}
	}, "columbina-gravity-on-dmg")
}

func (c *char) gravityAccum() {
	if !c.StatusIsActive(gravityKey) {
		c.gravityTask = false
		return
	}

	if !c.StatusIsActive(skillKey) {
		c.gravityTask = false
		return
	}
	c.gravityTask = true
	amt := 10 * 0.05 * (1 + c.c2GravityRate()) // 10 gravity per 1s
	switch c.gravityLastReaction {
	case info.ReactionTypeLunarCharged:
		c.gravity[LunarCharge] += amt
	case info.ReactionTypeLunarBloom:
		c.gravity[LunarBloom] += amt
	case info.ReactionTypeLunarCrystallize:
		c.gravity[LunarCrystallize] += amt
	}

	if c.totalGravity() >= gravityMax {
		c.gravityTick(true)
	}
	c.QueueCharTask(c.gravityAccum, 0.05*60)
}

func (c *char) totalGravity() float64 {
	sum := 0.0
	for _, g := range c.gravity {
		sum += g
	}
	return sum
}

func (c *char) clearGravity() {
	for i := range c.gravity {
		c.gravity[i] = 0
	}
}

// returns the index of the maximum value. In case of ties, the earlier index is returned
func argmax[T cmp.Ordered](x T, y ...T) int {
	max_ind := 0
	max_val := x
	for i, v := range y {
		if v > max_val {
			max_ind = i + 1
			max_val = v
		}
	}
	return max_ind
}

func (c *char) gravityTick(clearGravity bool) {
	// ties are broken by team member count for c1 for first cast
	// TODO: ties for non C1 tick is broken randomly by the first two reactions. https://discord.com/channels/763583452762734592/1460819588446163056/1463116888673616014
	electro := 0
	dendro := 0
	geo := 0
	for _, char := range c.Core.Player.Chars() {
		switch char.Base.Element {
		case attributes.Electro:
			electro += 1
		case attributes.Dendro:
			dendro += 1
		case attributes.Geo:
			geo += 1
		}
	}

	maxReaction := LunarCharge
	maxGravity := 0.0
	switch argmax(electro, dendro, geo) {
	case 0:
		maxReaction = LunarCharge
	case 1:
		maxReaction = LunarBloom
	case 2:
		maxReaction = LunarCrystallize
	}

	for i, g := range c.gravity {
		if g > maxGravity {
			maxGravity = g
			maxReaction = lunarReaction(i)
		}
	}

	if clearGravity {
		c.clearGravity()
	}

	var mult float64
	var atkTag attacks.AttackTag
	var elem attributes.Element
	var abil string
	radius := 6.0
	travel := 0
	switch maxReaction {
	case LunarCharge:
		mult = skillLC[c.TalentLvlSkill()]
		atkTag = attacks.AttackTagDirectLunarCharged
		elem = attributes.Electro
		abil = "Gravity Interference (Lunar-Charged)"

	case LunarBloom:
		mult = skillLB[c.TalentLvlSkill()]
		atkTag = attacks.AttackTagDirectLunarBloom
		elem = attributes.Dendro
		abil = "Gravity Interference (Lunar-Bloom)"
		radius = 0.5
		travel = 20
	case LunarCrystallize:
		mult = skillLCr[c.TalentLvlSkill()]
		atkTag = attacks.AttackTagDirectLunarCrystallize
		elem = attributes.Geo
		abil = "Gravity Interference (Lunar-Crystallize)"

	default:
		return
	}

	ai := info.AttackInfo{
		ActorIndex:       c.Index(),
		Abil:             abil,
		AttackTag:        atkTag,
		ICDTag:           attacks.ICDTagNone,
		StrikeType:       attacks.StrikeTypeDefault,
		Element:          elem,
		Durability:       0,
		UseHP:            true,
		IgnoreDefPercent: 1.0,
		Mult:             mult,
	}

	ap := combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, radius)
	c.a1OGravityTick()
	c.c1OnGravityTick(maxReaction)
	c.c2OnGravityTick(maxReaction)

	for _, delay := range skillHitmarks[maxReaction] {
		ai.FlatDmg = c.c4OnGravityTickFlatDMG(maxReaction)
		c.Core.QueueAttack(ai, ap, delay, delay+travel)
	}
}

func (c *char) Skill(p map[string]int) (action.Info, error) {
	c.QueueCharTask(func() {
		c.skillSrc = c.Core.F
		ai := info.AttackInfo{
			ActorIndex: c.Index(),
			Abil:       "Skill",
			AttackTag:  attacks.AttackTagElementalArt,
			ICDTag:     attacks.ICDTagNone,
			ICDGroup:   attacks.ICDGroupDefault,
			StrikeType: attacks.StrikeTypeDefault,
			Element:    attributes.Hydro,
			Durability: 25,
			UseHP:      true,
			Mult:       skill[c.TalentLvlSkill()],
		}
		ap := combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 6)
		c.Core.QueueAttack(ai, ap, 0, 0)
		if !c.StatusIsActive(skillKey) {
			c.clearGravity()
		}
		c.AddStatus(skillKey, 25*60+1, true)
		c.QueueCharTask(c.skillTickTask(c.skillSrc), 126)
		c.SetCDWithDelay(action.ActionSkill, 17*60, 0)
		c.c1OnSkill()
	}, skillHitmark)

	return action.Info{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillHitmark,
		State:           action.SkillState,
	}, nil
}

func (c *char) particleCB(a info.AttackCB) {
	if a.Target.Type() != info.TargettableEnemy {
		return
	}
	if c.StatusIsActive(particleICDKey) {
		return
	}
	c.AddStatus(particleICDKey, 3*60, false)
	if c.Core.Rand.Float64() < 0.67 {
		c.Core.QueueParticle(c.Base.Key.String(), 1, attributes.Hydro, c.ParticleDelay)
	} else {
		c.Core.QueueParticle(c.Base.Key.String(), 2, attributes.Hydro, c.ParticleDelay)
	}
}

// Helper function that handles damage, healing, and particle components of every tick of her E
func (c *char) skillTick() {
	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Skill (DoT)",
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Hydro,
		Durability: 25,
		UseHP:      true,
		Mult:       skillDoT[c.TalentLvlSkill()],
	}

	radius := 4.0
	if c.Core.Player.GetMoonsignLevel() >= 2 {
		radius = 6
	}
	ap := combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, radius)
	c.Core.QueueAttack(ai, ap, 0, 0, c.particleCB)
}

func (c *char) skillTickTask(src int) func() {
	return func() {
		if c.skillSrc > src {
			return
		}

		if !c.StatusIsActive(skillKey) {
			return
		}

		c.skillTick()

		c.Core.Tasks.Add(c.skillTickTask(src), 120)
	}
}

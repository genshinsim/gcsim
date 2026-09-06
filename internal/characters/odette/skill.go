package odette

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var (
	skillFrames       []int
	skillRecastFrames []int
)

const (
	skillHitmark              = 23
	particleICDKey            = "odette-particle-icd"
	danceDoubleKey            = "odette-dance-double"
	danceDoubleUpgradeKey     = "odette-dance-double-upgrade"
	skillFirstTickDelay       = 134
	skillRecastKey            = "odette-skill-recast"
	skillRecastFinalHitmark   = 11 + 9 + 8 + 34
	skillRecastFirstTickDelay = 52
)

var (
	skillTickDelay         = []int{109, 125}
	skillRecastDoTHitmarks = []int{11, 11 + 9, 11 + 9 + 8}
)

func init() {
	skillFrames = frames.InitAbilSlice(42) // E -> J
	skillFrames[action.ActionAttack] = 41
	skillFrames[action.ActionSkill] = 41
	skillFrames[action.ActionBurst] = 41
	skillFrames[action.ActionDash] = 40
	skillFrames[action.ActionWalk] = 41
	skillFrames[action.ActionSwap] = 40

	skillRecastFrames = frames.InitAbilSlice(76) // E -> Q
	skillRecastFrames[action.ActionAttack] = 75
	skillRecastFrames[action.ActionBurst] = 76
	skillRecastFrames[action.ActionDash] = 74
	skillRecastFrames[action.ActionJump] = 75
	skillRecastFrames[action.ActionWalk] = 75
	skillRecastFrames[action.ActionSwap] = 74
}

// With slow, graceful dance steps, Odette deals AoE Cryo DMG to the opponent, and also summons her
// Solo Dance Double to the field.
// If a Dance Double summoned by Odette is already on the field, this will re-summon the Dance
// Double and reset its duration.

// Solo Dance Double
// Alternates between the Plume and Wing dance moves, periodically attacking nearby opponents and
// dealing to them AoE Cryo DMG.
func (c *char) Skill(p map[string]int) (action.Info, error) {
	if c.useSpecialSkill() {
		return c.skillRecast(p)
	}

	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Adagio: Phantom Night Dancers",
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagElementalArt,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Cryo,
		Durability: 25,
		Mult:       skill[c.TalentLvlSkill()],
	}

	ap := combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 6)

	c.Core.QueueAttack(
		ai,
		ap,
		skillHitmark,
		skillHitmark,
		c.particleCB,
	)
	c.summonDanceDouble(skillFirstTickDelay)
	c.AddStatus(skillRecastKey, 6*60+skillHitmark, false)
	c.SetCDWithDelay(action.ActionSkill, 15*60, 14)
	return action.Info{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionSwap], // earliest cancel
		State:           action.SkillState,
	}, nil
}

// Additionally, for 6s after unleashing the Elemental Skill Adagio: Phantom Night Dancers, Odette's
// Elemental Skill Adagio: Phantom Night Dancers will become the special Elemental Skill Adagio:
// Coda at Dawn's Tolling instead, where a dance duet deals AoE Cryo DMG to nearby opponents over
// time. Then, when the duet ends, she deals another instance of AoE Cryo DMG that is considered
// Stellar-Conduct or Stellar Swirl DMG.
func (c *char) skillRecast(_ map[string]int) (action.Info, error) {
	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Coda at Dawn's Tolling",
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagElementalArt,
		ICDGroup:   attacks.ICDGroupOdetteDanceDuo,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Cryo,
		Durability: 25,
		Mult:       codaDoT[c.TalentLvlSkill()],
	}
	ap := combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 6)
	for _, delay := range skillRecastDoTHitmarks {
		c.Core.QueueAttack(ai, ap, delay, delay, c.particleCB)
	}

	c.QueueCharTask(func() {
		aiFinal := info.AttackInfo{
			ActorIndex:       c.Index(),
			Abil:             "Coda at Dawn's Tolling" + stellarConductText,
			AttackTag:        attacks.AttackTagDirectStellarConduct,
			ICDTag:           attacks.ICDTagNone,
			ICDGroup:         attacks.ICDGroupDefault,
			StrikeType:       attacks.StrikeTypeDefault,
			Element:          attributes.Cryo,
			Mult:             codaSSC[c.TalentLvlSkill()] * c.a4StellarGlimmerMult(),
			IgnoreDefPercent: 1,
		}
		if c.getRadiance() == radianceStellarSwirl {
			aiFinal.Abil = "Coda at Dawn's Tolling" + stellarSwirlText
			aiFinal.AttackTag = attacks.AttackTagDirectStellarSwirl
			aiFinal.Mult = codaSSw[c.TalentLvlSkill()] * c.a4StellarGlimmerMult()
		}

		c.Core.QueueAttack(aiFinal, ap, 0, 0, c.particleCB)
		c.AddStatus(danceDoubleUpgradeKey, c.StatusDuration(danceDoubleKey), false)

		c.c1OnSkillRecast(aiFinal.AttackTag)
	}, skillRecastFinalHitmark)

	// cancel existing dance tickers during the recast
	src := c.Core.F
	c.danceDoubleSrc = src
	// restart dance double at the end of the recast
	c.Core.Tasks.Add(func() {
		c.danceDoubleTicker(src, 0)
	}, skillRecastFinalHitmark+skillRecastFirstTickDelay)

	c.SetCD(action.ActionSpecialSkill, 15*60)
	return action.Info{
		Frames:          frames.NewAbilFunc(skillRecastFrames),
		AnimationLength: skillRecastFrames[action.InvalidAction],
		CanQueueAfter:   skillRecastFrames[action.ActionSwap], // earliest cancel
		State:           action.SkillState,
	}, nil
}

func (c *char) summonDanceDouble(firstTickDelay int) {
	src := c.Core.F
	c.danceDoubleSrc = src
	c.AddStatus(danceDoubleKey, 20*60, false)
	c.Core.Tasks.Add(func() { c.danceDoubleTicker(src, 0) }, firstTickDelay)

	c.a1OnDanceSummon()
	c.c2OnDanceSummon()
}

func (c *char) danceDoubleTicker(src, count int) {
	if c.danceDoubleSrc != src {
		return
	}

	if !c.StatusIsActive(danceDoubleKey) {
		return
	}

	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Cryo,
		Durability: 25,
	}
	if count%2 == 0 {
		ai.Abil = "\"Plume\" Dance Move"
		ai.Mult = plume[c.TalentLvlSkill()]
	} else {
		ai.Abil = "\"Wing\" Dance Move"
		ai.Mult = wing[c.TalentLvlSkill()]
	}

	ap := combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 4)

	c.Core.QueueAttack(ai, ap, 0, 0)
	c.Core.Tasks.Add(func() { c.danceDoubleTicker(src, count+1) }, skillTickDelay[count%2])

	if !c.StatusIsActive(danceDoubleUpgradeKey) {
		return
	}

	radiance := c.getRadiance()
	if radiance == radianceNone {
		return
	}

	aiStellar := info.AttackInfo{
		ActorIndex:       c.Index(),
		ICDTag:           attacks.ICDTagNone,
		ICDGroup:         attacks.ICDGroupDefault,
		StrikeType:       attacks.StrikeTypeDefault,
		Element:          attributes.Cryo,
		IgnoreDefPercent: 1,
	}

	baseAbil := "\"Plume\" Dance Move"
	mults := map[radianceState][]float64{radianceStellarConduct: plumeSSC, radianceStellarSwirl: plumeSSw}

	if count%2 != 0 {
		baseAbil = "\"Wing\" Dance Move"
		mults[radianceStellarConduct] = wingSSC
		mults[radianceStellarSwirl] = wingSSw
	}

	aiStellar.Mult = mults[radiance][c.TalentLvlSkill()] * c.a4StellarGlimmerMult()

	if radiance == radianceStellarConduct {
		aiStellar.AttackTag = attacks.AttackTagDirectStellarConduct
		aiStellar.Abil = baseAbil + stellarConductText
	} else {
		aiStellar.AttackTag = attacks.AttackTagDirectStellarSwirl
		aiStellar.Abil = baseAbil + stellarSwirlText
	}

	c.Core.QueueAttack(aiStellar, ap, 0, 0)
}

func (c *char) particleCB(a info.AttackCB) {
	if a.Target.Type() != info.TargettableEnemy {
		return
	}

	if c.StatusIsActive(particleICDKey) {
		return
	}

	c.AddStatus(particleICDKey, 6*60, true)
	c.Core.QueueParticle(c.Base.Key.String(), 5, attributes.Cryo, c.ParticleDelay)
}

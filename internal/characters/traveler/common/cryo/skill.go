package cryo

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var (
	skillFrames       [][]int
	skillTickHitmarks = []int{5, 8}
)

const (
	skillTapHitmark     = 24
	skillFirstTickDelay = 60
	skillInterval       = 3 * 60
	skillSscICD         = 0.2 * 60
	particleICDKey      = "travelercryo-particle-icd"
	skillKey            = "travelercryo-e"
	skillICDKey         = "travelercryo-e-icd"
	skillStacksMax      = 8
)

func init() {
	skillFrames = make([][]int, 2)

	// TODO: Placeholder using DMC frames

	// Male
	skillFrames[0] = frames.InitAbilSlice(37) // E -> N1
	skillFrames[0][action.ActionDash] = 29    // E -> D
	skillFrames[0][action.ActionJump] = 29    // E -> J
	skillFrames[0][action.ActionSwap] = 36    // E -> Swap

	// Female
	skillFrames[1] = frames.InitAbilSlice(37) // E -> N1/Q
	skillFrames[1][action.ActionDash] = 28    // E -> D
	skillFrames[1][action.ActionJump] = 28    // E -> J
	skillFrames[1][action.ActionSwap] = 35    // E -> Swap
}

// Stabs at the opponent with the Traveler's weapon, which releases an ice-cold fog. This deals Cryo
// DMG to opponents up ahead and forms a Frostpierce Star next to the Traveler.
//
// Frostpierce Star
//   - Follows your own party members around the field and periodically fires ice crystals at nearby
//     opponents, dealing Cryo DMG.
//     Radiance: Stellar-Conduct: The Frostpierce Star will no longer fire ice crystals at opponents
//     at intervals but will instead fire an ice crystal at opponents in a coordinated attack when
//     the Traveler hits an opponent with a Normal Attack, Charged Attack, or Plunging Attack. This
//     effect can trigger once every 0.2s.
//   - When an opponent is hit with an ice crystal, the Traveler gains 1 stack of
//     Frostglow (max 8 stacks).
//   - When not in combat, Frostglow expires after 30s.
func (c *Traveler) Skill(p map[string]int) (action.Info, error) {
	travel, ok := p["travel"]
	if !ok {
		travel = 10
	}
	c.skillTravel = travel

	ai := info.AttackInfo{
		ActorIndex:     c.Index(),
		Abil:           "Ice Fog Piercer",
		AttackTag:      attacks.AttackTagElementalArt,
		AdditionalTags: []attacks.AdditionalTag{attacks.AdditionalTagNightsoul},
		ICDTag:         attacks.ICDTagTravelerHoldDMG,
		ICDGroup:       attacks.ICDGroupDefault,
		StrikeType:     attacks.StrikeTypeDefault,
		Element:        attributes.Cryo,
		Durability:     25,
		Mult:           skill[c.TalentLvlSkill()],
	}

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 3.0),
		skillTapHitmark,
		skillTapHitmark,
		c.particleCB,
	)

	src := c.Core.F
	c.skillSrc = src
	c.QueueCharTask(func() { c.skillTicker(src) }, skillFirstTickDelay)
	c.AddStatus(skillKey, 12*60+skillFirstTickDelay+c.c4SkillBonusDur(), false)
	c.SetCD(action.ActionSkill, 15*60)

	return action.Info{
		Frames:          frames.NewAbilFunc(skillFrames[c.gender]),
		AnimationLength: skillFrames[c.gender][action.InvalidAction],
		CanQueueAfter:   skillFrames[c.gender][action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}, nil
}

func (c *Traveler) skillTicker(src int) {
	if c.skillSrc != src {
		return
	}

	if !c.StatusIsActive(skillKey) {
		return
	}

	if c.getRadiance() != radianceStellarConduct {
		c.queueCrystal(skillTickHitmarks...)
	}

	c.QueueCharTask(func() { c.skillTicker(src) }, skillInterval)
}

func (c *Traveler) queueCrystal(hitmarks ...int) {
	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Frostpierce Star",
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagElementalArt,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Cryo,
		Durability: 25,
		Mult:       crystal[c.TalentLvlSkill()],
	}

	for _, delay := range hitmarks {
		c.Core.QueueAttack(
			ai,
			combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 3.0),
			delay,
			delay+c.skillTravel,
			c.crystalCB,
			c.c2CB,
		)
	}
}

func (c *Traveler) naCaPlungeCB(a info.AttackCB) {
	if a.Target.Type() != info.TargettableEnemy {
		return
	}

	if c.getRadiance() != radianceStellarConduct {
		return
	}

	if !c.StatusIsActive(skillKey) {
		return
	}

	if c.StatusIsActive(skillICDKey) {
		return
	}

	c.AddStatus(skillICDKey, skillSscICD, true)
	c.queueCrystal(3)
}

func (c *Traveler) crystalCB(a info.AttackCB) {
	if a.Target.Type() != info.TargettableEnemy {
		return
	}

	c.flostglowStacks = min(c.flostglowStacks+1, skillStacksMax)
}

func (c *Traveler) particleCB(a info.AttackCB) {
	if a.Target.Type() != info.TargettableEnemy {
		return
	}
	if c.StatusIsActive(particleICDKey) {
		return
	}
	c.AddStatus(particleICDKey, 0.3*60, true)

	c.Core.QueueParticle(c.Base.Key.String(), 3, attributes.Cryo, c.ParticleDelay)
}

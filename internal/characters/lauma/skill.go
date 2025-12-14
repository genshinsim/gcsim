package lauma

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

var skillFrames [][]int

const (
	skillPressHitmark                 = 16
	frostgroveSanctuaryFirstHitPress  = 31
	frostgroveSanctuaryFirstHitHold   = 24
	frostgroveSanctuaryInterval       = 117
	skillHoldHitmark                  = 45
	skillConsumeDew                   = 29
	skillOffset                       = 0
	frostgroveSanctuaryKey            = "lauma-frostgrove-sanctuary"
	frostgroveSanctuaryParticleICDKey = "lauma-frostgrove-sanctuary-particle-icd"
	laumaC4RefundKey                  = "lauma-c4-refund"
	verdantDewKey                     = "verdant-dew"
)

func init() {
	skillFrames = make([][]int, 2)

	// skill (press) -> x
	skillFrames[0] = frames.InitAbilSlice(38)
	skillFrames[0][action.ActionSkill] = 42
	// skillFrames[0][action.ActionSkillHoldFramesOnly] = 40
	skillFrames[0][action.ActionDash] = 37
	skillFrames[0][action.ActionSwap] = 36

	// skill (hold=1) -> x
	skillFrames[1] = frames.InitAbilSlice(58)
	skillFrames[1][action.ActionAttack] = 57
	skillFrames[1][action.ActionSkill] = 64
	skillFrames[1][action.ActionBurst] = 59
	skillFrames[1][action.ActionJump] = 57
	skillFrames[1][action.ActionSwap] = 56
}

func (c *char) Skill(p map[string]int) (action.Info, error) {
	c.AddStatus(c1Key, 20*60, true)
	c.AddStatus(frostgroveSanctuaryKey, 15*60, true)
	c.AddStatus(a1Key, 20*60, true)
	c.c6SkillPaleHymnCount = 8

	// this piece of code is causing it idk why
	if p["hold"] == 0 {
		if int(c.Core.Flags.Custom[verdantDewKey]) == 0 {
			return c.skillPress()
		}
	}
	return c.skillHold()
}

func (c *char) skillPress() (action.Info, error) {
	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Runo: Dawnless Rest of Karsikko (Press)",
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Dendro,
		PoiseDMG:   50,
		Durability: 25,
		Mult:       skillPress[c.TalentLvlSkill()],
	}

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHit(
			c.Core.Combat.Player(),
			info.Point{Y: skillOffset},
			info.Point{Y: skillOffset},
			6,
		),
		skillPressHitmark,
		skillPressHitmark,
		c.frostgroveSantuary,
		c.applySkillShred,
	)

	c.SetCDWithDelay(action.ActionSkill, 12*60, 13)

	return action.Info{
		Frames:          frames.NewAbilFunc(skillFrames[0]),
		AnimationLength: skillFrames[0][action.InvalidAction],
		CanQueueAfter:   skillFrames[0][action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}, nil
}

func (c *char) skillHold() (action.Info, error) {
	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Runo: Dawnless Rest of Karsikko (Hold Hit 1)",
		AttackTag:  attacks.AttackTagElementalArtHold,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Dendro,
		Durability: 25,
		Mult:       skillHold1[c.TalentLvlSkill()],
	}

	aiDirectLB := info.AttackInfo{
		ActorIndex:       c.Index(),
		Abil:             "Runo: Dawnless Rest of Karsikko (Hold Hit 2)",
		AttackTag:        attacks.AttackTagDirectLunarBloom,
		ICDTag:           attacks.ICDTagNone,
		ICDGroup:         attacks.ICDGroupDefault,
		StrikeType:       attacks.StrikeTypeDefault,
		Element:          attributes.Dendro,
		UseEM:            true,
		IgnoreDefPercent: 1,
	}

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHit(
			c.Core.Combat.Player(),
			info.Point{Y: skillOffset},
			info.Point{Y: skillOffset},
			6,
		),
		skillHoldHitmark,
		skillHoldHitmark,
		c.frostgroveSantuary,
		c.applySkillShred,
	)

	c.Core.Tasks.Add(func() {
		dewCount := c.Core.Flags.Custom[verdantDewKey]

		aiDirectLB.Mult = skillHold2[c.TalentLvlSkill()] * dewCount

		c.Core.Flags.Custom[verdantDewKey] -= dewCount

		c.moonSong(int(dewCount))

		c.Core.QueueAttack(
			aiDirectLB,
			combat.NewCircleHit(
				c.Core.Combat.Player(),
				info.Point{Y: skillOffset},
				info.Point{Y: skillOffset},
				6,
			),
			skillHoldHitmark-skillConsumeDew,
			skillHoldHitmark-skillConsumeDew,
		)
	}, skillConsumeDew)

	c.SetCDWithDelay(action.ActionSkill, 12*60, 13)

	return action.Info{
		Frames:          frames.NewAbilFunc(skillFrames[1]),
		AnimationLength: skillFrames[1][action.InvalidAction],
		CanQueueAfter:   skillFrames[1][action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}, nil
}

func (c *char) moonSong(dewCount int) {
	if !c.StatusIsActive(burstKey) {
		return
	}

	c.addPaleHymn(dewCount*6, false)
}

func (c *char) applySkillShred(a info.AttackCB) {
	e, ok := a.Target.(*enemy.Enemy)
	if !ok {
		return
	}
	shredAmount := skillResShred[c.TalentLvlSkill()]
	e.AddResistMod(info.ResistMod{
		Base:  modifier.NewBaseWithHitlag("lauma-skill-shred", 10*60),
		Ele:   attributes.Dendro,
		Value: -shredAmount,
	})
	e.AddResistMod(info.ResistMod{
		Base:  modifier.NewBaseWithHitlag("lauma-skill-shred", 10*60),
		Ele:   attributes.Hydro,
		Value: -shredAmount,
	})
}

func (c *char) frostgroveSantuary(a info.AttackCB) {
	// TODO: check if frostgrove sanctuary gets overriden by skill while it's active

	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Frostgrove Sanctuary",
		AttackTag:  attacks.AttackTagElementalArtHold,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Dendro,
		Durability: 25,
		Mult:       frostgroveSanctuaryAtk[c.TalentLvlSkill()],
	}

	frostgroveSanctuaryFirstHit := frostgroveSanctuaryFirstHitPress
	if a.AttackEvent.Info.Abil == "Runo: Dawnless Rest of Karsikko (Hold Hit 1)" {
		frostgroveSanctuaryFirstHit = frostgroveSanctuaryFirstHitHold
	}

	c.Core.Tasks.Add(func() {
		c.frostgroveSantuary(a)
	}, frostgroveSanctuaryInterval)

	em := c.Stat(attributes.EM)
	ai.FlatDmg = em * frostgroveSanctuaryEM[c.TalentLvlSkill()]

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHit(
			c.Core.Combat.Player(),
			info.Point{Y: skillOffset},
			info.Point{Y: skillOffset},
			6,
		),
		frostgroveSanctuaryFirstHit,
		frostgroveSanctuaryFirstHit,
		c.particleCB,
		c.applySkillShred,
		c.c4Refund,
	)

	if c.Base.Cons >= 6 {
		if c.c6SkillPaleHymnCount <= 0 {
			return
		}

		c.c6SkillPaleHymnCount--

		aiC6 := info.AttackInfo{
			ActorIndex:       c.Index(),
			Abil:             "Frostgrove Sanctuary C6",
			AttackTag:        attacks.AttackTagDirectLunarBloom,
			ICDTag:           attacks.ICDTagNone,
			ICDGroup:         attacks.ICDGroupDefault,
			StrikeType:       attacks.StrikeTypeDefault,
			Element:          attributes.Dendro,
			Durability:       25,
			UseEM:            true,
			Mult:             1.85,
			IgnoreDefPercent: 1,
		}
		c.Core.QueueAttack(
			aiC6,
			combat.NewCircleHit(
				c.Core.Combat.Player(),
				info.Point{Y: skillOffset},
				info.Point{Y: skillOffset},
				6,
			),
			frostgroveSanctuaryFirstHit,
			frostgroveSanctuaryFirstHit,
			c.c6PaleHymn,
		)
	}
}

func (c *char) particleCB(a info.AttackCB) {
	if a.Target.Type() != info.TargettableEnemy {
		return
	}
	if c.StatusIsActive(frostgroveSanctuaryParticleICDKey) {
		return
	}
	c.AddStatus(frostgroveSanctuaryParticleICDKey, 3.3*60, true)

	count := 1.0
	if c.Core.Rand.Float64() < 0.3 {
		count = 2
	}
	c.Core.QueueParticle(c.Base.Key.String(), count, attributes.Dendro, c.ParticleDelay)
}

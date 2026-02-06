package lauma

import (
	"math"

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
	skillTicks                        = 8
	skillFirstTickDelay               = 62
	skillPressHitmark                 = 16
	frostgroveSanctuaryInterval       = 117
	skillHoldHitmark                  = 45
	skillConsumeDew                   = 29
	frostgroveSanctuaryKey            = "lauma-frostgrove-sanctuary"
	frostgroveSanctuaryParticleICDKey = "lauma-frostgrove-sanctuary-particle-icd"
	laumaC4RefundKey                  = "lauma-c4-refund"
	c6SkillHitName                    = "Frostgrove Sanctuary C6"
	moonSongIcdKey                    = "moonsong-icd"
)

func init() {
	skillFrames = make([][]int, 2)

	// skill (press) -> x
	skillFrames[0] = frames.InitAbilSlice(42)
	skillFrames[0][action.ActionAttack] = 38
	skillFrames[0][action.ActionCharge] = 38
	skillFrames[0][action.ActionBurst] = 38
	// skillFrames[0][action.ActionSkillHoldFramesOnly] = 40
	skillFrames[0][action.ActionDash] = 37
	skillFrames[0][action.ActionJump] = 38
	skillFrames[0][action.ActionWalk] = 38
	skillFrames[0][action.ActionSwap] = 36

	// skill (hold=1) -> x
	skillFrames[1] = frames.InitAbilSlice(64)
	skillFrames[1][action.ActionAttack] = 57
	skillFrames[1][action.ActionCharge] = 57
	skillFrames[1][action.ActionBurst] = 59
	skillFrames[1][action.ActionBurst] = 58
	skillFrames[1][action.ActionJump] = 57
	skillFrames[1][action.ActionWalk] = 58
	skillFrames[1][action.ActionSwap] = 56
}

func ceil(x float64) int {
	return int(math.Ceil(x))
}

func (c *char) Skill(p map[string]int) (action.Info, error) {
	c.AddStatus(c1Key, 20*60, true)
	c.AddStatus(a1Key, 20*60, true)
	c.c6OnSkill()
	c.AddStatus(frostgroveSanctuaryKey, 15*60, true)
	c.skillSrc = c.Core.F
	for i := 0.0; i < skillTicks; i++ {
		c.QueueCharTask(c.frostgroveSantuaryTick(c.skillSrc), skillFirstTickDelay+ceil(frostgroveSanctuaryInterval*i))
	}

	if p["hold"] == 0 || c.Core.Player.VerdantDew() > 0 {
		return c.skillPress()
	}
	return c.skillHold()
}

func (c *char) skillPress() (action.Info, error) {
	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Hymn of Hunting (Press)",
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Dendro,
		Durability: 25,
		Mult:       skillPress[c.TalentLvlSkill()],
	}

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHit(
			c.Core.Combat.Player(),
			c.Core.Combat.Player().Pos(),
			nil,
			6,
		),
		skillPressHitmark,
		skillPressHitmark,
		c.applySkillShredCB,
	)

	c.SetCDWithDelay(action.ActionSkill, 12*60, 13)

	return action.Info{
		Frames:          frames.NewAbilFunc(skillFrames[0]),
		AnimationLength: skillFrames[0][action.InvalidAction],
		CanQueueAfter:   skillFrames[0][action.ActionSwap], // earliest cancel
		State:           action.SkillState,
	}, nil
}

func (c *char) skillHold() (action.Info, error) {
	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Hymn of Eternal Rest (Hold)",
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
		Abil:             "Hymn of Eternal Rest (Hold) (Lunar-Bloom)",
		AttackTag:        attacks.AttackTagDirectLunarBloom,
		ICDTag:           attacks.ICDTagNone,
		ICDGroup:         attacks.ICDGroupDefault,
		StrikeType:       attacks.StrikeTypeDefault,
		Element:          attributes.Dendro,
		UseEM:            true,
		IgnoreDefPercent: 1,
		Mult:             skillHold2[c.TalentLvlSkill()],
	}

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHit(
			c.Core.Combat.Player(),
			c.Core.Combat.Player().Pos(),
			nil,
			6,
		),
		skillHoldHitmark,
		skillHoldHitmark,
		c.applySkillShredCB,
	)

	c.QueueCharTask(func() {
		dewCount := c.Core.Player.VerdantDew()
		c.Core.Player.ConsumeVerdantDew(dewCount)

		aiDirectLB.Mult = skillHold2[c.TalentLvlSkill()] * float64(dewCount)
		c.addMoonSong(dewCount)

		c.Core.QueueAttack(
			aiDirectLB,
			combat.NewCircleHitOnTarget(
				c.Core.Combat.Player(),
				nil,
				6,
			),
			skillHoldHitmark-skillConsumeDew,
			skillHoldHitmark-skillConsumeDew,
			c.applySkillShredCB,
		)
	}, skillConsumeDew)

	c.SetCDWithDelay(action.ActionSkill, 12*60, 13)

	return action.Info{
		Frames:          frames.NewAbilFunc(skillFrames[1]),
		AnimationLength: skillFrames[1][action.InvalidAction],
		CanQueueAfter:   skillFrames[1][action.ActionSwap], // earliest cancel
		State:           action.SkillState,
	}, nil
}

func (c *char) addMoonSong(moonsong int) {
	if moonsong <= 0 {
		return
	}

	if c.StatusIsActive(burstKey) && !c.StatusIsActive(moonSongIcdKey) {
		c.addPaleHymnMoonsong(moonsong * 6)
		c.AddStatus(moonSongIcdKey, c.StatusDuration(burstKey), true)
		c.moonSong = 0
		c.moonSongSrc = -1
		return
	}

	c.moonSong = moonsong
	src := c.Core.F
	c.moonSongSrc = src

	// remove moonsong stacks after 15s if not refreshed
	c.QueueCharTask(func() {
		if c.moonSongSrc == src {
			c.moonSong = 0
		}
	}, 15*60)
}

func (c *char) moonSongOnBurst() {
	if c.moonSong <= 0 {
		return
	}

	if !c.StatusIsActive(burstKey) {
		return
	}

	if c.StatusIsActive(moonSongIcdKey) {
		return
	}

	c.addPaleHymnMoonsong(c.moonSong * 6)
	c.AddStatus(moonSongIcdKey, c.StatusDuration(burstKey), true)
	c.moonSong = 0
	c.moonSongSrc = -1
}

func (c *char) applySkillShredCB(a info.AttackCB) {
	e, ok := a.Target.(*enemy.Enemy)
	if !ok {
		return
	}
	shredAmount := skillResShred[c.TalentLvlSkill()]
	e.AddResistMod(info.ResistMod{
		Base:  modifier.NewBaseWithHitlag("lauma-skill-shred-dendro", 10*60),
		Ele:   attributes.Dendro,
		Value: -shredAmount,
	})
	e.AddResistMod(info.ResistMod{
		Base:  modifier.NewBaseWithHitlag("lauma-skill-shred-hydro", 10*60),
		Ele:   attributes.Hydro,
		Value: -shredAmount,
	})
}

func (c *char) frostgroveSantuaryTick(src int) func() {
	return func() {
		if src != c.skillSrc {
			return
		}

		em := c.Stat(attributes.EM)
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
			FlatDmg:    em * frostgroveSanctuaryEM[c.TalentLvlSkill()],
		}

		c.Core.QueueAttack(
			ai,
			combat.NewCircleHitOnTarget(
				c.Core.Combat.Player(),
				nil,
				6,
			),
			0,
			0,
			c.particleCB,
			c.applySkillShredCB,
			c.c4RefundCB,
		)

		c.c6OnFrostgroveTick()
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

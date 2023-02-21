package kazuha

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var skillPressFrames [][]int
var skillHoldFrames [][]int

const (
	skillPressHitmark = 10
	skillPressRadius  = 5
	skillPressCDStart = 8
	skillHoldHitmark  = 33
	skillHoldRadius   = 9
	skillHoldCDStart  = 31
	particleICDKey    = "kazuha-particle-icd"
)

func init() {
	// Tap E
	skillPressFrames = make([][]int, 2)
	// Tap E -> X
	skillPressFrames[0] = frames.InitAbilSlice(77) // averaged all abils
	skillPressFrames[0][action.ActionHighPlunge] = 24
	// Tap E (Glide Cancel) -> X
	skillPressFrames[1] = frames.InitAbilSlice(69)
	skillPressFrames[1][action.ActionBurst] = 61
	skillPressFrames[1][action.ActionDash] = 61
	skillPressFrames[1][action.ActionJump] = 59
	skillPressFrames[1][action.ActionSwap] = 60

	// Hold E
	skillHoldFrames = make([][]int, 2)
	// Hold E -> X
	skillHoldFrames[0] = frames.InitAbilSlice(175) // averaged all abils
	skillHoldFrames[0][action.ActionHighPlunge] = 58
	// Hold E (Glide Cancel) -> X
	skillHoldFrames[1] = frames.InitAbilSlice(160)
	skillHoldFrames[1][action.ActionAttack] = 158
	skillHoldFrames[1][action.ActionBurst] = 159
	skillHoldFrames[1][action.ActionSwap] = 155
}

func (c *char) Skill(p map[string]int) action.ActionInfo {
	hold := p["hold"]
	glide := p["glide_cancel"]
	if glide < 0 {
		glide = 0
	}
	if glide > 1 {
		glide = 1
	}

	c.a1Absorb = attributes.NoElement

	if hold == 0 {
		return c.skillPress(glide)
	}
	return c.skillHold(glide)
}

func (c *char) makeParticleCB(count float64) combat.AttackCBFunc {
	return func(a combat.AttackCB) {
		if a.Target.Type() != combat.TargettableEnemy {
			return
		}
		if c.StatusIsActive(particleICDKey) {
			return
		}
		c.AddStatus(particleICDKey, 0.2*60, true)
		c.Core.QueueParticle(c.Base.Key.String(), count, attributes.Anemo, c.ParticleDelay)
	}
}

func (c *char) skillPress(glide int) action.ActionInfo {
	c.a1AbsorbCheckLocation = combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, skillPressRadius)
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Chihayaburu",
		AttackTag:  combat.AttackTagElementalArt,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeDefault,
		Element:    attributes.Anemo,
		Durability: 25,
		Mult:       skill[c.TalentLvlSkill()],
	}
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, skillPressRadius),
		0,
		skillPressHitmark,
		c.makeParticleCB(3),
	)
	if c.Base.Ascension >= 1 {
		c.Core.Tasks.Add(c.absorbCheckA1(c.Core.F, 0, int(skillPressHitmark/6)), 1)
	}

	cd := 360
	if c.Base.Cons >= 1 {
		cd = 324
	}
	if c.Base.Cons >= 6 {
		// TODO: when does the infusion kick in?
		// -> For now, assume that it starts on hitmark.
		c.Core.Tasks.Add(func() {
			c.c6()
		}, skillPressHitmark)
	}

	c.SetCDWithDelay(action.ActionSkill, cd, skillPressCDStart)

	shortestAction := action.ActionHighPlunge
	if glide == 1 {
		shortestAction = action.ActionJump
	}
	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillPressFrames[glide]),
		AnimationLength: skillPressFrames[glide][action.InvalidAction],
		CanQueueAfter:   skillPressFrames[glide][shortestAction], // earliest cancel
		State:           action.SkillState,
	}
}

func (c *char) skillHold(glide int) action.ActionInfo {
	c.a1AbsorbCheckLocation = combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, skillHoldRadius)
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Chihayaburu",
		AttackTag:  combat.AttackTagElementalArt,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeDefault,
		Element:    attributes.Anemo,
		Durability: 50,
		Mult:       skillHold[c.TalentLvlSkill()],
	}
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, skillHoldRadius),
		0,
		skillHoldHitmark,
		c.makeParticleCB(4),
	)
	if c.Base.Ascension >= 1 {
		c.Core.Tasks.Add(c.absorbCheckA1(c.Core.F, 0, int(skillHoldHitmark/6)), 1)
	}

	cd := 540
	if c.Base.Cons >= 1 {
		cd = 486
	}
	if c.Base.Cons >= 6 {
		// TODO: when does the infusion kick in?
		// -> For now, assume that it starts on hitmark.
		c.Core.Tasks.Add(func() {
			c.c6()
		}, skillHoldHitmark)
	}

	c.SetCDWithDelay(action.ActionSkill, cd, skillHoldCDStart)

	shortestAction := action.ActionHighPlunge
	if glide == 1 {
		shortestAction = action.ActionSwap
	}
	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillHoldFrames[glide]),
		AnimationLength: skillHoldFrames[glide][action.InvalidAction],
		CanQueueAfter:   skillHoldFrames[glide][shortestAction], // earliest cancel
		State:           action.SkillState,
	}
}

package varka

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

const defaultConversionElement = attributes.Physical

var (
	skillFrames             []int
	specialSkillFrames      []int
	fourWindsHitmark        = []int{34, 34 + 9}
	fourWindsHitHaltFrame   = []float64{0.09, 0}
	fourWindsCanBeDefHalted = []bool{true, false}
	fourWindsPoiseDmg       = []float64{65, 35}
)

const (
	skillHitmark   = 40
	particleICDKey = "varka-particle-icd"
	skillKey       = "sturm-und-drang"
	skillCD        = 16 * 60

	fourWindsCD = 11 * 60
)

func init() {
	skillFrames = frames.InitAbilSlice(55)
	skillFrames[action.ActionAttack] = 49
	skillFrames[action.ActionCharge] = 49 // Assumed same as NA
	skillFrames[action.ActionBurst] = 44
	skillFrames[action.ActionSkill] = 44 // Assumed same as Q
	skillFrames[action.ActionDash] = 65 - 19
	skillFrames[action.ActionJump] = 73 - 30
	skillFrames[action.ActionSwap] = 49

	specialSkillFrames = frames.InitAbilSlice(68)
	specialSkillFrames[action.ActionAttack] = 55
	specialSkillFrames[action.ActionCharge] = 64
	specialSkillFrames[action.ActionSkill] = 56
	specialSkillFrames[action.ActionBurst] = 55
	specialSkillFrames[action.ActionDash] = 75 - 19
	specialSkillFrames[action.ActionJump] = 85 - 30
	specialSkillFrames[action.ActionWalk] = 65
}

func (c *char) onExitField() {
	c.Core.Events.Subscribe(event.OnCharacterSwap, func(args ...any) {
		// do nothing if previous char wasn't varka
		prev := args[0].(int)
		if prev != c.Index() {
			return
		}
		if !c.StatusIsActive(skillKey) {
			return
		}

		c.DeleteStatus(skillKey)
	}, "varka-exit")
}

func (c *char) convertToFourWinds() bool {
	return c.conversionElem != defaultConversionElement
}

func (c *char) getConversionElem(prio ...attributes.Element) attributes.Element {
	for _, ele := range prio {
		for _, char := range c.Core.Player.Chars() {
			if char.Base.Element == ele {
				return ele
			}
		}
	}
	return defaultConversionElement
}

func (c *char) Skill(p map[string]int) (action.Info, error) {
	if c.useSpecialSkill() {
		return c.fourWinds(c.c6FreeSkill())
	}

	ai := info.AttackInfo{
		ActorIndex:         c.Index(),
		Abil:               "Windbound Execution",
		AttackTag:          attacks.AttackTagElementalArt,
		ICDTag:             attacks.ICDTagNone,
		ICDGroup:           attacks.ICDGroupDefault,
		StrikeType:         attacks.StrikeTypeBlunt,
		PoiseDMG:           100,
		Element:            attributes.Anemo,
		Durability:         25,
		Mult:               skillInitial[c.TalentLvlSkill()],
		HitlagHaltFrames:   0.09 * 60,
		HitlagFactor:       0.01,
		CanBeDefenseHalted: true,
	}
	c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 5), skillHitmark, skillHitmark, c.particleCB)

	c.QueueCharTask(func() {
		c.AddStatus(skillKey, 12*60, true)

		if c.convertToFourWinds() {
			c.fourWindsCDStacks = 0

			// discard any in progress CD queues
			c.DiscardActionCooldown(action.ActionSpecialSkill, fourWindsCD)

			// discard any ready special skill charges and start a new CD for each charge
			for range c.AvailableCDCharge[action.ActionSpecialSkill] {
				c.SetCD(action.ActionSpecialSkill, fourWindsCD)
			}

			// we specifically don't discard any previously queued CDs if they haven't started yet.
			// this aligns with in game tested behaviour with Chongyun C2

			// must be called after the CDs are reset
			c.c1OnSkill()
		}
	}, skillHitmark-1) // converts to skill state before the skill hitmark, relevant for Sac GS resetting cooldowns
	c.SetCDWithDelay(action.ActionSkill, skillCD, 39)
	return action.Info{
		Frames:          func(next action.Action) int { return skillFrames[next] },
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionJump], // earliest cancel
		State:           action.SkillState,
	}, nil
}

func (c *char) fourWinds(c6Free bool) (action.Info, error) {
	ele := []attributes.Element{c.conversionElem, attributes.Anemo}

	c1Mult := c.c1OnSpecialSkill()
	for i := range 2 {
		ai := info.AttackInfo{
			ActorIndex:         c.Index(),
			Abil:               "Four Winds' Ascension",
			AttackTag:          attacks.AttackTagElementalArt,
			ICDTag:             attacks.ICDTagNone,
			ICDGroup:           attacks.ICDGroupDefault,
			StrikeType:         attacks.StrikeTypeBlunt,
			PoiseDMG:           fourWindsPoiseDmg[i],
			Element:            ele[i],
			Durability:         25,
			Mult:               skillAscension[i][c.TalentLvlSkill()] * c.a1SkillMulti() * c1Mult,
			AdditionalTags:     []attacks.AdditionalTag{attacks.AdditionalTagVarkaSpecial},
			HitlagHaltFrames:   fourWindsHitHaltFrame[i] * 60,
			HitlagFactor:       0.01,
			CanBeDefenseHalted: fourWindsCanBeDefHalted[i],
		}
		ap := combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 6.5)
		c.Core.QueueAttack(ai, ap, fourWindsHitmark[i], fourWindsHitmark[i])
	}

	if !c6Free {
		c.QueueCharTask(func() {
			c.SetCD(action.ActionSpecialSkill, fourWindsCD)
			c.c6OnSkill()
		}, 39)
	}

	c.c2OnSpecialSkill()

	return action.Info{
		Frames:          func(next action.Action) int { return specialSkillFrames[next] },
		AnimationLength: specialSkillFrames[action.InvalidAction],
		CanQueueAfter:   specialSkillFrames[action.ActionBurst], // earliest cancel
		State:           action.SkillState,
	}, nil
}

func (c *char) fourWindsCDRedCB() func(ac info.AttackCB) {
	done := false
	return func(ac info.AttackCB) {
		if ac.Target.Type() != info.TargettableEnemy {
			return
		}

		if c.fourWindsCDStacks >= 15 {
			return
		}

		if !c.StatusIsActive(skillKey) {
			return
		}

		if done {
			return
		}
		done = true
		c.fourWindsCDStacks++

		amt := c.hexSkillCDReduction()
		c.ReduceActionCooldown(action.ActionSpecialSkill, amt)
	}
}

func (c *char) particleCB(ac info.AttackCB) {
	if ac.Target.Type() != info.TargettableEnemy {
		return
	}

	if c.StatusIsActive(particleICDKey) {
		return
	}

	c.AddStatus(particleICDKey, 0.3*60, true)
	c.Core.QueueParticle(c.Base.Key.String(), 6, attributes.Anemo, c.ParticleDelay)
}

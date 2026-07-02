package varka

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var (
	skillFrames      []int
	spearStormFrames []int
	fourWindsHitmark = []int{23, 25}
)

const (
	skillHitmark   = 19
	particleICDKey = "varka-particle-icd"
	skillKey       = "sturm-und-drang"
	skillCD        = 16 * 60

	fourWindsCD             = 11 * 60
	fourWindsCDReduceICDKey = "four-winds-cd-reduction-icd"
)

func init() {
	skillFrames = frames.InitAbilSlice(44)
	skillFrames[action.ActionAttack] = 19
	skillFrames[action.ActionSkill] = 22
	skillFrames[action.ActionBurst] = 19
	skillFrames[action.ActionDash] = 17
	skillFrames[action.ActionJump] = 18
	skillFrames[action.ActionSwap] = 17

	spearStormFrames = frames.InitAbilSlice(42)
	spearStormFrames[action.ActionAttack] = 28
	spearStormFrames[action.ActionBurst] = 28
	spearStormFrames[action.ActionDash] = 26
	spearStormFrames[action.ActionJump] = 26
	spearStormFrames[action.ActionWalk] = 32
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
	return c.conversionElem != attributes.NoElement
}

func (c *char) getConversionElem(prio ...attributes.Element) attributes.Element {
	for _, ele := range prio {
		for _, char := range c.Core.Player.Chars() {
			if char.Base.Element == ele {
				return ele
			}
		}
	}
	return attributes.Anemo
}

func (c *char) Skill(p map[string]int) (action.Info, error) {
	if c.StatModIsActive(skillKey) && c.convertToFourWinds() {
		return c.fourWinds(c.c6FreeSkill())
	}

	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Windbound Execution",
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Anemo,
		Durability: 25,
		Mult:       skillInitial[c.TalentLvlSkill()],
	}
	c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 3), skillHitmark, skillHitmark, c.particleCB)

	// c.SetCDWithDelay(action.ActionSkill, 16*60, skillHitmark)
	c.QueueCharTask(func() {
		c.AddStatus(skillKey, 12*60, true)

		if c.convertToFourWinds() {
			c.fourWindsCDStacks = 0
			c.fourWindsChargesStarted = 0
			c.fourWindsChargesAva = 0
			c.c1OnSkill()
			c.startFourWindsCD()
		}
		c.SetCD(action.ActionSkill, skillCD)
	}, skillHitmark)

	return action.Info{
		Frames:          func(next action.Action) int { return skillFrames[next] },
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}, nil
}

func (c *char) fourWinds(c6Free bool) (action.Info, error) {
	ele := []attributes.Element{c.conversionElem, attributes.Anemo}

	c1Mult := c.c1OnSpecialSkill()
	for i := range 2 {
		ai := info.AttackInfo{
			ActorIndex:     c.Index(),
			Abil:           "Four Winds' Ascension",
			AttackTag:      attacks.AttackTagElementalArt,
			ICDTag:         attacks.ICDTagNone,
			ICDGroup:       attacks.ICDGroupDefault,
			StrikeType:     attacks.StrikeTypeDefault,
			Element:        ele[i],
			Durability:     25,
			Mult:           skillAscension[i][c.TalentLvlSkill()] * c.a1SkillMulti() * c1Mult,
			AdditionalTags: []attacks.AdditionalTag{attacks.AdditionalTagVarkaSpecial},
		}
		ap := combat.NewBoxHitOnTarget(c.Core.Combat.Player(), nil, 4, 6)
		c.Core.QueueAttack(ai, ap, fourWindsHitmark[i], fourWindsHitmark[i])
	}

	if !c6Free {
		c.useFourWindsCharge()
		c.c6OnSkill()
	}

	c.c2OnSpecialSkill()

	return action.Info{
		Frames:          func(next action.Action) int { return spearStormFrames[next] },
		AnimationLength: spearStormFrames[action.InvalidAction],
		CanQueueAfter:   spearStormFrames[action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}, nil
}

func (c *char) fourWindsCDRedCB(ac info.AttackCB) {
	if ac.Target.Type() != info.TargettableEnemy {
		return
	}

	if c.fourWindsCDStacks >= 15 {
		return
	}

	if !c.StatusIsActive(skillKey) {
		return
	}

	if c.StatusIsActive(fourWindsCDReduceICDKey) {
		return
	}

	if c.fourWindsCDDoneF < 0 {
		return
	}

	c.AddStatus(fourWindsCDReduceICDKey, 0.1*60, true)
	c.fourWindsCDStacks++

	amt := c.hexSkillCDReduction()
	c.reduceFourWindsCD(amt)
}

func (c *char) reduceFourWindsCD(amt int) {
	if c.fourWindsCDDoneF < c.Core.F {
		panic("c.fourWindsCDDoneF should not be less than c.Core.F")
	}
	c.fourWindsCDDoneF = max(c.fourWindsCDDoneF-amt, c.Core.F)

	c.Core.Log.NewEventBuildMsg(glog.LogCooldownEvent, c.Index(), action.ActionSkill.String(), " (four winds) cooldown forcefully reduced").
		Write("type", action.ActionSkill.String()).
		Write("expiry", c.fourWindsCDDoneF-amt).
		Write("charges_remain", c.fourWindsCharges()).
		Write("four_winds_charges_started", c.fourWindsChargesStarted)

	c.queueCDTask()
}

func (c *char) startFourWindsCD() {
	modified := c.CDReduction(action.ActionSkill, fourWindsCD)
	c.fourWindsCDDoneF = c.Core.F + modified
	c.fourWindsChargesStarted += 1
	c.queueCDTask()
}

func (c *char) queueCDTask() {
	src := c.fourWindsCDDoneF
	c.Core.Tasks.Add(func() {
		if c.fourWindsCDDoneF != src {
			return
		}
		c.fourWindsChargesAva += 1
		c.Core.Log.NewEventBuildMsg(glog.LogCooldownEvent, c.Index(), action.ActionSkill.String(), " (four winds) cooldown ready").
			Write("type", action.ActionSkill.String()).
			Write("charges_remain", c.fourWindsCharges())

		if c.fourWindsChargesStarted < 2 {
			c.startFourWindsCD()
			return
		}

		c.fourWindsCDDoneF = -1
	}, src-c.Core.F)
}

func (c *char) resetFourWindsCD() {
	if c.fourWindsCDDoneF > c.Core.F {
		c.fourWindsCDDoneF = c.Core.F
		c.queueCDTask()
	}
}

func (c *char) fourWindsCD() int {
	if !c.StatusIsActive(skillKey) {
		return -1
	}

	if c.fourWindsCharges() > 0 {
		return 0
	}

	if c.fourWindsCDDoneF < 0 {
		return -1
	}

	return c.fourWindsCDDoneF - c.Core.F
}

func (c *char) fourWindsCharges() int {
	if !c.StatusIsActive(skillKey) {
		return 0
	}
	return c.fourWindsChargesAva
}

func (c *char) useFourWindsCharge() {
	if !c.StatusIsActive(skillKey) {
		return
	}
	if c.fourWindsChargesAva <= 0 {
		panic("unexpected charges less than 0")
	}
	c.fourWindsChargesAva -= 1

	// with C1, it's possible to stack up to 2, so we need to start a CD charge here
	if c.fourWindsCDDoneF < 0 && c.fourWindsChargesStarted < 2 {
		c.startFourWindsCD()
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

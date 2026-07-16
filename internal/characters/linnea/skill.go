package linnea

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/construct"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

const (
	skillPoundAttack int = iota
	skillHeavyAttack
)

var (
	skillFrames        []int
	skillRecastFrames  []int
	skillSuperInterval = []int{111, 63, 144} // time from pound -> heavy, heavy -> pound, pound -> pound
	skillSequence      = []int{skillPoundAttack, skillHeavyAttack, skillPoundAttack}
)

const (
	skillRecastKey      = "linnea-skill-recast"
	skillRecastDur      = 0.6 * 60
	skillDur            = 25 * 60
	skillSecondHitDelay = 16 // I have observed 16-19f most of the time, and occisionally 22-24
	skillSuperPower     = "linnea-super-power"
	skillSuperStart     = 114 - skillRecastDur // 114 since E cast, but 78f after recast expiry

	skillStandardPower     = "linnea-standard-power"
	skillMillionTonHitmark = 50
	skillStandardStart     = 48
	skillStandardInterval  = 340
	skillMillionAbil       = "Lumi Million Ton Crush"
	skillOverdriveAbil     = "Lumi Heavy Overdrive Hammer"

	skillCD      = 18 * 60
	skillCDStart = 1

	particleICDKey = "linnea-particle-icd"
	particleICD    = 9 * 60
	particleCount  = 3
)

func init() {
	skillFrames = frames.InitAbilSlice(33) // E -> D
	skillFrames[action.ActionAttack] = 20  // E -> recast
	skillFrames[action.ActionAim] = 20
	skillFrames[action.ActionSkill] = 20 // E -> recast
	skillFrames[action.ActionBurst] = 20
	skillFrames[action.ActionJump] = 20
	skillFrames[action.ActionWalk] = 33
	skillFrames[action.ActionSwap] = 20

	skillRecastFrames = frames.InitAbilSlice(10) // E -> D
}

func (c *char) Skill(p map[string]int) (action.Info, error) {
	if c.StatusIsActive(skillRecastKey) {
		return c.skillRecast(), nil
	}
	c.recastCount = 0
	c.DeleteStatus(skillStandardPower)
	c.DeleteStatus(skillSuperPower)

	c.AddStatus(skillRecastKey, skillRecastDur, true)
	src := c.Core.F
	c.skillRecastSrc = src
	c.QueueCharTask(func() {
		if c.skillRecastSrc != src {
			return
		}
		if !c.StatusIsActive(skillSuperPower) && c.Core.Constructs.CountByType(construct.GeoConstructLunarCrystallize) > 0 {
			c.skillHitNum = 1
		}
		c.AddStatus(skillSuperPower, skillDur, false)
		c.skillSrc = src
		c.a1OnLumi(src)
		c.Core.Tasks.Add(func() { c.lumiAttack(src) }, skillSuperStart)
	}, skillRecastDur)

	c.c1OnSkill()
	c.SetCDWithDelay(action.ActionSkill, skillCD, skillCDStart)

	return action.Info{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionSwap],
		State:           action.SkillState,
	}, nil
}

func (c *char) lumiAttack(src int) {
	if c.skillSrc != src {
		return
	}
	if !c.StatusIsActive(skillStandardPower) && !c.StatusIsActive(skillSuperPower) {
		return
	}

	attack := skillPoundAttack
	if c.StatusIsActive(skillSuperPower) && skillSequence[c.skillHitNum] == skillHeavyAttack {
		attack = skillHeavyAttack
	}

	switch attack {
	case skillPoundAttack:
		c.poundPoundHammer(src)
	case skillHeavyAttack:
		c.heavyOverdriveHammer()
	}

	delay := skillSuperInterval[c.skillHitNum]
	if c.StatusIsActive(skillStandardPower) {
		delay = skillStandardInterval
	} else {
		c.advanceSkillIndex()
	}

	c.Core.Tasks.Add(func() { c.lumiAttack(src) }, delay)
}

func (c *char) advanceSkillIndex() {
	if c.Core.Constructs.CountByType(construct.GeoConstructLunarCrystallize) == 0 {
		c.skillHitNum = 0
		return
	}

	c.skillHitNum++
	if c.skillHitNum >= 3 {
		c.skillHitNum = 0
	}
}

func (c *char) skillRecast() action.Info {
	c.recastCount++
	src := c.Core.F
	c.skillRecastSrc = src

	if c.recastCount == 3 {
		c.DeleteStatus(skillRecastKey)
		// do the big hit and then lower frequency
		c.skillHitNum = 0
		c.skillSrc = src
		c.a1OnLumi(src)
		c.Core.Tasks.Add(func() {
			c.millionTonHammer()
		}, skillMillionTonHitmark)
		c.AddStatus(skillStandardPower, skillDur, false)
		c.Core.Tasks.Add(func() { c.lumiAttack(src) }, skillMillionTonHitmark+skillStandardStart)
	} else {
		c.AddStatus(skillRecastKey, skillRecastDur, true)
		c.QueueCharTask(func() {
			if c.skillRecastSrc != src {
				return
			}
			if !c.StatusIsActive(skillSuperPower) && c.Core.Constructs.CountByType(construct.GeoConstructLunarCrystallize) > 0 {
				c.skillHitNum = 1
			}
			c.AddStatus(skillSuperPower, skillDur, false)
			c.skillSrc = src
			c.a1OnLumi(src)
			c.Core.Tasks.Add(func() { c.lumiAttack(src) }, skillSuperStart)
		}, skillRecastDur)
	}

	return action.Info{
		Frames:          frames.NewAbilFunc(skillRecastFrames),
		AnimationLength: skillRecastFrames[action.InvalidAction],
		CanQueueAfter:   skillRecastFrames[action.ActionJump],
		State:           action.SkillState, // TODO: does this matter?
	}
}

func (c *char) millionTonHammer() {
	ai := info.AttackInfo{
		Abil:             skillMillionAbil,
		ActorIndex:       c.Index(),
		AttackTag:        attacks.AttackTagDirectLunarCrystallize,
		ICDTag:           attacks.ICDTagNone,
		ICDGroup:         attacks.ICDGroupDefault,
		StrikeType:       attacks.StrikeTypeBlunt,
		Element:          attributes.Geo,
		Mult:             skillMillion[c.TalentLvlSkill()],
		UseDef:           true,
		IgnoreDefPercent: 1,
	}
	snap := c.Snapshot(&ai)
	c.c2MillionTonCDBonus(&snap)
	ap := combat.NewCircleHitOnTarget(c.LumiPos(), nil, 5)

	ae := info.AttackEvent{
		Info:        ai,
		Pattern:     ap,
		Snapshot:    snap,
		SourceFrame: c.Core.F,
	}

	ae.Callbacks = append(ae.Callbacks, c.particleCB)
	c.Core.QueueAttackEvent(&ae, 0)
	c.c2TriggerMoonDrift()
}

func (c *char) poundPoundHammer(src int) {
	ai := info.AttackInfo{
		Abil:       "Lumi Pound-Pound Pummeler 1",
		ActorIndex: c.Index(),
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagElementalArt,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeBlunt,
		Element:    attributes.Geo,
		Durability: 25,
		Mult:       skillPound[c.TalentLvlSkill()],
		UseDef:     true,
	}
	ap := combat.NewCircleHitOnTarget(c.LumiPos(), nil, 3)
	c.Core.QueueAttack(ai, ap, 0, 0, c.particleCB)
	c.Core.Tasks.Add(func() {
		if c.skillSrc != src {
			return
		}
		ai.Abil = "Lumi Pound-Pound Pummeler 2"
		c.Core.QueueAttack(ai, ap, 0, 0, c.particleCB)
	}, skillSecondHitDelay)
}

func (c *char) heavyOverdriveHammer() {
	ai := info.AttackInfo{
		Abil:             skillOverdriveAbil,
		ActorIndex:       c.Index(),
		AttackTag:        attacks.AttackTagDirectLunarCrystallize,
		ICDTag:           attacks.ICDTagNone,
		ICDGroup:         attacks.ICDGroupDefault,
		StrikeType:       attacks.StrikeTypeBlunt,
		Element:          attributes.Geo,
		Mult:             skillHeavy[c.TalentLvlSkill()],
		UseDef:           true,
		IgnoreDefPercent: 1,
	}
	snap := c.Snapshot(&ai)
	ap := combat.NewCircleHitOnTarget(c.LumiPos(), nil, 4)

	ae := info.AttackEvent{
		Info:        ai,
		Pattern:     ap,
		Snapshot:    snap,
		SourceFrame: c.Core.F,
	}
	ae.Callbacks = append(ae.Callbacks, c.particleCB)
	c.Core.QueueAttackEvent(&ae, 0)

	c.c2TriggerMoonDrift()
}

func (c *char) particleCB(a info.AttackCB) {
	if a.Target.Type() != info.TargettableEnemy {
		return
	}
	if c.StatusIsActive(particleICDKey) {
		return
	}
	c.AddStatus(particleICDKey, particleICD, true)

	c.Core.QueueParticle(c.Base.Key.String(), particleCount, attributes.Geo, c.ParticleDelay)
}

// NOTE: Lumi is not a Target in this implementation
func (c *char) LumiPos() info.Point {
	return c.Core.Combat.PrimaryTarget().Pos()
}

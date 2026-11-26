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
	skillSuperInterval = []int{62, 122, 90}
	skillSequence      = []int{skillPoundAttack, skillPoundAttack, skillHeavyAttack}
)

const (
	skillRecastKey  = "linnea-skill-recast"
	skillRecastDur  = 0.6 * 60
	skillDur        = 25 * 60
	skillSuperPower = "linnea-super-power"

	skillStandardPower       = "linnea-standard-power"
	skillMillionTonHitmark   = 50
	skillStandardStart       = 48
	skillStandardInterval    = 120
	skillPoundSecondHitDelay = 21
	skillMillionAbil         = "Lumi Million Ton Crush"
	skillOverdriveAbil       = "Lumi Heavy Overdrive Hammer"

	skillCD      = 18 * 60
	skillCDStart = 1

	particleICDKey = "linnea-particle-icd"
	particleICD    = 9 * 60
	particleCount  = 3
)

func init() {
	skillFrames = frames.InitAbilSlice(44)
	skillFrames[action.ActionAttack] = 33
	skillFrames[action.ActionAim] = 44
	skillFrames[action.ActionSkill] = 20
	skillFrames[action.ActionBurst] = 15
	skillFrames[action.ActionJump] = 14
	skillFrames[action.ActionWalk] = 32
	skillFrames[action.ActionDash] = 14
	skillFrames[action.ActionSwap] = 13

	skillRecastFrames = frames.InitAbilSlice(10) // ?
}

func (c *char) Skill(p map[string]int) (action.Info, error) {
	if c.StatusIsActive(skillRecastKey) {
		return c.skillRecast(), nil
	}
	c.recastCount = 0
	c.DeleteStatus(skillStandardPower)
	c.DeleteStatus(skillSuperPower)

	src := c.Core.F
	c.skillRecastSrc = src
	c.QueueCharTask(func() {
		if c.skillRecastSrc != src {
			return
		}
		c.summonLumiSuper(src)
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

func (c *char) skillRecast() action.Info {
	c.recastCount++
	src := c.Core.F
	c.skillRecastSrc = src

	if c.recastCount == 3 {
		c.DeleteStatus(skillRecastKey)
		c.summonLumiUltimate(src)
	} else {
		c.AddStatus(skillRecastKey, skillRecastDur, true)
		c.QueueCharTask(func() {
			if c.skillRecastSrc != src {
				return
			}
			c.summonLumiSuper(src)
		}, skillRecastDur)
	}

	return action.Info{
		Frames:          frames.NewAbilFunc(skillRecastFrames),
		AnimationLength: skillRecastFrames[action.InvalidAction],
		CanQueueAfter:   skillRecastFrames[action.ActionJump],
		State:           action.SkillState, // TODO: does this matter?
	}
}

func (c *char) lumiAttackStandard(src int) {
	if c.skillSrc != src {
		return
	}
	if !c.StatusIsActive(skillStandardPower) {
		return
	}
	c.poundPoundHammer(src)
	c.Core.Tasks.Add(func() { c.lumiAttackStandard(src) }, skillStandardInterval)
}

func (c *char) lumiAttackSuper(src int) {
	if c.skillSrc != src {
		return
	}
	if !c.StatusIsActive(skillSuperPower) {
		return
	}

	attack := skillPoundAttack

	moondrifts, _ := c.Core.Constructs.ConstructsByType(construct.GeoConstructLunarCrystallize)
	moondriftnearby := false
	playerPos := c.Core.Combat.Player().Pos()
	for _, moondrift := range moondrifts {
		if playerPos.Distance(moondrift.Pos()) < 16 {
			moondriftnearby = true
			break
		}
	}
	if skillSequence[c.skillHitNum] == skillHeavyAttack && moondriftnearby {
		attack = skillHeavyAttack
	}

	switch attack {
	case skillPoundAttack:
		c.poundPoundHammer(src)
	case skillHeavyAttack:
		c.heavyOverdriveHammer()
	}

	c.advanceSkillIndex()
	c.Core.Tasks.Add(func() { c.lumiAttackSuper(src) }, skillSuperInterval[c.skillHitNum])
}

func (c *char) advanceSkillIndex() {
	moondrifts, _ := c.Core.Constructs.ConstructsByType(construct.GeoConstructLunarCrystallize)
	moondriftnearby := false
	playerPos := c.Core.Combat.Player().Pos()
	for _, moondrift := range moondrifts {
		if playerPos.Distance(moondrift.Pos()) < 16 {
			moondriftnearby = true
			break
		}
	}

	if !moondriftnearby {
		return
	}

	c.skillHitNum++
	if c.skillHitNum >= 3 {
		c.skillHitNum = 0
	}
}

func (c *char) summonLumiStandard(src int) {
	c.DeleteStatus(skillStandardPower)
	c.DeleteStatus(skillSuperPower)
	c.AddStatus(skillStandardPower, skillDur, false)
	c.skillHitNum = 0
	c.skillSrc = src
	c.a1OnLumi(src)
	c.Core.Tasks.Add(func() { c.lumiAttackStandard(src) }, skillStandardStart)
}

func (c *char) summonLumiSuper(src int) {
	c.DeleteStatus(skillStandardPower)
	c.DeleteStatus(skillSuperPower)
	c.AddStatus(skillSuperPower, skillDur, false)
	// attempts to use millionTonHammer right after summoning in super form if there are moondrifts nearby
	c.skillHitNum = 1
	c.advanceSkillIndex()
	c.skillSrc = src
	c.a1OnLumi(src)
	c.Core.Tasks.Add(func() { c.lumiAttackSuper(src) }, skillSuperInterval[1])
}

func (c *char) summonLumiUltimate(src int) {
	c.DeleteStatus(skillStandardPower)
	c.DeleteStatus(skillSuperPower)
	c.AddStatus(skillStandardPower, skillDur, false)
	c.skillHitNum = 0
	c.skillSrc = src
	c.a1OnLumi(src)
	c.Core.Tasks.Add(func() {
		c.millionTonHammer()
	}, skillMillionTonHitmark)
	c.Core.Tasks.Add(func() { c.lumiAttackStandard(src) }, skillMillionTonHitmark+skillStandardStart)
}

func (c *char) millionTonHammer() {
	ai := info.AttackInfo{
		Abil:             skillMillionAbil,
		ActorIndex:       c.Index(),
		AttackTag:        attacks.AttackTagDirectLunarCrystallize,
		ICDTag:           attacks.ICDTagNone,
		ICDGroup:         attacks.ICDGroupDefault,
		StrikeType:       attacks.StrikeTypeBlunt,
		PoiseDMG:         150,
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
	c.c2TriggerMoonDrift(&ae)
}

func (c *char) poundPoundHammer(src int) {
	ai := info.AttackInfo{
		Abil:       "Lumi Pound-Pound Pummeler 1",
		ActorIndex: c.Index(),
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagElementalArt,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeBlunt,
		PoiseDMG:   35,
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
	}, skillPoundSecondHitDelay)
}

func (c *char) heavyOverdriveHammer() {
	ai := info.AttackInfo{
		Abil:             skillOverdriveAbil,
		ActorIndex:       c.Index(),
		AttackTag:        attacks.AttackTagDirectLunarCrystallize,
		ICDTag:           attacks.ICDTagNone,
		ICDGroup:         attacks.ICDGroupDefault,
		StrikeType:       attacks.StrikeTypeBlunt,
		PoiseDMG:         75,
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
	c.c2TriggerMoonDrift(&ae)
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

package mavuika

import (
	"errors"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/targets"
)

var (
	skillFrames             []int
	skillFramesHold         []int
	skillRecastFramesToBike []int
	skillRecastFramesToRing []int
	skillBikeRefreshFrames  []int
)

const (
	skillHitmark     = 16
	particleICDKey   = "mavuika-particle-icd"
	skillRecastCD    = 60
	skillRecastCDKey = "mavuika-skill-recast-cd"
)

func init() {
	skillFrames = frames.InitAbilSlice(28) // E -> Dash/Jump/Walk
	skillFrames[action.ActionAttack] = 18  // E -> N1
	skillFrames[action.ActionCharge] = 19  // E -> CA
	skillFrames[action.ActionSkill] = 19   // E -> Skill Recast
	skillFrames[action.ActionBurst] = 18   // E -> Burst
	skillFrames[action.ActionSwap] = 24    // E -> Swap

	skillFramesHold = frames.InitAbilSlice(43) // E -> N1
	skillFramesHold[action.ActionDash] = 42    // E -> Dash
	skillFramesHold[action.ActionJump] = 42    // E -> Jump
	skillFramesHold[action.ActionSwap] = 34    // E -> Swap

	skillRecastFramesToBike = frames.InitAbilSlice(24) // E -> Swap
	skillRecastFramesToBike[action.ActionAttack] = 13  // E -> N1
	skillRecastFramesToBike[action.ActionCharge] = 13  // E -> CA
	skillRecastFramesToBike[action.ActionBurst] = 13   // E -> Burst
	skillRecastFramesToBike[action.ActionDash] = 12    // E -> Dash
	skillRecastFramesToBike[action.ActionJump] = 13    // E -> Jump

	skillRecastFramesToRing = frames.InitAbilSlice(38) // E -> Jump
	skillRecastFramesToRing[action.ActionAttack] = 27  // E -> N1
	skillRecastFramesToRing[action.ActionCharge] = 28  // E -> CA
	skillRecastFramesToRing[action.ActionBurst] = 28   // E -> Burst
	skillRecastFramesToRing[action.ActionDash] = 37    // E -> Dash
	skillRecastFramesToRing[action.ActionSwap] = 24    // E -> Swap

	skillBikeRefreshFrames = frames.InitAbilSlice(39) // E -> E
	skillBikeRefreshFrames[action.ActionAttack] = 27  // E -> N1
	skillBikeRefreshFrames[action.ActionCharge] = 27  // E -> CA
	skillBikeRefreshFrames[action.ActionBurst] = 27   // E -> Burst
	skillBikeRefreshFrames[action.ActionDash] = 25    // E -> Dash
	skillBikeRefreshFrames[action.ActionJump] = 27    // E -> Jump
	skillBikeRefreshFrames[action.ActionSwap] = 24    // E -> Swap
}

func (c *char) nightsoulPointReduceTask(src int) {
	c.QueueCharTask(func() {
		if c.nightsoulSrc != src {
			return
		}
		val := 0.5
		if c.armamentState == bike {
			val += 0.4
			if c.Core.Player.CurrentState() == action.ChargeAttackState {
				val += 0.2
			}
		}
		c.reduceNightsoulPoints(val)
		c.nightsoulPointReduceTask(src)
	}, 0.1*60)
}

func (c *char) reduceNightsoulPoints(val float64) {
	val *= c.nightsoulConsumptionMul()
	if val == 0 {
		return
	}
	c.nightsoulState.ConsumePoints(val)

	if c.nightsoulState.Points() < 0.001 {
		c.exitNightsoul()
	}
}

func (c *char) exitNightsoul() {
	if !c.nightsoulState.HasBlessing() {
		return
	}
	c.nightsoulState.ExitBlessing()
	c.nightsoulState.ClearPoints()
	c.nightsoulSrc = -1
	if c.armamentState == bike && c.Core.Player.CurrentState() == action.NormalAttackState {
		c.NormalCounter = min(3, c.savedNormalCounter)
	} else {
		c.NormalCounter = 0
	}
	c.NormalHitNum = normalHitNum
}

func (c *char) enterNightsoulOrRegenerate(points float64) {
	if !c.nightsoulState.HasBlessing() {
		c.nightsoulState.EnterBlessing(points)
		c.nightsoulSrc = c.Core.F
		c.nightsoulPointReduceTask(c.nightsoulSrc)
		return
	}
	c.nightsoulState.GeneratePoints(points)
}

func (c *char) Skill(p map[string]int) (action.Info, error) {
	h := p["hold"]
	recast := p["recast"]
	if recast != 0 {
		if h > 0 {
			return action.Info{}, errors.New("cannot hold E while recasting")
		}
		if !c.nightsoulState.HasBlessing() {
			return action.Info{}, errors.New("cannot recast E while not in nightsoul blessing")
		}
		return c.skillRecast(), nil
	}

	var ai action.Info
	switch {
	case c.armamentState == bike && c.nightsoulState.HasBlessing():
		ai = c.skillBikeRefresh()
	case h > 0:
		ai = c.skillHold()
	default:
		ai = c.skillPress()
	}

	c.enterNightsoulOrRegenerate(c.nightsoulState.MaxPoints)
	return ai, nil
}

func (c *char) enterBike() {
	c.Core.Log.NewEvent("switching to bike state", glog.LogCharacterEvent, c.Index)
	c.armamentState = bike
	c.NormalHitNum = bikeHitNum
	c.NormalCounter = 0
	c.c6Bike()
}

func (c *char) exitBike() {
	c.Core.Log.NewEvent("switching to ring state", glog.LogCharacterEvent, c.Index)
	c.armamentState = ring
	c.NormalHitNum = normalHitNum
	c.ringSrc = c.Core.F

	c.skillRingTask(c.ringSrc)
	c.c2Ring()
}

func (c *char) skillRecast() action.Info {
	c.AddStatus(skillRecastCDKey, skillRecastCD, false)
	switch c.armamentState {
	case ring:
		c.enterBike()
		return action.Info{
			Frames:          frames.NewAbilFunc(skillRecastFramesToBike),
			AnimationLength: skillRecastFramesToBike[action.InvalidAction],
			CanQueueAfter:   skillRecastFramesToBike[action.ActionAttack],
			State:           action.SkillState,
		}

	default:
		c.exitBike()
		return action.Info{
			Frames:          frames.NewAbilFunc(skillRecastFramesToRing),
			AnimationLength: skillRecastFramesToRing[action.InvalidAction],
			CanQueueAfter:   skillRecastFramesToRing[action.ActionAttack],
			State:           action.SkillState,
		}
	}
}

func (c *char) skillHold() action.Info {
	ai := combat.AttackInfo{
		ActorIndex:     c.Index,
		Abil:           "The Named Moment (Flamestrider)",
		AttackTag:      attacks.AttackTagElementalArt,
		ICDTag:         attacks.ICDTagNone,
		AdditionalTags: []attacks.AdditionalTag{attacks.AdditionalTagNightsoul},
		ICDGroup:       attacks.ICDGroupDefault,
		StrikeType:     attacks.StrikeTypeBlunt,
		PoiseDMG:       75,
		Element:        attributes.Pyro,
		Durability:     25,
		Mult:           skill[c.TalentLvlSkill()],
		HitlagFactor:   0.05,
	}
	ap := combat.NewCircleHitOnTarget(
		c.Core.Combat.Player(),
		geometry.Point{Y: 1.0},
		6,
	)
	c.Core.QueueAttack(ai, ap, skillHitmark, skillHitmark, c.particleCB)
	c.enterBike()
	c.SetCDWithDelay(action.ActionSkill, 15*60, 18)
	c.Core.Tasks.Add(func() {
		c.AddStatus(skillRecastCDKey, skillRecastCD, false)
	}, 24)

	return c.getSkillCastActionInfo(skillFramesHold)
}

func (c *char) skillPress() action.Info {
	ai := combat.AttackInfo{
		ActorIndex:     c.Index,
		Abil:           "The Named Moment",
		AttackTag:      attacks.AttackTagElementalArt,
		ICDTag:         attacks.ICDTagNone,
		AdditionalTags: []attacks.AdditionalTag{attacks.AdditionalTagNightsoul},
		ICDGroup:       attacks.ICDGroupDefault,
		StrikeType:     attacks.StrikeTypeDefault,
		Element:        attributes.Pyro,
		Durability:     25,
		Mult:           skill[c.TalentLvlSkill()],
	}
	ap := combat.NewCircleHitOnTarget(
		c.Core.Combat.Player(),
		geometry.Point{Y: 0.5},
		5,
	)
	c.Core.QueueAttack(ai, ap, skillHitmark, skillHitmark, c.particleCB)
	c.exitBike()
	c.SetCDWithDelay(action.ActionSkill, 15*60, 18)

	return c.getSkillCastActionInfo(skillFrames)
}

// Recasting E while on bike, occurs with Sac or Burst allowing E to come off of cd
func (c *char) skillBikeRefresh() action.Info {
	ai := combat.AttackInfo{
		ActorIndex:     c.Index,
		Abil:           "The Named Moment (Flamestrider)",
		AttackTag:      attacks.AttackTagElementalArt,
		ICDTag:         attacks.ICDTagNone,
		AdditionalTags: []attacks.AdditionalTag{attacks.AdditionalTagNightsoul},
		ICDGroup:       attacks.ICDGroupDefault,
		StrikeType:     attacks.StrikeTypeBlunt,
		PoiseDMG:       75,
		Element:        attributes.Pyro,
		Durability:     25,
		Mult:           skill[c.TalentLvlSkill()],
		HitlagFactor:   0.05,
	}
	ap := combat.NewCircleHitOnTarget(
		c.Core.Combat.Player(),
		geometry.Point{Y: 1.0},
		6,
	)
	c.Core.QueueAttack(ai, ap, skillHitmark, skillHitmark, c.particleCB)
	c.SetCDWithDelay(action.ActionSkill, 15*60, 18)

	return c.getSkillCastActionInfo(skillBikeRefreshFrames)
}

// Recast can occur earlier out of Plunge, this extends a normal skill use to match total frames
func (c *char) getSkillCastActionInfo(f []int) action.Info {
	plungeFrames := 0

	// If using skill out of plunge, extend animation for non-recast skill
	if c.Core.Player.CurrentState() == action.PlungeAttackState {
		switch c.armamentState {
		case bike:
			plungeFrames = 19
		default:
			plungeFrames = 14
		}
	}

	return action.Info{
		Frames:          func(next action.Action) int { return f[next] + plungeFrames },
		AnimationLength: f[action.InvalidAction] + plungeFrames,
		CanQueueAfter:   f[action.ActionSwap] + plungeFrames,
		State:           action.SkillState,
	}
}

func (c *char) skillRingTask(src int) {
	c.QueueCharTask(func() {
		if c.ringSrc != src {
			return
		}
		if c.armamentState != ring {
			return
		}
		if !c.nightsoulState.HasBlessing() {
			return
		}
		ai := combat.AttackInfo{
			ActorIndex:     c.Index,
			Abil:           "Rings of Searing Radiance",
			AttackTag:      attacks.AttackTagElementalArt,
			ICDTag:         attacks.ICDTagNone,
			AdditionalTags: []attacks.AdditionalTag{attacks.AdditionalTagNightsoul},
			ICDGroup:       attacks.ICDGroupDefault,
			StrikeType:     attacks.StrikeTypeDefault,
			Element:        attributes.Pyro,
			Durability:     25,
			Mult:           skillRing[c.TalentLvlSkill()],
		}
		ap := combat.NewCircleHitOnTarget(
			c.Core.Combat.Player(),
			geometry.Point{Y: 1.0},
			6,
		)
		c.Core.QueueAttack(ai, ap, 0, 0, c.c6RingCB())
		c.reduceNightsoulPoints(3)
		c.skillRingTask(src)
	}, 2*60)
}

func (c *char) particleCB(a combat.AttackCB) {
	if a.Target.Type() != targets.TargettableEnemy {
		return
	}
	if c.StatusIsActive(particleICDKey) {
		return
	}
	c.AddStatus(particleICDKey, 0.5*60, true)
	c.Core.QueueParticle(c.Base.Key.String(), 5, attributes.Pyro, c.ParticleDelay)
}

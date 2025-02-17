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
)

const (
	skillHitmark     = 16
	particleICDKey   = "mavuika-particle-icd"
	skillRecastCD    = 60
	skillRecastCDKey = "mavuika-skill-recast-cd"
)

func init() {
	skillFrames = frames.InitAbilSlice(29) // E -> Dash/Jump
	skillFrames[action.ActionAttack] = 18  // E -> N1
	skillFrames[action.ActionCharge] = 18  // E -> CA
	skillFrames[action.ActionSkill] = 18   // E -> Skill Recast
	skillFrames[action.ActionBurst] = 18   // E -> Burst
	skillFrames[action.ActionWalk] = 28    // E -> Walk
	skillFrames[action.ActionSwap] = 24    // E -> Swap

	skillFramesHold = frames.InitAbilSlice(44) // E -> N1
	skillFramesHold[action.ActionSwap] = 34    // E -> Swap

	skillRecastFramesToBike = frames.InitAbilSlice(24) // E -> Swap
	skillRecastFramesToBike[action.ActionAttack] = 12  // E -> N1
	skillRecastFramesToBike[action.ActionCharge] = 12  // E -> CA

	skillRecastFramesToRing = frames.InitAbilSlice(27) // E -> N1
	skillRecastFramesToBike[action.ActionSwap] = 24    // E -> N1

}

func (c *char) nightsoulPointReduceFunc(src int) func() {
	return func() {
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
		c.Core.Tasks.Add(c.nightsoulPointReduceFunc(src), 6)
	}
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
		c.Core.Tasks.Add(c.nightsoulPointReduceFunc(c.nightsoulSrc), 6)
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

	if c.nightsoulState.HasBlessing() {
		c.QueueCharTask(c.skillRing(c.ringSrc), 120)
		c.c2Ring()
	}
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

	return action.Info{
		Frames:          frames.NewAbilFunc(skillFramesHold),
		AnimationLength: skillFramesHold[action.InvalidAction],
		CanQueueAfter:   skillFramesHold[action.ActionSwap],
		State:           action.SkillState,
	}
}

func (c *char) skillPress() action.Info {
	ai := combat.AttackInfo{
		ActorIndex:     c.Index,
		Abil:           "The Named Moment",
		AttackTag:      attacks.AttackTagElementalArt,
		ICDTag:         attacks.ICDTagNone,
		AdditionalTags: []attacks.AdditionalTag{attacks.AdditionalTagNightsoul},
		ICDGroup:       attacks.ICDGroupDefault,
		StrikeType:     attacks.StrikeTypePierce,
		Element:        attributes.Pyro,
		Durability:     25,
		Mult:           skill[c.TalentLvlSkill()],
	}
	ap := combat.NewCircleHitOnTarget(
		c.Core.Combat.Player(),
		geometry.Point{Y: 1.0},
		6,
	)
	c.Core.QueueAttack(ai, ap, skillHitmark, skillHitmark, c.particleCB)
	c.exitBike()
	c.SetCDWithDelay(action.ActionSkill, 15*60, 18)

	return action.Info{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionSwap],
		State:           action.SkillState,
	}
}

// Tried redirecting skillPress to skillHold, but it ran into errors
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
	}
	ap := combat.NewCircleHitOnTarget(
		c.Core.Combat.Player(),
		geometry.Point{Y: 1.0},
		6,
	)
	c.Core.QueueAttack(ai, ap, skillHitmark, skillHitmark, c.particleCB)
	c.SetCDWithDelay(action.ActionSkill, 15*60, 18)

	return action.Info{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionSwap],
		State:           action.SkillState,
	}
}

func (c *char) skillRing(src int) func() {
	return func() {
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
			StrikeType:     attacks.StrikeTypePierce,
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
		c.QueueCharTask(c.skillRing(src), 120)
	}
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

package pyro

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/targets"
)

var skillFrames [][]int
var skillHoldFrames [][]int

const (
	skillTapHitmark                = 26
	skillHoldHitmark               = 52
	blazingThresholdInterval       = 65
	scorchingThresholdHitmarkDelay = 12
	tapCdStart                     = 25
	holdCdStart                    = 47
	particleICDKey                 = "travelerpyro-particle-icd"
	scoringThresholdKey            = "travelerpyro-e"
)

func init() {
	skillFrames = make([][]int, 2)
	skillHoldFrames = make([][]int, 2)

	// Male
	// Tap
	skillFrames[0] = frames.InitAbilSlice(50)
	// Hold
	skillHoldFrames[0] = frames.InitAbilSlice(58) // E -> Ungray + 2

	// Female
	// Tap
	skillFrames[1] = frames.InitAbilSlice(50)
	// Hold
	skillHoldFrames[1] = frames.InitAbilSlice(58) // E -> Ungray + 2
}

func (c *Traveler) Skill(p map[string]int) (action.Info, error) {
	hold, ok := p["hold"]
	if !ok {
		hold = 0
	}
	switch hold {
	case 0:
	case 1:
	default:
		return action.Info{}, fmt.Errorf("invalid hold param supplied, got %v", hold)
	}

	// Enter Nightsoul and start reducing Points
	if c.nightsoulState.HasBlessing() {
		c.nightsoulState.GeneratePoints(42)
	} else {
		c.nightsoulState.EnterBlessing(c.nightsoulState.Points() + 42)
	}
	c.nightsoulSrc = c.Core.F
	c.QueueCharTask(c.nightsoulPointReduceFunc(c.Core.F), 10)
	c.c1AddMod()
	c.c2()
	c.c6AddMod()

	if hold == 0 {
		return c.SkillTap(p)
	}
	return c.SkillHold(p)
}

func (c *Traveler) SkillTap(p map[string]int) (action.Info, error) {
	c.SetCDWithDelay(action.ActionSkill, 18*60, tapCdStart)
	c.QueueCharTask(c.blazingThresholdHit(c.Core.F), skillTapHitmark)
	c.DeleteStatus(scoringThresholdKey)
	return action.Info{
		Frames:          frames.NewAbilFunc(skillFrames[c.gender]),
		AnimationLength: skillFrames[c.gender][action.InvalidAction],
		CanQueueAfter:   skillFrames[c.gender][action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}, nil
}

func (c *Traveler) SkillHold(p map[string]int) (action.Info, error) {
	c.SetCDWithDelay(action.ActionSkill, 18*60, holdCdStart)
	ai := combat.AttackInfo{
		ActorIndex:     c.Index,
		Abil:           "Flowfire Blade (Hold DMG)",
		AttackTag:      attacks.AttackTagElementalArt,
		AdditionalTags: []attacks.AdditionalTag{attacks.AdditionalTagNightsoul},
		ICDTag:         attacks.ICDTagTravelerHoldDMG,
		ICDGroup:       attacks.ICDGroupDefault,
		StrikeType:     attacks.StrikeTypeDefault,
		Element:        attributes.Pyro,
		Durability:     25,
		Mult:           holdDMG[c.TalentLvlSkill()],
	}

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 6.5),
		skillHoldHitmark,
		skillHoldHitmark,
	)

	c.AddStatus(scoringThresholdKey, -1, false)

	return action.Info{
		Frames:          frames.NewAbilFunc(skillHoldFrames[c.gender]),
		AnimationLength: skillHoldFrames[c.gender][action.InvalidAction],
		CanQueueAfter:   skillHoldFrames[c.gender][action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}, nil
}

func (c *Traveler) nightsoulPointReduceFunc(src int) func() {
	return func() {
		if c.nightsoulSrc != src {
			return
		}
		if !c.nightsoulState.HasBlessing() {
			return
		}
		val := 1.
		c.reduceNightsoulPoints(src, val)
		c.QueueCharTask(c.nightsoulPointReduceFunc(src), 10)
	}
}

func (c *Traveler) reduceNightsoulPoints(src int, val float64) {
	c.nightsoulState.ConsumePoints(val)
	if c.nightsoulState.Points() <= 0.00001 || c.Core.F >= src+12*60 {
		if !c.nightsoulState.HasBlessing() {
			return
		}
		c.DeleteStatus(scoringThresholdKey)
		for _, char := range c.Core.Player.Chars() {
			char.DeleteAttackMod(c1AttackModKey)
		}
		c.DeleteAttackMod(c6AttackModKey)
		c.nightsoulState.ExitBlessing()
		c.nightsoulState.ClearPoints()
		c.nightsoulSrc = -1
	}
}

func (c *Traveler) blazingThresholdHit(src int) func() {
	return func() {
		if src != c.nightsoulSrc {
			return
		}
		if !c.nightsoulState.HasBlessing() {
			return
		}
		ai := combat.AttackInfo{
			ActorIndex:     c.Index,
			Abil:           "Blazing Threshold DMG",
			AttackTag:      attacks.AttackTagElementalArt,
			AdditionalTags: []attacks.AdditionalTag{attacks.AdditionalTagNightsoul},
			ICDTag:         attacks.ICDTagTravelerBlazingThreshold,
			ICDGroup:       attacks.ICDGroupDefault,
			StrikeType:     attacks.StrikeTypeDefault,
			Element:        attributes.Pyro,
			Durability:     25,
			Mult:           blazingThreshold[c.TalentLvlSkill()],
		}
		// TODO: change hurtbox
		radius := 4.
		if c.Base.Ascension >= 1 && c.nightsoulState.Points() >= 20 {
			radius = 6.
		}
		c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, radius), 0, 0, c.particleCB)
		c.QueueCharTask(c.blazingThresholdHit(src), blazingThresholdInterval)
	}
}

func (c *Traveler) scorchingThresholdOnDamage() {
	c.Core.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		ae := args[1].(*combat.AttackEvent)
		dmg := args[2].(float64)
		if c.scorchingThresholdICD > c.Core.F {
			return false
		}
		if !c.StatusIsActive(scoringThresholdKey) {
			return false
		}
		// ignore EC, hydro swirl, and burning damage
		// this clause is here since these damage types are sourced to the target rather than character
		if ae.Info.AttackTag == attacks.AttackTagECDamage || ae.Info.AttackTag == attacks.AttackTagBurningDamage ||
			ae.Info.AttackTag == attacks.AttackTagSwirlHydro {
			return false
		}
		// ignore self dmg
		if ae.Info.ActorIndex == c.Index &&
			ae.Info.AttackTag == attacks.AttackTagElementalArt &&
			ae.Info.StrikeType == attacks.StrikeTypeSlash {
			return false
		}
		// ignore 0 damage
		if dmg == 0 {
			return false
		}

		ai := combat.AttackInfo{
			ActorIndex:     c.Index,
			Abil:           "Scorching Threshold DMG",
			AttackTag:      attacks.AttackTagElementalArt,
			AdditionalTags: []attacks.AdditionalTag{attacks.AdditionalTagNightsoul},
			ICDTag:         attacks.ICDTagTravelerScorchingThreshold,
			ICDGroup:       attacks.ICDGroupDefault,
			StrikeType:     attacks.StrikeTypeDefault,
			Element:        attributes.Pyro,
			Durability:     25,
			Mult:           scorchingThreshold[c.TalentLvlSkill()],
		}
		// TODO: change hurtbox
		radius := 4.
		if c.Base.Ascension >= 1 && c.nightsoulState.Points() >= 20 {
			radius = 6.
		}
		c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, radius),
			scorchingThresholdHitmarkDelay, scorchingThresholdHitmarkDelay, c.particleCB)

		c.scorchingThresholdICD = c.Core.F + 180 // 3 sec icd
		return false
	}, "travelerpyro-scorching-threshold")
}

func (c *Traveler) particleCB(a combat.AttackCB) {
	if a.Target.Type() != targets.TargettableEnemy {
		return
	}
	if c.StatusIsActive(particleICDKey) {
		return
	}
	c.AddStatus(particleICDKey, int(2.9*60), true)

	count := 1.0
	c.Core.QueueParticle(c.Base.Key.String(), count, attributes.Pyro, c.ParticleDelay)
}

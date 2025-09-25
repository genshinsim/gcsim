package pyro

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/enemy"
)

var (
	skillFrames     [][]int
	skillHoldFrames [][]int
)

const (
	skillTapHitmark                = 24
	skillHoldHitmark               = 50
	blazingThresholdInterval       = 66
	scorchingThresholdHitmarkDelay = 12
	tapCdStart                     = 19
	holdCdStart                    = 48
	enterNightsoulDelay            = 19
	nightsoulReduceDelay           = 8 // From wiki: consumption is 7.5 points per second -> 1 per 8f
	scorchingThresholdICD          = 180
	particleICDKey                 = "travelerpyro-particle-icd"
	scoringThresholdKey            = "travelerpyro-e"
	scorchingThresholdICDKey       = "travelerpyro-scorching-threshold-icd"
)

func init() {
	skillFrames = make([][]int, 2)
	skillHoldFrames = make([][]int, 2)

	// Male
	// Tap
	skillFrames[0] = frames.InitAbilSlice(49) // E -> N1
	skillFrames[0][action.ActionDash] = 31
	skillFrames[0][action.ActionJump] = 31
	skillFrames[0][action.ActionSwap] = 48
	// Hold
	skillHoldFrames[0] = frames.InitAbilSlice(77) // E -> E
	skillHoldFrames[0][action.ActionAttack] = 61
	skillHoldFrames[0][action.ActionBurst] = 59
	skillHoldFrames[0][action.ActionDash] = 60
	skillHoldFrames[0][action.ActionJump] = 61
	skillHoldFrames[0][action.ActionWalk] = 81
	skillHoldFrames[0][action.ActionSwap] = 59

	// Female
	// Tap
	skillFrames[1] = frames.InitAbilSlice(49) // E -> N1
	skillFrames[1][action.ActionDash] = 31
	skillFrames[1][action.ActionJump] = 31
	skillFrames[1][action.ActionSwap] = 48
	// Hold
	skillHoldFrames[1] = frames.InitAbilSlice(77) // E -> E
	skillHoldFrames[1][action.ActionAttack] = 61
	skillHoldFrames[1][action.ActionBurst] = 59
	skillHoldFrames[1][action.ActionDash] = 60
	skillHoldFrames[1][action.ActionJump] = 61
	skillHoldFrames[1][action.ActionWalk] = 81
	skillHoldFrames[1][action.ActionSwap] = 59
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

	c.c2OnSkill()

	if hold == 0 {
		return c.SkillTap(p)
	}
	return c.SkillHold(p)
}

func (c *Traveler) SkillTap(p map[string]int) (action.Info, error) {
	// Enter Nightsoul and start reducing Points
	skillSrc := c.Core.F + enterNightsoulDelay
	c.QueueCharTask(func() {
		c.enterNightsoul(skillSrc)
	}, enterNightsoulDelay)
	c.SetCDWithDelay(action.ActionSkill, 18*60, tapCdStart)
	c.QueueCharTask(c.blazingThresholdHit(skillSrc), skillTapHitmark)
	c.DeleteStatus(scoringThresholdKey)
	return action.Info{
		Frames:          frames.NewAbilFunc(skillFrames[c.gender]),
		AnimationLength: skillFrames[c.gender][action.InvalidAction],
		CanQueueAfter:   skillFrames[c.gender][action.ActionSwap], // earliest cancel
		State:           action.SkillState,
	}, nil
}

func (c *Traveler) SkillHold(p map[string]int) (action.Info, error) {
	c.QueueCharTask(func() {
		c.enterNightsoul(c.Core.F)
	}, 48)
	c.SetCDWithDelay(action.ActionSkill, 18*60, holdCdStart)
	ai := info.AttackInfo{
		ActorIndex:     c.Index(),
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
		combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 3.0),
		skillHoldHitmark,
		skillHoldHitmark,
	)

	c.AddStatus(scoringThresholdKey, -1, false)

	return action.Info{
		Frames:          frames.NewAbilFunc(skillHoldFrames[c.gender]),
		AnimationLength: skillHoldFrames[c.gender][action.InvalidAction],
		CanQueueAfter:   skillHoldFrames[c.gender][action.ActionSwap], // earliest cancel
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
		c.reduceNightsoulPoints(val)
		c.QueueCharTask(c.nightsoulPointReduceFunc(src), nightsoulReduceDelay)
	}
}

func (c *Traveler) reduceNightsoulPoints(val float64) {
	c.nightsoulState.ConsumePoints(val)
	if c.nightsoulState.Points() <= 0.00001 && c.nightsoulState.HasBlessing() {
		c.exitNightsoul()
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
		ai := info.AttackInfo{
			ActorIndex:     c.Index(),
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
		radius := 0.5
		if c.Base.Ascension >= 1 && c.nightsoulState.Points() >= 20 {
			radius = 3.
		}
		c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, radius), 0, 0, c.particleCB)
		c.QueueCharTask(c.blazingThresholdHit(src), blazingThresholdInterval)
	}
}

func (c *Traveler) scorchingThresholdOnDamage() {
	c.Core.Events.Subscribe(event.OnEnemyDamage, func(args ...any) bool {
		_, ok := args[0].(*enemy.Enemy)
		ae := args[1].(*info.AttackEvent)
		dmg := args[2].(float64)
		if !ok {
			return false
		}
		if c.StatusIsActive(scorchingThresholdICDKey) {
			return false
		}
		if !c.StatusIsActive(scoringThresholdKey) {
			return false
		}
		// ignore burning damage
		if ae.Info.AttackTag == attacks.AttackTagBurningDamage ||
			ae.Info.AttackTag == attacks.AttackTagSwirlHydro {
			return false
		}
		// ignore 0 damage
		if dmg == 0 {
			return false
		}

		ai := info.AttackInfo{
			ActorIndex:     c.Index(),
			Abil:           "Scorching Threshold",
			AttackTag:      attacks.AttackTagElementalArt,
			AdditionalTags: []attacks.AdditionalTag{attacks.AdditionalTagNightsoul},
			ICDTag:         attacks.ICDTagTravelerScorchingThreshold,
			ICDGroup:       attacks.ICDGroupDefault,
			StrikeType:     attacks.StrikeTypeDefault,
			Element:        attributes.Pyro,
			Durability:     25,
			Mult:           scorchingThreshold[c.TalentLvlSkill()],
		}
		radius := 1.5
		if c.Base.Ascension >= 1 && c.nightsoulState.Points() >= 20 {
			radius = 4.
		}
		c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, radius),
			scorchingThresholdHitmarkDelay, scorchingThresholdHitmarkDelay, c.particleCB)

		c.AddStatus(scorchingThresholdICDKey, scorchingThresholdICD, true)
		return false
	}, "travelerpyro-scorching-threshold")
}

func (c *Traveler) particleCB(a info.AttackCB) {
	if a.Target.Type() != info.TargettableEnemy {
		return
	}
	if c.StatusIsActive(particleICDKey) {
		return
	}
	c.AddStatus(particleICDKey, int(2.9*60), true)

	count := 1.0
	c.Core.QueueParticle(c.Base.Key.String(), count, attributes.Pyro, c.ParticleDelay)
}

func (c *Traveler) enterNightsoul(src int) {
	points := 42.0
	if !c.nightsoulState.HasBlessing() {
		points += c.nightsoulState.Points()
	}

	c.nightsoulSrc = src
	c.nightsoulState.EnterTimedBlessing(points, 12*60, c.exitNightsoul)
	c.QueueCharTask(c.nightsoulPointReduceFunc(c.nightsoulSrc), nightsoulReduceDelay)
}

func (c *Traveler) exitNightsoul() {
	c.DeleteStatus(scoringThresholdKey)
	c.nightsoulState.ExitBlessing()
	c.nightsoulState.ClearPoints()
	c.nightsoulSrc = -1
}

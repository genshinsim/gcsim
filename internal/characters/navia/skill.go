package navia

import (
	"fmt"
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/targets"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
	"math"
)

var skillPressFrames []int
var skillHoldFrames []int
var crystallise = []event.Event{event.OnCrystallizeElectro, event.OnCrystallizeCryo, event.OnCrystallizeHydro,
	event.OnCrystallizePyro}
var skillMultiplier = []float64{0, 1, 1.05, 1.1, 1.15, 1.2, 1.36, 1.4, 1.6, 1.666, 1.9, 2}

const (
	skillPressCDStart = 30
	skillPressHitmark = 30
	skillHoldHitmark  = 48
	arkheDelay        = 12
	particleICDKey    = "navia-particle-icd"
	arkheICDKey       = "navia-arkhe-icd"
)

func init() {
	skillPressFrames = frames.InitAbilSlice(38) // E -> N1/Q
	skillPressFrames[action.ActionDash] = 38
	skillPressFrames[action.ActionJump] = 38
	skillPressFrames[action.ActionSwap] = 38

	// skill (hold) -> x
	skillHoldFrames = frames.InitAbilSlice(51) // E -> Swap
	skillHoldFrames[action.ActionAttack] = 51
	skillHoldFrames[action.ActionBurst] = 51
	skillHoldFrames[action.ActionDash] = 51
	skillHoldFrames[action.ActionJump] = 51
}

func (c *char) Skill(p map[string]int) (action.Info, error) {
	c.c2ready = true
	hold := p["hold"]
	hitmark := 0
	if hold > 0 {
		hitmark += skillHoldHitmark + hold
		c.SetCDWithDelay(action.ActionSkill, 9*60, hitmark)
	} else {
		hitmark += skillPressHitmark + hold
		c.SetCDWithDelay(action.ActionSkill, 9*60, skillPressCDStart)
	}

	if p["shrapnel"] != 0 {
		c.shrapnel = int(math.Min(float64(p["shrapnel"]), 6))
	}

	c.Core.Log.NewEvent(fmt.Sprintf("%v crystal shrapnel", c.shrapnel), glog.LogCharacterEvent, c.Index)

	shots := 5 + int(math.Max(float64(c.shrapnel-3), 0))*2

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Rosula Shardshot",
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeBlunt,
		Element:    attributes.Geo,
		Durability: 25,
		Mult:       skillshotgun[c.TalentLvlSkill()],
	}

	excess := math.Max(float64(c.shrapnel-3), 0)
	m := make([]float64, attributes.EndStatType)
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("navia-skill-dmgup", hitmark+6),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			if atk.Info.AttackTag != attacks.AttackTagElementalArt {
				return nil, false
			}
			m[attributes.DmgP] = 0.15 * excess
			if c.Base.Cons >= 2 {
				m[attributes.CR] = 0.12 * excess
			}
			if c.Base.Cons >= 6 {
				m[attributes.CD] = 0.45 * excess
			}
			return m, true
		},
	})

	// Looks for enemies in the path of each bullet
	// Initially trims enemies to check by scanning only the hit zone
	for _, t := range c.Core.Combat.EnemiesWithinArea(
		combat.NewCircleHitOnTargetFanAngle(c.Core.Combat.Player(), geometry.Point{Y: 0}, 25, 15),
		nil,
	) {
		// Tallies up the hits
		hits := 0
		for i := 0; i < shots; i++ {
			if ok, _ := t.AttackWillLand(combat.NewCircleHitOnTargetFanAngle(c.Core.Combat.Player(),
				geometry.Point{Y: 0},
				25, 15)); ok {
				hits++
			}
		}
		// Applies damage based on the hits
		ai.Mult = skillshotgun[c.TalentLvlSkill()] * skillMultiplier[hits]
		c.Core.QueueAttack(
			ai,
			combat.NewSingleTargetHit(t.Key()),
			hitmark,
			hitmark+5,
			c.SkillCB(hitmark),
		)
	}

	// remove the shrapnel after firing and action C1
	c.QueueCharTask(
		func() {
			c.c1(c.shrapnel)
			if c.Base.Cons < 6 {
				c.shrapnel = 0
			} else {
				c.shrapnel = c.shrapnel - 3 // C6 keeps any more than the three
			}
			return
		},
		hitmark,
	)
	// Add Geo Infusion
	c.QueueCharTask(
		c.a1,
		hitmark+30,
	)
	c.c2ready = false
	if hold > 1 {
		return action.Info{
			Frames:          frames.NewAbilFunc(skillHoldFrames),
			AnimationLength: skillHoldFrames[action.ActionDash] + hold,
			CanQueueAfter:   skillHoldFrames[action.ActionDash] + hold,
			State:           action.SkillState,
		}, nil
	}
	return action.Info{
		Frames:          frames.NewAbilFunc(skillPressFrames),
		AnimationLength: skillPressFrames[action.ActionDash],
		CanQueueAfter:   skillPressFrames[action.ActionDash],
		State:           action.SkillState,
	}, nil

}

func (c *char) SkillCB(hitmark int) combat.AttackCBFunc {
	done := false
	return func(a combat.AttackCB) {
		e := a.Target.(*enemy.Enemy)
		if e.Type() != targets.TargettableEnemy {
			return
		}

		// When firing, attack with the Surging Blade
		c.SurgingBlade(hitmark)

		if done {
			return
		}
		c.c2(a)
		c.c2ready = false
		if !c.StatusIsActive(particleICDKey) {
			count := 3.0
			if c.Core.Rand.Float64() < 0.5 {
				count = 4
			}
			c.Core.QueueParticle(c.Base.Key.String(), count, attributes.Geo, c.ParticleDelay)
			c.AddStatus(particleICDKey, 0.2*60, true)
		}

		done = true

	}
}

// ShrapnelGain adds Shrapnel Stacks when crystallise occurs. Stacks should last 300s but this is way to long to bother
// When a character in the party obtains an Elemental Shard created from the Crystallize reaction,
// Navia will gain 1 Crystal Shrapnel charge. Navia can hold up to 6 charges of Crystal Shrapnel at once.
// Each time Crystal Shrapnel gain is triggered, the duration of the Shards you have already will be reset.
func (c *char) ShrapnelGain() {

	for _, crystal := range crystallise {
		c.Core.Events.Subscribe(crystal, func(args ...interface{}) bool {
			if c.shrapnel < 6 {
				c.shrapnel++
				c.Core.Log.NewEvent("Crystal Shrapnel gained from Crystallise", glog.LogCharacterEvent, c.Index)
			}
			return false
		}, "shrapnel-gain")
	}

}

func (c *char) SurgingBlade(hitmark int) {
	if c.StatusIsActive(arkheICDKey) {
		return
	}
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Surging Blade",
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeBlunt,
		Element:    attributes.Geo,
		Durability: 0,
		Mult:       skillblade[c.TalentLvlSkill()],
	}
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), geometry.Point{Y: 0}, 3),
		arkheDelay+hitmark,
		arkheDelay+hitmark,
		nil,
	)
	c.QueueCharTask(
		func() {
			c.AddStatus(arkheICDKey, 7*60, true)
		},
		hitmark,
	)

	return
}

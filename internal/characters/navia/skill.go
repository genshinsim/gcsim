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
	"github.com/genshinsim/gcsim/pkg/core/player/shield"
	"github.com/genshinsim/gcsim/pkg/core/targets"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
	"github.com/genshinsim/gcsim/pkg/reactable"
	"math"
)

var skillPressFrames []int
var skillPowerPressFrames []int
var skillHoldFrames []int
var skillPowerHoldFrames []int

var skillMultiplier = []float64{
	0,                   // 0 hits
	1,                   // 1 hit
	1.05000000074505806, // 2 hit
	1.10000000149011612, // 3 hit etc
	1.15000000596046448,
	1.20000000298023224,
	1.36000001430511475,
	1.4000000059604645,
	1.6000000238418579,
	1.6660000085830688,
	1.8999999761581421,
	2,
}

var hitscans = [][]float64{
	{0.9, 0, 0, 0},
	{0.25, 7.911, 0, 0},
	{0.25, -1.826, 0, 0},
	{0.25, -4.325, 0, 0},
	{0.25, 0.773, 0, 0},
	{0.25, 6.209, 0, 0},
	{0.25, -2.752, 0, 0},
	{0.25, 7.845, 0.01, 0.01},
	{0.25, -7.933, -0.01, -0.01},
	{0.25, 2.626, 0, 0},
	{0.25, -5.43724, 0, 0},
	//width, angle, x offset, y offset
}

const (
	skillPressCDStart = 11

	travelDelay = 9

	arkheDelay = 65

	skillHoldCDStart  = 41
	skillHoldDuration = 241 //an additional 1f to account for hold being set to 1 to activate

	particleICDKey = "navia-particle-icd"
	arkheICDKey    = "navia-arkhe-icd"
)

func init() {
	skillPressFrames = frames.InitAbilSlice(40) // E -> E/Q
	skillPressFrames[action.ActionDash] = 24
	skillPressFrames[action.ActionJump] = 24
	skillPressFrames[action.ActionSwap] = 38
	skillPressFrames[action.ActionWalk] = 35
	skillPressFrames[action.ActionAttack] = 38

	// skill with >=3 shrapnel -> x
	skillPowerPressFrames = frames.InitAbilSlice(41) // E -> E/Q
	skillPowerPressFrames[action.ActionDash] = 26
	skillPowerPressFrames[action.ActionJump] = 24
	skillPowerPressFrames[action.ActionSwap] = 40
	skillPowerPressFrames[action.ActionWalk] = 39
	skillPowerPressFrames[action.ActionAttack] = 40

	// skill (hold) -> x
	skillHoldFrames = frames.InitAbilSlice(71) // E -> E/Q
	skillHoldFrames[action.ActionDash] = 54
	skillHoldFrames[action.ActionJump] = 54
	skillHoldFrames[action.ActionSwap] = 69
	skillHoldFrames[action.ActionWalk] = 65
	skillHoldFrames[action.ActionAttack] = 69

	// skill (hold) with >=3 shrapnel -> x
	skillPowerHoldFrames = frames.InitAbilSlice(73) // E -> E/Q
	skillPowerHoldFrames[action.ActionDash] = 56
	skillPowerHoldFrames[action.ActionJump] = 56
	skillPowerHoldFrames[action.ActionSwap] = 71
	skillPowerHoldFrames[action.ActionWalk] = 70
	skillPowerHoldFrames[action.ActionAttack] = 72
}

func (c *char) Skill(p map[string]int) (action.Info, error) {
	c.c2ready = true
	shrapnel := 0
	hold := p["hold"]
	firingTime := 0
	if hold > 0 {
		if hold > skillHoldDuration {
			hold = skillHoldDuration - 1
		}
		firingTime += skillHoldCDStart + hold - 1

		// At 0.25, tap is converted to hold, and the suction begins
		// TODO: Confirm suction begins at 15f
		for i := 15; i < firingTime; i += 30 {
			c.PullCrystals(firingTime, i)
		}
		c.SetCDWithDelay(action.ActionSkill, 9*60, skillHoldCDStart+hold)
	} else {
		firingTime += skillPressCDStart
		c.SetCDWithDelay(action.ActionSkill, 9*60, skillPressCDStart)
	}
	shots := 5
	
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

	c.QueueCharTask(
		func() {
			if p["shrapnel"] != 0 {
				c.shrapnel = int(math.Min(float64(p["shrapnel"]), 6))
			}
			c.Core.Log.NewEvent(fmt.Sprintf("%v crystal shrapnel", c.shrapnel), glog.LogCharacterEvent, c.Index)
			shots = 5 + int(math.Min(float64(c.shrapnel), 3))*2
			shrapnel = c.shrapnel
			// Calculate buffs based on excess shrapnel
			excess := math.Max(float64(c.shrapnel-3), 0)
			m := make([]float64, attributes.EndStatType)
			c.AddAttackMod(character.AttackMod{
				Base: modifier.NewBase("navia-skill-dmgup", travelDelay+1),
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
				combat.NewCircleHitOnTargetFanAngle(c.Core.Combat.Player(), geometry.Point{Y: 0}, 11.5, 15),
				nil,
			) {
				// Tallies up the hits
				hits := 0
				for i := 0; i < shots; i++ {
					if ok, _ := t.AttackWillLand(combat.NewBoxHitOnTarget(c.Core.Combat.Player(), geometry.Point{X: hitscans[i][2],
						Y: hitscans[i][3]}.Rotate(geometry.DegreesToDirection(hitscans[i][1])),
						hitscans[i][0], 11.5)); ok {
						hits++
					}
				}
				// Applies damage based on the hits
				ai.Mult = skillshotgun[c.TalentLvlSkill()] * skillMultiplier[hits]
				c.Core.QueueAttack(
					ai,
					combat.NewSingleTargetHit(t.Key()),
					0,
					travelDelay,
					c.particleCB(),
					c.c2(),
				)
			}
			// remove the shrapnel after firing and action C1 and A1
			c.c1(c.shrapnel)
			if c.Base.Cons < 6 {
				c.shrapnel = 0
			} else {
				c.shrapnel -= 3 // C6 keeps any more than the three
			}
			c.a1()
		},
		firingTime,
	)

	c.Core.Tasks.Add(c.SurgingBlade, firingTime)

	c.c2ready = false
	if hold > 1 {
		if shrapnel >= 3 {
			return action.Info{
				Frames:          func(next action.Action) int { return hold + skillPowerHoldFrames[next] },
				AnimationLength: skillPowerHoldFrames[action.InvalidAction] + hold,
				CanQueueAfter:   skillPowerHoldFrames[action.ActionDash] + hold,
				State:           action.SkillState,
			}, nil
		}
		return action.Info{
			Frames:          func(next action.Action) int { return hold + skillHoldFrames[next] },
			AnimationLength: skillHoldFrames[action.InvalidAction] + hold,
			CanQueueAfter:   skillHoldFrames[action.ActionDash] + hold,
			State:           action.SkillState,
		}, nil
	}
	if shrapnel >= 3 {
		return action.Info{
			Frames:          frames.NewAbilFunc(skillPowerPressFrames),
			AnimationLength: skillPowerPressFrames[action.InvalidAction],
			CanQueueAfter:   skillPowerPressFrames[action.ActionDash],
			State:           action.SkillState,
		}, nil
	}
	return action.Info{
		Frames:          frames.NewAbilFunc(skillPressFrames),
		AnimationLength: skillPressFrames[action.InvalidAction],
		CanQueueAfter:   skillPressFrames[action.ActionDash],
		State:           action.SkillState,
	}, nil

}

func (c *char) particleCB() combat.AttackCBFunc {
	done := false
	return func(a combat.AttackCB) {
		e := a.Target.(*enemy.Enemy)
		if e.Type() != targets.TargettableEnemy {
			return
		}

		if done {
			return
		}
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

// ShrapnelGain adds Shrapnel Stacks when a Crystallise Shield is picked up.
// Stacks should last 300s but this is way to long to bother
// When a character in the party obtains an Elemental Shard created from the Crystallize reaction,
// Navia will gain 1 Crystal Shrapnel charge. Navia can hold up to 6 charges of Crystal Shrapnel at once.
// Each time Crystal Shrapnel gain is triggered, the duration of the Shards you have already will be reset.
func (c *char) ShrapnelGain() {
	c.Core.Events.Subscribe(event.OnShielded, func(args ...interface{}) bool {
		// Check shield
		shd := args[0].(shield.Shield)

		if shd.Type() != shield.Crystallize {
			return false
		}
		if c.shrapnel < 6 {
			c.shrapnel++
			c.Core.Log.NewEvent("Crystal Shrapnel gained from Crystallise", glog.LogCharacterEvent, c.Index)
		}
		return false
	}, "shrapnel-gain")
}

func (c *char) SurgingBlade() {
	if c.StatusIsActive(arkheICDKey) {
		return
	}
	c.AddStatus(arkheICDKey, 7*60, false)
	e := c.Core.Combat.ClosestEnemyWithinArea(combat.NewCircleHitFanAngle(c.Core.Combat.Player(),
		geometry.Point{Y: 0}.Rotate(geometry.Point{Y: 0}), geometry.Point{Y: 0}, 11.5, 7.933+7.911), nil)
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
	if e != nil {
		c.Core.QueueAttack(
			ai,
			combat.NewCircleHitOnTarget(e, geometry.Point{Y: 0}, 3),
			0,
			arkheDelay,
			nil,
		)
	} else {
		c.Core.QueueAttack(
			ai,
			combat.NewCircleHitOnTarget(c.Core.Combat.Player(), geometry.Point{Y: 3}, 3),
			0,
			0,
			nil,
		)
	}
	return
}

// PullCrystals to Navia. This has a range of 12m from datamine, so
// check every every 30f. 
func (c *char) PullCrystals(firingTime, i int) {
	for j, k := 0, 0; j < c.Core.Combat.GadgetCount(); j++ {
		cs, ok := c.Core.Combat.Gadget(j).(*reactable.CrystallizeShard)
		// skip if not a shard
		if !ok {
			continue
		}

		// If shard is out of 12m range, skip
		if !cs.IsWithinArea(combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 12)) {
			continue
		}

		// approximate sucking in as 0.4m per frame (~8m distance took 20f to arrive at gorou)
		distance := cs.Pos().Distance(c.Core.Combat.Player().Pos())
		travel := int(math.Ceil(distance / 0.4))
		// if the crystal won't arrive before the shot is fired, skip
		if firingTime-i < travel {
			continue
		}
		// special check to account for edge case if shard just spawned and will arrive before it can be picked up
		if c.Core.F+travel < cs.EarliestPickup {
			continue
		}
		c.Core.Tasks.Add(func() {
			cs.AddShieldKillShard()
		}, travel)
		// max three crystals
		k++
		if k >= 3 {
			break
		}
	}
}
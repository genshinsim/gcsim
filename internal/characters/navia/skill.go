package navia

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/targets"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"math"
)

var skillPressFrames []int
var skillHoldFrames []int
var crystallise = []event.Event{event.OnCrystallizeElectro, event.OnCrystallizeCryo, event.OnCrystallizeHydro,
	event.OnCrystallizePyro}

const (
	skillPressCDStart = 16
	skillPressHitmark = 17

	skillHoldCDStart = 11
	skillHoldHitmark = 12
)

func init() {
	skillPressFrames = frames.InitAbilSlice(39) // E -> N1/Q
	skillPressFrames[action.ActionDash] = 34
	skillPressFrames[action.ActionJump] = 35
	skillPressFrames[action.ActionSwap] = 37

	// skill (hold) -> x
	skillHoldFrames = frames.InitAbilSlice(46) // E -> Swap
	skillHoldFrames[action.ActionAttack] = 38
	skillHoldFrames[action.ActionBurst] = 37
	skillHoldFrames[action.ActionDash] = 30
	skillHoldFrames[action.ActionJump] = 30
}

func (c *char) Skill(p map[string]int) (action.Info, error) {
	hold := p["hold"]
	c.SurgingBlade()
	hitmark := 0
	if hold > 0 {
		hitmark += skillHoldHitmark + hold
	} else {
		hitmark += skillPressHitmark + hold
	}

	if p["shrapnel"] != 0 {
		c.shrapnel = int(math.Min(float64(p["shrapnel"]), 6))
	}
	travel := 5
	if p["travel"] != 0 {
		travel = p["travel"]
	}
	shots := 1.0
	switch c.shrapnel {
	case 0:
		shots = 1.2
	case 1:
		shots = 1.4
	case 2:
		shots = 1.66
	default:
		shots = 2.0
	}
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Rosula Shardshot",
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeBlunt,
		Element:    attributes.Geo,
		Durability: 25,
		Mult:       skillshotgun[c.TalentLvlSkill()] * shots,
	}
	// When firing, attack with the Surging Blade
	c.SurgingBlade(hitmark)

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 25),
		hitmark,
		travel+hitmark,
		c.SkillCB(),
	)
	// C1 Energy Restoration and CD reduction
	if c.Base.Cons >= 1 {
		c.QueueCharTask(func() {
			c.c1(c.shrapnel)
		}, hitmark)
	}
	// remove the shrapnel after firing
	c.QueueCharTask(
		func() {
			if c.Base.Cons < 6 {
				c.shrapnel = 0
			} else {
				c.shrapnel = c.shrapnel - 3 // C6 keeps any more than the three
			}
			return
		},
		hitmark,
	)

	return action.Info{}, nil
}

func (c *char) SkillCB() combat.AttackCBFunc {
	done := false
	return func(a combat.AttackCB) {
		e := a.Target.(*enemy.Enemy)
		if e.Type() != targets.TargettableEnemy {
			return
		}

		c.c4()

		if done {
			return
		}
		c.Core.QueueParticle(c.Base.Key.String(), 4, attributes.Geo, c.ParticleDelay)

		c.c2()

		done = true

	}
}

func (c *char) SurgingBlade(delay int) {
	if c.StatusIsActive("surging-blade-cd") {
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
		3+delay,
		3+delay,
		nil,
	)
	c.AddStatus("surging-blade-cd", 7*60, true)
	return
}

func (c *char) ShrapnelGain() {

	for _, crystal := range crystallise {
		c.Core.Events.Subscribe(crystal, func(args ...interface{}) bool {
			if c.shrapnel < 6 {
				c.shrapnel++
			}
			return false
		}, "shrapnel-gain")
	}

}

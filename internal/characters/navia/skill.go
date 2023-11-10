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

const (
	skillPressCDStart = 16
	skillPressHitmark = 17

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
				m[attributes.CR] = 0.08 * excess
			}
			if c.Base.Cons >= 6 {
				m[attributes.CD] = 0.35 * excess
			}
			return m, true
		},
	})

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 25),
		hitmark,
		5+hitmark,
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
	// Add Geo Infusion
	c.QueueCharTask(
		c.a1,
		hitmark+30,
	)

	return action.Info{
		Frames:          frames.NewAbilFunc(skillHoldFrames),
		AnimationLength: skillHoldFrames[action.ActionDash],
		CanQueueAfter:   skillHoldFrames[action.ActionDash],
		State:           action.SkillState,
	}, nil
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
		10+hitmark,
		10+hitmark,
		nil,
	)
	c.QueueCharTask(
		func() {
			c.AddStatus("surging-blade-cd", 7*60, true)
		},
		hitmark,
	)

	return
}

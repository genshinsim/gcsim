package aloy

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/targets"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

var skillFrames []int

const skillRelease = 20 // release frame for the Bomb, travel comes on top, bomb_delay comes on top afterwards

func init() {
	skillFrames = frames.InitAbilSlice(49) // E -> Dash
	skillFrames[action.ActionAttack] = 47  // E -> N1
	skillFrames[action.ActionBurst] = 48   // E -> Q
	skillFrames[action.ActionJump] = 47    // E -> J
	skillFrames[action.ActionSwap] = 66    // E -> Swap
}

const (
	rushingIceKey      = "rushingice"
	rushingIceDuration = 600
)

// Skill - Handles main damage, bomblet, and coil effects
//
// Has 3 parameters:
//
// - "travel" = Delay in frames until main damage, bomblets spawn on main damage
//
// - "bomblets" = Number of bomblets that hit
//
// - "bomb_delay" = Delay in frames before bomblets go off and coil stacks get added
//
// - too many potential bomblet hit variations to keep syntax short, so we simplify how they can be handled here
func (c *char) Skill(p map[string]int) action.ActionInfo {
	travel, ok := p["travel"]
	if !ok {
		travel = 5
	}

	bomblets, ok := p["bomblets"]
	if !ok {
		bomblets = 2
	}

	delay, ok := p["bomb_delay"]
	if !ok {
		delay = 0
	}

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Freeze Bomb",
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Cryo,
		Durability: 25,
		Mult:       skillMain[c.TalentLvlSkill()],
	}
	// TODO: accurate snapshot timing, assumes snapshot on release and not on hit/bomb creation
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHit(c.Core.Combat.Player(), c.Core.Combat.PrimaryTarget(), nil, 4),
		skillRelease,
		skillRelease+travel,
		c.makeParticleCB(),
	)

	// Bomblets snapshot on cast
	ai = combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               "Chillwater Bomblets",
		AttackTag:          attacks.AttackTagElementalArt,
		ICDTag:             attacks.ICDTagElementalArt,
		ICDGroup:           attacks.ICDGroupDefault,
		StrikeType:         attacks.StrikeTypeDefault,
		Element:            attributes.Cryo,
		Durability:         25,
		Mult:               skillBomblets[c.TalentLvlSkill()],
		CanBeDefenseHalted: true,
		IsDeployable:       true,
	}

	// Queue up bomblets
	for i := 0; i < bomblets; i++ {
		c.Core.QueueAttack(
			ai,
			combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 2),
			skillRelease+travel,
			skillRelease+travel+delay+((i+1)*6),
			c.coilStacks,
		)
	}

	c.SetCDWithDelay(action.ActionSkill, 20*60, 19)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillRelease,
		State:           action.SkillState,
	}
}

func (c *char) makeParticleCB() combat.AttackCBFunc {
	done := false
	return func(a combat.AttackCB) {
		if a.Target.Type() != targets.TargettableEnemy {
			return
		}
		if done {
			return
		}
		done = true
		c.Core.QueueParticle(c.Base.Key.String(), 5, attributes.Cryo, c.ParticleDelay)
	}
}

// Handles coil stacking and associated effects, including triggering rushing ice
func (c *char) coilStacks(a combat.AttackCB) {
	if a.Target.Type() != targets.TargettableEnemy {
		return
	}
	if c.coilICDExpiry > c.Core.F {
		return
	}
	// Can't gain coil stacks while in rushing ice
	if c.StatusIsActive(rushingIceKey) {
		return
	}
	c.coils++
	c.coilICDExpiry = c.Core.F + 6

	c.Core.Log.NewEvent("coil stack gained", glog.LogCharacterEvent, c.Index).
		Write("stacks", c.coils)

	c.a1()

	if c.coils == 4 {
		c.coils = 0
		c.rushingIce()
	}
}

// Handles rushing ice state
func (c *char) rushingIce() {
	c.AddStatus(rushingIceKey, rushingIceDuration, true)
	c.Core.Player.AddWeaponInfuse(c.Index, "aloy-rushing-ice", attributes.Cryo, 600, true, attacks.AttackTagNormal)

	// Rushing ice NA bonus
	val := make([]float64, attributes.EndStatType)
	val[attributes.DmgP] = skillRushingIceNABonus[c.TalentLvlSkill()]
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBaseWithHitlag("aloy-rushing-ice", rushingIceDuration),
		Amount: func(atk *combat.AttackEvent, _ combat.Target) ([]float64, bool) {
			if atk.Info.AttackTag == attacks.AttackTagNormal {
				return val, true
			}
			return nil, false
		},
	})

	c.a4()
}

// Add coil mod at the beginning of the sim
// Can't be made dynamic easily as coils last until 30s after when Aloy swaps off field
func (c *char) coilMod() {
	val := make([]float64, attributes.EndStatType)
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("aloy-coil-stacks", -1),
		Amount: func(atk *combat.AttackEvent, _ combat.Target) ([]float64, bool) {
			if atk.Info.AttackTag == attacks.AttackTagNormal && c.coils > 0 {
				val[attributes.DmgP] = skillCoilNABonus[c.coils-1][c.TalentLvlSkill()]
				return val, true
			}
			return nil, false
		},
	})

}

// Exit Field Hook to start timer to clear coil stacks
func (c *char) onExitField() {
	c.Core.Events.Subscribe(event.OnCharacterSwap, func(args ...interface{}) bool {
		prev := args[0].(int)
		if prev != c.Index {
			return false
		}
		c.lastFieldExit = c.Core.F

		c.Core.Tasks.Add(func() {
			if c.lastFieldExit != (c.Core.F - 30*60) {
				return
			}
			c.coils = 0
		}, 30*60)

		return false
	}, "aloy-exit")
}

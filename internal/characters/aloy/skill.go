package aloy

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

var skillFrames []int

const skillHitmark = 20 // release frame for the Bomb, travel comes on top, bomb_delay comes on top afterwards

func init() {
	skillFrames = frames.InitAbilSlice(49) // E -> Dash
	skillFrames[action.ActionAttack] = 47  // E -> N1
	skillFrames[action.ActionBurst] = 48   // E -> Q
	skillFrames[action.ActionJump] = 47    // E -> J
	skillFrames[action.ActionSwap] = 66    // E -> Swap
}

const (
	rushingIceKey = "rushingice"
)

// Skill - Handles main damage, bomblet, and coil effects
// Has 3 parameters, "bomblets" = Number of bomblets that hit
// "bomblet_coil_stacks" = Number of coil stacks gained
// "delay" - Delay in frames before bomblets go off and coil stacks get added
// Too many potential bomblet hit variations to keep syntax short, so we simplify how they can be handled here
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

	c.Core.Tasks.Add(func() {
		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Freeze Bomb",
			AttackTag:  combat.AttackTagElementalArt,
			ICDTag:     combat.ICDTagNone,
			ICDGroup:   combat.ICDGroupDefault,
			StrikeType: combat.StrikeTypeDefault,
			Element:    attributes.Cryo,
			Durability: 25,
			Mult:       skillMain[c.TalentLvlSkill()],
		}
		// TODO: accurate snapshot timing, assumes snapshot on release and not on hit/bomb creation
		c.Core.QueueAttack(
			ai,
			combat.NewCircleHit(c.Core.Combat.Player(), c.Core.Combat.PrimaryTarget(), nil, 4),
			0,
			travel,
		)
	}, skillHitmark)

	// Bomblets snapshot on cast
	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               "Chillwater Bomblets",
		AttackTag:          combat.AttackTagElementalArt,
		ICDTag:             combat.ICDTagElementalArt,
		ICDGroup:           combat.ICDGroupDefault,
		StrikeType:         combat.StrikeTypeDefault,
		Element:            attributes.Cryo,
		Durability:         25,
		Mult:               skillBomblets[c.TalentLvlSkill()],
		CanBeDefenseHalted: true,
		IsDeployable:       true,
	}

	// Queue up bomblets
	for i := 0; i < bomblets; i++ {
		c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 2), 0, skillHitmark+travel+delay+((i+1)*6), c.coilStacks)
	}

	c.Core.QueueParticle("aloy", 5, attributes.Cryo, skillHitmark+travel+c.ParticleDelay)
	c.SetCDWithDelay(action.ActionSkill, 20*60, 19)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillHitmark,
		State:           action.SkillState,
	}
}

// Handles coil stacking and associated effects, including triggering rushing ice
func (c *char) coilStacks(a combat.AttackCB) {
	if a.Target.Type() != combat.TargettableEnemy {
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

	// A1
	// When Aloy receives the Coil effect from Frozen Wilds, her ATK is increased by 16%, while nearby party members' ATK is increased by 8%. This effect lasts 10s.
	for _, char := range c.Core.Player.Chars() {
		valA1 := make([]float64, attributes.EndStatType)
		valA1[attributes.ATKP] = .08
		if char.Index == c.Index {
			valA1[attributes.ATKP] = .16
		}
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag("aloy-a1", 600),
			AffectedStat: attributes.NoStat,
			Amount: func() ([]float64, bool) {
				return valA1, true
			},
		})
	}

	if c.coils == 4 {
		c.coils = 0
		c.rushingIce()
	}
}

// Handles rushing ice state
func (c *char) rushingIce() {
	c.AddStatus(rushingIceKey, 600, true)
	c.Core.Player.AddWeaponInfuse(c.Index, "aloy-rushing-ice", attributes.Cryo, 600, true, combat.AttackTagNormal)

	// Rushing ice NA bonus
	val := make([]float64, attributes.EndStatType)
	val[attributes.DmgP] = skillRushingIceNABonus[c.TalentLvlSkill()]
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBaseWithHitlag("aloy-rushing-ice", 600),
		Amount: func(atk *combat.AttackEvent, _ combat.Target) ([]float64, bool) {
			if atk.Info.AttackTag == combat.AttackTagNormal {
				return val, true
			}
			return nil, false
		},
	})

	// A4 cryo damage increase
	valA4 := make([]float64, attributes.EndStatType)
	stacks := 1
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBaseWithHitlag("aloy-strong-strike", 600),
		AffectedStat: attributes.NoStat,
		Amount: func() ([]float64, bool) {
			if stacks > 10 {
				stacks = 10
			}
			valA4[attributes.CryoP] = float64(stacks) * 0.035
			return valA4, true
		},
	})

	for i := 0; i < 10; i++ {
		//every 1 s, affected by hitlag
		c.QueueCharTask(func() { stacks++ }, 60*(1+i))
	}
}

// Add coil mod at the beginning of the sim
// Can't be made dynamic easily as coils last until 30s after when Aloy swaps off field
func (c *char) coilMod() {
	val := make([]float64, attributes.EndStatType)
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("aloy-coil-stacks", -1),
		Amount: func(atk *combat.AttackEvent, _ combat.Target) ([]float64, bool) {
			if atk.Info.AttackTag == combat.AttackTagNormal && c.coils > 0 {
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

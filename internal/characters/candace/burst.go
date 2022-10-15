package candace

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

var burstFrames []int

const (
	burstHitmark = 33
	burstKey     = "candace-burst"
	waveHitmark  = 1
)

func init() {
	burstFrames = frames.InitAbilSlice(51)
	burstFrames[action.ActionAttack] = 50
	burstFrames[action.ActionDash] = 50
	burstFrames[action.ActionJump] = 50
	burstFrames[action.ActionSwap] = 49
}

func (c *char) Burst(p map[string]int) action.ActionInfo {
	c.waveCount = 0
	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               "Sacred Rite: Wagtail's Tide (Q)",
		AttackTag:          combat.AttackTagElementalBurst,
		ICDTag:             combat.ICDTagNone,
		ICDGroup:           combat.ICDGroupDefault,
		StrikeType:         combat.StrikeTypeBlunt,
		Element:            attributes.Hydro,
		Durability:         25,
		FlatDmg:            burstDmg[c.TalentLvlBurst()] * c.MaxHP(),
		HitlagFactor:       0.01,
		CanBeDefenseHalted: true,
	}
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHit(c.Core.Combat.Player(), 1, false, combat.TargettableEnemy),
		burstHitmark,
		burstHitmark,
	)

	duration := 540
	if c.Base.Cons >= 1 {
		duration = 720
	}

	m := make([]float64, attributes.EndStatType)
	m[attributes.DmgP] = 0.2
	// timer starts at hitmark
	c.Core.Tasks.Add(func() {
		for _, char := range c.Core.Player.Chars() {
			char.AddAttackMod(character.AttackMod{
				Base: modifier.NewBaseWithHitlag(burstKey, duration),
				Amount: func(atk *combat.AttackEvent, _ combat.Target) ([]float64, bool) {
					if atk.Info.AttackTag != combat.AttackTagNormal {
						return nil, false
					}
					if atk.Info.Element == attributes.Physical || atk.Info.Element == attributes.NoElement {
						return nil, false
					}
					return m, true
				},
			})
			switch char.Weapon.Class {
			case weapon.WeaponClassClaymore,
				weapon.WeaponClassSpear,
				weapon.WeaponClassSword:
				c.Core.Player.AddWeaponInfuse(
					c.Index,
					"candace-infuse",
					attributes.Hydro,
					duration,
					true,
					combat.AttackTagNormal, combat.AttackTagExtra, combat.AttackTagPlunge,
				) // TODO: does this refresh constantly or one time?
			}
			c.a4(char, duration)
		}
	}, burstHitmark)

	c.ConsumeEnergy(4)
	c.SetCD(action.ActionBurst, 15*60)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}
}

func (c *char) burstSwap() {
	c.Core.Events.Subscribe(event.OnCharacterSwap, func(args ...interface{}) bool {
		if c.waveCount > 2 {
			return false
		}
		next := args[1].(int)
		char := c.Core.Player.Chars()[next]
		if !char.StatusIsActive(burstKey) {
			return false
		}
		ai := combat.AttackInfo{
			ActorIndex:         c.Index,
			Abil:               "Sacred Rite: Wagtail's Tide (Wave)",
			AttackTag:          combat.AttackTagElementalBurst,
			ICDTag:             combat.ICDTagNone,
			ICDGroup:           combat.ICDGroupDefault,
			StrikeType:         combat.StrikeTypeBlunt,
			Element:            attributes.Hydro,
			Durability:         25,
			FlatDmg:            burstWaveDmg[c.TalentLvlBurst()] * c.MaxHP(),
			HitlagFactor:       0.01,
			CanBeDefenseHalted: true,
		}
		c.Core.QueueAttack(
			ai,
			combat.NewCircleHit(c.Core.Combat.Player(), 1, false, combat.TargettableEnemy),
			waveHitmark,
			waveHitmark,
		)
		c.waveCount++
		return false
	}, "candace-burst-swap")
}

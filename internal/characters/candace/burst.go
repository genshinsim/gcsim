package candace

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
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
	burstKey     = "candace-q"
	burstDmgKey  = "candace-q-dmg"
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
		AttackTag:          attacks.AttackTagElementalBurst,
		ICDTag:             attacks.ICDTagNone,
		ICDGroup:           attacks.ICDGroupDefault,
		StrikeType:         attacks.StrikeTypeDefault,
		Element:            attributes.Hydro,
		Durability:         25,
		FlatDmg:            burstDmg[c.TalentLvlBurst()] * c.MaxHP(),
		HitlagFactor:       0.01,
		CanBeDefenseHalted: true,
	}
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 3.5),
		burstHitmark,
		burstHitmark,
	)

	duration := 540
	if c.Base.Cons >= 1 {
		duration = 720
	}

	c.burstSrc = c.Core.F
	// timer starts at hitmark
	c.Core.Tasks.Add(func() {
		c.AddStatus(burstKey, duration, true)
		c.burstInfuseFn(c.CharWrapper, c.burstSrc)
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

func (c *char) burstInfuseFn(char *character.CharWrapper, src int) {
	if src != c.burstSrc {
		return
	}
	if c.Core.Player.Active() != char.Index {
		return
	}
	if !c.StatusIsActive(burstKey) {
		return
	}
	switch char.Weapon.Class {
	case weapon.WeaponClassClaymore,
		weapon.WeaponClassSpear,
		weapon.WeaponClassSword:
		c.Core.Player.AddWeaponInfuse(
			char.Index,
			"candace-q-infuse",
			attributes.Hydro,
			60,
			true,
			attacks.AttackTagNormal, attacks.AttackTagExtra, attacks.AttackTagPlunge,
		)
	}
	c.QueueCharTask(func() { c.burstInfuseFn(char, src) }, 30)
}

func (c *char) burstSwap() {
	c.Core.Events.Subscribe(event.OnCharacterSwap, func(args ...interface{}) bool {
		if !c.StatusIsActive(burstKey) {
			return false
		}
		next := args[1].(int)
		char := c.Core.Player.Chars()[next]
		c.burstInfuseFn(char, c.burstSrc)
		if c.waveCount > 2 {
			return false
		}
		ai := combat.AttackInfo{
			ActorIndex:         c.Index,
			Abil:               "Sacred Rite: Wagtail's Tide (Wave)",
			AttackTag:          attacks.AttackTagElementalBurst,
			ICDTag:             attacks.ICDTagNone,
			ICDGroup:           attacks.ICDGroupDefault,
			StrikeType:         attacks.StrikeTypeDefault,
			Element:            attributes.Hydro,
			Durability:         25,
			FlatDmg:            burstWaveDmg[c.TalentLvlBurst()] * c.MaxHP(),
			HitlagFactor:       0.01,
			CanBeDefenseHalted: true,
		}
		c.Core.QueueAttack(
			ai,
			combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 3.5),
			waveHitmark,
			waveHitmark,
		)
		c.waveCount++
		return false
	}, "candace-q-swap")
}

func (c *char) burstInit(char *character.CharWrapper) {
	m := make([]float64, attributes.EndStatType)
	m[attributes.DmgP] = 0.2
	char.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase(burstDmgKey, -1),
		Amount: func(atk *combat.AttackEvent, _ combat.Target) ([]float64, bool) {
			if !c.StatusIsActive(burstKey) {
				return nil, false
			}
			if atk.Info.AttackTag != attacks.AttackTagNormal {
				return nil, false
			}
			if atk.Info.Element == attributes.Physical || atk.Info.Element == attributes.NoElement {
				return nil, false
			}
			return m, true
		},
	})
}

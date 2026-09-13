package candace

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const c6ICDKey = "candace-c6-icd"

// When Sacred Rite: Heron's Sanctum hits opponents,
// Candace's Max HP will be increased by 20% for 15s.
func (c *char) c2() {
	m := make([]float64, attributes.EndStatType)
	m[attributes.HPP] = 0.2
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBaseWithHitlag("candace-c2", 15*60),
		AffectedStat: attributes.HPP,
		Amount: func() []float64 {
			return m
		},
	})
}

// When characters (excluding Candace herself) affected by the Prayer of the
// Crimson Crown caused by Sacred Rite: Wagtail's Tide deal Elemental DMG to
// opponents using Normal Attacks, an attack wave will be unleashed that deals
// AoE Hydro DMG equal to 15% of Candace's Max HP. This effect can trigger once
// every 2.3s and is considered Elemental Burst DMG.
func (c *char) c6() {
	c.Core.Events.Subscribe(event.OnEnemyDamage, func(args ...any) {
		atk := args[1].(*info.AttackEvent)
		dmg := args[2].(float64)
		if atk.Info.AttackTag != attacks.AttackTagNormal {
			return
		}
		if atk.Info.Element == attributes.Physical || atk.Info.Element == attributes.NoElement {
			return
		}
		if atk.Info.ActorIndex != c.Core.Player.Active() {
			return
		}
		if atk.Info.ActorIndex == c.Index() {
			return
		}
		if !c.StatusIsActive(burstKey) {
			return
		}
		if c.StatusIsActive(c6ICDKey) {
			return
		}
		if dmg == 0 {
			return
		}
		c.AddStatus(c6ICDKey, 138, true)
		ai := info.AttackInfo{
			ActorIndex:         c.Index(),
			Abil:               "The Overflow (C6)",
			AttackTag:          attacks.AttackTagElementalBurst,
			ICDTag:             attacks.ICDTagNone,
			ICDGroup:           attacks.ICDGroupDefault,
			StrikeType:         attacks.StrikeTypeDefault,
			Element:            attributes.Hydro,
			Durability:         25,
			FlatDmg:            0.15 * c.MaxHP(),
			CanBeDefenseHalted: true,
		}
		c.Core.QueueAttack(
			ai,
			combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 3.5),
			waveHitmark,
			waveHitmark,
		)
	}, "candace-c6")
}

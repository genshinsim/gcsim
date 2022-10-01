package candace

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const c6ICDKey = "candace-c6"

// When Sacred Rite: Heron's Sanctum hits opponents,
// Candace's Max HP will be increased by 20% for 15s.
func (c *char) c2() {
	m := make([]float64, attributes.EndStatType)
	m[attributes.HPP] = 0.2
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBaseWithHitlag("candace-c2", 15*60),
		AffectedStat: attributes.HPP,
		Amount: func() ([]float64, bool) {
			return m, true
		},
	})
}

// When characters (excluding Candace herself) affected by the Prayer of the
// Crimson Crown caused by Sacred Rite: Wagtail's Tide deal Elemental DMG to
// opponents using Normal Attacks, an attack wave will be unleashed that deals
// AoE Hydro DMG equal to 15% of Candace's Max HP. This effect can trigger once
// every 2.3s and is considered Elemental Burst DMG.
func (c *char) c6() {
	c.Core.Events.Subscribe(event.OnDamage, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.AttackTag != combat.AttackTagNormal {
			return false
		}
		if atk.Info.Element == attributes.Physical || atk.Info.Element == attributes.NoElement {
			return false
		}
		if atk.Info.ActorIndex != c.Core.Player.Active() {
			return false
		}
		if atk.Info.ActorIndex == c.Index {
			return false
		}
		if !c.StatusIsActive(burstKey) {
			return false
		}
		if c.StatusIsActive(c6ICDKey) {
			return false
		}
		c.AddStatus(c6ICDKey, 138, false) // TODO: is c6 hitlag affected?
		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "The Overflow (C6)",
			AttackTag:  combat.AttackTagElementalBurst,
			ICDTag:     combat.ICDTagNone,
			ICDGroup:   combat.ICDGroupDefault,
			StrikeType: combat.StrikeTypeBlunt,
			Element:    attributes.Hydro,
			Durability: 25,
			FlatDmg:    0.15 * c.MaxHP(),
		}
		c.Core.QueueAttack(
			ai,
			combat.NewCircleHit(c.Core.Combat.Player(), 1, false, combat.TargettableEnemy),
			waveHitmark, // TODO find correct timing
			waveHitmark,
		)
		return false
	}, "yunjin-burst")
}

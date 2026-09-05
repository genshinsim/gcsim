package cryo

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	c1Key        = "travelercryo-c1"
	c1ICDKey     = "travelercryo-c1-icd"
	c2Key        = "travelercryo-c2"
	c2UpgradeKey = "travelercryo-upgrade-c2"
	c6Key        = "travelercryo-c6"
)

// The Traveler regenerates 5 Elemental Energy when they deal Stellar Glimmer DMG. This effect can
// trigger once every 0.5s.
func (c *Traveler) c1Init() {
	if c.Base.Cons < 1 {
		return
	}
	c.Core.Events.Subscribe(event.OnEnemyDamage, func(args ...any) {
		atk := args[1].(*info.AttackEvent)

		if _, ok := args[0].(*enemy.Enemy); !ok {
			return
		}
		switch atk.Info.AttackTag {
		case attacks.AttackTagDirectStellarConduct:
		case attacks.AttackTagDirectStellarSwirl:
		case attacks.AttackTagReactionStellarSwirl:
		default:
			return
		}

		if atk.Info.ActorIndex != c.Index() {
			return
		}

		if c.StatusIsActive(c1ICDKey) {
			return
		}

		c.AddStatus(c1ICDKey, 0.5*60, true)

		c.AddEnergy(c1Key, 5)
	}, c1Key)
}

// The active character's Elemental Mastery is increased by 60 for 5s when an opponent is hit by an
// ice crystal fired by the Elemental Skill Ice Fog Piercer.
// Additionally, when a character affected by the aforementioned effect triggers a Stellar Glimmer
// reaction or deals Stellar Glimmer reaction DMG, the Elemental Mastery boost described above will
// be increased to 120 for its duration.
func (c *Traveler) c2Init() {
	if c.Base.Cons < 2 {
		return
	}
	fReactionHook := func(args ...any) {
		if _, ok := args[0].(*enemy.Enemy); !ok {
			return
		}

		atk := args[1].(*info.AttackEvent)

		char := c.Core.Player.Chars()[atk.Info.ActorIndex]

		if !char.StatusIsActive(c2Key) {
			return
		}

		char.AddStatus(c2UpgradeKey, 5*60, true)
	}

	fDMGHook := func(args ...any) {
		if _, ok := args[0].(*enemy.Enemy); !ok {
			return
		}

		atk := args[1].(*info.AttackEvent)

		switch atk.Info.AttackTag {
		case attacks.AttackTagDirectStellarConduct:
		case attacks.AttackTagDirectStellarSwirl:
		case attacks.AttackTagReactionStellarSwirl:
		default:
			return
		}

		char := c.Core.Player.Chars()[atk.Info.ActorIndex]
		if !char.StatModIsActive(c2Key) {
			return
		}

		c2Dur := char.StatusDuration(c2Key)
		char.AddStatus(c2UpgradeKey, c2Dur, true)
	}

	c.Core.Events.Subscribe(event.OnStellarConduct, fReactionHook, c2Key)
	c.Core.Events.Subscribe(event.OnStellarSwirl, fReactionHook, c2Key)
	c.Core.Events.Subscribe(event.OnEnemyDamage, fDMGHook, c2Key)
}

func (c *Traveler) c2CB(a info.AttackCB) {
	if c.Base.Cons < 2 {
		return
	}
	if a.Target.Type() != info.TargettableEnemy {
		return
	}
	m := make([]float64, attributes.EndStatType)
	char := c.Core.Player.ActiveChar()
	char.AddStatMod(character.StatMod{
		Base:         modifier.NewBaseWithHitlag(c2Key, 5*60),
		AffectedStat: attributes.EM,
		Amount: func() []float64 {
			if char.StatusIsActive(c2UpgradeKey) {
				m[attributes.EM] = 120
			} else {
				m[attributes.EM] = 60
			}
			return m
		},
	})
}

// The duration of the Frostpierce Star created by the Elemental Skill Ice Fog Piercer is increased
// by 25%.
func (c *Traveler) c4SkillBonusDur() int {
	if c.Base.Cons < 4 {
		return 0
	}

	return 3 * 60
}

// Every stack of Frostglow consumed when using the Elemental Burst Frostbound Javelin increases
// Stellar Glimmer reaction DMG dealt by other party members by 5% for 15s. A maximum 40% increase
// can be obtained in this way.
func (c *Traveler) c6OnBurst(stacks int) {
	if c.Base.Cons < 6 {
		return
	}
	for _, char := range c.Core.Player.Chars() {
		char.AddReactBonusMod(character.ReactBonusMod{
			Base: modifier.NewBaseWithHitlag(c6Key, 15*60),
			Amount: func(ai info.AttackInfo) float64 {
				switch ai.AttackTag {
				case
					attacks.AttackTagDirectStellarConduct,
					attacks.AttackTagDirectStellarSwirl,
					attacks.AttackTagReactionStellarSwirl:
					return 0.05 * float64(stacks)
				}
				return 0
			},
		})
	}
}

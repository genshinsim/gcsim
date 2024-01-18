package chevreuse

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

// When all party members are Pyro and Electro characters and there is at least
// one Pyro and one Electro character each in the party:
// Chevreuse grants "Coordinated Tactics" to nearby party members:
// After a character triggers the Overloaded reaction, the Pyro and Electro RES
// of the opponent(s) affected by this Overloaded reaction will be decreased by 40% for 6s.
// The "Coordinated Tactics" effect will be removed when the Elemental Types
// of the characters in the party do not meet the basic requirements for the Passive Talent.
func (c *char) a1() {
	if c.Base.Ascension < 1 {
		return
	}

	// check if only pyro + electro
	chars := c.Core.Player.Chars()
	count := make(map[attributes.Element]int)
	for _, this := range chars {
		count[this.Base.Element]++
	}
	c.onlyPyroElectro = count[attributes.Pyro] > 0 && count[attributes.Electro] > 0 && count[attributes.Electro]+count[attributes.Pyro] == len(chars)

	if !c.onlyPyroElectro {
		return
	}

	c.Core.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		// don't trigger if no overload dmg
		if atk.Info.AttackTag != attacks.AttackTagOverloadDamage {
			return false
		}

		t, ok := args[0].(*enemy.Enemy)
		if !ok {
			return false
		}
		t.AddResistMod(combat.ResistMod{
			Base:  modifier.NewBaseWithHitlag("chev-a1-pyro", 6*60),
			Ele:   attributes.Pyro,
			Value: -0.40,
		})
		t.AddResistMod(combat.ResistMod{
			Base:  modifier.NewBaseWithHitlag("chev-a1-electro", 6*60),
			Ele:   attributes.Electro,
			Value: -0.40,
		})

		return false
	}, "cheuv-a1")
}

// After Chevreuse fires an Overcharged Ball using Short-Range Rapid Interdiction Fire,
// nearby Pyro and Electro characters in the party gain 1% increased ATK for every 1,000 Max HP Chevreuse has for 30s.
// ATK can be increased by up to 40% in this way.
func (c *char) a4() {
	if c.Base.Ascension < 4 {
		return
	}

	m := make([]float64, attributes.EndStatType)
	m[attributes.ATKP] = min(c.MaxHP()/1000*0.01, 0.4)
	for _, char := range c.Core.Player.Chars() {
		if char.Base.Element != attributes.Pyro && char.Base.Element != attributes.Electro {
			continue
		}
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag("chev-a4", 30*60),
			AffectedStat: attributes.ATKP,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})
	}
}

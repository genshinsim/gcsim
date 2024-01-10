package chevreuse

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

/*
When the Elemental Type of all party members is Pyro or Electro and there is at least
one Pyro and one Electro Elemental Type each in the party:
Chevreuse grants "Coordinated Tactics" to nearby party members: After a character triggers the Overloaded reaction,

	the Pyro and Electro RES of the opponent(s) affected by this Overloaded reaction will be decreased by 40% for 6s.
*/
func (c *char) a1() {
	if c.onlyPyroElectro {
		c.Core.Events.Subscribe(event.OnOverload, c.OnOverload, "cheuv-a1")
	}
}

func (c *char) OnOverload(args ...interface{}) bool {
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
}

/*
After Chevreuse fires an Overcharged Ball using Short-Range Rapid Interdiction Fire,
nearby Pyro and Electro characters in the party gain 1% increased ATK for every 1000 Max HP Chevreuse has for 30s.
ATK can be increased by up to 40% in this way.
*/
func (c *char) a4() {
	if c.Base.Ascension < 4 {
		return
	}

	m := make([]float64, attributes.EndStatType)
	atkBuff := c.MaxHP() / 1000 * 0.01

	if atkBuff > 0.4 {
		atkBuff = 0.4
	}

	m[attributes.ATKP] = atkBuff

	buffDuration := 30

	for _, char := range c.Core.Player.Chars() {
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag("chev-a4", buffDuration*60),
			AffectedStat: attributes.NoStat,
			Extra:        true,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})
	}

}

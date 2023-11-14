package mika

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	a1Stacks = "detector-stacks"
	a1Buff   = "detector-buff"
)

func (c *char) addDetectorStack() {
	stacks := c.Tag(a1Stacks)

	if stacks < c.maxDetectorStacks {
		stacks++
		c.Core.Log.NewEvent("add detector stack", glog.LogCharacterEvent, c.Index).
			Write("stacks", stacks).
			Write("maxstacks", c.maxDetectorStacks)
	}
	c.SetTag(a1Stacks, stacks)
}

/*
When the Elemental Type of all party members is Pyro or Electro and there is at least
one Pyro and one Electro Elemental Type each in the party:
Chevreuse grants "Coordinated Tactics" to nearby party members: After a character triggers the Overloaded reaction,

	the Pyro and Electro RES of the opponent(s) affected by this Overloaded reaction will be decreased by 40% for 6s.
*/
func (c *char) a1(char *character.CharWrapper) {
	c.Core.Events.Subscribe(event.OnOverload, func(args ...interface{}) bool {
		c.OnOverload()
		return false
	}, "overload")
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

	return true
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
	atkBuff := c.MaxHP() * 0.01

	if atkBuff > 0.4 {
		atkBuff = 0.4
	}

	m[attributes.ATKP] = atkBuff

	buffDuration := 30

	for _, char := range c.Core.Player.Chars() {
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag("chev-a4", buffDuration),
			AffectedStat: attributes.NoStat,
			Extra:        true,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})
	}

}

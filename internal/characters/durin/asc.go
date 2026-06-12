package durin

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

var a1ReactToElements = map[event.Event][]attributes.Element{
	event.OnOverload:        {attributes.Electro, attributes.Pyro},
	event.OnSwirlPyro:       {attributes.Anemo, attributes.Pyro},
	event.OnCrystallizePyro: {attributes.Geo, attributes.Pyro},
	event.OnBurning:         {attributes.Dendro, attributes.Pyro},
}

const a1BlackKey = "durin-a1-black"

func (c *char) a1Init() {
	if c.Base.Ascension < 1 {
		return
	}

	for event, elements := range a1ReactToElements {
		c.Core.Events.Subscribe(event, c.a1MakeResShred(elements), fmt.Sprintf("durin-a1-hook-%v", event))
	}

	c.Core.Events.Subscribe(event.OnEnemyDamage, func(args ...any) {
		t, ok := args[0].(*enemy.Enemy)
		if !ok {
			return
		}

		atk := args[1].(*info.AttackEvent)

		if !t.IsBurning() {
			return
		}

		switch atk.Info.Element {
		case attributes.Dendro:
		case attributes.Pyro:
		default:
			return
		}

		if !c.StatusIsActive(burstKeyWhite) {
			return
		}

		t.AddResistMod(info.ResistMod{
			Base:  modifier.NewBaseWithHitlag("durin-a1-"+attributes.Dendro.String(), 6*60),
			Ele:   attributes.Dendro,
			Value: -0.20 * c.hexereiA1Bonus(),
		})

		t.AddResistMod(info.ResistMod{
			Base:  modifier.NewBaseWithHitlag("durin-a1-"+attributes.Pyro.String(), 6*60),
			Ele:   attributes.Pyro,
			Value: -0.20 * c.hexereiA1Bonus(),
		})
	}, "durin-a1-hook-on-dmg")
}

func (c *char) a1OnBurst(isWhite bool) {
	if c.Base.Ascension < 1 {
		return
	}

	if isWhite {
		c.DeleteReactBonusMod(a1BlackKey)
		return
	}
	reactMod := character.ReactBonusMod{
		Base: modifier.NewBaseWithHitlag(a1BlackKey, burstDuration),
		Amount: func(ai info.AttackInfo) float64 {
			if ai.Amped {
				return 0.40 * c.hexereiA1Bonus()
			}
			return 0
		},
	}
	c.AddReactBonusMod(reactMod)
}

func (c *char) a1MakeResShred(elements []attributes.Element) func(args ...any) {
	return func(args ...any) {
		t, ok := args[0].(*enemy.Enemy)
		if !ok {
			return
		}

		if !c.StatusIsActive(burstKeyWhite) {
			return
		}

		for _, ele := range elements {
			t.AddResistMod(info.ResistMod{
				Base:  modifier.NewBaseWithHitlag("durin-a1-"+ele.String(), 6*60),
				Ele:   ele,
				Value: -0.20 * c.hexereiA1Bonus(),
			})
		}
	}
}

func (c *char) a4OnBurst() {
	if c.Base.Ascension < 4 {
		return
	}
	c.a4stacks = 10
}

func (c *char) a4Dmg() float64 {
	if c.Base.Ascension < 4 {
		return 1.0
	}

	if c.a4stacks <= 0 {
		return 1.0
	}

	if !c.StatusIsActive(burstKeyWhite) && !c.StatusIsActive(burstKeyBlack) {
		return 1.0
	}

	c.a4stacks -= 1

	bonus := min(c.TotalAtk()/100*0.03, 0.75)

	return 1 + bonus
}

func (c *char) hexereiA1Bonus() float64 {
	if !c.IsHexerei {
		return 1
	}

	if c.Core.Player.GetHexereiCount() < 2 {
		return 1
	}

	return 1.75
}

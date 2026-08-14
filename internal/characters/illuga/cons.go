package illuga

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

var geoReactionEvents = []event.Event{event.OnCrystallizeCryo, event.OnCrystallizeElectro, event.OnCrystallizeHydro, event.OnCrystallizeElectro, event.OnLunarCrystallize}

const (
	c1EnergyKey      = "illuga-c1-energy"
	c1ICDKey         = "illuga-c1-icd"
	c2QuillThreshold = 7
	c2Hitmark        = 50
	c6CR             = 0.1
	c6CD             = 0.3
	c6EM             = 80
)

func (c *char) c1Init() {
	if c.Base.Cons < 1 {
		return
	}

	energyGain := func(args ...any) {
		if c.StatusIsActive(c1ICDKey) {
			return
		}
		char := args[0].(*character.CharWrapper)

		if char.Index() != c.Index() {
			return
		}

		if char.Index() != c.Core.Player.Active() {
			return
		}

		c.AddStatus(c1ICDKey, 15*60, true)
		c.AddEnergy(c1EnergyKey, 12)
	}

	for _, event := range geoReactionEvents {
		c.Core.Events.Subscribe(event, energyGain, "illuga-c1")
	}
}

func (c *char) c2Reset() {
	if c.Base.Cons < 2 {
		return
	}

	c.c2QuillCounter = 0
}

func (c *char) c2Increment() {
	if c.Base.Cons < 2 {
		return
	}

	c.c2QuillCounter++

	if c.c2QuillCounter >= c2QuillThreshold {
		c.c2QuillCounter -= c2QuillThreshold
		c.c2Hit()
	}
}

func (c *char) c2Hit() {
	if c.Base.Cons < 2 {
		return
	}

	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Aedon C2 Hit",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeBlunt,
		PoiseDMG:   75,
		Element:    attributes.Geo,
		Durability: 25,
	}

	ai.FlatDmg += c.Stat(attributes.EM) * 4
	ai.FlatDmg += c.TotalDef(false) * 2

	c.Core.QueueAttack(
		ai,
		combat.NewSingleTargetHit(c.Core.Combat.PrimaryTarget().Key()),
		c2Hitmark,
		c2Hitmark,
	)
}

func (c *char) c4(src int) func() {
	if c.Base.Cons < 4 {
		return func() {}
	}

	m := make([]float64, attributes.EndStatType)
	m[attributes.DEF] = 200

	return func() {
		if src < c.c4Src {
			return
		}

		// new source will alway be higher than the previous one
		if src > c.c4Src {
			c.c4Src = src
		}

		if !c.StatusIsActive(burstKey) {
			return
		}

		for _, char := range c.Core.Player.Chars() {
			if char.Index() == c.Index() {
				continue
			}
			char.AddStatMod(character.StatMod{
				Base: modifier.NewBaseWithHitlag("illuga-c4", 1.1*60),
				Amount: func() []float64 {
					return m
				},
			})
		}

		c.QueueCharTask(c.c4(src), 1*60)
	}
}

func (c *char) c6CR() float64 {
	if c.Base.Cons < 6 {
		return 0
	}

	return c6CR - a1CR
}

func (c *char) c6CD() float64 {
	if c.Base.Cons < 6 {
		return 0
	}

	return c6CD - a1CD
}

func (c *char) c6EM() float64 {
	if c.Base.Cons < 6 {
		return 0
	}

	return c6EM - a1EM
}

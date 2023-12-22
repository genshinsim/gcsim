package neuvillette

import (
	"math"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/gadget"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

type NeuvA1Keys struct {
	Evt event.Event
	Key string
}

var a1Multipliers = [4]float64{1, 1.1, 1.25, 1.6}

func (c *char) a1() {
	a1 := []NeuvA1Keys{
		{event.OnBloom, "neuvillette-a1-bloom"},
		{event.OnCrystallizeHydro, "neuvillette-a1-crystallize-hydro"},
		{event.OnElectroCharged, "neuvillette-a1-electro-charged"},
		{event.OnFrozen, "neuvillette-a1-frozen"},
		{event.OnSwirlHydro, "neuvillette-a1-swirl-hydro"},
		{event.OnVaporize, "neuvillette-a1-vaporize"},
	}

	c.a1Statuses = append(c.a1Statuses,
		a1...,
	)

	for _, val := range a1 {
		// need to make a copy of key for the status key
		key := val.Key
		c.Core.Events.Subscribe(val.Evt, func(args ...interface{}) bool {
			if _, ok := args[0].(*gadget.Gadget); ok {
				return false
			}
			c.AddStatus(key, 30*60, true)
			return false
		}, key)
	}
}

func (c *char) countA1() int {
	if c.Base.Ascension < 1 {
		return 0
	}
	a1TriggeredReactionsCount := 0
	for _, val := range c.a1Statuses {
		if c.StatusIsActive(val.Key) {
			a1TriggeredReactionsCount += 1
		}
		if a1TriggeredReactionsCount == 3 {
			break
		}
	}
	return a1TriggeredReactionsCount
}

func (c *char) a4() {
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("neuvillette-a4", -1),
		AffectedStat: attributes.HydroP,
		Extra:        true,
		Amount: func() ([]float64, bool) {
			return c.a4Buff, true
		},
	})

	c.Core.Events.Subscribe(event.OnPlayerHPDrain, func(args ...interface{}) bool {
		di := args[0].(player.DrainInfo)

		if di.Amount <= 0 {
			return false
		}

		if di.ActorIndex != c.Index {
			return false
		}

		c.updateA4()

		return false
	}, "neuv-a4-update-on-hp-drain")

	c.Core.Events.Subscribe(event.OnHeal, func(args ...interface{}) bool {
		target := args[1].(int)
		amount := args[2].(float64)
		overheal := args[3].(float64)

		if amount <= 0 {
			return false
		}

		if math.Abs(amount-overheal) <= 1e-9 {
			return false
		}

		if target != c.Index {
			return false
		}

		c.updateA4()

		return false
	}, "neuv-a4-update-on-heal")

	c.updateA4()
}

func (c *char) updateA4() {
	hpRatio := c.CurrentHPRatio()
	hydroDmgBuff := (hpRatio - 0.3) * 0.6

	if hydroDmgBuff < 0 {
		hydroDmgBuff = 0
	} else if hydroDmgBuff > 0.3 {
		hydroDmgBuff = 0.3
	}

	c.a4Buff[attributes.HydroP] = hydroDmgBuff
}

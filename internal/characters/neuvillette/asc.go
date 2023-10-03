package neuvillette

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

type NeuvA1Keys struct {
	Evt event.Event
	Key string
}

func init() {

}

var a1Multipliers = [4]float64{1, 1.1, 1.25, 1.6}

func (c *char) a1() {
	if c.Base.Ascension < 1 {
		return
	}
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
		c.Core.Events.Subscribe(val.Evt, func(args ...interface{}) bool {
			c.AddStatus(val.Key, 30*60, false)
			return false
		}, val.Key)
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
	if c.Base.Ascension < 4 {
		return
	}

	c.AddStatMod(character.StatMod{
		AffectedStat: attributes.HydroP,
		Extra:        true,
		Base:         modifier.NewBase("neuvillette-a4", -1),
		Amount: func() ([]float64, bool) {
			return c.a4Buff, true
		},
	})

	c.a4Tick()
}

func (c *char) a4Tick() {
	hpRatio := c.CurrentHPRatio()
	hydroDmgBuff := (hpRatio - 0.3) * 0.6

	if hydroDmgBuff < 0 {
		hydroDmgBuff = 0
	} else if hydroDmgBuff > 0.3 {
		hydroDmgBuff = 0.3
	}

	c.a4Buff[attributes.HydroP] = hydroDmgBuff

	// TODO: Is this on HP change or just on tick?
	c.Core.Tasks.Add(c.a4Tick, 30)
}

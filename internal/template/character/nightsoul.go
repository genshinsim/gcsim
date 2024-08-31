package character

import (
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

const NightsoulBlessingStatus = "nightsoul-blessing"

func (c *Character) EnterNightsoulBlessing(amount float64) {
	c.NightsoulPoints = amount
	c.AddStatus(NightsoulBlessingStatus, -1, true)
	c.Core.Log.NewEvent("enter nightsoul blessing", glog.LogCharacterEvent, c.Index).
		Write("points", c.NightsoulPoints)
}

func (c *Character) ExitNightsoulBlessing() {
	c.DeleteStatus(NightsoulBlessingStatus)
	c.Core.Log.NewEvent("exit nightsoul blessing", glog.LogCharacterEvent, c.Index)
}

func (c *Character) GenerateNightsoulPoints(amount float64) {
	prevPoints := c.NightsoulPoints
	c.NightsoulPoints += amount
	c.Core.Events.Emit(event.OnNightsoulGenerate, c.Index, amount)
	c.Core.Log.NewEvent("generate nightsoul points", glog.LogCharacterEvent, c.Index).
		Write("previous points", prevPoints).
		Write("amount", amount).
		Write("final", c.NightsoulPoints)
}

func (c *Character) ConsumeNightsoulPoints(amount float64) {
	prevPoints := c.NightsoulPoints
	c.NightsoulPoints -= amount
	c.Core.Events.Emit(event.OnNightsoulConsume, c.Index, amount)
	c.Core.Log.NewEvent("consume nightsoul points", glog.LogCharacterEvent, c.Index).
		Write("previous points", prevPoints).
		Write("amount", amount).
		Write("final", c.NightsoulPoints)
}

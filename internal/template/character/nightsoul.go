package character

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

const NightsoulBlessingStatus = "nightsoul-blessing"

type Nightsoul struct {
	char            *character.CharWrapper
	c               *core.Core
	nightsoulPoints float64
}

func NewNightsoul(c *core.Core, char *character.CharWrapper) *Nightsoul {
	t := &Nightsoul{
		char: char,
		c:    c,
	}
	return t
}

func (n *Nightsoul) EnterBlessing(amount float64) {
	n.nightsoulPoints = amount
	n.char.AddStatus(NightsoulBlessingStatus, -1, true)
	n.c.Log.NewEvent("enter nightsoul blessing", glog.LogCharacterEvent, n.char.Index).
		Write("points", n.nightsoulPoints)
}

func (n *Nightsoul) ExitBlessing() {
	n.char.DeleteStatus(NightsoulBlessingStatus)
	n.c.Log.NewEvent("exit nightsoul blessing", glog.LogCharacterEvent, n.char.Index)
}

func (n *Nightsoul) GeneratePoints(amount float64) {
	prevPoints := n.char.NightsoulPoints
	n.nightsoulPoints += amount
	n.c.Events.Emit(event.OnNightsoulGenerate, n.char.Index, amount)
	n.c.Log.NewEvent("generate nightsoul points", glog.LogCharacterEvent, n.char.Index).
		Write("previous points", prevPoints).
		Write("amount", amount).
		Write("final", n.nightsoulPoints)
}

func (n *Nightsoul) ConsumePoints(amount float64) {
	prevPoints := n.nightsoulPoints
	n.nightsoulPoints -= amount
	n.c.Events.Emit(event.OnNightsoulConsume, n.char.Index, amount)
	n.c.Log.NewEvent("consume nightsoul points", glog.LogCharacterEvent, n.char.Index).
		Write("previous points", prevPoints).
		Write("amount", amount).
		Write("final", n.nightsoulPoints)
}

func (n *Nightsoul) Points() float64 {
	return n.nightsoulPoints
}

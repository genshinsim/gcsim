package nefer

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/gadget"
)

const seedDuration = 12 * 60

type seedGadget struct {
	*gadget.Gadget
}

func newSeedGadget(c *core.Core, pos info.Point) *seedGadget {
	g := &seedGadget{}
	g.Gadget = gadget.New(c, pos, 2, info.GadgetTypDendroCore)
	g.Duration = seedDuration
	return g
}

func (g *seedGadget) HandleAttack(atk *info.AttackEvent) float64 {
	g.Core.Events.Emit(event.OnGadgetHit, g, atk)
	return 0
}

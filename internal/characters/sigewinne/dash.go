package sigewinne

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/template/sourcewaterdroplet"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

const (
	BoLPctPerDroplet = 0.1
)

func (c *char) Dash(p map[string]int) (action.Info, error) {
	c.burstEarlyCancelled = false
	dropletsToPickup, ok := p["pickup_droplets"]
	if !ok {
		return c.Character.Dash(p)
	}
	if dropletsToPickup == 0 {
		return c.Character.Dash(p)
	}
	droplets := c.getSourcewaterDroplets()
	dropletsToPickup = min(dropletsToPickup, len(droplets))

	// TODO: additional delay if >2 droplets picked up
	indices := c.Core.Combat.Rand.Perm(dropletsToPickup)
	for _, ind := range indices {
		g := droplets[ind]
		c.consumeDroplet(g)
	}
	c.Core.Combat.Log.NewEvent(fmt.Sprint("Picked up ", dropletsToPickup, " droplets"), glog.LogCharacterEvent, c.Index)

	return c.Character.Dash(p)
}

func (c *char) consumeDroplet(g *sourcewaterdroplet.Gadget) {
	g.Kill()
	c.ModifyHPDebtByAmount(c.MaxHP() * BoLPctPerDroplet)
}

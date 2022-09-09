// gadget provides a skeleton implementation for gadgets and is
// basically a thin wrapper around target.Target and adds some
// helper functionalities
package gadget

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/target"
)

type Gadget struct {
	*target.Target
	OnRemoved       func()
	ThinkInterval   int //should be > 0
	OnThinkInterval func()
	//internal helper
	sinceLastThink int
}

func New(core *core.Core, pos core.Coord) *Gadget {
	g := &Gadget{}
	g.Target = target.New(core, pos.X, pos.Y, pos.R)
	return g
}

func (g *Gadget) Kill() {
	g.Target.Kill()
	if g.OnRemoved != nil {
		g.OnRemoved()
	}
}

func (g *Gadget) Tick() {
	if g.OnThinkInterval != nil && g.ThinkInterval > 0 {
		if g.sinceLastThink < g.ThinkInterval {
			g.sinceLastThink++
		}
		if g.sinceLastThink != g.ThinkInterval {
			return
		}
		g.OnThinkInterval()
		g.sinceLastThink = 0
	}
}

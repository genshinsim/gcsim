// gadget provides a skeleton implementation for gadgets and is
// basically a thin wrapper around target.Target and adds some
// helper functionalities
package gadget

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/target"
)

type Gadget struct {
	*target.Target
	core            *core.Core
	OnRemoved       func()
	ThinkInterval   int //should be > 0
	OnThinkInterval func()
	Duration        int //how long gadget should live for; use -1 for infinite
	//internal helper
	sinceLastThink int
}

func New(core *core.Core, pos core.Coord) *Gadget {
	g := &Gadget{
		core: core,
	}
	g.Target = target.New(core, pos.X, pos.Y, pos.R)
	return g
}

func (g *Gadget) Kill() {
	g.Target.Kill()
	if g.OnRemoved != nil {
		g.OnRemoved()
	}
	g.core.Combat.RemoveGadget(g.Index())
}

func (g *Gadget) Type() combat.TargettableType { return combat.TargettableGadget }

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
	if g.Duration != -1 && g.Duration > 0 {
		g.Duration--
		if g.Duration == 0 {
			g.Kill()
		}
	}
}

// gadget provides a skeleton implementation for gadgets and is
// basically a thin wrapper around target.Target and adds some
// helper functionalities
package gadget

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/targets"
	"github.com/genshinsim/gcsim/pkg/target"
)

type Gadget struct {
	*target.Target
	src             int
	gadgetTyp       combat.GadgetTyp
	core            *core.Core
	OnKill          func()
	OnExpiry        func() // only called if gadget dies from expiry
	ThinkInterval   int    //should be > 0
	OnThinkInterval func()
	Duration        int //how long gadget should live for; use -1 for infinite
	//internal helper
	sinceLastThink int
}

func New(core *core.Core, p geometry.Point, r float64, typ combat.GadgetTyp) *Gadget {
	g := &Gadget{
		core:      core,
		src:       core.F,
		gadgetTyp: typ,
	}
	g.Target = target.New(core, p, r)
	return g
}

func (g *Gadget) Kill() {
	g.Target.Kill()
	if g.OnKill != nil {
		g.OnKill()
	}
	g.core.Combat.RemoveGadget(g.Key())
}

func (g *Gadget) Type() targets.TargettableType { return targets.TargettableGadget }
func (g *Gadget) Src() int                      { return g.src }
func (g *Gadget) GadgetTyp() combat.GadgetTyp   { return g.gadgetTyp }

func (g *Gadget) Tick() {
	if g.OnThinkInterval != nil && g.ThinkInterval > 0 {
		if g.sinceLastThink < g.ThinkInterval {
			g.sinceLastThink++
		}
		if g.sinceLastThink == g.ThinkInterval {
			g.OnThinkInterval()
			g.sinceLastThink = 0
		}
	}
	if g.Duration != -1 && g.Duration > 0 {
		g.Duration--
		if g.Duration == 0 {
			if g.OnExpiry != nil {
				g.OnExpiry()
			}
			g.core.Combat.RemoveGadget(g.Key())
		}
	}
}

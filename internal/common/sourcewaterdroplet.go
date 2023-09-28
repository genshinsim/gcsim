package common

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/targets"
	"github.com/genshinsim/gcsim/pkg/gadget"
)

type SourcewaterDroplet struct {
	*gadget.Gadget
}

func NewSourcewaterDroplet(core *core.Core, pos geometry.Point) *SourcewaterDroplet {
	p := &SourcewaterDroplet{}
	p.Gadget = gadget.New(core, pos, 0.3, combat.GadgetTypSourcewaterDroplet)
	p.Gadget.Duration = 914
	core.Combat.AddGadget(p)
	return p
}

func (s *SourcewaterDroplet) Tick() {
	// this is needed since gadget tick
	s.Gadget.Tick()
}

func (s *SourcewaterDroplet) Type() targets.TargettableType                          { return targets.TargettableGadget }
func (s *SourcewaterDroplet) HandleAttack(*combat.AttackEvent) float64               { return 0 }
func (s *SourcewaterDroplet) Attack(*combat.AttackEvent, glog.Event) (float64, bool) { return 0, false }

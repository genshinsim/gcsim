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

func NewSourcewaterDroplet(core *core.Core, pos geometry.Point, typ combat.GadgetTyp) *SourcewaterDroplet {
	p := &SourcewaterDroplet{}
	p.Gadget = gadget.New(core, pos, 1, typ)
	p.Gadget.Duration = 878
	core.Combat.AddGadget(p)
	return p
}

func (s *SourcewaterDroplet) HandleAttack(*combat.AttackEvent) float64 { return 0 }
func (s *SourcewaterDroplet) SetDirection(trg geometry.Point)          {}
func (s *SourcewaterDroplet) SetDirectionToClosestEnemy()              {}
func (s *SourcewaterDroplet) CalcTempDirection(trg geometry.Point) geometry.Point {
	return geometry.DefaultDirection()
}

func (s *SourcewaterDroplet) Type() targets.TargettableType                          { return targets.TargettableGadget }
func (s *SourcewaterDroplet) Attack(*combat.AttackEvent, glog.Event) (float64, bool) { return 0, false }

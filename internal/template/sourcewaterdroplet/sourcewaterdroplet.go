package sourcewaterdroplet

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/targets"
	"github.com/genshinsim/gcsim/pkg/gadget"
)

type Gadget struct {
	*gadget.Gadget
}

func New(core *core.Core, pos geometry.Point, typ combat.GadgetTyp) *Gadget {
	p := &Gadget{}
	p.Gadget = gadget.New(core, pos, 1, typ)
	p.Gadget.Duration = 878
	core.Combat.AddGadget(p)
	return p
}

func (s *Gadget) HandleAttack(*combat.AttackEvent) float64 { return 0 }
func (s *Gadget) SetDirection(trg geometry.Point)          {}
func (s *Gadget) SetDirectionToClosestEnemy()              {}
func (s *Gadget) CalcTempDirection(trg geometry.Point) geometry.Point {
	return geometry.DefaultDirection()
}

func (s *Gadget) Type() targets.TargettableType                          { return targets.TargettableGadget }
func (s *Gadget) Attack(*combat.AttackEvent, glog.Event) (float64, bool) { return 0, false }

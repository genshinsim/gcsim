package sourcewaterdroplet

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/gadget"
)

type Gadget struct {
	*gadget.Gadget
}

func New(core *core.Core, pos info.Point, typ info.GadgetTyp) *Gadget {
	p := &Gadget{}
	p.Gadget = gadget.New(core, pos, 1, typ)
	p.Duration = 878
	core.Combat.AddGadget(p)
	return p
}

func (s *Gadget) HandleAttack(*info.AttackEvent) float64 { return 0 }
func (s *Gadget) SetDirection(trg info.Point)            {}
func (s *Gadget) SetDirectionToClosestEnemy()            {}
func (s *Gadget) CalcTempDirection(trg info.Point) info.Point {
	return info.DefaultDirection()
}

func (s *Gadget) Type() info.TargettableType                           { return info.TargettableGadget }
func (s *Gadget) Attack(*info.AttackEvent, glog.Event) (float64, bool) { return 0, false }

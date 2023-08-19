package hydro

import (
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/targets"
	"github.com/genshinsim/gcsim/pkg/gadget"
)

type sourcewaterDroplet struct {
	*gadget.Gadget
	c *char
}

func (c *char) newDropblet() *sourcewaterDroplet {
	p := &sourcewaterDroplet{c: c}
	player := c.Core.Combat.Player()
	pos := geometry.CalcOffsetPoint(
		player.Pos(),
		geometry.Point{X: 0.7, Y: 3.5},
		player.Direction(),
	)
	p.Gadget = gadget.New(c.Core, pos, 0.3, combat.GadgetTypSourcewaterDroplet)
	p.Gadget.Duration = 914

	return p
}

func (s *sourcewaterDroplet) Tick() {
	// this is needed since gadget tick
	s.Gadget.Tick()
}

func (s *sourcewaterDroplet) Type() targets.TargettableType                          { return targets.TargettableGadget }
func (s *sourcewaterDroplet) HandleAttack(*combat.AttackEvent) float64               { return 0 }
func (s *sourcewaterDroplet) Attack(*combat.AttackEvent, glog.Event) (float64, bool) { return 0, false }

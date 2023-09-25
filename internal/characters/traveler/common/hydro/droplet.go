package hydro

import (
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/gadget"
)

type sourcewaterDroplet struct {
	*gadget.Gadget
	c *char
}

func (c *char) newDropblet() *sourcewaterDroplet {
	p := &sourcewaterDroplet{c: c}
	player := c.Core.Combat.Player()
	pos := geometry.CalcRandomPointFromCenter(
		geometry.CalcOffsetPoint(
			player.Pos(),
			geometry.Point{Y: 3.5},
			player.Direction(),
		),
		0.3,
		3,
		c.Core.Rand,
	)
	p.Gadget = gadget.New(c.Core, pos, 0.3, combat.GadgetTypSourcewaterDroplet)
	p.Gadget.Duration = 878

	return p
}

func (s *sourcewaterDroplet) HandleAttack(*combat.AttackEvent) float64 { return 0 }
func (s *sourcewaterDroplet) SetDirection(trg geometry.Point)          {}
func (s *sourcewaterDroplet) SetDirectionToClosestEnemy()              {}
func (s *sourcewaterDroplet) CalcTempDirection(trg geometry.Point) geometry.Point {
	return geometry.DefaultDirection()
}

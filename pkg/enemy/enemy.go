package enemy

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/target"
	"github.com/genshinsim/gcsim/pkg/reactable"
)

type Enemy struct {
	*target.Tmpl
}

func New(index int, c *core.Core, p core.EnemyProfile) *Enemy {
	e := &Enemy{}
	e.Tmpl = &target.Tmpl{}
	e.Reactable = &reactable.Reactable{}
	e.TargetIndex = index
	e.Level = p.Level
	e.Res = p.Resist
	e.Core = c
	e.HPMax = p.HP
	e.HPCurrent = p.HP

	e.Reactable.Init(e, c)
	e.Tmpl.Init(p.CoordX, p.CoordY, p.Size)

	return e
}

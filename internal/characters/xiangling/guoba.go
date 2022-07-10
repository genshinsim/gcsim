package xiangling

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/reactable"
	"github.com/genshinsim/gcsim/pkg/target"
)

type panda struct {
	*target.Target
	*reactable.Reactable
	pyroWindowStart int
	pyroWindowEnd   int
}

func newGuoba(c *core.Core) *panda {
	p := &panda{}
	p.Target = target.New(c, 0, 0, 0.5)
	p.Reactable = &reactable.Reactable{}
	p.Reactable.Init(p, c)

	p.Target.HPCurrent = 1
	p.Target.HPMax = 1

	return p
}

func (p *panda) Type() combat.TargettableType { return combat.TargettableObject }

func (p *panda) Attack(atk *combat.AttackEvent, evt glog.Event) (float64, bool) {
	//don't take damage, trigger swirl reaction only on sucrose E
	if p.Core.Player.Chars()[atk.Info.ActorIndex].Base.Key != keys.Sucrose {
		return 0, false
	}
	if atk.Info.AttackTag != combat.AttackTagElementalArt {
		return 0, false
	}
	//check pyro window
	if p.Core.F < p.pyroWindowStart || p.Core.F > p.pyroWindowEnd {
		return 0, false
	}

	//cheat a bit, set the durability just enough to match incoming sucrose E gauge
	p.Durability[attributes.Pyro] = 25
	p.React(atk)
	//wipe out the durability after
	p.Durability[attributes.Pyro] = 0

	return 0, false
}

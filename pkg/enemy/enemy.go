// Package enemy implements an enemey target
package enemy

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/targets"
	"github.com/genshinsim/gcsim/pkg/core/task"
	"github.com/genshinsim/gcsim/pkg/modifier"
	"github.com/genshinsim/gcsim/pkg/reactable"
	"github.com/genshinsim/gcsim/pkg/target"
)

type Enemy struct {
	*target.Target
	*reactable.Reactable

	Level   int
	resists map[attributes.Element]float64
	prof    info.EnemyProfile
	hp      float64
	maxhp   float64

	damageTaken       float64
	lastParticleDrop  int
	particleDropIndex int // for custom HP drops

	// mods
	mods []modifier.Mod

	// hitlag stuff
	timePassed   int
	frozenFrames int
	queue        *task.Handler
}

func New(core *core.Core, p info.EnemyProfile) *Enemy {
	e := &Enemy{}
	e.queue = task.New(&e.timePassed)
	e.Level = p.Level
	//TODO: do we need to clone this map isntead?
	e.resists = p.Resist
	//TODO: this is kinda redundant to keep both profile and lvl/resist
	e.prof = p
	e.Target = target.New(core, geometry.Point{X: p.Pos.X, Y: p.Pos.Y}, p.Pos.R)
	e.Reactable = &reactable.Reactable{}
	e.Reactable.Init(e, core)
	e.Reactable.FreezeResist = e.prof.FreezeResist
	e.mods = make([]modifier.Mod, 0, 10)
	if core.Combat.DamageMode {
		e.hp = p.HP
		e.maxhp = p.HP
	}
	return e
}

func (e *Enemy) Type() targets.TargettableType { return targets.TargettableEnemy }

func (e *Enemy) MaxHP() float64 { return e.maxhp }
func (e *Enemy) HP() float64    { return e.hp }
func (e *Enemy) Kill() {
	e.Alive = false
	if e.Key() == e.Core.Combat.DefaultTarget {
		player := e.Core.Combat.Player()
		// try setting default target to closest enemy to player if target died
		enemy := e.Core.Combat.ClosestEnemy(player.Pos())
		if enemy == nil {
			// all enemies dead, do nothing for now
			return
		}
		e.Core.Combat.DefaultTarget = enemy.Key()
		e.Core.Combat.Log.NewEvent("default target changed on enemy death", glog.LogWarnings, -1)
		player.SetDirection(enemy.Pos())
	}
}

func (e *Enemy) SetDirection(trg geometry.Point) {}
func (e *Enemy) SetDirectionToClosestEnemy()     {}
func (e *Enemy) CalcTempDirection(trg geometry.Point) geometry.Point {
	return geometry.DefaultDirection()
}

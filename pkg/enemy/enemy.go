// Package enemy implements an enemey target
package enemy

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/targets"
	"github.com/genshinsim/gcsim/pkg/modifier"
	"github.com/genshinsim/gcsim/pkg/queue"
	"github.com/genshinsim/gcsim/pkg/reactable"
	"github.com/genshinsim/gcsim/pkg/target"
)

type EnemyProfile struct {
	Level                 int                   `json:"level"`
	HP                    float64               `json:"hp"`
	Resist                attributes.ElementMap `json:"resist"`
	Pos                   core.Coord            `json:"-"`
	ParticleDropThreshold float64               `json:"particle_drop_threshold"` // drop particle every x dmg dealt
	ParticleDropCount     float64               `json:"particle_drop_count"`
	ParticleElement       attributes.Element    `json:"particle_element"`
}

func (e *EnemyProfile) Clone() EnemyProfile {
	r := *e
	r.Resist = make(map[attributes.Element]float64)
	for k, v := range e.Resist {
		r.Resist[k] = v
	}
	return r
}

type Enemy struct {
	*target.Target
	*reactable.Reactable

	Level   int
	resists map[attributes.Element]float64
	prof    EnemyProfile
	hp      float64
	maxhp   float64

	damageTaken      float64
	lastParticleDrop int

	//mods
	mods []modifier.Mod

	//hitlag stuff
	timePassed   int
	frozenFrames int
	queue        []queue.Task
}

func New(core *core.Core, p EnemyProfile) *Enemy {
	e := &Enemy{}
	e.Level = p.Level
	//TODO: do we need to clone this map isntead?
	e.resists = p.Resist
	//TODO: this is kinda redundant to keep both profile and lvl/resist
	e.prof = p
	e.Target = target.New(core, geometry.Point{X: p.Pos.X, Y: p.Pos.Y}, p.Pos.R)
	e.Reactable = &reactable.Reactable{}
	e.Reactable.Init(e, core)
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

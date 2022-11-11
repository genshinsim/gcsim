// Package enemy implements an enemey target
package enemy

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/modifier"
	"github.com/genshinsim/gcsim/pkg/queue"
	"github.com/genshinsim/gcsim/pkg/reactable"
	"github.com/genshinsim/gcsim/pkg/target"
)

type EnemyProfile struct {
	Level                 int                            `json:"level"`
	HP                    float64                        `json:"hp"`
	Resist                map[attributes.Element]float64 `json:"resist"`
	Pos                   core.Coord                     `json:"pos"`
	ParticleDropThreshold float64                        `json:"particleDropThreshold"` // drop particle every x dmg dealt
	ParticleDropCount     float64                        `json:"particleDropCount"`
	ParticleElement       attributes.Element             `json:"particleElement"`
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

	Level  int
	resist map[attributes.Element]float64
	prof   EnemyProfile
	hp     float64
	maxhp  float64

	damageTaken      float64
	lastParticleDrop int

	//mods
	mods []modifier.Mod

	//hitlag stuff
	timePassed   float64
	frozenFrames float64
	queue        []queue.Task
}

func New(core *core.Core, p EnemyProfile) *Enemy {
	e := &Enemy{}
	e.Level = p.Level
	//TODO: do we need to clone this map isntead?
	e.resist = p.Resist
	//TODO: this is kinda redundant to keep both profile and lvl/resist
	e.prof = p
	e.Target = target.New(core, p.Pos.X, p.Pos.Y, p.Pos.R)
	e.Reactable = &reactable.Reactable{}
	e.Reactable.Init(e, core)
	e.mods = make([]modifier.Mod, 0, 10)
	if core.Combat.DamageMode {
		e.hp = p.HP
		e.maxhp = p.HP
	}
	return e
}

func (e *Enemy) Type() combat.TargettableType { return combat.TargettableEnemy }

func (t *Enemy) MaxHP() float64 { return t.maxhp }
func (t *Enemy) HP() float64    { return t.hp }

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

const MaxTeamSize = 4

type EnemyProfile struct {
	Level                 int                            `json:"level"`
	HP                    float64                        `json:"-"`
	Resist                map[attributes.Element]float64 `json:"-"`
	Pos                   core.Coord                     `json:"-"`
	ParticleDropThreshold float64                        `json:"-"` // drop particle every x dmg dealt
	ParticleDropCount     float64                        `json:"-"`
	ParticleElement       attributes.Element             `json:"-"`
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

	damageTaken      float64
	lastParticleDrop int

	//mods
	mods []modifier.Mod

	//hitlag stuff
	timePassed   float64
	frozenFrames float64
	queue        []queue.Task

	//icd related
	icdTagOnTimer       [MaxTeamSize][combat.ICDTagLength]bool
	icdTagCounter       [MaxTeamSize][combat.ICDTagLength]int
	icdDamageTagOnTimer [MaxTeamSize][combat.ICDTagLength]bool
	icdDamageTagCounter [MaxTeamSize][combat.ICDTagLength]int
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
		e.Target.HPCurrent = p.HP
		e.Target.HPMax = p.HP
	}
	return e
}

func (e *Enemy) Type() combat.TargettableType { return combat.TargettableEnemy }

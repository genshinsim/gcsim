// Package enemy implements an enemey target
package enemy

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
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
	resistMods  []*resistMod
	defenseMods []*defenseMod

	//icd related
	icdTagOnTimer       [MaxTeamSize][combat.ICDTagLength]bool
	icdTagCounter       [MaxTeamSize][combat.ICDTagLength]int
	icdDamageTagOnTimer [MaxTeamSize][combat.ICDTagLength]bool
	icdDamageTagCounter [MaxTeamSize][combat.ICDTagLength]int
}

func New(core *core.Core, p EnemyProfile) *Enemy {
	e := &Enemy{}
	e.Target = target.New(core, p.Pos.X, p.Pos.Y, p.Pos.R)
	e.Reactable = &reactable.Reactable{}
	e.Reactable.Init(e, core)
	return e
}

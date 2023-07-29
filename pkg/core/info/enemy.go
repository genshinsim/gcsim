package info

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
)

type EnemyProfile struct {
	Level                 int                   `json:"level"`
	HP                    float64               `json:"hp"`
	Resist                attributes.ElementMap `json:"resist"`
	Pos                   Coord                 `json:"-"`
	ParticleDropThreshold float64               `json:"particle_drop_threshold"` // drop particle every x dmg dealt
	ParticleDropCount     float64               `json:"particle_drop_count"`
	ParticleElement       attributes.Element    `json:"particle_element"`
	FreezeResist          float64               `json:"freeze_resist"`
}

func (e *EnemyProfile) Clone() EnemyProfile {
	r := *e
	r.Resist = make(map[attributes.Element]float64)
	for k, v := range e.Resist {
		r.Resist[k] = v
	}
	return r
}

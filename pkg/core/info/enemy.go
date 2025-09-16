package info

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/model"
)

type EnemyProfile struct {
	Level                 int                    `json:"level"`
	HP                    float64                `json:"hp"`
	Resist                attributes.ElementMap  `json:"resist"`
	Pos                   Coord                  `json:"-"`
	ParticleDropThreshold float64                `json:"particle_drop_threshold"` // drop particle every x dmg dealt
	ParticleDropCount     float64                `json:"particle_drop_count"`
	ParticleElement       attributes.Element     `json:"particle_element"`
	FreezeResist          float64                `json:"freeze_resist"`
	ParticleDrops         []*model.MonsterHPDrop `json:"-"`
	HpBase                float64                `json:"-"`
	HpGrowCurve           model.MonsterCurveType `json:"-"`
	Id                    int                    `json:"-"`
	MonsterName           string                 `json:"monster_name"`
	Modified              bool                   `json:"modified"`
}

func (e *EnemyProfile) Clone() EnemyProfile {
	r := *e
	r.Resist = make(map[attributes.Element]float64)
	for k, v := range e.Resist {
		r.Resist[k] = v
	}
	return r
}

type Enemy interface {
	Target
	// hp related
	MaxHP() float64
	HP() float64
	// hitlag related
	ApplyHitlag(factor, dur float64)
	QueueEnemyTask(f func(), delay int)
	// modifier related
	// add
	AddStatus(key string, dur int, hitlag bool)
	AddResistMod(mod ResistMod)
	AddDefMod(mod DefMod)
	// delete
	DeleteStatus(key string)
	DeleteResistMod(key string)
	DeleteDefMod(key string)
	// active
	StatusIsActive(key string) bool
	ResistModIsActive(key string) bool
	DefModIsActive(key string) bool
	StatusExpiry(key string) int
}

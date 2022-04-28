package simulation

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/player"
	"github.com/genshinsim/gcsim/pkg/simulation/queue"
)

type SimulationConfig struct {
	//these settings relate to each simulation iteration
	Duration   int            `json:"duration"`
	DamageMode bool           `json:"damage_mode"`
	Targets    []EnemyProfile `json:"targets"`
	Characters struct {
		Initial player.CharKey            `json:"initial"`
		Profile []player.CharacterProfile `json:"profile"`
	} `json:"characters"`
	Rotation []queue.ActionBlock `json:"-"`
	Hurt     HurtEvent           `json:"-"`
	Energy   EnergyEvent         `json:"-"`
}

func (c *SimulationConfig) Clone() SimulationConfig {
	r := *c

	r.Targets = make([]EnemyProfile, len(c.Targets))
	for i, v := range c.Targets {
		r.Targets[i] = v.Clone()
	}

	r.Characters.Profile = make([]player.CharacterProfile, len(c.Characters.Profile))
	for i, v := range c.Characters.Profile {
		r.Characters.Profile[i] = v.Clone()
	}

	r.Rotation = make([]queue.ActionBlock, len(c.Rotation))
	for i, v := range c.Rotation {
		r.Rotation[i] = v.Clone()
	}

	return r
}

type EnemyProfile struct {
	Level          int                            `json:"level"`
	HP             float64                        `json:"-"`
	Resist         map[attributes.Element]float64 `json:"-"`
	Size           float64                        `json:"-"`
	CoordX, CoordY float64                        `json:"-"`
}

func (e *EnemyProfile) Clone() EnemyProfile {
	r := EnemyProfile{
		Level:  e.Level,
		Resist: make(map[attributes.Element]float64),
	}
	for k, v := range e.Resist {
		r.Resist[k] = v
	}
	return r
}

type EnergyEvent struct {
	Active    bool
	Once      bool //how often
	Start     int
	End       int
	Particles int
}

type HurtEvent struct {
	Active bool
	Once   bool //how often
	Start  int  //
	End    int
	Min    float64
	Max    float64
	Ele    attributes.Element
}

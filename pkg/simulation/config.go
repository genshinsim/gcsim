package simulation

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/simulation/queue"
)

type SimulationConfig struct {
	//these settings relate to each simulation iteration
	Duration    int                          `json:"duration"`
	DamageMode  bool                         `json:"damage_mode"`
	Targets     []enemy.EnemyProfile         `json:"targets"`
	PlayerPos   Pos                          `json:"player_initial_pos"`
	Characters  []character.CharacterProfile `json:"characters"`
	InitialChar keys.Char                    `json:"initial"`
	Rotation    []queue.ActionBlock          `json:"-"`
	Hurt        HurtEvent                    `json:"-"`
	Energy      EnergyEvent                  `json:"-"`
}

type Pos struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	R float64 `json:"r"`
}

func (c *SimulationConfig) Clone() SimulationConfig {
	r := *c

	r.Targets = make([]enemy.EnemyProfile, len(c.Targets))
	for i, v := range c.Targets {
		r.Targets[i] = v.Clone()
	}

	r.Characters = make([]character.CharacterProfile, len(c.Characters))
	for i, v := range c.Characters {
		r.Characters[i] = v.Clone()
	}

	r.Rotation = make([]queue.ActionBlock, len(c.Rotation))
	for i, v := range c.Rotation {
		r.Rotation[i] = v.Clone()
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

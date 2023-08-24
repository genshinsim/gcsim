package info

import (
	"encoding/json"
	"log"

	"github.com/genshinsim/gcsim/pkg/core/keys"
)

type ActionList struct {
	Targets     []EnemyProfile     `json:"targets"`
	PlayerPos   Coord              `json:"player_initial_pos"`
	Characters  []CharacterProfile `json:"characters"`
	InitialChar keys.Char          `json:"initial"`
	Energy      EnergySettings     `json:"energy_settings"`
	Settings    SimulatorSettings  `json:"settings"`
	Errors      []error            `json:"-"` //These represents errors preventing ActionList from being executed
	ErrorMsgs   []string           `json:"errors"`
}

type EnergySettings struct {
	Active         bool `json:"active"`
	Once           bool `json:"once"` //how often
	Start          int  `json:"start"`
	End            int  `json:"end"`
	Amount         int  `json:"amount"`
	LastEnergyDrop int  `json:"last_energy_drop"`
}

type SimulatorSettings struct {
	Duration     float64 `json:"-"`
	DamageMode   bool    `json:"damage_mode"`
	EnableHitlag bool    `json:"enable_hitlag"`
	DefHalt      bool    `json:"def_halt"` // for hitlag
	//other stuff
	NumberOfWorkers int    `json:"-"`          // how many workers to run the simulation
	Iterations      int    `json:"iterations"` // how many iterations to run
	Delays          Delays `json:"delays"`
}

type Delays struct {
	Skill  int `json:"skill"`
	Burst  int `json:"burst"`
	Attack int `json:"attack"`
	Charge int `json:"charge"`
	Aim    int `json:"aim"`
	Dash   int `json:"dash"`
	Jump   int `json:"jump"`
	Swap   int `json:"swap"`
}

func (c *ActionList) Copy() *ActionList {

	r := *c

	r.Targets = make([]EnemyProfile, len(c.Targets))
	for i, v := range c.Targets {
		r.Targets[i] = v.Clone()
	}

	r.Characters = make([]CharacterProfile, len(c.Characters))
	for i, v := range c.Characters {
		r.Characters[i] = v.Clone()
	}

	return &r
}

func (a *ActionList) PrettyPrint() string {
	prettyJson, err := json.MarshalIndent(a, "", "  ")
	if err != nil {
		log.Fatalf(err.Error())
	}
	return string(prettyJson)
}

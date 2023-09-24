package info

import (
	"encoding/json"
	"log"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/keys"
)

type ActionList struct {
	Targets          []EnemyProfile     `json:"targets"`
	InitialPlayerPos Coord              `json:"initial_player_pos"`
	Characters       []CharacterProfile `json:"characters"`
	InitialChar      keys.Char          `json:"initial_char"`
	EnergySettings   EnergySettings     `json:"energy_settings"`
	HurtSettings     HurtSettings       `json:"hurt_settings"`
	Settings         SimulatorSettings  `json:"settings"`
	Errors           []error            `json:"-"` // These represents errors preventing ActionList from being executed
	ErrorMsgs        []string           `json:"error_msgs"`
}

type EnergySettings struct {
	Active         bool `json:"active"`
	Once           bool `json:"once"` // how often
	Start          int  `json:"start"`
	End            int  `json:"end"`
	Amount         int  `json:"amount"`
	LastEnergyDrop int  `json:"last_energy_drop"`
}

type HurtSettings struct {
	Active   bool               `json:"active"`
	Once     bool               `json:"once"`
	Start    int                `json:"start"`
	End      int                `json:"end"`
	Min      float64            `json:"min"`
	Max      float64            `json:"max"`
	Element  attributes.Element `json:"element"`
	LastHurt int                `json:"last_hurt"`
}

type SimulatorSettings struct {
	Duration        float64 `json:"-"`
	DamageMode      bool    `json:"damage_mode"`
	EnableHitlag    bool    `json:"enable_hitlag"`
	DefHalt         bool    `json:"def_halt"` // for hitlag
	ErCalc          bool    `json:"er_calc"`
	ExpectedCritDmg bool    `json:"expected_dmg"`
	// other stuff
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

func (a *ActionList) Copy() *ActionList {
	r := *a

	r.Targets = make([]EnemyProfile, len(a.Targets))
	for i, v := range a.Targets {
		r.Targets[i] = v.Clone()
	}

	r.Characters = make([]CharacterProfile, len(a.Characters))
	for i := range a.Characters {
		r.Characters[i] = a.Characters[i].Clone()
	}

	return &r
}

func (a *ActionList) PrettyPrint() string {
	prettyJSON, err := json.MarshalIndent(a, "", "  ")
	if err != nil {
		log.Fatalf(err.Error())
	}
	return string(prettyJSON)
}

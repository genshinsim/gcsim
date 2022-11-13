package result

import (
	"fmt"
	"strings"

	"github.com/genshinsim/gcsim/pkg/agg"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/gcs/ast"
	"github.com/genshinsim/gcsim/pkg/simulation"
)

type Summary struct {
	// sim metadata
	SchemaVersion Version `json:"schema_version"`
	SimVersion    string  `json:"sim_version"`
	BuildDate     string  `json:"build_date"`
	Modified      bool    `json:"modified"`

	// character & enemy metadata
	InitialCharacter  string                       `json:"initial_character"`
	CharacterDetails  []simulation.CharacterDetail `json:"character_details"`
	TargetDetails     []enemy.EnemyProfile         `json:"target_details"`
	SimulatorSettings ast.SimulatorSettings        `json:"simulator_settings"`
	EnergySettings    ast.EnergySettings           `json:"energy_settings"`

	Config     string `json:"config_file"`
	SampleSeed string `json:"sample_seed"`

	// calculations/simulation data
	Statistics agg.Result `json:"statistics"`
}

type Version struct {
	Major int `json:"major"`
	Minor int `json:"minor"`
}

// TODO: repopulate or delete
func (r *Summary) PrettyPrint() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf(
		"Average duration of %.2f seconds (min: %.2f max: %.2f std: %.2f)\n",
		r.Statistics.Duration.Mean, r.Statistics.Duration.Min,
		r.Statistics.Duration.Max, r.Statistics.Duration.SD))
	sb.WriteString(fmt.Sprintf(
		"Average %.2f damage over %.2f seconds, resulting in %.0f dps (min: %.2f max: %.2f std: %.2f) \n",
		r.Statistics.TotalDamage.Mean, r.Statistics.Duration.Mean,
		r.Statistics.DPS.Mean, r.Statistics.DPS.Min, r.Statistics.DPS.Max, r.Statistics.DPS.SD))
	sb.WriteString(fmt.Sprintf(
		"Simulation completed %v iterations in %.3f seconds\n", r.Statistics.Iterations, r.Statistics.Runtime/1000000000))

	return sb.String()
}

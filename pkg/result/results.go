package result

import (
	"fmt"
	"strings"

	"github.com/genshinsim/gcsim/pkg/agg"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/gcs/ast"
	"github.com/genshinsim/gcsim/pkg/model"
	"github.com/genshinsim/gcsim/pkg/simulation"
)

type Summary struct {
	// sim metadata
	SchemaVersion Version `json:"schema_version"`
	SimVersion    string  `json:"sim_version"`
	BuildDate     string  `json:"build_date"`
	Modified      bool    `json:"modified"`
	Mode          int     `json:"mode"`

	// character & enemy metadata
	// TODO: TargetDetails as map. Need changes to how target keys work
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

// TODO: this function is incomplete. not all data copied
func (s *Summary) ToPBModel() *model.SimulationResult {
	r := &model.SimulationResult{
		SchemaVersion: &model.Version{
			Major: int64(s.SchemaVersion.Major),
			Minor: int64(s.SchemaVersion.Minor),
		},
		SimVersion:       s.SimVersion,
		BuildDate:        s.BuildDate,
		Modified:         s.Modified,
		InitialCharacter: s.InitialCharacter,
		Config:           s.Config,
		SampleSeed:       s.SampleSeed,
	}
	for _, v := range s.CharacterDetails {
		next := &model.Character{
			Key:      v.Name, //TODO: to be updated when we rekey characters
			Name:     v.Name,
			Element:  v.Element,
			Level:    int64(v.Level),
			MaxLevel: int64(v.MaxLevel),
			Cons:     int64(v.Cons),
			Weapon: &model.Weapon{
				Name:     v.Weapon.Name,
				Refine:   int64(v.Weapon.Refine),
				Level:    int64(v.Weapon.Level),
				MaxLevel: int64(v.MaxLevel),
			},
			Talents: &model.CharacterTalents{
				Attack: int64(v.Talents.Attack),
				Skill:  int64(v.Talents.Skill),
				Burst:  int64(v.Talents.Burst),
			},
		}
		next.Sets = make(map[string]int64)
		for k, x := range v.Sets {
			next.Sets[k] = int64(x)
		}
		next.Stats = make(map[string]float64)
		for i, x := range v.Stats {
			next.Stats[attributes.StatTypeString[i]] = x
		}
		next.Snapshot = make(map[string]float64)
		for i, x := range v.SnapshotStats {
			next.Snapshot[attributes.StatTypeString[i]] = x
		}
		r.CharacterDetails = append(r.CharacterDetails, next)
	}
	for _, v := range s.TargetDetails {
		next := &model.Enemy{
			Level: int64(v.Level),
			HP:    v.HP,
			Pos: &model.Coord{
				X: v.Pos.X,
				Y: v.Pos.Y,
				R: v.Pos.R,
			},
			ParticleDropThreshold: v.ParticleDropThreshold,
			ParticleDropCount:     v.ParticleDropCount,
			ParticleElement:       v.ParticleElement.String(),
		}
		next.Resist = make(map[string]float64)
		for k, x := range v.Resist {
			next.Resist[k.String()] = x
		}
		r.TargetDetails = append(r.TargetDetails, next)
	}

	r.Statistics = &model.SimulationStatistics{
		MinSeed: s.Statistics.MinSeed,
		MaxSeed: s.Statistics.MaxSeed,
		Duration: &model.OverviewStats{
			Min:  s.Statistics.Duration.Min,
			Max:  s.Statistics.Duration.Max,
			Mean: s.Statistics.Duration.Mean,
		},
		Runtime:    s.Statistics.Runtime,
		Iterations: int64(s.Statistics.Iterations),
	}

	return r
}

func (s *Summary) ToPBDBEntry() *model.DBEntry {
	r := &model.DBEntry{
		SimDuration: &model.DescriptiveStats{
			Min:  s.Statistics.Duration.Min,
			Max:  s.Statistics.Duration.Max,
			Mean: s.Statistics.Duration.Mean,
			SD:   s.Statistics.Duration.SD,
		},
		TotalDamage: &model.DescriptiveStats{
			Min:  s.Statistics.TotalDamage.Min,
			Max:  s.Statistics.TotalDamage.Max,
			Mean: s.Statistics.TotalDamage.Mean,
			SD:   s.Statistics.TotalDamage.SD,
		},
		TargetCount:      int32(len(s.TargetDetails)),
		Hash:             s.SimVersion,
		Config:           s.Config,
		MeanDpsPerTarget: s.Statistics.TotalDamage.Mean / (float64(len(s.TargetDetails)) * s.Statistics.Duration.Mean),
	}
	if s.Mode == 1 {
		r.Mode = model.SimMode_TTK_MODE
	}
	for _, v := range s.CharacterDetails {
		next := &model.Character{
			Key:      v.Name, //TODO: to be updated when we rekey characters
			Name:     v.Name,
			Element:  v.Element,
			Level:    int64(v.Level),
			MaxLevel: int64(v.MaxLevel),
			Cons:     int64(v.Cons),
			Weapon: &model.Weapon{
				Name:     v.Weapon.Name,
				Refine:   int64(v.Weapon.Refine),
				Level:    int64(v.Weapon.Level),
				MaxLevel: int64(v.MaxLevel),
			},
			Talents: &model.CharacterTalents{
				Attack: int64(v.Talents.Attack),
				Skill:  int64(v.Talents.Skill),
				Burst:  int64(v.Talents.Burst),
			},
		}
		next.Sets = make(map[string]int64)
		for k, x := range v.Sets {
			next.Sets[k] = int64(x)
		}
		next.Stats = make(map[string]float64)
		for i, x := range v.Stats {
			next.Stats[attributes.StatTypeString[i]] = x
		}
		next.Snapshot = make(map[string]float64)
		for i, x := range v.SnapshotStats {
			next.Snapshot[attributes.StatTypeString[i]] = x
		}
		r.Team = append(r.Team, next)
		r.CharNames = append(r.CharNames, v.Name)
	}

	return r
}

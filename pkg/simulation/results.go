package simulation

import "github.com/genshinsim/gcsim/pkg/core"

type Result struct {
	IsDamageMode          bool                      `json:"is_damage_mode"`
	CharNames             []string                  `json:"char_names"`
	CharDetails           []CharDetail              `json:"char_details"`
	DamageByChar          []map[string]float64      `json:"damage_by_char"`
	DamageInstancesByChar []map[string]int          `json:"damage_instances_by_char"`
	DamageByCharByTargets []map[int]float64         `json:"damage_by_char_by_targets"`
	DamageDetailByTime    map[DamageDetails]float64 `json:"-"`
	CharActiveTime        []int                     `json:"char_active_time"`
	AbilUsageCountByChar  []map[string]int          `json:"abil_usage_count_by_char"`
	ParticleCount         map[string]int            `json:"particle_count"`
	ReactionsTriggered    map[core.ReactionType]int `json:"reactions_triggered"`
	Duration              int                       `json:"sim_duration"`
	ElementUptime         []map[core.EleType]int    `json:"ele_uptime"`
	EnergyWhenBurst       [][]float64               `json:"energy_when_burst"`
	//final result
	Damage float64 `json:"damage"`
	DPS    float64 `json:"dps"`
	//for tracking min/max run
	Seed int64 `json:"seed"`
}

type CharDetail struct {
	Name          string         `json:"name"`
	Element       string         `json:"element"`
	Level         int            `json:"level"`
	MaxLevel      int            `json:"max_level"`
	Cons          int            `json:"cons"`
	Weapon        WeaponDetail   `json:"weapon"`
	Talents       TalentDetail   `json:"talents"`
	Sets          map[string]int `json:"sets"`
	Stats         []float64      `json:"stats"`
	SnapshotStats []float64      `json:"snapshot"`
}

type DamageDetails struct {
	FrameBucket int
	Char        int
	Target      int
}

type WeaponDetail struct {
	Name     string `json:"name"`
	Refine   int    `json:"refine"`
	Level    int    `json:"level"`
	MaxLevel int    `json:"max_level"`
}

type TalentDetail struct {
	Attack int `json:"attack"`
	Skill  int `json:"skill"`
	Burst  int `json:"burst"`
}

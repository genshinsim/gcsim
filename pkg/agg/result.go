package agg

type Result struct {
	// metadata
	MinSeed    string  `json:"min_seed"`
	MaxSeed    string  `json:"max_seed"`
	P25Seed    string  `json:"p25_seed"`
	P50Seed    string  `json:"p50_seed"`
	P75Seed    string  `json:"p75_seed"`
	Runtime    float64 `json:"runtime"`
	Iterations int     `json:"iterations"`

	// global overview (global/no group by)
	Duration    SummaryStat `json:"duration"`
	DPS         SummaryStat `json:"dps"`
	RPS         SummaryStat `json:"rps"`
	EPS         SummaryStat `json:"eps"`
	HPS         SummaryStat `json:"hps"`
	SHP         SummaryStat `json:"shp"`
	TotalDamage SummaryStat `json:"total_damage"`

	Warnings      Warnings        `json:"warnings"`
	FailedActions []FailedActions `json:"failed_actions"`

	// damage
	ElementDPS            map[string]FloatStat   `json:"element_dps"`
	TargetDPS             map[int]FloatStat      `json:"target_dps"`
	CharacterDPS          []FloatStat            `json:"character_dps"`
	BreakdownByElementDPS []map[string]FloatStat `json:"dps_by_element"`
	BreakdownByTargetDPS  []map[int]FloatStat    `json:"dps_by_target"`

	DamageBucketSize             int           `json:"damage_bucket_size"`
	DamageBuckets                []FloatStat   `json:"damage_buckets"`
	CumulativeDamageContribution [][]FloatStat `json:"cumu_damage_contrib"`

	// shield
	Shields map[string]ShieldInfo `json:"shields"`
}

// TODO: OverviewResult w/ Histogram data for distribution graphs?
type SummaryStat struct {
	Min  float64 `json:"min"`
	Max  float64 `json:"max"`
	Mean float64 `json:"mean"`
	SD   float64 `json:"sd"`

	// Only use if necessary.
	// w/o quantile can use StreamStats
	// O(1) vs O(n) space for stream vs sample
	// O(1) vs O(nlogn) time for stream vs sample
	Q1 float64 `json:"q1"`
	Q2 float64 `json:"q2"`
	Q3 float64 `json:"q3"`

	// requires storing all iteration values in buffer to properly compute
	// lower bound inclusive, upper exclusive (except on last bin)
	Hist []uint `json:"histogram,omitempty"`
}

type FloatStat struct {
	Min  float64 `json:"min"`
	Max  float64 `json:"max"`
	Mean float64 `json:"mean"`
	SD   float64 `json:"sd"`
}

type Warnings struct {
	TargetOverlap       bool `json:"target_overlap"`
	InsufficientEnergy  bool `json:"insufficient_energy"`
	InsufficientStamina bool `json:"insufficient_stamina"`
	SwapCD              bool `json:"swap_cd"`
	SkillCD             bool `json:"skill_cd"`
}

// TODO: add more statistics?
type FailedActions struct {
	InsufficientEnergy  FloatStat `json:"insufficient_energy"`
	InsufficientStamina FloatStat `json:"insufficient_stamina"`
	SwapCD              FloatStat `json:"swap_cd"`
	SkillCD             FloatStat `json:"skill_cd"`
}

type ShieldInfo struct {
	HP     map[string]FloatStat `json:"hp"`
	Uptime FloatStat            `json:"uptime"`
}

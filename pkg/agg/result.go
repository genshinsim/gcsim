package agg

type Result struct {
	// metadata
	MinSeed  int64     `json:"min_seed"`
	MaxSeed  int64     `json:"max_seed"`
	Duration FloatStat `json:"duration"`

	// overview
	TotalDamage FloatStat `json:"total_damage"`
	DPS         FloatStat `json:"dps"`
}

type FloatStat struct {
	Min  float64 `json:"min,omitempty"`
	Max  float64 `json:"max,omitempty"`
	Mean float64 `json:"mean,omitempty"`
	SD   float64 `json:"sd,omitempty"`
}

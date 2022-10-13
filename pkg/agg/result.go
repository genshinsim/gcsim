package agg

type Result struct {
	// metadata
	MinSeed uint64 `json:"min_seed"`
	MaxSeed uint64 `json:"max_seed"`

	// global overview (global/no group by)
	Duration    FloatStat `json:"duration"`
	DPS         FloatStat `json:"dps"`
	RPS         FloatStat `json:"rps"`
	EPS         FloatStat `json:"eps"`
	HPS         FloatStat `json:"hps"`
	SPS         FloatStat `json:"sps"`
	TotalDamage FloatStat `json:"total_damage"`
}

// TODO: OverviewResult w/ Histogram data for distribution graphs?
type FloatStat struct {
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
}

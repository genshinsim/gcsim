package agg

type Result struct {
	// metadata
	MinSeed  uint64    `json:"min_seed"`
	MaxSeed  uint64    `json:"max_seed"`
	Duration FloatStat `json:"duration"`

	// overview
	TotalDamage FloatStat `json:"total_damage"`
	DPS         FloatStat `json:"dps"`

	// legacy. needs to be removed
	Legacy LegacyStats
}

// TODO: remove. just backporting of the old stats collection for supporting current UI and tools
type LegacyStats struct {
	DamageByChar          []map[string]FloatStat
	DamageInstancesByChar []map[string]IntStat
	DamageByCharByTargets []map[int]FloatStat
	CharActiveTime        []IntStat
	AbilUsageCountByChar  []map[string]IntStat
	ParticleCount         map[string]FloatStat
	ReactionsTriggered    map[string]IntStat
	ElementUptime         []map[string]IntStat
	DPSByTarget           map[int]FloatStat
	DamageOverTime        map[string]FloatStat
}

type FloatStat struct {
	Min  float64 `json:"min"`
	Max  float64 `json:"max"`
	Mean float64 `json:"mean"`
	SD   float64 `json:"sd"`
}

type IntStat struct {
	Min  int     `json:"min"`
	Max  int     `json:"max"`
	Mean float64 `json:"mean"`
	SD   float64 `json:"sd"`
}

package stats

//go:generate msgp

// NOTE: all maps MUST use a string key. This is a requirement of the MessagePack spec
//   any string type aliasing must be defined in this module or msgp will not know that the type
//   is a string

type FieldStatus string
type ReactionModifier string

const (
	OnField  FieldStatus = "on_field"
	OffField FieldStatus = "off_field"

	Melt      ReactionModifier = "melt"
	Vaporize  ReactionModifier = "vaporize"
	Spread    ReactionModifier = "spread"
	Aggravate ReactionModifier = "aggravate"
)

type Result struct {
	Seed          uint64    `json:"seed"           msg:"seed"`
	Duration      int       `json:"duration"       msg:"duration"`
	TotalDamage   float64   `json:"total_damage"   msg:"total_damage"`
	DPS           float64   `json:"dps"            msg:"dps"`
	DamageBuckets []float64 `json:"damage_buckets" msg:"damage_buckets"`

	ActiveCharacters []ActiveCharacterInterval `json:"active_characters" msg:"active_characters"`
	DamageMitigation []float64                 `json:"damage_mitigation" msg:"damage_mitigation"`
	ShieldResults    ShieldResult              `json:"shield_results"    msg:"shield_results"`

	Characters    []CharacterResult `json:"characters"     msg:"characters"`
	Enemies       []EnemyResult     `json:"enemies"        msg:"enemies"`
	TargetOverlap bool              `json:"target_overlap" msg:"target_overlap"`
}

type CharacterResult struct {
	// For raw data usage outside of gcsim only
	Name string `json:"name" msg:"name"`

	DamageEvents   []DamageEvent   `json:"damage_events"   msg:"damage_events"`
	ReactionEvents []ReactionEvent `json:"reaction_events" msg:"reaction_events"`
	ActionEvents   []ActionEvent   `json:"action_events"   msg:"action_events"`
	EnergyEvents   []EnergyEvent   `json:"energy_events"   msg:"energy_events"`
	HealEvents     []HealEvent     `json:"heal_events"     msg:"heal_events"`

	// TODO: Move to Result since only active character can perform actions?
	FailedActions []ActionFailInterval `json:"failed_actions" msg:"failed_actions"`

	EnergyStatus []float64 `json:"energy_status" msg:"energy_status"` // can be completely replaced by EnergyEvents?
	HealthStatus []float64 `json:"health_status" msg:"health_status"`

	DamageCumulativeContrib []float64 `json:"damage_cumulative_contrib" msg:"damage_cumulative_contrib"`

	ActiveTime  int     `json:"active_time"  msg:"active_time"`
	EnergySpent float64 `json:"energy_spent" msg:"energy_spent"`
	ErNeeded    float64 `json:"er_needed"    msg:"er_needed"`
	WeightedER  float64 `json:"weighted_er"  msg:"weight_er"`
}

type EnemyResult struct {
	ReactionStatus []ReactionStatusInterval `json:"reaction_status" msg:"reaction_status"`
	ReactionUptime map[string]int           `json:"reaction_uptime" msg:"reaction_uptime"` // can calculate from intervals?
}

type ShieldResult struct {
	Shields         []ShieldStats                     `json:"shields"          msg:"shields"`
	EffectiveShield map[string][]ShieldSingleInterval `json:"effective_shield" msg:"effective_shield"`
}

type DamageEvent struct {
	Frame              int              `json:"frame"               msg:"frame"`
	ActionId           int              `json:"action_id"           msg:"action_id"`
	Source             string           `json:"source"              msg:"source"`
	Target             int              `json:"target"              msg:"target"`
	Element            string           `json:"element"             msg:"element"`
	ReactionModifier   ReactionModifier `json:"reaction_modifier"   msg:"reaction_modifier"`
	Crit               bool             `json:"crit"                msg:"crit"`
	Modifiers          []string         `json:"modifiers"           msg:"modifiers"`
	MitigationModifier float64          `json:"mitigation_modifier" msg:"mitigation_modifier"`
	Damage             float64          `json:"damage"              msg:"damage"`
}

type ActionEvent struct {
	Frame    int    `json:"frame"     msg:"frame"`
	ActionId int    `json:"action_id" msg:"action_id"`
	Action   string `json:"action"    msg:"action"`
}

type ReactionEvent struct {
	Frame    int    `json:"frame"    msg:"frame"`
	Source   string `json:"source"   msg:"source"`
	Target   int    `json:"target"   msg:"target"`
	Reaction string `json:"reaction" msg:"reaction"`
}

type EnergyEvent struct {
	Frame       int         `json:"frame"        msg:"frame"`
	Source      string      `json:"source"       msg:"source"`
	FieldStatus FieldStatus `json:"field_status" msg:"field_status"`
	Gained      float64     `json:"gained"       msg:"gained"`
	Wasted      float64     `json:"wasted"       msg:"wasted"`
	Current     float64     `json:"current"      msg:"current"` // this is pre + gained
}

// Heal events are stored in the source character
type HealEvent struct {
	Frame  int     `json:"frame"  msg:"frame"`
	Source string  `json:"source" msg:"source"`
	Target int     `json:"target" msg:"target"`
	Heal   float64 `json:"heal"   msg:"heal"`
}

type ActionFailInterval struct {
	Start  int    `json:"start"  msg:"start"`
	End    int    `json:"end"    msg:"end"`
	Reason string `json:"reason" msg:"reason"`
}

type ReactionStatusInterval struct {
	Start int    `json:"start" msg:"start"`
	End   int    `json:"end"   msg:"end"`
	Type  string `json:"type"  msg:"type"`
}

type ActiveCharacterInterval struct {
	Start     int `json:"start"     msg:"start"`
	End       int `json:"end"       msg:"end"`
	Character int `json:"character" msg:"character"`
}

type ShieldStats struct {
	Name      string           `json:"name"      msg:"name"`
	Intervals []ShieldInterval `json:"intervals" msg:"intervals"`
}

type ShieldInterval struct {
	Start int                `json:"start" msg:"start"`
	End   int                `json:"end"   msg:"end"`
	HP    map[string]float64 `json:"hp"    msg:"hp"`
}

type ShieldSingleInterval struct {
	Start int     `json:"start" msg:"start"`
	End   int     `json:"end"   msg:"end"`
	HP    float64 `json:"hp"    msg:"hp"`
}

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

	Amp15     ReactionModifier = "amp_1_5"
	Amp20     ReactionModifier = "amp_2_0"
	Spread    ReactionModifier = "spread"
	Aggravate ReactionModifier = "aggravate"
)

type Result struct {
	Seed          uint64    `msg:"seed" json:"seed"`
	Duration      int       `msg:"duration" json:"sim_duration"`
	TotalDamage   float64   `msg:"total_damage" json:"total_damage"`
	DPS           float64   `msg:"dps" json:"dps"`
	DamageBuckets []float64 `msg:"damage_buckets" json:"damage_buckets"`

	ActiveCharacters []ActiveCharacterInterval `msg:"active_characters" json:"active_characters"`
	DamageMitigation []float64                 `msg:"damage_mitigation" json:"damage_mitigation"`
	ShieldResults    ShieldResult              `msg:"shield_results" json:"shield_results"`

	Characters    []CharacterResult `msg:"characters" json:"characters"`
	Enemies       []EnemyResult     `msg:"enemies" json:"enemies"`
	TargetOverlap bool              `msg:"target_overlap" json:"target_overlap"`
}

type CharacterResult struct {
	// For raw data usage outside of gcsim only
	Name string `msg:"name" json:"name"`

	DamageEvents   []DamageEvent   `msg:"damage_events" json:"damage_events"`
	ReactionEvents []ReactionEvent `msg:"reaction_events" json:"reaction_events"`
	ActionEvents   []ActionEvent   `msg:"action_events" json:"action_events"`
	EnergyEvents   []EnergyEvent   `msg:"energy_events" json:"energy_events"`
	HealEvents     []HealEvent     `msg:"heal_events" json:"heal_events"`

	// TODO: Move to Result since only active character can perform actions?
	FailedActions []ActionFailInterval `msg:"failed_actions" json:"failed_actions"`

	EnergyStatus []float64 `msg:"energy" json:"energy"` // can be completely replaced by EnergyEvents?
	HealthStatus []float64 `msg:"health" json:"health"`

	DamageCumulativeContrib []float64 `msg:"damage_cumu_contrib" json:"damage_cumu_contrib"`

	ActiveTime  int     `msg:"active_time" json:"active_time"`
	EnergySpent float64 `msg:"energy_spent" json:"energy_spent"`
}

type EnemyResult struct {
	ReactionStatus []ReactionStatusInterval `msg:"reaction_status" json:"reaction_status"`
	ReactionUptime map[string]int           `msg:"reaction_uptime" json:"reaction_uptime"` // can calculate from intervals?
}

type ShieldResult struct {
	Shields         []ShieldStats                     `msg:"shields" json:"shields"`
	EffectiveShield map[string][]ShieldSingleInterval `msg:"effective_shield" json:"effective_shield"`
}

type DamageEvent struct {
	Frame            int              `msg:"frame" json:"frame"`
	ActionId         int              `msg:"action_id" json:"action_id"`
	Source           string           `msg:"src" json:"src"`
	Target           int              `msg:"target" json:"target"`
	Element          string           `msg:"element" json:"element"`
	ReactionModifier ReactionModifier `msg:"reaction_modifier" json:"reaction_modifier"`
	Crit             bool             `msg:"crit" json:"crit"`
	Modifiers        []string         `msg:"modifiers" json:"modifiers"`
	Mitigation       float64          `msg:"mitigation_modifier" json:"mitigation_modifier"`
	Damage           float64          `msg:"damage" json:"damage"`
}

type ActionEvent struct {
	Frame    int    `msg:"frame" json:"frame"`
	ActionId int    `msg:"action_id" json:"action_id"`
	Action   string `msg:"action" json:"action"`
}

type ReactionEvent struct {
	Frame    int    `msg:"frame" json:"frame"`
	Source   string `msg:"src" json:"src"`
	Target   int    `msg:"target" json:"target"`
	Reaction string `msg:"reaction" json:"reaction"`
}

type EnergyEvent struct {
	Frame   int         `msg:"frame" json:"frame"`
	Source  string      `msg:"src" json:"src"`
	Status  FieldStatus `msg:"field_status" json:"field_status"`
	Gained  float64     `msg:"gained" json:"gained"`
	Wasted  float64     `msg:"wasted" json:"wasted"`
	Current float64     `msg:"current" json:"current"` // this is pre + gained
}

// Heal events are stored in the source character
type HealEvent struct {
	Frame  int     `msg:"frame" json:"frame"`
	Source string  `msg:"src" json:"src"`
	Target int     `msg:"target" json:"target"`
	Heal   float64 `msg:"heal" json:"heal"`
}

type ActionFailInterval struct {
	Start  int    `msg:"start" json:"start"`
	End    int    `msg:"end" json:"end"`
	Reason string `msg:"reason" json:"reason"`
}

type ReactionStatusInterval struct {
	Start int    `msg:"start" json:"start"`
	End   int    `msg:"end" json:"end"`
	Type  string `msg:"type" json:"type"`
}

type ActiveCharacterInterval struct {
	Start     int `msg:"start" json:"start"`
	End       int `msg:"end" json:"end"`
	Character int `msg:"character" json:"character"`
}

type ShieldStats struct {
	Name      string           `msg:"name" json:"name"`
	Intervals []ShieldInterval `msg:"intervals" json:"intervals"`
}

type ShieldInterval struct {
	Start int                `msg:"start" json:"start"`
	End   int                `msg:"end" json:"end"`
	HP    map[string]float64 `msg:"hp" json:"hp"`
}

type ShieldSingleInterval struct {
	Start int     `msg:"start" json:"start"`
	End   int     `msg:"end" json:"end"`
	HP    float64 `msg:"hp" json:"hp"`
}

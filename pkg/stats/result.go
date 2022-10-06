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
	Seed        uint64  `msg:"seed"`
	Duration    int     `msg:"duration"`
	TotalDamage float64 `msg:"total_damage"`
	DPS         float64 `msg:"dps"`

	// TODO: Remove. just here for backwards compatibility
	Legacy LegacyResult `msg:"legacy"`

	ActiveCharacters []ActiveCharacterInterval `msg:"active_characters"`
	Shields          []ShieldInterval          `msg:"shields"`
	DamageMitigation []float64                 `msg:"damage_mitigation"`

	Characters []CharacterResult `msg:"characters"`
	Enemies    []EnemyResult     `msg:"enemies"`
}

type LegacyResult struct {
	DamageOverTime        map[string]float64   `msg:"damage_over_time"`
	DamageByChar          []map[string]float64 `msg:"damage_by_char"`
	DamageByCharByTargets []map[string]float64 `msg:"damage_by_char_by_targets"`
	DamageInstancesByChar []map[string]int     `msg:"damage_instances_by_char"`
	AbilUsageCountByChar  []map[string]int     `msg:"abil_usage_count_by_char"`
	CharActiveTime        []int                `msg:"char_active_time"`
	ElementUptime         []map[string]int     `msg:"element_uptime"`
	ParticleCount         map[string]float64   `msg:"particle_count"`
	ReactionsTriggered    map[string]int       `msg:"reactions_triggered"`
}

type CharacterResult struct {
	// For raw data usage outside of gcsim only
	Name string `msg:"name"`

	DamageEvents   []DamageEvent   `msg:"damage_events"`
	ReactionEvents []ReactionEvent `msg:"reaction_events"`
	ActionEvents   []ActionEvent   `msg:"action_events"`
	EnergyEvents   []EnergyEvent   `msg:"energy_events"`
	HealEvents     []HealEvent     `msg:"heal_events"`

	// TODO: Move to Result since only active character can perform actions?
	FailedActions []ActionFailInterval `msg:"failed_actions"`

	EnergyStatus []float64 `msg:"energy"` // can be completely replaced by EnergyEvents?
	HealthStatus []float64 `msg:"health"`

	ActiveTime  int     `msg:"active_time"`
	EnergySpent float64 `msg:"energy_spent"`
}

type EnemyResult struct {
	ReactionStatus []ReactionStatusInterval `msg:"reaction_status"`
	ReactionUptime map[string]int           `msg:"reaction_uptime"` // can calculate from intervals?
}

type DamageEvent struct {
	Frame            int              `msg:"frame"`
	ActionId         int              `msg:"action_id"`
	Source           string           `msg:"src"`
	Target           int              `msg:"target"`
	Element          string           `msg:"element"`
	ReactionModifier ReactionModifier `msg:"reaction_modifier"`
	Crit             bool             `msg:"crit"`
	Modifiers        []string         `msg:"modifiers"`
	Mitigation       float64          `msg:"mitigation_modifier"`
	Damage           float64          `msg:"damage"`
}

type ActionEvent struct {
	Frame    int    `msg:"frame"`
	ActionId int    `msg:"action_id"`
	Action   string `msg:"action"`
}

type ReactionEvent struct {
	Frame    int    `msg:"frame"`
	Source   string `msg:"src"`
	Target   int    `msg:"target"`
	Reaction string `msg:"reaction"`
}

type EnergyEvent struct {
	Frame   int         `msg:"frame"`
	Source  string      `msg:"src"`
	Status  FieldStatus `msg:"field_status"`
	Gained  float64     `msg:"gained"`
	Wasted  float64     `msg:"wasted"`
	Current float64     `msg:"current"` // this is pre + gained
}

// Heal events are stored in the source character
type HealEvent struct {
	Frame  int     `msg:"frame"`
	Source string  `msg:"src"`
	Target int     `msg:"target"`
	Heal   float64 `msg:"heal"`
}

type ActionFailInterval struct {
	Start  int    `msg:"start"`
	End    int    `msg:"end"`
	Reason string `msg:"reason"`
}

type ReactionStatusInterval struct {
	Start int    `msg:"start"`
	End   int    `msg:"end"`
	Type  string `msg:"type"`
}

type ActiveCharacterInterval struct {
	Start     int `msg:"start"`
	End       int `msg:"end"`
	Character int `msg:"character"`
}

type ShieldInterval struct {
	Start      int     `msg:"start"`
	End        int     `msg:"end"`
	Name       string  `msg:"name"`
	Absorption float64 `msg:"absoption"`
}

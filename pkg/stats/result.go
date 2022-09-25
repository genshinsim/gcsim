package stats

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/reactable"
)

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
	Seed        int64   `json:"seed"`
	Duration    int     `json:"duration"`
	TotalDamage float64 `json:"total_damage"`
	DPS         float64 `json:"dps"`
	/*
	 * TODO: evaluate HPS, eHPS, TotalHealing, etc (Overview/Summary struct?)
	 */

	ActiveCharacters []ActiveCharacterInterval `json:"active_characters"`
	Shields          []ShieldInterval          `json:"shields"`
	DamageMitigation []float64                 `json:"damage_mitigation"`

	Characters []CharacterResult `json:"characters"`
	Enemies    []EnemyResult     `json:"enemies"`
}

type CharacterResult struct {
	// For raw data usage outside of gcsim only
	Name string `json:"name"`

	DamageEvents   []DamageEvent   `json:"damage_events"`
	ReactionEvents []ReactionEvent `json:"reaction_events"`
	ActionEvents   []ActionEvent   `json:"action_events"`
	EnergyEvents   []EnergyEvent   `json:"energy_events"`
	HealEvents     []HealEvent     `json:"heal_events"`

	// TODO: Move to Result since only active character can perform actions?
	FailedActions []ActionFailInterval `json:"failed_actions"`

	EnergyStatus []float64 `json:"energy"` // can be completely replaced by EnergyEvents?
	HealthStatus []float64 `json:"health"`

	ActiveTime  int     `json:"active_time"`
	EnergySpent float64 `json:"energy_spent"`
}

type EnemyResult struct {
	ReactionStatus []ReactionStatusInterval `json:"reaction_status"`
	ReactionUptime map[string]int           `json:"reaction_uptime"` // can calculate from intervals?
}

type DamageEvent struct {
	Frame            int                `json:"frame"`
	ActionId         int                `json:"action_id"`
	Source           string             `json:"src"`
	Target           int                `json:"target"`
	Element          attributes.Element `json:"element"`
	ReactionModifier ReactionModifier   `json:"reaction_modifier"`
	Crit             bool               `json:"crit"`
	Modifiers        []string           `json:"modifiers"`
	Mitigation       float64            `json:"mitigation_modifier"`
	Damage           float64            `json:"damage"`
}

type ActionEvent struct {
	Frame    int           `json:"frame"`
	ActionId int           `json:"action_id"`
	Action   action.Action `json:"action"`
}

type ReactionEvent struct {
	Frame    int                 `json:"frame"`
	Source   string              `json:"src"`
	Target   int                 `json:"target"`
	Reaction combat.ReactionType `json:"reaction"`
}

type EnergyEvent struct {
	Frame   int         `json:"frame"`
	Source  string      `json:"src"`
	Status  FieldStatus `json:"field_status"`
	Gained  float64     `json:"gained"`
	Wasted  float64     `json:"wasted"`
	Current float64     `json:"current"` // this is pre + gained
}

// Heal events are stored in the source character
type HealEvent struct {
	Frame  int     `json:"frame"`
	Source string  `json:"src"`
	Target int     `json:"target"`
	Heal   float64 `json:"heal"`
}

type ActionFailInterval struct {
	Start  int                  `json:"start"`
	End    int                  `json:"end"`
	Reason action.ActionFailure `json:"reason"`
}

type ReactionStatusInterval struct {
	Start int                         `json:"start"`
	End   int                         `json:"end"`
	Type  reactable.ReactableModifier `json:"type"`
}

type ActiveCharacterInterval struct {
	Start     int `json:"start"`
	End       int `json:"end"`
	Character int `json:"character"`
}

type ShieldInterval struct {
	Start      int     `json:"start"`
	End        int     `json:"end"`
	Name       string  `json:"name"`
	Absorption float64 `json:"absoption"`
}

package character

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

type CharWrapper struct {
	Character

	//base characteristics
	Base     CharacterBase
	Weapon   WeaponProfile
	Talents  TalentProfile
	CharZone ZoneType

	//current status
	Energy    float64
	EnergyMax float64
	HPCurrent float64

	//tags
	Tags map[string]int
}

func (m *CharWrapper) Tag(key string) int {
	return m.Tags[key]
}

func (m *CharWrapper) AddTag(key string, val int) {
	m.Tags[key] = val
}

func (m *CharWrapper) RemoveTag(key string) {
	delete(m.Tags, key)
}

type Character interface {
	Attack(p map[string]int) action.ActionInfo
	Aimed(p map[string]int) action.ActionInfo
	ChargeAttack(p map[string]int) action.ActionInfo
	HighPlungeAttack(p map[string]int) action.ActionInfo
	LowPlungeAttack(p map[string]int) action.ActionInfo
	Skill(p map[string]int) action.ActionInfo
	Burst(p map[string]int) action.ActionInfo
	Dash(p map[string]int) action.ActionInfo

	ActionReady(a action.Action, p map[string]int) bool
	ActionStam(a action.Action, p map[string]int) float64

	SetCD(a action.Action, dur int)
	Cooldown(a action.Action) int
	ResetActionCooldown(a action.Action)
	ReduceActionCooldown(a action.Action, v int)
	Charges(a action.Action) int

	Snapshot(a *combat.AttackInfo) combat.Snapshot
}

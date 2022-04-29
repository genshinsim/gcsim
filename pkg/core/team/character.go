package team

import "github.com/genshinsim/gcsim/pkg/core/combat"

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
	Tag map[string]int
}

type Character interface {
	Attack(p map[string]int) ActionInfo
	Aimed(p map[string]int) ActionInfo
	ChargeAttack(p map[string]int) ActionInfo
	HighPlungeAttack(p map[string]int) ActionInfo
	LowPlungeAttack(p map[string]int) ActionInfo
	Skill(p map[string]int) ActionInfo
	Burst(p map[string]int) ActionInfo
	Dash(p map[string]int) ActionInfo

	ActionReady(a Action, p map[string]int) bool
	ActionStam(a Action, p map[string]int) float64

	SetCD(a Action, dur int)
	Cooldown(a Action) int
	ResetActionCooldown(a Action)
	ReduceActionCooldown(a Action, v int)
	Charges(Action) int

	Snapshot(a *combat.AttackInfo) combat.Snapshot
}

package character

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

type CharWrapper struct {
	Index int
	f     *int //current frame
	debug bool //debug mode?
	Character
	events event.Eventter
	log    glog.Logger

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

	stats [attributes.EndStatType]float64

	//mods
	statsMod            []*statMod
	attackMods          []*attackMod
	reactionBonusMods   []*reactionBonusMod
	cooldownMods        []*cooldownMod
	healBonusMods       []*healBonusMod
	damageReductionMods []*damageReductionMod
}

func (c *CharWrapper) Tag(key string) int {
	return c.Tags[key]
}

func (c *CharWrapper) AddTag(key string, val int) {
	c.Tags[key] = val
}

func (c *CharWrapper) RemoveTag(key string) {
	delete(c.Tags, key)
}

func (c *CharWrapper) ModifyHP(amt float64) {
	c.HPCurrent += amt
	if c.HPCurrent < 0 {
		c.HPCurrent = -1
	}
	maxhp := c.MaxHP()
	if c.HPCurrent > maxhp {
		c.HPCurrent = maxhp
	}
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

	ActionStam(a action.Action, p map[string]int) float64

	ActionReady(a action.Action, p map[string]int) bool
	SetCD(a action.Action, dur int)
	Cooldown(a action.Action) int
	ResetActionCooldown(a action.Action)
	ReduceActionCooldown(a action.Action, v int)
	Charges(a action.Action) int

	Snapshot(a *combat.AttackInfo) combat.Snapshot

	AddEnergy(src string, amt float64)
	ReceiveParticle(p Particle, isActive bool, partyCount int)
}

type Particle struct {
	Source string
	Num    float64
	Ele    attributes.Element
}

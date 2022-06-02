package character

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/artifact"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
	"github.com/genshinsim/gcsim/pkg/core/task"
)

type Character interface {
	Init() error //init function built into every char to setup any variables etc.

	Attack(p map[string]int) action.ActionInfo
	Aimed(p map[string]int) action.ActionInfo
	ChargeAttack(p map[string]int) action.ActionInfo
	HighPlungeAttack(p map[string]int) action.ActionInfo
	LowPlungeAttack(p map[string]int) action.ActionInfo
	Skill(p map[string]int) action.ActionInfo
	Burst(p map[string]int) action.ActionInfo
	Dash(p map[string]int) action.ActionInfo
	Walk(p map[string]int) action.ActionInfo
	Jump(p map[string]int) action.ActionInfo

	ActionStam(a action.Action, p map[string]int) float64

	ActionReady(a action.Action, p map[string]int) bool
	SetCD(a action.Action, dur int)
	Cooldown(a action.Action) int
	ResetActionCooldown(a action.Action)
	ReduceActionCooldown(a action.Action, v int)
	Charges(a action.Action) int

	Snapshot(a *combat.AttackInfo) combat.Snapshot

	AddEnergy(src string, amt float64)

	ApplyHitlag(factor float64, dur int)

	Tick()
}

type CharWrapper struct {
	Index int
	f     *int //current frame
	debug bool //debug mode?
	Character
	events event.Eventter
	log    glog.Logger
	tasks  task.Tasker

	//base characteristics
	Base     CharacterBase
	Weapon   weapon.WeaponProfile
	Talents  TalentProfile
	CharZone ZoneType
	SkillCon int
	BurstCon int

	Equip struct {
		Weapon weapon.Weapon
		Sets   map[keys.Set]artifact.Set
	}

	//current status
	Energy    float64
	EnergyMax float64
	HPCurrent float64

	//normal attack counter
	NormalHitNum  int //how many hits in a normal combo
	NormalCounter int

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

func New(
	p CharacterProfile,
	f *int, //current frame
	debug bool, //are we running in debug mode
	log glog.Logger, //logging, can be nil
	events event.Eventter, //event emitter
	task task.Tasker,
) *CharWrapper {
	c := &CharWrapper{
		Base:                p.Base,
		Weapon:              p.Weapon,
		Talents:             p.Talents,
		log:                 log,
		events:              events,
		Tags:                make(map[string]int),
		statsMod:            make([]*statMod, 0, 10),
		attackMods:          make([]*attackMod, 0, 10),
		reactionBonusMods:   make([]*reactionBonusMod, 0, 10),
		cooldownMods:        make([]*cooldownMod, 0, 10),
		healBonusMods:       make([]*healBonusMod, 0, 10),
		damageReductionMods: make([]*damageReductionMod, 0, 10),
	}
	s := (*[attributes.EndStatType]float64)(p.Stats)
	c.stats = *s
	c.Equip.Sets = make(map[keys.Set]artifact.Set)

	return c
}

func (c *CharWrapper) SetIndex(index int) {
	c.Index = index
}

func (c *CharWrapper) SetWeapon(w weapon.Weapon) {
	c.Equip.Weapon = w
}

func (c *CharWrapper) SetArtifactSet(key keys.Set, set artifact.Set) {
	c.Equip.Sets[key] = set
}

func (c *CharWrapper) Tag(key string) int {
	return c.Tags[key]
}

func (c *CharWrapper) SetTag(key string, val int) {
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

func (c *CharWrapper) TalentLvlAttack() int {
	//TODO: make sure childe on init sets this tag
	if c.Tags[keys.ChildePassive] > 0 {
		return c.Talents.Attack
	}
	return c.Talents.Attack - 1
}
func (c *CharWrapper) TalentLvlSkill() int {
	if c.Base.Cons >= c.SkillCon {
		return c.Talents.Skill + 2
	}
	return c.Talents.Skill - 1
}
func (c *CharWrapper) TalentLvlBurst() int {
	if c.Base.Cons >= c.BurstCon {
		return c.Talents.Burst + 2
	}
	return c.Talents.Burst - 1
}

type Particle struct {
	Source string
	Num    float64
	Ele    attributes.Element
}

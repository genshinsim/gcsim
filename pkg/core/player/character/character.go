package character

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/artifact"
	"github.com/genshinsim/gcsim/pkg/core/player/character/profile"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
	"github.com/genshinsim/gcsim/pkg/core/task"
	"github.com/genshinsim/gcsim/pkg/modifier"
	"github.com/genshinsim/gcsim/pkg/queue"
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

	ActionReady(a action.Action, p map[string]int) (bool, action.ActionFailure)
	SetCD(a action.Action, dur int)
	Cooldown(a action.Action) int
	ResetActionCooldown(a action.Action)
	ReduceActionCooldown(a action.Action, v int)
	Charges(a action.Action) int

	Snapshot(a *combat.AttackInfo) combat.Snapshot

	AddEnergy(src string, amt float64)

	ApplyHitlag(factor, dur float64)

	Condition([]string) (any, error)

	ResetNormalCounter()
	NextNormalCounter() int
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
	Base     profile.CharacterBase
	Weapon   weapon.WeaponProfile
	Talents  profile.TalentProfile
	CharZone profile.ZoneType
	CharBody profile.BodyType
	SkillCon int
	BurstCon int

	Equip struct {
		Weapon weapon.Weapon
		Sets   map[keys.Set]artifact.Set
	}

	//current status
	ParticleDelay int // character custom particle delay
	Energy        float64
	EnergyMax     float64
	HPCurrent     float64

	//normal attack counter
	NormalHitNum  int //how many hits in a normal combo
	NormalCounter int

	//tags
	Tags      map[string]int
	BaseStats [attributes.EndStatType]float64

	//mods
	mods []modifier.Mod

	//dash cd: keeps track of remaining cd frames for off-field chars
	RemainingDashCD int
	DashLockout     bool

	//hitlag stuff
	timePassed   int //how many frames have passed since start of sim
	frozenFrames int //how many frames are we still frozen for
	queue        []queue.Task
}

type charTask struct {
	f     func()
	delay float64
}

func New(
	p profile.CharacterProfile,
	f *int, //current frame
	debug bool, //are we running in debug mode
	log glog.Logger, //logging, can be nil
	events event.Eventter, //event emitter
	task task.Tasker,
) (*CharWrapper, error) {
	c := &CharWrapper{
		Base:          p.Base,
		Weapon:        p.Weapon,
		Talents:       p.Talents,
		ParticleDelay: 100, //default particle delay
		log:           log,
		events:        events,
		tasks:         task,
		Tags:          make(map[string]int),
		mods:          make([]modifier.Mod, 0, 20),
		f:             f,
		debug:         debug,
	}
	s := (*[attributes.EndStatType]float64)(p.Stats)
	c.BaseStats = *s
	c.Equip.Sets = make(map[keys.Set]artifact.Set)
	//default cons
	c.SkillCon = 3
	c.BurstCon = 5
	//check talents
	if c.Talents.Attack < 1 || c.Talents.Attack > 10 {
		return nil, fmt.Errorf("invalid talent lvl: attack - %v", c.Talents.Attack)
	}
	if c.Talents.Skill < 1 || c.Talents.Skill > 10 {
		return nil, fmt.Errorf("invalid talent lvl: skill - %v", c.Talents.Skill)
	}
	if c.Talents.Burst < 1 || c.Talents.Burst > 10 {
		return nil, fmt.Errorf("invalid talent lvl: burst - %v", c.Talents.Burst)
	}
	//setup base hp
	if p.Base.StartHP > -1 {
		c.HPCurrent = p.Base.StartHP
	} else {
		c.HPCurrent = -1 //to be cleared up in init
	}

	return c, nil
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

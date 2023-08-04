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
	Base      profile.CharacterBase
	Weapon    weapon.WeaponProfile
	Talents   profile.TalentProfile
	CharZone  profile.ZoneType
	CharBody  profile.BodyType
	NormalCon int
	SkillCon  int
	BurstCon  int

	Equip struct {
		Weapon weapon.Weapon
		Sets   map[keys.Set]artifact.Set
	}

	//current status
	ParticleDelay  int // character custom particle delay
	Energy         float64
	EnergyMax      float64
	currentHPRatio float64
	// needed so that start hp is not influenced by hp mods added during team initialization
	StartHP int

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

	//set to -1 by default and let each char specify normal/skill/burst cons
	c.NormalCon = -1
	c.SkillCon = -1
	c.BurstCon = -1

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

func (c *CharWrapper) clampHPRatio() {
	if c.currentHPRatio > 1 {
		c.currentHPRatio = 1
	} else if c.currentHPRatio < 0 {
		c.currentHPRatio = 0
	}
}

func (c *CharWrapper) SetHPByAmount(amt float64) {
	c.currentHPRatio = amt / c.MaxHP()
	c.clampHPRatio()
}

func (c *CharWrapper) SetHPByRatio(r float64) {
	c.currentHPRatio = r
	c.clampHPRatio()
}

func (c *CharWrapper) ModifyHPByAmount(amt float64) {
	newHP := c.CurrentHP() + amt
	c.SetHPByAmount(newHP)
}

func (c *CharWrapper) ModifyHPByRatio(r float64) {
	newHPRatio := c.currentHPRatio + r
	c.SetHPByRatio(newHPRatio)
}

func (c *CharWrapper) consCheck() {
	consUnset := 0
	if c.NormalCon < 0 {
		consUnset++
	}
	if c.SkillCon < 0 {
		consUnset++
	}
	if c.BurstCon < 0 {
		consUnset++
	}
	if consUnset != 1 {
		panic(fmt.Sprintf("cons not set properly for %v, please set two out of three values:\nNormalCon: %v\nSkillCon: %v\nBurstCon: %v", c.Base.Key.String(), c.NormalCon, c.SkillCon, c.BurstCon))
	}
}

func (c *CharWrapper) TalentLvlAttack() int {
	c.consCheck()
	add := -1
	if c.Tags[keys.ChildePassive] > 0 {
		add++
	}
	if c.NormalCon > 0 && c.Base.Cons >= c.NormalCon {
		add += 3
	}
	if add >= 4 {
		add = 4
	}
	return c.Talents.Attack + add
}
func (c *CharWrapper) TalentLvlSkill() int {
	c.consCheck()
	add := -1
	if c.SkillCon > 0 && c.Base.Cons >= c.SkillCon {
		add += 3
	}
	if add >= 4 {
		add = 4
	}
	return c.Talents.Skill + add
}
func (c *CharWrapper) TalentLvlBurst() int {
	c.consCheck()
	add := -1
	if c.BurstCon > 0 && c.Base.Cons >= c.BurstCon {
		add += 3
	}
	if add >= 4 {
		add = 4
	}
	return c.Talents.Burst + add
}

type Particle struct {
	Source string
	Num    float64
	Ele    attributes.Element
}

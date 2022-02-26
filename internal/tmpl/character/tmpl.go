package character

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/genshinsim/gcsim/pkg/core"
)

type Tmpl struct {
	Core  *core.Core
	Rand  *rand.Rand
	Index int
	//this should describe the frame in which the abil becomes available
	//if frame > current then it's available. no need to decrement this way
	// CD        map[string]int
	ActionCD      []int
	Mods          []core.CharStatMod
	PreDamageMods []core.PreDamageMod
	ReactMod      []core.ReactionBonusMod
	Tags          map[string]int
	//Profile info
	Base     core.CharacterBase
	Weapon   core.WeaponProfile
	Stats    [core.EndStatType]float64
	Talents  core.TalentProfile
	SkillCon int
	BurstCon int
	CharZone core.ZoneType

	CDReductionFuncs []core.CDAdjust

	Energy    float64
	EnergyMax float64

	HPCurrent float64
	HPMax     float64

	//counters
	NormalHitNum  int //how many hits in a normal combo
	NormalCounter int

	//map to track frames
	normalCancelFrames map[int]map[core.ActionType]int             //this maps normal hit number into
	cancelFrames       map[core.ActionType]map[core.ActionType]int //this maps all other skills

	//infusion
	Infusion core.WeaponInfusion //TODO currently just overides the old; disregarding any existing
}

func NewTemplateChar(x *core.Core, p core.CharacterProfile) (*Tmpl, error) {
	t := Tmpl{}
	t.Core = x
	t.Rand = x.Rand

	t.ActionCD = make([]int, core.EndActionType)
	t.Mods = make([]core.CharStatMod, 0, 10)
	t.Tags = make(map[string]int)
	t.CDReductionFuncs = make([]core.CDAdjust, 0, 5)
	t.Base = p.Base
	t.Weapon = p.Weapon
	t.Talents = p.Talents
	t.SkillCon = 3
	t.BurstCon = 5
	if t.Talents.Attack < 1 || t.Talents.Attack > 15 {
		return nil, fmt.Errorf("invalid talent lvl: attack - %v", t.Talents.Attack)
	}
	if t.Talents.Attack < 1 || t.Talents.Attack > 12 {
		return nil, fmt.Errorf("invalid talent lvl: skill - %v", t.Talents.Skill)
	}
	if t.Talents.Attack < 1 || t.Talents.Attack > 12 {
		return nil, fmt.Errorf("invalid talent lvl: burst - %v", t.Talents.Burst)
	}
	for i, v := range p.Stats {
		t.Stats[i] = v
	}
	if p.Base.StartHP > -1 {
		t.Core.Log.NewEvent("setting starting hp", core.LogCharacterEvent, t.Index, "character", p.Base.Key.String(), "hp", p.Base.StartHP)
		t.HPCurrent = p.Base.StartHP
	} else {
		t.HPCurrent = math.MaxInt64
	}

	t.normalCancelFrames = make(map[int]map[core.ActionType]int)
	t.cancelFrames = make(map[core.ActionType]map[core.ActionType]int)

	return &t, nil
}

func (t *Tmpl) SetIndex(index int) {
	t.Index = index
}

// Character initialization function. Occurs AFTER all char/weapons are initially loaded
func (t *Tmpl) Init() {
	hpp := t.Stats[core.HPP]
	hp := t.Stats[core.HP]

	for _, m := range t.Mods {
		if m.Expiry > t.Core.F || m.Expiry == -1 {
			a, ok := m.Amount()
			if ok {
				hpp += a[core.HPP]
				hp += a[core.HP]
			}
		}
	}

	t.HPMax = t.Base.HP*(1+hpp) + hp
	// c.HPCurrent = 1
	if t.HPCurrent > t.HPMax {
		t.HPCurrent = t.HPMax
	}
}

func (c *Tmpl) SetWeaponKey(k string) {
	c.Weapon.Key = k
}

func (c *Tmpl) WeaponKey() string {
	return c.Weapon.Key
}

func (c *Tmpl) AddWeaponInfuse(inf core.WeaponInfusion) {
	c.Infusion = inf
}

func (c *Tmpl) ModIsActive(key string) bool {
	ind := -1
	for i, v := range c.Mods {
		if v.Key == key {
			ind = i
		}
	}
	//mod doesnt exist
	if ind == -1 {
		return false
	}
	//check expiry
	if c.Mods[ind].Expiry < c.Core.F && c.Mods[ind].Expiry > -1 {
		return false
	}
	_, ok := c.Mods[ind].Amount()
	return ok
}

// TODO: This design pattern feels wrong... I think we should have mods be their own separate interface
// Each mod would then have an "IsActive" function that we pass a character to?
func (c *Tmpl) PreDamageModIsActive(key string) bool {
	ind := -1
	for i, v := range c.PreDamageMods {
		if v.Key == key {
			ind = i
		}
	}
	//mod doesnt exist
	if ind == -1 {
		return false
	}
	//check expiry
	if c.PreDamageMods[ind].Expiry < c.Core.F && c.PreDamageMods[ind].Expiry > -1 {
		return false
	}
	return true
}

func (c *Tmpl) ReactBonusModIsActive(key string) bool {
	ind := -1
	for i, v := range c.ReactMod {
		if v.Key == key {
			ind = i
		}
	}
	//mod doesnt exist
	if ind == -1 {
		return false
	}
	//check expiry
	if c.ReactMod[ind].Expiry < c.Core.F && c.ReactMod[ind].Expiry > -1 {
		return false
	}
	return true
}

func (c *Tmpl) Tag(key string) int {
	return c.Tags[key]
}

func (c *Tmpl) AddTag(key string, val int) {
	c.Tags[key] = val
}

func (c *Tmpl) RemoveTag(key string) {
	delete(c.Tags, key)
}

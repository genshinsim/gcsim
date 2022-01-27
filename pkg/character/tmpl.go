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

	//infusion
	Infusion core.WeaponInfusion //TODO currently just overides the old; disregarding any existing
}

func NewTemplateChar(x *core.Core, p core.CharacterProfile) (*Tmpl, error) {
	c := Tmpl{}
	c.Core = x
	c.Rand = x.Rand

	c.ActionCD = make([]int, core.EndActionType)
	c.Mods = make([]core.CharStatMod, 0, 10)
	c.Tags = make(map[string]int)
	c.CDReductionFuncs = make([]core.CDAdjust, 0, 5)
	c.Base = p.Base
	c.Weapon = p.Weapon
	c.Talents = p.Talents
	c.SkillCon = 3
	c.BurstCon = 5
	if c.Talents.Attack < 1 || c.Talents.Attack > 15 {
		return nil, fmt.Errorf("invalid talent lvl: attack - %v", c.Talents.Attack)
	}
	if c.Talents.Attack < 1 || c.Talents.Attack > 12 {
		return nil, fmt.Errorf("invalid talent lvl: skill - %v", c.Talents.Skill)
	}
	if c.Talents.Attack < 1 || c.Talents.Attack > 12 {
		return nil, fmt.Errorf("invalid talent lvl: burst - %v", c.Talents.Burst)
	}
	for i, v := range p.Stats {
		c.Stats[i] = v
	}
	if p.Base.StartHP > -1 {
		c.Core.Log.Debugw("setting starting hp", "frame", x.F, "event", core.LogCharacterEvent, "character", p.Base.Key.String(), "hp", p.Base.StartHP)
		c.HPCurrent = p.Base.StartHP
	} else {
		c.HPCurrent = math.MaxInt64
	}

	return &c, nil
}

// Character initialization function. Occurs AFTER all char/weapons are initially loaded
func (t *Tmpl) Init(index int) {
	t.Index = index
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

func (c *Tmpl) AddPreDamageMod(mod core.PreDamageMod) {
	ind := -1
	for i, v := range c.PreDamageMods {
		if v.Key == mod.Key {
			ind = i
		}
	}
	if ind > -1 {
		c.Core.Log.Debugw("mod refreshed", "frame", c.Core.F, "event", core.LogStatusEvent, "char", c.Index, "overwrite", true, "key", mod.Key, "expiry", mod.Expiry)
		c.PreDamageMods[ind] = mod
	} else {
		c.PreDamageMods = append(c.PreDamageMods, mod)
		c.Core.Log.Debugw("mod added", "frame", c.Core.F, "event", core.LogStatusEvent, "char", c.Index, "overwrite", false, "key", mod.Key, "expiry", mod.Expiry)
	}

	// Add task to check for mod expiry in debug instances
	if c.Core.Flags.LogDebug && mod.Expiry > -1 {
		c.AddTask(func() {
			if c.PreDamageModIsActive(mod.Key) {
				return
			}
			c.Core.Log.Debugw("mod expired", "frame", c.Core.F, "event", core.LogStatusEvent, "char", c.Index, "overwrite", false, "key", mod.Key, "expiry", mod.Expiry)
		}, "check-mod-expiry", mod.Expiry+1-c.Core.F)
	}
}

func (c *Tmpl) AddMod(mod core.CharStatMod) {
	ind := -1
	for i, v := range c.Mods {
		if v.Key == mod.Key {
			ind = i
		}
	}
	if ind > -1 {
		c.Core.Log.Debugw("mod refreshed", "frame", c.Core.F, "event", core.LogStatusEvent, "char", c.Index, "overwrite", true, "key", mod.Key, "expiry", mod.Expiry)
		c.Mods[ind] = mod
	} else {
		c.Mods = append(c.Mods, mod)
		c.Core.Log.Debugw("mod added", "frame", c.Core.F, "event", core.LogStatusEvent, "char", c.Index, "overwrite", false, "key", mod.Key, "expiry", mod.Expiry)
	}

	// Add task to check for mod expiry in debug instances
	if c.Core.Flags.LogDebug && mod.Expiry > -1 {
		c.AddTask(func() {
			if c.ModIsActive(mod.Key) {
				return
			}
			c.Core.Log.Debugw("mod expired", "frame", c.Core.F, "event", core.LogStatusEvent, "char", c.Index, "overwrite", false, "key", mod.Key, "expiry", mod.Expiry)
		}, "check-mod-expiry", mod.Expiry+1-c.Core.F)
	}
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

func (t *Tmpl) AddReactBonusMod(mod core.ReactionBonusMod) {
	ind := -1
	for i, v := range t.ReactMod {
		if v.Key == mod.Key {
			ind = i
		}
	}
	if ind != -1 {
		t.ReactMod[ind] = mod
		t.Core.Log.Debugw("mod refreshed", "frame", t.Core.F, "char", t.Index, "event", core.LogStatusEvent, "char", t.Index, "overwrite", true, "key", mod.Key, "expiry", mod.Expiry)
		return
	}
	t.ReactMod = append(t.ReactMod, mod)
	t.Core.Log.Debugw("mod added", "frame", t.Core.F, "char", t.Index, "event", core.LogStatusEvent, "char", t.Index, "key", mod.Key, "expiry", mod.Expiry)

	// Add task to check for mod expiry in debug instances
	if t.Core.Flags.LogDebug && mod.Expiry > -1 {
		t.AddTask(func() {
			if t.ReactBonusModIsActive(mod.Key) {
				return
			}
			t.Core.Log.Debugw("mod expired", "frame", t.Core.F, "char", t.Index, "event", core.LogStatusEvent, "char", t.Index, "key", mod.Key, "expiry", mod.Expiry)
		}, "check-mod-expiry", mod.Expiry+1-t.Core.F)
	}
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

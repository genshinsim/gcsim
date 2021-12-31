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

func (t *Tmpl) Init(index int) {
	t.Index = index
	hpp := t.Stats[core.HPP]
	hp := t.Stats[core.HP]

	for _, m := range t.Mods {
		if m.Expiry > t.Core.F || m.Expiry == -1 {
			a, ok := m.Amount(core.AttackTagNone)
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

func (c *Tmpl) AddWeaponInfuse(inf core.WeaponInfusion) {
	c.Infusion = inf
}

func (c *Tmpl) AddPreDamageMod(mod core.PreDamageMod) {
	ind := len(c.PreDamageMods)
	for i, v := range c.PreDamageMods {
		if v.Key == mod.Key {
			ind = i
		}
	}
	if ind != 0 && ind != len(c.PreDamageMods) {
		c.Core.Log.Debugw("char pre damage mod added", "frame", c.Core.F, "event", core.LogCharacterEvent, "overwrite", true, "key", mod.Key, "expiry", mod.Expiry)
		c.PreDamageMods[ind] = mod
	} else {
		c.PreDamageMods = append(c.PreDamageMods, mod)
		c.Core.Log.Debugw("char pre damage mod added", "frame", c.Core.F, "event", core.LogCharacterEvent, "overwrite", true, "key", mod.Key, "expiry", mod.Expiry)
	}

}

func (c *Tmpl) AddMod(mod core.CharStatMod) {
	ind := len(c.Mods)
	for i, v := range c.Mods {
		if v.Key == mod.Key {
			ind = i
		}
	}
	if ind != 0 && ind != len(c.Mods) {
		c.Core.Log.Debugw("char mod added", "frame", c.Core.F, "char", c.Index, "event", core.LogCharacterEvent, "overwrite", true, "key", mod.Key, "expiry", mod.Expiry)
		c.Mods[ind] = mod
	} else {
		c.Mods = append(c.Mods, mod)
		c.Core.Log.Debugw("char mod added", "frame", c.Core.F, "char", c.Index, "event", core.LogCharacterEvent, "overwrite", true, "key", mod.Key, "expiry", mod.Expiry)
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
	_, ok := c.Mods[ind].Amount(core.AttackTagNone)
	return ok
}

func (t *Tmpl) AddReactBonusMod(mod core.ReactionBonusMod) {
	ind := -1
	for i, v := range t.ReactMod {
		if v.Key == mod.Key {
			ind = i
		}
	}
	if ind != -1 {
		t.Core.Log.Debugw("react bonus mod overwritten", "frame", t.Core.F, "event", core.LogEnemyEvent, "count", len(t.ReactMod), "char", t.Index)
		// LogEnemyEvent
		t.ReactMod[ind] = mod
		return
	}
	t.ReactMod = append(t.ReactMod, mod)
	t.Core.Log.Debugw("react bonus mod added", "frame", t.Core.F, "event", core.LogEnemyEvent, "count", len(t.ReactMod), "char", t.Index)
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

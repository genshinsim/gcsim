package character

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/genshinsim/gcsim/pkg/coretype"
)

type Core interface {
	coretype.Framer
	coretype.EventEmitter
	coretype.Logger
	coretype.RandomGenerator
	coretype.TaskHandler
}

type Tmpl struct {
	Core  Core
	Rand  *rand.Rand
	Index int

	//cooldown related
	ActionCD               []int
	cdQueueWorkerStartedAt []int
	cdCurrentQueueWorker   []*func()
	cdQueue                [][]int
	AvailableCDCharge      []int
	additionalCDCharge     []int

	//mods
	Mods          []coretype.CharStatMod
	PreDamageMods []coretype.PreDamageMod
	ReactMod      []coretype.ReactionBonusMod
	Tags          map[string]int
	//Profile info
	Base     coretype.CharacterBase
	Weapon   coretype.WeaponProfile
	Stats    [coretype.EndStatType]float64
	Talents  coretype.TalentProfile
	SkillCon int
	BurstCon int
	CharZone coretype.ZoneType

	CDReductionFuncs []coretype.CDAdjust

	Energy    float64
	EnergyMax float64

	HPCurrent float64
	HPMax     float64

	//counters
	NormalHitNum  int //how many hits in a normal combo
	NormalCounter int

	//map to track frames
	normalCancelFrames map[int]map[coretype.ActionType]int                 //this maps normal hit number into
	cancelFrames       map[coretype.ActionType]map[coretype.ActionType]int //this maps all other skills

	//infusion
	Infusion coretype.WeaponInfusion //TODO currently just overides the old; disregarding any existing
}

func NewTemplateChar(x Core, p coretype.CharacterProfile) (*Tmpl, error) {
	t := Tmpl{}
	t.Core = x
	t.Rand = x.R()

	t.ActionCD = make([]int, coretype.EndActionType)
	t.Mods = make([]coretype.CharStatMod, 0, 10)
	t.Tags = make(map[string]int)
	t.CDReductionFuncs = make([]coretype.CDAdjust, 0, 5)
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
		t.Core.NewEvent("setting starting hp", coretype.LogCharacterEvent, t.Index, "character", p.Base.Key.String(), "hp", p.Base.StartHP)
		t.HPCurrent = p.Base.StartHP
	} else {
		t.HPCurrent = math.MaxInt64
	}

	t.normalCancelFrames = make(map[int]map[coretype.ActionType]int)
	t.cancelFrames = make(map[coretype.ActionType]map[coretype.ActionType]int)

	t.cdQueueWorkerStartedAt = make([]int, coretype.EndActionType)
	t.cdCurrentQueueWorker = make([]*func(), coretype.EndActionType)
	t.cdQueue = make([][]int, coretype.EndActionType)
	t.additionalCDCharge = make([]int, coretype.EndActionType)
	t.AvailableCDCharge = make([]int, coretype.EndActionType)

	for i := 0; i < len(t.cdQueue); i++ {
		t.cdQueue[i] = make([]int, 0, 4)
		t.AvailableCDCharge[i] = 1
	}

	return &t, nil
}

func (t *Tmpl) SetIndex(index int) {
	t.Index = index
}

// Character initialization function. Occurs AFTER all char/weapons are initially loaded
func (t *Tmpl) Init() {
	hpp := t.Stats[coretype.HPP]
	hp := t.Stats[coretype.HP]

	for _, m := range t.Mods {
		if m.Expiry > t.Core.F() || m.Expiry == -1 {
			a, ok := m.Amount()
			if ok {
				hpp += a[coretype.HPP]
				hp += a[coretype.HP]
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

func (c *Tmpl) AddWeaponInfuse(inf coretype.WeaponInfusion) {
	c.Infusion = inf
}

func (c *Tmpl) AddPreDamageMod(mod coretype.PreDamageMod) {
	ind := -1
	for i, v := range c.PreDamageMods {
		if v.Key == mod.Key {
			ind = i
		}
	}

	// check if mod exists and has not expired
	if ind != -1 && (c.PreDamageMods[ind].Expiry > c.Core.F() || c.PreDamageMods[ind].Expiry == -1) {
		c.Core.NewEvent("mod refreshed", coretype.LogStatusEvent, c.Index, "overwrite", true, "key", mod.Key, "expiry", mod.Expiry)
		mod.Event = c.PreDamageMods[ind].Event
	} else {
		mod.Event = c.Core.NewEvent("mod added", coretype.LogStatusEvent, c.Index, "overwrite", false, "key", mod.Key, "expiry", mod.Expiry)
		// append empty mod if we can not reuse mods[ind]
		if ind == -1 {
			c.PreDamageMods = append(c.PreDamageMods, coretype.PreDamageMod{})
			ind = len(c.PreDamageMods) - 1
		}
	}
	mod.Event.SetEnded(mod.Expiry)
	c.PreDamageMods[ind] = mod
}

func (c *Tmpl) AddMod(mod coretype.CharStatMod) {
	ind := -1
	for i, v := range c.Mods {
		if v.Key == mod.Key {
			ind = i
		}
	}

	// check if mod exists and has not expired
	if ind != -1 && (c.Mods[ind].Expiry > c.Core.F() || c.Mods[ind].Expiry == -1) {
		c.Core.NewEvent("mod refreshed", coretype.LogStatusEvent, c.Index, "overwrite", true, "key", mod.Key, "expiry", mod.Expiry)
		mod.Event = c.Mods[ind].Event
	} else {
		mod.Event = c.Core.NewEvent("mod added", coretype.LogStatusEvent, c.Index, "overwrite", false, "key", mod.Key, "expiry", mod.Expiry)
		// append empty mod if we can not reuse mods[ind]
		if ind == -1 {
			c.Mods = append(c.Mods, coretype.CharStatMod{})
			ind = len(c.Mods) - 1
		}
	}
	mod.Event.SetEnded(mod.Expiry)
	c.Mods[ind] = mod
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
	if c.Mods[ind].Expiry < c.Core.F() && c.Mods[ind].Expiry > -1 {
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
	if c.PreDamageMods[ind].Expiry < c.Core.F() && c.PreDamageMods[ind].Expiry > -1 {
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
	if c.ReactMod[ind].Expiry < c.Core.F() && c.ReactMod[ind].Expiry > -1 {
		return false
	}
	return true
}

func (c *Tmpl) WeaponInfuseIsActive(key string) bool {
	if c.Infusion.Key != key {
		return false
	}
	//check expiry
	if c.Infusion.Expiry < c.Core.F() && c.Infusion.Expiry > -1 {
		return false
	}
	return true
}

func (t *Tmpl) AddReactBonusMod(mod coretype.ReactionBonusMod) {
	ind := -1
	for i, v := range t.ReactMod {
		if v.Key == mod.Key {
			ind = i
		}
	}

	// check if mod exists and has not expired
	if ind != -1 && (t.ReactMod[ind].Expiry > t.Core.F() || t.ReactMod[ind].Expiry == -1) {
		t.Core.NewEvent("mod refreshed", coretype.LogStatusEvent, t.Index, "overwrite", true, "key", mod.Key, "expiry", mod.Expiry)
		mod.Event = t.ReactMod[ind].Event
	} else {
		mod.Event = t.Core.NewEvent("mod added", coretype.LogStatusEvent, t.Index, "overwrite", false, "key", mod.Key, "expiry", mod.Expiry)
		// append empty mod if we can not reuse mods[ind]
		if ind == -1 {
			t.ReactMod = append(t.ReactMod, coretype.ReactionBonusMod{})
			ind = len(t.ReactMod) - 1
		}
	}
	mod.Event.SetEnded(mod.Expiry)
	t.ReactMod[ind] = mod
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

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

	//cooldown related
	ActionCD               []int
	cdQueueWorkerStartedAt []int
	cdCurrentQueueWorker   []*func()
	cdQueue                [][]int
	AvailableCDCharge      []int
	additionalCDCharge     []int

	//mods
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

	// TODO: maybe should change this to % of max hp
	HPCurrent float64

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

	t.cdQueueWorkerStartedAt = make([]int, core.EndActionType)
	t.cdCurrentQueueWorker = make([]*func(), core.EndActionType)
	t.cdQueue = make([][]int, core.EndActionType)
	t.additionalCDCharge = make([]int, core.EndActionType)
	t.AvailableCDCharge = make([]int, core.EndActionType)

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
	maxhp := t.MaxHP()
	if t.HP() > maxhp {
		t.HPCurrent = maxhp
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

	// check if mod exists and has not expired
	if ind != -1 && (c.PreDamageMods[ind].Expiry > c.Core.F || c.PreDamageMods[ind].Expiry == -1) {
		c.Core.Log.NewEvent("mod refreshed", core.LogStatusEvent, c.Index, "overwrite", true, "key", mod.Key, "expiry", mod.Expiry)
		mod.Event = c.PreDamageMods[ind].Event
	} else {
		mod.Event = c.Core.Log.NewEvent("mod added", core.LogStatusEvent, c.Index, "overwrite", false, "key", mod.Key, "expiry", mod.Expiry)
		// append empty mod if we can not reuse mods[ind]
		if ind == -1 {
			c.PreDamageMods = append(c.PreDamageMods, core.PreDamageMod{})
			ind = len(c.PreDamageMods) - 1
		}
	}
	mod.Event.SetEnded(mod.Expiry)
	c.PreDamageMods[ind] = mod
}

func (c *Tmpl) DeletePreDamageMod(key string) {
	n := 0
	for _, v := range c.PreDamageMods {
		if v.Key == key {
			v.Event.SetEnded(c.Core.F)
			c.Core.Log.NewEvent("mod deleted", core.LogStatusEvent, c.Index, "key", key)
		} else {
			c.PreDamageMods[n] = v
			n++
		}
	}
	c.PreDamageMods = c.PreDamageMods[:n]
}

func (c *Tmpl) AddMod(mod core.CharStatMod) {
	ind := -1
	for i, v := range c.Mods {
		if v.Key == mod.Key {
			ind = i
		}
	}

	// check if mod exists and has not expired
	if ind != -1 && (c.Mods[ind].Expiry > c.Core.F || c.Mods[ind].Expiry == -1) {
		c.Core.Log.NewEvent("mod refreshed", core.LogStatusEvent, c.Index, "overwrite", true, "key", mod.Key, "expiry", mod.Expiry)
		mod.Event = c.Mods[ind].Event
	} else {
		mod.Event = c.Core.Log.NewEvent("mod added", core.LogStatusEvent, c.Index, "overwrite", false, "key", mod.Key, "expiry", mod.Expiry)
		// append empty mod if we can not reuse mods[ind]
		if ind == -1 {
			c.Mods = append(c.Mods, core.CharStatMod{})
			ind = len(c.Mods) - 1
		}
	}
	mod.Event.SetEnded(mod.Expiry)
	c.Mods[ind] = mod
}

func (c *Tmpl) DeleteMod(key string) {
	n := 0
	for _, v := range c.Mods {
		if v.Key == key {
			v.Event.SetEnded(c.Core.F)
			c.Core.Log.NewEvent("mod deleted", core.LogStatusEvent, c.Index, "key", key)
		} else {
			c.Mods[n] = v
			n++
		}
	}
	c.Mods = c.Mods[:n]
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

func (c *Tmpl) WeaponInfuseIsActive(key string) bool {
	if c.Infusion.Key != key {
		return false
	}
	//check expiry
	if c.Infusion.Expiry < c.Core.F && c.Infusion.Expiry > -1 {
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

	// check if mod exists and has not expired
	if ind != -1 && (t.ReactMod[ind].Expiry > t.Core.F || t.ReactMod[ind].Expiry == -1) {
		t.Core.Log.NewEvent("mod refreshed", core.LogStatusEvent, t.Index, "overwrite", true, "key", mod.Key, "expiry", mod.Expiry)
		mod.Event = t.ReactMod[ind].Event
	} else {
		mod.Event = t.Core.Log.NewEvent("mod added", core.LogStatusEvent, t.Index, "overwrite", false, "key", mod.Key, "expiry", mod.Expiry)
		// append empty mod if we can not reuse mods[ind]
		if ind == -1 {
			t.ReactMod = append(t.ReactMod, core.ReactionBonusMod{})
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

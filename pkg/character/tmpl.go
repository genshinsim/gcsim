package character

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/genshinsim/gsim/pkg/core"
	"go.uber.org/zap"
)

type Tmpl struct {
	Core  *core.Core
	Rand  *rand.Rand
	Log   *zap.SugaredLogger
	Index int
	//this should describe the frame in which the abil becomes available
	//if frame > current then it's available. no need to decrement this way
	// CD        map[string]int
	ActionCD []int
	Mods     []core.CharStatMod
	Tags     map[string]int
	//Profile info
	Base     core.CharacterBase
	Weapon   core.WeaponProfile
	Stats    []float64
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
	c.Log = x.Log
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
	c.Stats = make([]float64, core.EndStatType)
	for i, v := range p.Stats {
		c.Stats[i] = v
	}
	if p.Base.StartHP > -1 {
		c.Log.Debugw("setting starting hp", "frame", x.F, "event", core.LogCharacterEvent, "character", p.Base.Name, "hp", p.Base.StartHP)
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

func (c *Tmpl) AddWeaponInfuse(inf core.WeaponInfusion) {
	c.Infusion = inf
}

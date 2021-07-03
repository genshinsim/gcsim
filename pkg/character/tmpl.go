package character

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/genshinsim/gsim/pkg/def"
	"go.uber.org/zap"
)

type Tmpl struct {
	Sim   def.Sim
	Rand  *rand.Rand
	Log   *zap.SugaredLogger
	Index int
	//this should describe the frame in which the abil becomes available
	//if frame > current then it's available. no need to decrement this way
	// CD        map[string]int
	ActionCD []int
	Mods     []def.CharStatMod
	Tags     map[string]int
	//Profile info
	Base     def.CharacterBase
	Weapon   def.WeaponProfile
	Stats    []float64
	Talents  def.TalentProfile
	SkillCon int
	BurstCon int
	CharZone def.ZoneType

	CDReductionFuncs []def.CDAdjust

	Energy    float64
	MaxEnergy float64

	HPCurrent float64
	HPMax     float64

	//Tasks specific to the character to be executed at set frames
	Tasks map[int][]CharTask
	//counters
	NormalHitNum  int //how many hits in a normal combo
	NormalCounter int

	//infusion
	Infusion def.WeaponInfusion //TODO currently just overides the old; disregarding any existing
}

type CharTask struct {
	Name        string
	F           func()
	Delay       int
	originFrame int
}

func NewTemplateChar(s def.Sim, log *zap.SugaredLogger, p def.CharacterProfile) (*Tmpl, error) {
	c := Tmpl{}
	c.Sim = s
	c.Log = log
	c.Rand = s.Rand()

	c.ActionCD = make([]int, def.EndActionType)
	c.Mods = make([]def.CharStatMod, 0, 10)
	c.Tags = make(map[string]int)
	c.Tasks = make(map[int][]CharTask)
	c.CDReductionFuncs = make([]def.CDAdjust, 0, 5)
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
	c.Stats = make([]float64, len(def.StatTypeString))
	for i, v := range p.Stats {
		c.Stats[i] = v
	}
	if p.Base.StartHP > -1 {
		c.Log.Debugw("setting starting hp", "frame", s.Frame(), "event", def.LogCharacterEvent, "character", p.Base.Name, "hp", p.Base.StartHP)
		c.HPCurrent = p.Base.StartHP
	} else {
		c.HPCurrent = math.MaxInt64
	}

	return &c, nil
}

func (t *Tmpl) Init(index int) {
	t.Index = index
	hpp := t.Stats[def.HPP]
	hp := t.Stats[def.HP]

	for _, m := range t.Mods {
		if m.Expiry > t.Sim.Frame() || m.Expiry == -1 {
			a, ok := m.Amount(def.AttackTagNone)
			if ok {
				hpp += a[def.HPP]
				hp += a[def.HP]
			}
		}
	}

	t.HPMax = t.Base.HP*(1+hpp) + hp
	// c.HPCurrent = 1
	if t.HPCurrent > t.HPMax {
		t.HPCurrent = t.HPMax
	}
}

func (c *Tmpl) AddWeaponInfuse(inf def.WeaponInfusion) {
	c.Infusion = inf
}

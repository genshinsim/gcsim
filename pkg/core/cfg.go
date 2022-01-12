package core

type Config struct {
	Label      string
	DamageMode bool
	Targets    []EnemyProfile
	Characters struct {
		Initial CharKey
		Profile []CharacterProfile
	}
	Rotation []ActionBlock

	Hurt      HurtEvent
	Energy    EnergyEvent
	FixedRand bool //if this is true then use the same seed
}

type RunOpt struct {
	LogDetails bool `json:"log_details"`
	Iteration  int  `json:"iter"`
	Workers    int  `json:"workers"`
	Duration   int  `json:"seconds"`
	Debug      bool `json:"debug"`
	ERCalcMode bool `json:"er_calc_mode"`
	DebugPaths []string
}

type CharacterProfile struct {
	Base      CharacterBase             `json:"base"`
	Weapon    WeaponProfile             `json:"weapon"`
	Talents   TalentProfile             `json:"talents"`
	Stats     []float64                 `json:"stats"`
	Sets      map[string]int            `json:"sets"`
	SetParams map[string]map[string]int `json:"-"`
	Params    map[string]int            `json:"-"`
}

type CharacterBase struct {
	Key      CharKey `json:"-"`
	Name     string  `json:"name"`
	Element  EleType `json:"element"`
	Level    int     `json:"level"`
	MaxLevel int     `json:"max_level"`
	HP       float64 `json:"-"`
	Atk      float64 `json:"-"`
	Def      float64 `json:"-"`
	Cons     int     `json:"cons"`
	StartHP  float64 `json:"-"`
}

type WeaponProfile struct {
	Name     string         `json:"name"`
	Key      string         `json:""` //use this to match with weapon curve mapping
	Class    WeaponClass    `json:"-"`
	Refine   int            `json:"refine"`
	Level    int            `json:"level"`
	MaxLevel int            `json:"max_level"`
	Atk      float64        `json:"-"`
	Params   map[string]int `json:"-"`
}

type TalentProfile struct {
	Attack int `json:"attack"`
	Skill  int `json:"skill"`
	Burst  int `json:"burst"`
}

type EnemyProfile struct {
	Level          int                 `json:"level"`
	HP             float64             `json:"-"`
	Resist         map[EleType]float64 `json:"-"`
	Size           float64             `json:"-"`
	CoordX, CoordY float64             `json:"-"`
}

type EnergyEvent struct {
	Active    bool
	Once      bool //how often
	Start     int
	End       int
	Particles int
}

type HurtEvent struct {
	Active bool
	Once   bool //how often
	Start  int  //
	End    int
	Min    float64
	Max    float64
	Ele    EleType
}

func (e *EnemyProfile) Clone() EnemyProfile {
	r := EnemyProfile{
		Level:  e.Level,
		Resist: make(map[EleType]float64),
	}
	for k, v := range e.Resist {
		r.Resist[k] = v
	}
	return r
}

func (c *CharacterProfile) Clone() CharacterProfile {
	r := *c
	r.Weapon.Params = make(map[string]int)
	for k, v := range c.Weapon.Params {
		r.Weapon.Params[k] = v
	}
	r.Stats = make([]float64, len(c.Stats))
	copy(r.Stats, c.Stats)
	r.Sets = make(map[string]int)
	for k, v := range c.Sets {
		r.Sets[k] = v
	}

	return r
}

func (c *Config) Clone() Config {
	r := *c

	r.Targets = make([]EnemyProfile, len(c.Targets))

	for i, v := range c.Targets {
		r.Targets[i] = v.Clone()
	}

	return r
}

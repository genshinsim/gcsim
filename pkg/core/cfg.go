package core

//SimulationConfig describes the required settings to run an simulation
type SimulationConfig struct {
	//these settings relate to each simulation iteration
	DamageMode bool           `json:"damage_mode"`
	Targets    []EnemyProfile `json:"targets"`
	Characters struct {
		Initial CharKey            `json:"initial"`
		Profile []CharacterProfile `json:"profile"`
	} `json:"characters"`
	Rotation []ActionBlock     `json:"-"`
	Hurt     HurtEvent         `json:"-"`
	Energy   EnergyEvent       `json:"-"`
	Settings SimulatorSettings `json:"-"`
}

func (c *SimulationConfig) Clone() SimulationConfig {
	r := *c

	r.Targets = make([]EnemyProfile, len(c.Targets))
	for i, v := range c.Targets {
		r.Targets[i] = v.Clone()
	}

	r.Characters.Profile = make([]CharacterProfile, len(c.Characters.Profile))
	for i, v := range c.Characters.Profile {
		r.Characters.Profile[i] = v.Clone()
	}

	r.Rotation = make([]ActionBlock, len(c.Rotation))
	for i, v := range c.Rotation {
		r.Rotation[i] = v.Clone()
	}

	return r
}

type SimulatorSettings struct {
	Duration   int
	DamageMode bool
	Delays     Delays

	//modes
	QueueMode  SimulationQueueMode
	ERCalcMode bool

	//other stuff
	NumberOfWorkers int // how many workers to run the simulation
	Iterations      int // how many iterations to run
}

type SimulationQueueMode int

const (
	ActionPriorityList SimulationQueueMode = iota
	SequentialList
)

// type RunOpt struct {
// 	LogDetails bool `json:"log_details"`
// 	Iteration  int  `json:"iter"`
// 	Workers    int  `json:"workers"`
// 	Duration   int  `json:"seconds"`
// 	Debug      bool `json:"debug"`
// 	ERCalcMode bool `json:"er_calc_mode"`
// 	DebugPaths []string
// }

type CharacterProfile struct {
	Base         CharacterBase             `json:"base"`
	Weapon       WeaponProfile             `json:"weapon"`
	Talents      TalentProfile             `json:"talents"`
	Stats        []float64                 `json:"stats"`
	StatsByLabel map[string][]float64      `json:"stats_by_label"`
	Sets         map[string]int            `json:"sets"`
	SetParams    map[string]map[string]int `json:"-"`
	Params       map[string]int            `json:"-"`
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

type CharacterBase struct {
	Key      CharKey `json:"key"`
	Name     string  `json:"name"`
	Element  EleType `json:"element"`
	Level    int     `json:"level"`
	MaxLevel int     `json:"max_level"`
	HP       float64 `json:"base_hp"`
	Atk      float64 `json:"base_atk"`
	Def      float64 `json:"base_def"`
	Cons     int     `json:"cons"`
	StartHP  float64 `json:"start_hp"`
}

type WeaponProfile struct {
	Name     string         `json:"name"`
	Key      string         `json:"key"` //use this to match with weapon curve mapping
	Class    WeaponClass    `json:"-"`
	Refine   int            `json:"refine"`
	Level    int            `json:"level"`
	MaxLevel int            `json:"max_level"`
	Atk      float64        `json:"base_atk"`
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

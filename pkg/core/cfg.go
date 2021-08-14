package core

type Config struct {
	Label      string
	RunOptions struct {
		Debug      bool
		Duration   int
		Iteration  int
		Workers    int
		DamageMode bool
	}
	Targets    []EnemyProfile
	Characters struct {
		Initial string
		Profile []CharacterProfile
	}
	Rotation []Action

	Hurt      HurtEvent
	FixedRand bool //if this is true then use the same seed
}

type LogConfig struct {
	LogLevel      string
	LogFile       string
	LogShowCaller bool
}

type CharacterProfile struct {
	Base    CharacterBase
	Weapon  WeaponProfile
	Talents TalentProfile
	Stats   []float64
	Sets    map[string]int
}

type CharacterBase struct {
	Name    string
	Element EleType
	Level   int
	HP      float64
	Atk     float64
	Def     float64
	Cons    int
	StartHP float64
}

type WeaponProfile struct {
	Name   string
	Class  WeaponClass
	Refine int
	Atk    float64
	Param  map[string]int
}

type TalentProfile struct {
	Attack int
	Skill  int
	Burst  int
}

type EnemyProfile struct {
	Level  int
	HP     float64
	Resist map[EleType]float64
}

type HurtEvent struct {
	WillHurt bool
	Once     bool //how often
	Start    int  //
	End      int
	Min      float64
	Max      float64
	Ele      EleType
}

func CloneEnemy(e EnemyProfile) EnemyProfile {
	r := EnemyProfile{
		Level:  e.Level,
		Resist: make(map[EleType]float64),
	}
	for k, v := range e.Resist {
		r.Resist[k] = v
	}
	return r
}

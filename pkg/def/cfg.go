package def

type Config struct {
	Label string
	Mode  struct {
		HPMode     bool
		HP         float64
		FrameLimit int
	}
	Targets    []EnemyProfile
	Characters struct {
		Initial string
		Profile []CharacterProfile
	}
	Rotation []Action

	Hurt      HurtEvent
	FixedRand bool //if this is true then use the same seed
	LogConfig LogConfig
}

type LogConfig struct {
	LogLevel      string
	LogFile       string
	LogShowCaller bool
	LogEvents     []bool
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
	Resist map[EleType]float64
}

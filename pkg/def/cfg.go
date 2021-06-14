package def

type CharacterProfile struct {
	Base    CharacterBase
	Weapon  WeaponProfile
	Talents TalentProfile
	Stats   []float64
	Sets    map[string]int
}

type CharacterBase struct {
	Name    string  `yaml:"Name"`
	Element EleType `yaml:"Element"`
	Level   int     `yaml:"Level"`
	HP      float64 `yaml:"BaseHP"`
	Atk     float64 `yaml:"BaseAtk"`
	Def     float64 `yaml:"BaseDef"`
	Cons    int     `yaml:"Constellation"`
	StartHP float64
}

type WeaponProfile struct {
	Name   string      `yaml:"WeaponName"`
	Class  WeaponClass `yaml:"WeaponClass"`
	Refine int         `yaml:"WeaponRefinement"`
	Atk    float64     `yaml:"WeaponBaseAtk"`
	Param  map[string]int
}

type TalentProfile struct {
	Attack int
	Skill  int
	Burst  int
}

type EnemyProfile struct {
	Level  int                 `yaml:"Level"`
	Resist map[EleType]float64 `yaml:"Resist"`
}

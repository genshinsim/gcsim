package character

import "github.com/genshinsim/gcsim/pkg/core/attributes"

type CharacterProfile struct {
	Base      CharacterBase             `json:"base"`
	Weapon    WeaponProfile             `json:"weapon"`
	Talents   TalentProfile             `json:"talents"`
	Stats     []float64                 `json:"stats"`
	Sets      map[string]int            `json:"sets"`
	SetParams map[string]map[string]int `json:"-"`
	Params    map[string]int            `json:"-"`
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
	Key      CharKey            `json:"key"`
	Name     string             `json:"name"`
	Element  attributes.Element `json:"element"`
	Level    int                `json:"level"`
	MaxLevel int                `json:"max_level"`
	HP       float64            `json:"base_hp"`
	Atk      float64            `json:"base_atk"`
	Def      float64            `json:"base_def"`
	Cons     int                `json:"cons"`
	StartHP  float64            `json:"start_hp"`
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

type ZoneType int

const (
	ZoneMondstadt ZoneType = iota
	ZoneLiyue
	ZoneInazuma
)

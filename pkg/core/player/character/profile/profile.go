package profile

import (
	"strconv"
	"strings"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
)

type CharacterProfile struct {
	Base         CharacterBase               `json:"base"`
	Weapon       weapon.WeaponProfile        `json:"weapon"`
	Talents      TalentProfile               `json:"talents"`
	Stats        []float64                   `json:"stats"`
	StatsByLabel map[string][]float64        `json:"stats_by_label"`
	Sets         Sets                        `json:"sets"`
	SetParams    map[keys.Set]map[string]int `json:"-"`
	Params       map[string]int              `json:"-"`
}

type Sets map[keys.Set]int

func (s Sets) MarshalJSON() ([]byte, error) {
	//we'll use a custom string builder i guess
	var sb strings.Builder
	sb.WriteString("{")
	for k, v := range s {
		sb.WriteString(`"`)
		sb.WriteString(k.String())
		sb.WriteString(`":"`)
		sb.WriteString(strconv.Itoa(v))
		sb.WriteString(`",`)
	}
	str := sb.String()
	str = strings.TrimRight(str, ",")
	str = str + "}"
	return []byte(str), nil
}

func (c *CharacterProfile) Clone() CharacterProfile {
	r := *c
	r.Weapon.Params = make(map[string]int)
	for k, v := range c.Weapon.Params {
		r.Weapon.Params[k] = v
	}
	r.Stats = make([]float64, len(c.Stats))
	copy(r.Stats, c.Stats)
	r.Sets = make(map[keys.Set]int)
	for k, v := range c.Sets {
		r.Sets[k] = v
	}

	return r
}

type CharacterBase struct {
	Key       keys.Char          `json:"key"`
	Rarity    int                `json:"rarity"`
	Element   attributes.Element `json:"element"`
	Level     int                `json:"level"`
	MaxLevel  int                `json:"max_level"`
	Ascension int                `json:"ascension"`
	HP        float64            `json:"base_hp"`
	Atk       float64            `json:"base_atk"`
	Def       float64            `json:"base_def"`
	Cons      int                `json:"cons"`
}

type TalentProfile struct {
	Attack int `json:"attack"`
	Skill  int `json:"skill"`
	Burst  int `json:"burst"`
}

type BodyType int

const (
	BodyBoy BodyType = iota
	BodyGirl
	BodyMale
	BodyLady
	BodyLoli
)

type ZoneType int

const (
	ZoneUnknown ZoneType = iota
	ZoneMondstadt
	ZoneLiyue
	ZoneInazuma
	ZoneSumeru
	ZoneFontaine
	ZoneSnezhnaya
)

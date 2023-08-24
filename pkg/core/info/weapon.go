package info

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/genshinsim/gcsim/pkg/core/keys"
)

type WeaponClass int

const (
	WeaponClassSword WeaponClass = iota
	WeaponClassClaymore
	WeaponClassSpear
	WeaponClassBow
	WeaponClassCatalyst
	EndWeaponClass
)

var weaponName = []string{
	"sword",
	"claymore",
	"polearm",
	"bow",
	"catalyst",
}

func (w WeaponClass) String() string {
	return weaponName[w]
}

func (e WeaponClass) MarshalJSON() ([]byte, error) {
	return json.Marshal(weaponName[e])
}

func (e *WeaponClass) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	s = strings.ToLower(s)
	for i, v := range weaponName {
		if v == s {
			*e = WeaponClass(i)
			return nil
		}
	}
	return errors.New("unrecognized element")
}

type WeaponProfile struct {
	Name     string         `json:"name"`
	Key      keys.Weapon    `json:"key"` //use this to match with weapon curve mapping
	Class    WeaponClass    `json:"-"`
	Refine   int            `json:"refine"`
	Level    int            `json:"level"`
	MaxLevel int            `json:"max_level"`
	Atk      float64        `json:"base_atk"`
	Params   map[string]int `json:"-"`
}
type Weapon interface {
	SetIndex(int)
	Init() error
}

package weapon

import "github.com/genshinsim/gcsim/pkg/core/keys"

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
}

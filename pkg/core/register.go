package core

import (
	"sync"

	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

var (
	mu             sync.RWMutex
	NewCharFuncMap = make(map[keys.Char]NewCharacterFunc)
	setMap         = make(map[keys.Set]NewSetFunc)
	weaponMap      = make(map[keys.Weapon]NewWeaponFunc)
)

type NewCharacterFunc func(core *Core, char *character.CharWrapper, p info.CharacterProfile) error
type NewSetFunc func(core *Core, char *character.CharWrapper, count int, param map[string]int) (info.Set, error)
type NewWeaponFunc func(core *Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error)

func RegisterCharFunc(char keys.Char, f NewCharacterFunc) {
	mu.Lock()
	defer mu.Unlock()
	if _, dup := NewCharFuncMap[char]; dup {
		panic("combat: RegisterChar called twice for character " + char.String())
	}
	NewCharFuncMap[char] = f
}

func RegisterSetFunc(set keys.Set, f NewSetFunc) {
	mu.Lock()
	defer mu.Unlock()
	if _, dup := setMap[set]; dup {
		panic("combat: RegisterSetBonus called twice for character " + set.String())
	}
	setMap[set] = f
}

func RegisterWeaponFunc(weap keys.Weapon, f NewWeaponFunc) {
	mu.Lock()
	defer mu.Unlock()
	if _, dup := weaponMap[weap]; dup {
		panic("combat: RegisterWeapon called twice for character " + weap.String())
	}
	weaponMap[weap] = f
}

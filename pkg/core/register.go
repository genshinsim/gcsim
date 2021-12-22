package core

import (
	"sync"

	"github.com/genshinsim/gcsim/pkg/core/keys"
)

var (
	mu        sync.RWMutex
	charMap   = make(map[keys.Char]NewCharacterFunc)
	setMap    = make(map[string]NewSetFunc)
	weaponMap = make(map[string]NewWeaponFunc)
)

type NewCharacterFunc func(core *Core, p CharacterProfile) (Character, error)
type NewSetFunc func(c Character, core *Core, count int)
type NewWeaponFunc func(c Character, core *Core, r int, param map[string]int)

func RegisterCharFunc(char keys.Char, f NewCharacterFunc) {
	mu.Lock()
	defer mu.Unlock()
	if _, dup := charMap[char]; dup {
		panic("combat: RegisterChar called twice for character " + char.String())
	}
	charMap[char] = f
}

func RegisterSetFunc(name string, f NewSetFunc) {
	mu.Lock()
	defer mu.Unlock()
	if _, dup := setMap[name]; dup {
		panic("combat: RegisterSetBonus called twice for character " + name)
	}
	setMap[name] = f
}

func RegisterWeaponFunc(name string, f NewWeaponFunc) {
	mu.Lock()
	defer mu.Unlock()
	if _, dup := weaponMap[name]; dup {
		panic("combat: RegisterWeapon called twice for character " + name)
	}
	weaponMap[name] = f
}

package core

import (
	"sync"

	"github.com/genshinsim/gcsim/pkg/coretype"
)

var (
	mu        sync.RWMutex
	charMap   = make(map[coretype.CharKey]NewCharacterFunc)
	setMap    = make(map[string]NewSetFunc)
	weaponMap = make(map[string]NewWeaponFunc)
)

type NewCharacterFunc func(core *Core, p coretype.CharacterProfile) (coretype.Character, error)
type NewSetFunc func(c coretype.Character, core *Core, count int, param map[string]int)
type NewWeaponFunc func(c coretype.Character, core *Core, r int, param map[string]int) string

func RegisterCharFunc(char coretype.CharKey, f NewCharacterFunc) {
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

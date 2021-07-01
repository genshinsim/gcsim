package combat

import (
	"sync"

	"github.com/genshinsim/gsim/pkg/def"
	"go.uber.org/zap"
)

var (
	mu        sync.RWMutex
	charMap   = make(map[string]NewCharacterFunc)
	setMap    = make(map[string]NewSetFunc)
	weaponMap = make(map[string]NewWeaponFunc)
)

type NewCharacterFunc func(s def.Sim, log *zap.SugaredLogger, p def.CharacterProfile) (def.Character, error)
type NewSetFunc func(c def.Character, s def.Sim, count int)
type NewWeaponFunc func(c def.Character, s def.Sim, log def.Logger, r int, param map[string]int)

func RegisterCharFunc(name string, f NewCharacterFunc) {
	mu.Lock()
	defer mu.Unlock()
	if _, dup := charMap[name]; dup {
		panic("combat: RegisterChar called twice for character " + name)
	}
	charMap[name] = f
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

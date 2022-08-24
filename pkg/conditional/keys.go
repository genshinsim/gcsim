package conditional

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/shortcut"
)

func evalKeys(c *core.Core, fields []string) (int, error) {
	// .keys.weapon.polarstar
	if err := fieldsCheck(fields, 3, "keys"); err != nil {
		return 0, err
	}

	name := fields[2]
	switch typ := fields[1]; typ {
	case "weapon":
		return evalWeaponKey(name)
	case "set":
		return evalSetKey(name)
	case "char": // is this necessary? :pepela:
		return evalCharacterKey(name)
	default:
		return 0, fmt.Errorf("bad key condition: invalid type %v", typ)
	}
}

func evalWeaponKey(name string) (int, error) {
	key, ok := shortcut.WeaponNameToKey[name]
	if !ok {
		return 0, fmt.Errorf("bad key condition: invalid weapon %v", name)
	}
	return int(key), nil
}

func evalSetKey(name string) (int, error) {
	key, ok := shortcut.SetNameToKey[name]
	if !ok {
		return 0, fmt.Errorf("bad key condition: invalid artifact set %v", name)
	}
	return int(key), nil
}

func evalCharacterKey(name string) (int, error) {
	key, ok := shortcut.CharNameToKey[name]
	if !ok {
		return 0, fmt.Errorf("bad key condition: invalid character %v", name)
	}
	return int(key), nil
}

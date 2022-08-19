package conditional

import (
	"fmt"
	"strings"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/shortcut"
)

func evalWeaponKey[V core.Number](c *core.Core, fields []string) (V, error) {
	name := strings.TrimPrefix(fields[2], ".")
	key, ok := shortcut.WeaponNameToKey[name]
	if !ok {
		return 0, fmt.Errorf("bad key condition: invalid weapon %v", name)
	}
	return V(key), nil
}

func evalSetKey[V core.Number](c *core.Core, fields []string) (V, error) {
	name := strings.TrimPrefix(fields[2], ".")
	key, ok := shortcut.SetNameToKey[name]
	if !ok {
		return 0, fmt.Errorf("bad key condition: invalid artifact set %v", name)
	}
	return V(key), nil
}

func evalCharacterKey[V core.Number](c *core.Core, fields []string) (V, error) {
	name := strings.TrimPrefix(fields[2], ".")
	key, ok := shortcut.CharNameToKey[name]
	if !ok {
		return 0, fmt.Errorf("bad key condition: invalid character %v", name)
	}
	return V(key), nil
}

func evalKeys[V core.Number](c *core.Core, fields []string) (V, error) {
	// .keys.weapon.polarstar
	if err := fieldsCheck(fields, 3, "keys"); err != nil {
		return 0, err
	}

	name := strings.TrimPrefix(fields[1], ".")
	switch name {
	case "weapon":
		return evalWeaponKey[V](c, fields)
	case "set":
		return evalSetKey[V](c, fields)
	case "char": // is this necessary? :pepela:
		return evalCharacterKey[V](c, fields)
	default:
		return 0, fmt.Errorf("bad key condition: invalid type %v", name)
	}
}

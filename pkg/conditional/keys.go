package conditional

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/shortcut"
)

func evalKeys(fields []string) (int, error) {
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
	case "element":
		return evalElementKey(name)
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

func evalElementKey(name string) (int, error) {
	key := attributes.StringToEle(name)
	switch key {
	case attributes.UnknownElement, attributes.EndEleType, attributes.NoElement:
		return 0, fmt.Errorf("bad key condition: invalid element %v", name)
	}
	return int(key), nil
}

func evalAction(fields []string) (int, error) {
	if err := fieldsCheck(fields, 1, "action"); err != nil {
		return 0, err
	}

	a := action.StringToAction(fields[1])
	if a == action.InvalidAction {
		return 0, fmt.Errorf("bad action condition: invalid action %v", fields[1])
	}
	return int(a), nil
}

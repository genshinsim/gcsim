package conditional

import (
	"strings"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/shortcut"
)

func evalWeaponKey(c *core.Core, fields []string) int64 {
	name := strings.TrimPrefix(fields[2], ".")
	key, ok := shortcut.WeaponNameToKey[name]
	if !ok {
		c.Log.NewEvent("bad keys conditon: invalid weapon", glog.LogWarnings, -1).Write("fields", fields)
		return -1
	}
	return int64(key)
}

func evalSetKey(c *core.Core, fields []string) int64 {
	name := strings.TrimPrefix(fields[2], ".")
	key, ok := shortcut.SetNameToKey[name]
	if !ok {
		c.Log.NewEvent("bad keys conditon: invalid set", glog.LogWarnings, -1).Write("fields", fields)
		return -1
	}
	return int64(key)
}

func evalCharacterKey(c *core.Core, fields []string) int64 {
	name := strings.TrimPrefix(fields[2], ".")
	key, ok := shortcut.CharNameToKey[name]
	if !ok {
		c.Log.NewEvent("bad keys conditon: invalid character", glog.LogWarnings, -1).Write("fields", fields)
		return -1
	}
	return int64(key)
}

func evalKeys(c *core.Core, fields []string) int64 {
	// .keys.weapon.polarstar
	if len(fields) < 3 {
		c.Log.NewEvent("bad keys conditon: invalid num of fields", glog.LogWarnings, -1).Write("fields", fields)
		return -1
	}

	name := strings.TrimPrefix(fields[1], ".")
	switch name {
	case "weapon":
		return evalWeaponKey(c, fields)
	case "set":
		return evalSetKey(c, fields)
	case "char": // is this necessary? :pepela:
		return evalCharacterKey(c, fields)
	default:
		c.Log.NewEvent("bad keys conditon: invalid type", glog.LogWarnings, -1).Write("fields", fields)
		return -1
	}
}

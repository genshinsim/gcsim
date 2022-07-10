package conditional

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/shortcut"
)

func Eval(c *core.Core, fields []string) int64 {
	switch fields[0] {
	case ".debuff":
		return evalDebuff(c, fields)
	case ".element":
		return evalElement(c, fields)
	case ".cd":
		return evalCD(c, fields)
	case ".energy":
		return evalEnergy(c, fields)
	case ".status":
		return evalStatus(c, fields)
	case ".tags":
		return evalTags(c, fields)
	case ".stam":
		return evalStam(c, fields)
	case ".ready":
		return evalReady(c, fields)
	case ".mods":
		return evalCharMods(c, fields)
	case ".infusion":
		return evalInfusion(c, fields)
	case ".construct":
		return evalConstruct(c, fields)
	case ".normal":
		return evalNormalCounter(c, fields)
	default:
		//check if it's a char name; if so check char custom eval func
		if key, ok := shortcut.CharNameToKey[fields[0]]; ok {
			return evalCharCustom(c, key, fields)
		}
		return 0
	}
}

package conditional

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

func evalCharacter(c *core.Core, key keys.Char, fields []string) (any, error) {
	if err := fieldsCheck(fields, 2, "character"); err != nil {
		return 0, err
	}
	char, ok := c.Player.ByKey(key)
	if !ok {
		return 0, fmt.Errorf("character %v not in team when evaluating condition", key)
	}

	typ := fields[1]
	switch typ {
	case "cons":
		return char.Base.Cons, nil
	case "energy":
		return char.Energy, nil
	case "energymax":
		return char.EnergyMax, nil
	case "normal":
		return char.NextNormalCounter(), nil
	case "onfield":
		return c.Player.Active() == char.Index, nil
	case "weapon":
		return int(char.Weapon.Key), nil
	}

	// call character condition early
	if err := fieldsCheck(fields, 3, "character "+fields[1]); err != nil {
		// .kokomi.<cond>
		return char.Condition(fields[1:])
	}
	val := fields[2]

	// special case for ability conditions. since typ/val are swapped
	// .kokomi.<abil>.<cond>
	act := action.StringToAction(typ)
	if act != action.InvalidAction {
		return evalCharacterAbil(char, act, val)
	}

	switch typ {
	case "status", "mods":
		return char.StatusIsActive(val), nil
	case "infusion":
		return c.Player.WeaponInfuseIsActive(char.Index, val), nil
	case "tags":
		return char.Tag(val), nil
	case "stats":
		return evalCharacterStats(char, val)
	default: // .kokomi.<cond>.*
		return char.Condition(fields[1:])
	}
}

func evalCharacterStats(char *character.CharWrapper, stat string) (float64, error) {
	key := attributes.StrToStatType(stat)
	if key == -1 {
		return 0, fmt.Errorf("invalid stat key %v in character stat condition", stat)
	}
	return char.Stat(key), nil
}

func evalCharacterAbil(char *character.CharWrapper, act action.Action, typ string) (any, error) {
	switch typ {
	case "cd":
		return char.Cooldown(act), nil
	case "ready":
		//TODO: nil map may cause problems here??
		return char.ActionReady(act, nil), nil
	default:
		return 0, fmt.Errorf("bad character ability condition: invalid type %v", typ)
	}
}

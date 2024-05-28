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

	// special case for ability conditions. since fields are swapped
	// .kokomi.<abil>.<cond>
	typ := fields[1]
	act := action.StringToAction(typ)
	if act != action.InvalidAction {
		if err := fieldsCheck(fields, 3, "character ability"); err != nil {
			return 0, err
		}
		return evalCharacterAbil(c, char, act, fields[2])
	}

	charCat := "character " + typ

	switch typ {
	case "id":
		return int(char.Base.Key), nil
	case "cons":
		return char.Base.Cons, nil
	case "energy":
		return char.Energy, nil
	case "energymax":
		return char.EnergyMax, nil
	case "hp":
		return char.CurrentHP(), nil
	case "hpmax":
		return char.MaxHP(), nil
	case "hpratio":
		return char.CurrentHPRatio(), nil
	case "normal":
		return char.NextNormalCounter(), nil
	case "onfield":
		return c.Player.Active() == char.Index, nil
	case "weapon":
		return int(char.Weapon.Key), nil
	case "status":
		if err := fieldsCheck(fields, 3, charCat); err != nil {
			return 0, err
		}
		return char.StatusDuration(fields[2]), nil
	case "mods":
		if err := fieldsCheck(fields, 3, charCat); err != nil {
			return 0, err
		}
		return char.StatusDuration(fields[2]), nil
	case "infusion":
		if err := fieldsCheck(fields, 3, charCat); err != nil {
			return 0, err
		}
		return c.Player.WeaponInfuseIsActive(char.Index, fields[2]), nil
	case "tags":
		if err := fieldsCheck(fields, 3, charCat); err != nil {
			return 0, err
		}
		return char.Tag(fields[2]), nil
	case "stats":
		if err := fieldsCheck(fields, 3, charCat); err != nil {
			return 0, err
		}
		return evalCharacterStats(char, fields[2])
	case "bol":
		return char.CurrentHPDebt(), nil
	case "bolratio":
		return char.CurrentHPDebt() / char.MaxHP(), nil
	default: // .kokomi.*
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

func evalCharacterAbil(c *core.Core, char *character.CharWrapper, act action.Action, typ string) (any, error) {
	switch typ {
	case "cd":
		if act == action.ActionSwap {
			return c.Player.SwapCD, nil
		}
		return char.Cooldown(act), nil
	case "charge":
		return char.Charges(act), nil
	case "ready":
		if act == action.ActionSwap {
			return c.Player.SwapCD == 0 || c.Player.Active() == char.Index, nil
		}
		// TODO: nil map may cause problems here??
		ok, _ := char.ActionReady(act, nil)
		return ok, nil
	default:
		return 0, fmt.Errorf("bad character ability condition: invalid type %v", typ)
	}
}

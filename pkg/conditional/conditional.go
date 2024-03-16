package conditional

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/shortcut"
)

const countField = "count"

func fieldsCheck(fields []string, expecting int, category string) error {
	if len(fields) < expecting {
		return fmt.Errorf(
			"bad %v condition; invalid num of fields; expecting at least %v; got %v",
			category,
			expecting,
			len(fields),
		)
	}
	return nil
}

func Eval(c *core.Core, fields []string) (any, error) {
	switch fields[0] {
	case "debuff":
		return evalDebuff(c, fields)
	case "element":
		return evalElement(c, fields)
	case "status":
		if err := fieldsCheck(fields, 2, "status"); err != nil {
			return 0, err
		}
		return c.Status.Duration(fields[1]), nil
	case "stam":
		return c.Player.Stam, nil
	case "construct":
		return evalConstruct(c, fields)
	case "gadgets":
		return evalGadgets(c, fields)
	case "keys":
		return evalKeys(fields)
	case "state":
		return int(c.Player.CurrentState()), nil
	case "action":
		return evalAction(fields)
	case "previous-action":
		return int(c.Player.LastAction.Type), nil
	case "previous-char":
		return int(c.Player.ByIndex(c.Player.LastAction.Char).Base.Key), nil
	case "airborne":
		return c.Player.Airborne() != player.Grounded, nil
	default:
		// check if it's a char name; if so check char custom eval func
		name := fields[0]
		if key, ok := shortcut.CharNameToKey[name]; ok {
			return evalCharacter(c, key, fields)
		}
		return 0, fmt.Errorf("invalid character %v in character condition", name)
	}
}

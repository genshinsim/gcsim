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
	switch typ {
	case "cd":
		return evalCharacterCooldown(char, val)
	case "status", "mods":
		return char.StatusIsActive(val), nil
	case "infusion":
		return c.Player.WeaponInfuseIsActive(char.Index, val), nil
	case "tags":
		return char.Tag(val), nil
	case "ready":
		return evalCharacterAbilReady(char, val)
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

func evalCharacterAbilReady(char *character.CharWrapper, abil string) (bool, error) {
	ak := action.StringToAction(abil)
	if ak == action.InvalidAction {
		return false, fmt.Errorf("invalid abil %v in ready condition", abil)
	}
	//TODO: nil map may cause problems here??
	return char.ActionReady(ak, nil), nil
}

func evalCharacterCooldown(char *character.CharWrapper, abil string) (int, error) {
	switch abil {
	case "skill":
		return char.Cooldown(action.ActionSkill), nil
	case "burst":
		return char.Cooldown(action.ActionBurst), nil
	default:
		return 0, fmt.Errorf("invalid ability %v in character cooldown condition", abil)
	}
}

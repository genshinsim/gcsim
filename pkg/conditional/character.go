package conditional

import (
	"fmt"
	"strings"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

func fieldsCheck(fields []string, expecting int, category string) error {
	if len(fields) < expecting {
		return fmt.Errorf("bad %v condition; invalid num of fields; expecting at least %v; got %v", category, expecting, len(fields))
	}
	return nil
}

func evalCharacter[V core.Number](c *core.Core, key keys.Char, fields []string) (V, error) {
	if err := fieldsCheck(fields, 2, "character"); err != nil {
		return 0, err
	}
	char, ok := c.Player.ByKey(key)
	if !ok {
		return 0, fmt.Errorf("character %v not in team when evaluating condition", key)
	}

	switch fields[1] {
	case ".cd":
		if err := fieldsCheck(fields, 3, "character cooldown"); err != nil {
			return 0, err
		}
		return evalCharacterCooldown[V](char, fields[2])
	case ".energy":
		return V(char.Energy), nil
	case ".status":
		if err := fieldsCheck(fields, 3, "character status"); err != nil {
			return 0, err
		}
		return evalCharacterMods[V](char, fields[2])
	case ".mods":
		if err := fieldsCheck(fields, 3, "character mods"); err != nil {
			return 0, err
		}
		return evalCharacterStatus[V](char, fields[2])
	case ".infusion":
		if err := fieldsCheck(fields, 3, "character stats"); err != nil {
			return 0, err
		}
		if c.Player.WeaponInfuseIsActive(char.Index, strings.TrimPrefix(fields[2], ".")) {
			return 1, nil
		}
		return 0, nil
	case ".normal":
		return V(char.NextNormalCounter()), nil
	case ".onfield":
		if c.Player.Active() == char.Index {
			return 1, nil
		}
		return 0, nil
	case ".weapon":
		return V(char.Weapon.Key), nil
	case ".cons":
		return V(char.Base.Cons), nil
	case ".tags":
		return V(char.Tag(strings.TrimPrefix(fields[2], "."))), nil
	case ".ready":
		if err := fieldsCheck(fields, 3, "character abil ready"); err != nil {
			return 0, err
		}
		return evalCharacterAbilReady[V](char, fields[2])
	case ".stats":
		if err := fieldsCheck(fields, 3, "character stats"); err != nil {
			return 0, err
		}
		return evalCharacterStats[V](char, fields[2])
	default:
		return V(char.Condition(strings.TrimPrefix(fields[1], "."))), nil
	}
}

func evalCharacterStats[V core.Number](char *character.CharWrapper, stat string) (V, error) {
	key := attributes.StrToStatType(strings.TrimPrefix(stat, "."))
	if key == -1 {
		return 0, fmt.Errorf("invalid stat key %v in character stat condition", stat)
	}
	return V(char.Stat(key)), nil
}

func evalCharacterMods[V core.Number](char *character.CharWrapper, mod string) (V, error) {
	mod = strings.TrimPrefix(mod, ".")
	if char.StatModIsActive(mod) {
		return 1, nil
	}
	return 0, nil
}

func evalCharacterStatus[V core.Number](char *character.CharWrapper, mod string) (V, error) {
	mod = strings.TrimPrefix(mod, ".")
	if char.StatusIsActive(mod) {
		return 1, nil
	}
	return 0, nil
}

func evalCharacterAbilReady[V core.Number](char *character.CharWrapper, abil string) (V, error) {
	abil = strings.TrimPrefix(abil, ".")
	ak := action.StringToAction(abil)
	if ak == action.InvalidAction {
		return 0, fmt.Errorf("invalid abil %v in ready condition", abil)
	}

	//TODO: nil map may cause problems here??
	if char.ActionReady(ak, nil) {
		return 1, nil
	}

	return 0, nil
}

func evalCharacterCooldown[V core.Number](char *character.CharWrapper, abil string) (V, error) {
	switch abil {
	case ".skill":
		return V(char.Cooldown(action.ActionSkill)), nil
	case ".burst":
		return V(char.Cooldown(action.ActionBurst)), nil
	default:
		return 0, fmt.Errorf("invalid ability %v in character cooldown condition", abil)
	}
}

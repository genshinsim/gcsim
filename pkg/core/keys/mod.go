package keys

import (
	"encoding/json"
	"errors"
	"strings"
)

type Modifier int

func (m *Modifier) MarshalJSON() ([]byte, error) {
	return json.Marshal(modNames[*m])
}

func (m *Modifier) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	s = strings.ToLower(s)
	for i := range modNames {
		if modNames[i] == s {
			*m = Modifier(i)
			return nil
		}
	}
	return errors.New("unrecognized character key")
}

func (m Modifier) String() string {
	return modNames[m]
}

func (m Modifier) IsSpecialDecay() bool {
	switch m {
	case Dendro:
	case Quicken:
	case Frozen:
	case Anemo:
	case Geo:
	case Burning:
	default:
		return false
	}
	return true
}

const (
	InvalidModifier Modifier = iota
	TestingMod
	Electro
	Pyro
	Cryo
	Hydro
	BurningFuel
	Dendro
	Quicken
	Frozen
	Anemo
	Geo
	Burning
	BuiltinModifierDelim // delim
	// TODO: everything below here to EndModifierKeys should be generated
	EndModifierKeys
)

var modNames = [EndModifierKeys]string{
	"invalid",
	"testing",
	"electro",
	"pyro",
	"cryo",
	"hydro",
	"burning_fuel",
	"dendro",
	"quicken",
	"frozen",
	"anemo",
	"geo",
	"burning",
	"invalid", // delim
}

package info

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
)

// TODO: this needs to be eventually refactored into just generic modifier keys but for now
// keeping this to make it easier to find/access reactions
type ReactionModKey int

const (
	ReactionModKeyInvalid ReactionModKey = iota
	ReactionModKeyElectro
	ReactionModKeyPyro
	ReactionModKeyCryo
	ReactionModKeyHydro
	ReactionModKeyBurningFuel
	ReactionModKeySpecialDecayDelim
	ReactionModKeyDendro
	ReactionModKeyQuicken
	ReactionModKeyFrozen
	ReactionModKeyAnemo
	ReactionModKeyGeo
	ReactionModKeyBurning
	ReactionModKeyEnd
)

var ModifierString = []string{
	"",
	"electro",
	"pyro",
	"cryo",
	"hydro",
	"dendro-fuel",
	"",
	"dendro",
	"quicken",
	"frozen",
	"anemo",
	"geo",
	"burning",
	"",
}

var modifierElement = []attributes.Element{
	attributes.UnknownElement,
	attributes.Electro,
	attributes.Pyro,
	attributes.Cryo,
	attributes.Hydro,
	attributes.Dendro,
	attributes.UnknownElement,
	attributes.Dendro,
	attributes.Quicken,
	attributes.Frozen,
	attributes.Anemo,
	attributes.Geo,
	attributes.Pyro,
	attributes.UnknownElement,
}

func (r ReactionModKey) Element() attributes.Element { return modifierElement[r] }
func (r ReactionModKey) String() string              { return ModifierString[r] }

func (r ReactionModKey) MarshalJSON() ([]byte, error) {
	return json.Marshal(ModifierString[r])
}

func (r *ReactionModKey) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	s = strings.ToLower(s)
	for i, v := range ModifierString {
		if v == s {
			*r = ReactionModKey(i)
			return nil
		}
	}
	return errors.New("unrecognized ReactableModifier")
}

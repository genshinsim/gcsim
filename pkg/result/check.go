package result

import (
	"slices"

	"github.com/genshinsim/gcsim/pkg/core/keys"
)

// adding a character to this list will make the "incomplete warning" appear on the viewer
var incompleteCharacters = []keys.Char{
	keys.TestCharDoNotUse,
}

func IsCharacterComplete(char keys.Char) bool {
	return !slices.Contains(incompleteCharacters, char)
}

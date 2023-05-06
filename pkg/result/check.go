package result

import (
	"github.com/genshinsim/gcsim/pkg/core/keys"
)

// adding a character to this list will make the "incomplete warning" appear on the viewer
var incompleteCharacters = []keys.Char{
	keys.TestCharDoNotUse,
}

func IsCharacterComplete(char keys.Char) bool {
	for _, v := range incompleteCharacters {
		if v == char {
			return false
		}
	}
	return true
}

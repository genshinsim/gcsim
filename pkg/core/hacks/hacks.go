package hacks

import "github.com/genshinsim/gcsim/pkg/core/keys"

var noblesseSpecialChars = [keys.EndCharKeys]bool{}

func RegisterNOSpecialChar(k keys.Char) {
	noblesseSpecialChars[k] = true
}

func NOCharIsSpecial(k keys.Char) bool {
	return noblesseSpecialChars[k]
}

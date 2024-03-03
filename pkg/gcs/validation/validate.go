package validation

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/keys"
)

func ValidateCharParamKeys(c keys.Char, a action.Action, keys []string) error {
	f, ok := charValidParamKeys[c]
	if !ok {
		// all is ok if no validation function registered
		return nil
	}
	return f(a, keys)
}

package validation

import (
	"slices"

	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/keys"
)

// generic params that can be used for any character
var ignoreParams = []string{
	// iansan burst
	"movement",
}

func ValidateCharParamKeys(c keys.Char, a action.Action, keys []string) error {
	f, ok := charValidParamKeys[c]
	if !ok {
		// all is ok if no validation function registered
		return nil
	}

	filtered := make([]string, 0, len(keys))
	for _, v := range keys {
		if !slices.Contains(ignoreParams, v) {
			filtered = append(filtered, v)
		}
	}

	return f(a, filtered)
}

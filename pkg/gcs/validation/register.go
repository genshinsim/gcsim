package validation

import (
	"sync"

	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/keys"
)

var (
	mu                 sync.RWMutex
	charValidParamKeys = make(map[keys.Char]CharParamKeysValidationFunc)
)

type CharParamKeysValidationFunc func(a action.Action, p []string) error

func RegisterCharParamValidationFunc(char keys.Char, f CharParamKeysValidationFunc) {
	mu.Lock()
	defer mu.Unlock()
	if _, dup := charValidParamKeys[char]; dup {
		panic("validation: RegisterCharParamValidationFunc called twice for character " + char.String())
	}
	charValidParamKeys[char] = f
}

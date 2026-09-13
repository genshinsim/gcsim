package info

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

// THESE MODIFIERS SHOULD EVENTUALLY BE DEPRECATED

type Status struct {
	modifier.Base
}
type ResistMod struct {
	Ele   attributes.Element
	Value float64
	modifier.Base
}

type DefMod struct {
	Value float64
	Dur   int
	modifier.Base
}

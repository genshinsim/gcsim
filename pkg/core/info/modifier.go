package info

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/gmod"
)

// THESE MODIFIERS SHOULD EVENTUALLY BE DEPRECATED

type Status struct {
	gmod.Base
}
type ResistMod struct {
	Ele   attributes.Element
	Value float64
	gmod.Base
}

type DefMod struct {
	Value float64
	Dur   int
	gmod.Base
}

package info

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/keys"
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

// THESE ARE IDEALLY THE NEW MODIFIERS TO BE USED
type Modifier struct {
	ModifierListeners
	ModiferMixins

	Name       keys.Modifier
	Duration   int // duration in frames
	Durability Durability

	DecayRate Durability // always calculated on Add, but exposed so it can be modified

	Source   keys.TargetID
	Target   keys.TargetID
	Stacking StackingType
}

type ModifierListeners struct {
	OnAdd    func(*Modifier)
	OnRemove func(*Modifier)
	// OnBeingHit func(*Modifier, *AttackEvent)
	// OnHeal func(*Modifier))
	// OnBeingHealed func(*Modifier)
	OnThinkInterval func(*Modifier)

	PreTick  func(*Modifier)
	PostTick func(*Modifier)
}

type ModiferMixins struct {
	ModifyDamageMixin func(*Modifier, *AttackEvent)
	ModifyStatsMixin  func(*Modifier) *[]attributes.Stats
}

type StackingType int

const (
	InvalidStacking         StackingType = iota
	Refresh                              // single instance. re-apply resets durability and doesn't trigger onAdded/onRemoved or reset onThinkInterval
	Unique                               // single instance. can't be re-applied unless expired
	Prolong                              // same as Refresh. can only be re-applied within the initial duration
	Multiple                             // same as Unique. can hold multiple instances
	MultipleRefresh                      // same as Refresh. can hold multiple instances, a re-apply will trigger onAdded/onRemoved and reset onThinkInterval on the oldest instance
	MultipleRefreshNoRemove              // same as Refresh. can hold multiple instances
	MultipleAllRefresh                   // same as Refresh. can hold multiple instances, a re-apply resets durability for all active instances
	RefreshAndAddDurability              // not used
	GlobalUnique                         // not used
	RefreshUniqueDurability              // unknown behaviour
	Overlap                              // used for "auras"
	OverlapRefreshDuration               // used for "auras"
)

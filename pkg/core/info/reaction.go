package info

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
)

type ReactionType string

const (
	ReactionTypeOverload           ReactionType = "overload"
	ReactionTypeSuperconduct       ReactionType = "superconduct"
	ReactionTypeFreeze             ReactionType = "freeze"
	ReactionTypeMelt               ReactionType = "melt"
	ReactionTypeVaporize           ReactionType = "vaporize"
	ReactionTypeCrystallizeElectro ReactionType = "crystallize-electro"
	ReactionTypeCrystallizeHydro   ReactionType = "crystallize-hydro"
	ReactionTypeCrystallizePyro    ReactionType = "crystallize-pyro"
	ReactionTypeCrystallizeCryo    ReactionType = "crystallize-cryo"
	ReactionTypeSwirlElectro       ReactionType = "swirl-electro"
	ReactionTypeSwirlHydro         ReactionType = "swirl-hydro"
	ReactionTypeSwirlPyro          ReactionType = "swirl-pyro"
	ReactionTypeSwirlCryo          ReactionType = "swirl-cryo"
	ReactionTypeElectroCharged     ReactionType = "electrocharged"
	ReactionTypeShatter            ReactionType = "shatter"
	ReactionTypeBurning            ReactionType = "burning"
	ReactionTypeQuicken            ReactionType = "quicken"
	ReactionTypeAggravate          ReactionType = "aggravate"
	ReactionTypeSpread             ReactionType = "spread"
	ReactionTypeBloom              ReactionType = "bloom"
	ReactionTypeBurgeon            ReactionType = "burgeon"
	ReactionTypeHyperbloom         ReactionType = "hyperbloom"
	ReactionTypeNoReaction         ReactionType = ""
	ReactionTypeFreezeExtend       ReactionType = "freeze-extend"
)

const SelfDamageSuffix = " (self damage)"

type Durability float64

const ZeroDur Durability = 0.00000000001

type Reactable interface {
	Tick()

	React(a *AttackEvent)
	AttachOrRefill(a *AttackEvent) bool
	SetAuraDurability(mod ReactionModKey, dur Durability)
	SetAuraDecayRate(mod ReactionModKey, decay Durability)
	GetAuraDurability(mod ReactionModKey) Durability
	GetAuraDecayRate(mod ReactionModKey) Durability

	ActiveAuraString() []string
	AuraCount() int
	GetDurability() []Durability
	AuraContains(e ...attributes.Element) bool

	ReactableBloom
	ReactableBurning
	ReactableCatalyze
	ReactableCrystallize
	ReactableEC
	ReactableFreeze
	ReactableMelt
	ReactableOverload
	ReactableSuperconduct
	ReactableSwirl
	ReactableVaporize
}

type ReactableBloom interface {
	TryBloom(a *AttackEvent) bool
}

type ReactableBurning interface {
	TryBurning(a *AttackEvent) bool
	IsBurning() bool
}

type ReactableCatalyze interface {
	TryAggravate(a *AttackEvent) bool
	TrySpread(a *AttackEvent) bool
	TryQuicken(a *AttackEvent) bool
}

type ReactableCrystallize interface {
	TryCrystallizeElectro(a *AttackEvent) bool
	TryCrystallizeHydro(a *AttackEvent) bool
	TryCrystallizeCryo(a *AttackEvent) bool
	TryCrystallizePyro(a *AttackEvent) bool
	TryCrystallizeFrozen(a *AttackEvent) bool
}

type ReactableEC interface {
	TryAddEC(a *AttackEvent) bool
}

type ReactableFreeze interface {
	TryFreeze(a *AttackEvent) bool
	PoiseDMGCheck(a *AttackEvent) bool
	ShatterCheck(a *AttackEvent) bool
	SetFreezeResist(resist float64)
}

type ReactableMelt interface {
	TryMelt(a *AttackEvent) bool
}

type ReactableOverload interface {
	TryOverload(a *AttackEvent) bool
}

type ReactableSuperconduct interface {
	TrySuperconduct(a *AttackEvent) bool
	TryFrozenSuperconduct(a *AttackEvent) bool
}

type ReactableSwirl interface {
	TrySwirlElectro(a *AttackEvent) bool
	TrySwirlHydro(a *AttackEvent) bool
	TrySwirlCryo(a *AttackEvent) bool
	TrySwirlPyro(a *AttackEvent) bool
	TrySwirlFrozen(a *AttackEvent) bool
}

type ReactableVaporize interface {
	TryVaporize(a *AttackEvent) bool
}

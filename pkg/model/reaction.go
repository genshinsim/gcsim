package model

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

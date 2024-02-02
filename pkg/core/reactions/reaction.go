package reactions

type ReactionType string

const (
	Overload           ReactionType = "overload"
	Superconduct       ReactionType = "superconduct"
	Freeze             ReactionType = "freeze"
	Melt               ReactionType = "melt"
	Vaporize           ReactionType = "vaporize"
	CrystallizeElectro ReactionType = "crystallize-electro"
	CrystallizeHydro   ReactionType = "crystallize-hydro"
	CrystallizePyro    ReactionType = "crystallize-pyro"
	CrystallizeCryo    ReactionType = "crystallize-cryo"
	SwirlElectro       ReactionType = "swirl-electro"
	SwirlHydro         ReactionType = "swirl-hydro"
	SwirlPyro          ReactionType = "swirl-pyro"
	SwirlCryo          ReactionType = "swirl-cryo"
	ElectroCharged     ReactionType = "electrocharged"
	Shatter            ReactionType = "shatter"
	Burning            ReactionType = "burning"
	Quicken            ReactionType = "quicken"
	Aggravate          ReactionType = "aggravate"
	Spread             ReactionType = "spread"
	Bloom              ReactionType = "bloom"
	Burgeon            ReactionType = "burgeon"
	Hyperbloom         ReactionType = "hyperbloom"
	NoReaction         ReactionType = ""
	FreezeExtend       ReactionType = "freeze-extend"
)

const SelfDamageSuffix = " (self damage)"

type Durability float64

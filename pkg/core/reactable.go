package core

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/reactions"
	"github.com/genshinsim/gcsim/pkg/model"
)

type Reactable interface {
	Init(self combat.Target, c *Core)
	Tick()

	React(a *combat.AttackEvent)
	AttachOrRefill(a *combat.AttackEvent) bool
	SetAuraDurability(mod model.Element, dur reactions.Durability, decay reactions.Durability)

	ActiveAuraString() []string
	AuraCount() int
	GetAuraDurability(mod model.Element) reactions.Durability
	GetDurability() []reactions.Durability
	GetAuraDecayRate(mod model.Element) reactions.Durability
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
	TryBloom(a *combat.AttackEvent) bool
}

type ReactableBurning interface {
	TryBurning(a *combat.AttackEvent) bool
	IsBurning() bool
}

type ReactableCatalyze interface {
	TryAggravate(a *combat.AttackEvent) bool
	TrySpread(a *combat.AttackEvent) bool
	TryQuicken(a *combat.AttackEvent) bool
}

type ReactableCrystallize interface {
	TryCrystallizeElectro(a *combat.AttackEvent) bool
	TryCrystallizeHydro(a *combat.AttackEvent) bool
	TryCrystallizeCryo(a *combat.AttackEvent) bool
	TryCrystallizePyro(a *combat.AttackEvent) bool
	TryCrystallizeFrozen(a *combat.AttackEvent) bool
}

type ReactableEC interface {
	TryAddEC(a *combat.AttackEvent) bool
}

type ReactableFreeze interface {
	TryFreeze(a *combat.AttackEvent) bool
	PoiseDMGCheck(a *combat.AttackEvent) bool
	ShatterCheck(a *combat.AttackEvent) bool
	SetFreezeResist(resist float64)
}

type ReactableMelt interface {
	TryMelt(a *combat.AttackEvent) bool
}

type ReactableOverload interface {
	TryOverload(a *combat.AttackEvent) bool
}

type ReactableSuperconduct interface {
	TrySuperconduct(a *combat.AttackEvent) bool
	TryFrozenSuperconduct(a *combat.AttackEvent) bool
}

type ReactableSwirl interface {
	TrySwirlElectro(a *combat.AttackEvent) bool
	TrySwirlHydro(a *combat.AttackEvent) bool
	TrySwirlCryo(a *combat.AttackEvent) bool
	TrySwirlPyro(a *combat.AttackEvent) bool
	TrySwirlFrozen(a *combat.AttackEvent) bool
}

type ReactableVaporize interface {
	TryVaporize(a *combat.AttackEvent) bool
}

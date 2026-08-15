package frostbreath

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	buffKey      = "frostbreath-atk"
	energyKey    = "frostbreath-energy"
	buffICDKey   = "frostbreath-buff-icd"
	energyICDKey = "frostbreath-energy-icd"
)

type Weapon struct {
	Index  int
	char   *character.CharWrapper
	core   *core.Core
	refine int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error {
	m := make([]float64, attributes.EndStatType)
	m[attributes.ATKP] = 0.15 + float64(w.refine)*0.05

	energy := 6 + float64(w.refine)*1.5

	refundEnergy := func(char *character.CharWrapper) {
		if char.StatusIsActive(energyICDKey) {
			return
		}
		char.AddStatus(energyICDKey, 0.3*60, true)
		char.AddEnergy(energyKey, energy)
	}

	onHydroOrCryoReaction := func(args ...any) {
		if w.char.StatusIsActive(buffICDKey) {
			return
		}

		atk := args[1].(*info.AttackInfo)
		if atk.ActorIndex != w.char.Index() {
			return
		}

		w.char.AddStatus(buffICDKey, 16*60, true)

		w.char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag(buffKey, 15*60),
			AffectedStat: attributes.ATKP,
			Amount: func() []float64 {
				return m
			},
		})

		for _, char := range w.core.Player.Chars() {
			refundEnergy(char)
		}
	}

	w.core.Events.Subscribe(event.OnMelt, onHydroOrCryoReaction, "frostbreath-on-hydro-or-cryo-reaction")
	w.core.Events.Subscribe(event.OnFrozen, onHydroOrCryoReaction, "frostbreath-on-hydro-or-cryo-reaction")
	w.core.Events.Subscribe(event.OnShatter, onHydroOrCryoReaction, "frostbreath-on-hydro-or-cryo-reaction")
	w.core.Events.Subscribe(event.OnSuperconduct, onHydroOrCryoReaction, "frostbreath-on-hydro-or-cryo-reaction")
	w.core.Events.Subscribe(event.OnSwirlCryo, onHydroOrCryoReaction, "frostbreath-on-hydro-or-cryo-reaction")
	w.core.Events.Subscribe(event.OnSwirlHydro, onHydroOrCryoReaction, "frostbreath-on-hydro-or-cryo-reaction")
	w.core.Events.Subscribe(event.OnCrystallizeCryo, onHydroOrCryoReaction, "frostbreath-on-hydro-or-cryo-reaction")
	w.core.Events.Subscribe(event.OnCrystallizeHydro, onHydroOrCryoReaction, "frostbreath-on-hydro-or-cryo-reaction")
	w.core.Events.Subscribe(event.OnVaporize, onHydroOrCryoReaction, "frostbreath-on-hydro-or-cryo-reaction")
	w.core.Events.Subscribe(event.OnElectroCharged, onHydroOrCryoReaction, "frostbreath-on-hydro-or-cryo-reaction")
	w.core.Events.Subscribe(event.OnBloom, onHydroOrCryoReaction, "frostbreath-on-hydro-or-cryo-reaction")
	w.core.Events.Subscribe(event.OnLunarCharged, onHydroOrCryoReaction, "frostbreath-on-hydro-or-cryo-reaction")
	w.core.Events.Subscribe(event.OnLunarBloom, onHydroOrCryoReaction, "frostbreath-on-hydro-or-cryo-reaction")
	w.core.Events.Subscribe(event.OnLunarCrystallize, onHydroOrCryoReaction, "frostbreath-on-hydro-or-cryo-reaction")
	w.core.Events.Subscribe(event.OnStellarConduct, onHydroOrCryoReaction, "frostbreath-on-hydro-or-cryo-reaction")
	w.core.Events.Subscribe(event.OnStellarSwirl, onHydroOrCryoReaction, "frostbreath-on-hydro-or-cryo-reaction")
	// w.core.Events.Subscribe(event.OnStellarSwirlAnemo, onHydroOrCryoReaction, "frostbreath-on-hydro-or-cryo-reaction")

	return nil
}

func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{
		char:   char,
		core:   c,
		refine: p.Refine,
	}
	return w, nil
}

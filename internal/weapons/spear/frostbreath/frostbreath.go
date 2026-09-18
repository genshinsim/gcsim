package frostbreath

import (
	"fmt"

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

		atk := args[1].(*info.AttackEvent)
		if atk.Info.ActorIndex != w.char.Index() {
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

	w.core.Events.Subscribe(event.OnMelt, onHydroOrCryoReaction, fmt.Sprintf("frostbreath-on-melt-%v", w.char.Base.Key.String()))
	w.core.Events.Subscribe(event.OnFrozen, onHydroOrCryoReaction, fmt.Sprintf("frostbreath-on-frozen-%v", w.char.Base.Key.String()))
	w.core.Events.Subscribe(event.OnShatter, onHydroOrCryoReaction, fmt.Sprintf("frostbreath-on-shatter-%v", w.char.Base.Key.String()))
	w.core.Events.Subscribe(event.OnSuperconduct, onHydroOrCryoReaction, fmt.Sprintf("frostbreath-on-superconduct-%v", w.char.Base.Key.String()))
	w.core.Events.Subscribe(event.OnSwirlCryo, onHydroOrCryoReaction, fmt.Sprintf("frostbreath-on-swirl-cryo-%v", w.char.Base.Key.String()))
	w.core.Events.Subscribe(event.OnSwirlHydro, onHydroOrCryoReaction, fmt.Sprintf("frostbreath-on-swirl-hydro-%v", w.char.Base.Key.String()))
	w.core.Events.Subscribe(event.OnCrystallizeCryo, onHydroOrCryoReaction, fmt.Sprintf("frostbreath-on-crystallize-cryo-%v", w.char.Base.Key.String()))
	w.core.Events.Subscribe(event.OnCrystallizeHydro, onHydroOrCryoReaction, fmt.Sprintf("frostbreath-on-crystallize-hydro-%v", w.char.Base.Key.String()))
	w.core.Events.Subscribe(event.OnVaporize, onHydroOrCryoReaction, fmt.Sprintf("frostbreath-on-vaporize-%v", w.char.Base.Key.String()))
	w.core.Events.Subscribe(event.OnElectroCharged, onHydroOrCryoReaction, fmt.Sprintf("frostbreath-on-electro-charged-%v", w.char.Base.Key.String()))
	w.core.Events.Subscribe(event.OnBloom, onHydroOrCryoReaction, fmt.Sprintf("frostbreath-on-bloom-%v", w.char.Base.Key.String()))
	w.core.Events.Subscribe(event.OnLunarCharged, onHydroOrCryoReaction, fmt.Sprintf("frostbreath-on-lunar-charged-%v", w.char.Base.Key.String()))
	w.core.Events.Subscribe(event.OnLunarBloom, onHydroOrCryoReaction, fmt.Sprintf("frostbreath-on-lunar-bloom-%v", w.char.Base.Key.String()))
	w.core.Events.Subscribe(event.OnLunarCrystallize, onHydroOrCryoReaction, fmt.Sprintf("frostbreath-on-lunar-crystallize-%v", w.char.Base.Key.String()))
	w.core.Events.Subscribe(event.OnStellarConduct, onHydroOrCryoReaction, fmt.Sprintf("frostbreath-on-stellar-conduct-%v", w.char.Base.Key.String()))
	w.core.Events.Subscribe(event.OnStellarSwirl, onHydroOrCryoReaction, fmt.Sprintf("frostbreath-on-stellar-swirl-%v", w.char.Base.Key.String()))

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

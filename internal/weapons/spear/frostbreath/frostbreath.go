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

	energy := 4.5 + float64(w.refine)*1.5

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

	key := "frostbreath-on-hydro-or-cryo-reaction-" + w.char.Base.Key.String()
	w.core.Events.Subscribe(event.OnMelt, onHydroOrCryoReaction, key)
	w.core.Events.Subscribe(event.OnFrozen, onHydroOrCryoReaction, key)
	w.core.Events.Subscribe(event.OnShatter, onHydroOrCryoReaction, key)
	w.core.Events.Subscribe(event.OnSuperconduct, onHydroOrCryoReaction, key)
	w.core.Events.Subscribe(event.OnSwirlCryo, onHydroOrCryoReaction, key)
	w.core.Events.Subscribe(event.OnSwirlHydro, onHydroOrCryoReaction, key)
	w.core.Events.Subscribe(event.OnCrystallizeCryo, onHydroOrCryoReaction, key)
	w.core.Events.Subscribe(event.OnCrystallizeHydro, onHydroOrCryoReaction, key)
	w.core.Events.Subscribe(event.OnVaporize, onHydroOrCryoReaction, key)
	w.core.Events.Subscribe(event.OnElectroCharged, onHydroOrCryoReaction, key)
	w.core.Events.Subscribe(event.OnBloom, onHydroOrCryoReaction, key)
	w.core.Events.Subscribe(event.OnLunarCharged, onHydroOrCryoReaction, key)
	w.core.Events.Subscribe(event.OnLunarBloom, onHydroOrCryoReaction, key)
	w.core.Events.Subscribe(event.OnLunarCrystallize, onHydroOrCryoReaction, key)
	w.core.Events.Subscribe(event.OnStellarConduct, onHydroOrCryoReaction, key)
	w.core.Events.Subscribe(event.OnStellarSwirl, onHydroOrCryoReaction, key)

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

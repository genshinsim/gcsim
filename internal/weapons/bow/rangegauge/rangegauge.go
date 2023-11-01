package rangegauge

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const icdKey = "range-gauge-struggle-icd"

var symbol = []string{"unity-symbol-0", "unity-symbol-1", "unity-symbol-2"}

func init() {
	core.RegisterWeaponFunc(keys.RangeGauge, NewWeapon)
}

type Weapon struct {
	core    *core.Core
	char    *character.CharWrapper
	refine  int
	atkp    []float64
	eleDMGP []float64
	Index   int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

// When the wielder is healed or heals others, they will gain a Unity's Symbol that lasts 30s, up to a maximum of 3 Symbols.
// When using their Elemental Skill or Burst, all Symbols will be consumed and the Struggle effect will be granted for 10s.
// For each Symbol consumed, gain 3/4/5/6/7% ATK and 7/8.5/10/11.5/13% All Elemental DMG Bonus.
// The Struggle effect can be triggered once every 15s,
// and Symbols can be gained even when the character is not on the field.

func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{
		core:   c,
		char:   char,
		refine: p.Refine,
		buff:   make([]float64, attributes.EndStatType),
	}

	atkp := 0.02 + 0.01*float64(r)
	baseEleDMGP := 0.055 + 0.015*float64(r)
	duration := 10 * 60

	val := make([]float64, attributes.EndStatType)
	m := make([]float64, attributes.EndStatType)

	c.Events.Subscribe(event.OnHeal, func(args ...interface{}) bool {
		source := args[0].(*player.HealInfo)
		index := args[1].(int)
		amount := args[2].(float64)
		if source.Caller != char.Index && index != char.Index { // heal others and get healed including wielder
			return false
		}
		if amount <= 0 {
			return false
		}

		// override oldest symbol
		idx := 0
		for i, s := range symbol {
			if char.StatusExpiry(s) < char.StatusExpiry(symbol[idx]) {
				idx = i
			}
		}
		char.AddStatus((symbol[idx]), 30*60, true)

		c.Log.NewEvent("dockhands-assistant proc'd", glog.LogWeaponEvent, char.Index).
			Write("index", idx)

		return false
	}, fmt.Sprintf("dockhands-assistant-heal-%v", char.Base.Key.String()))

	c.Events.Subscribe(event.OnBurst, func(args ...interface{}) bool {
		// check for active before deleting symbol
		if char.StatusIsActive(icdKey) {
			return false
		}
		if c.Player.Active() != char.Index {
			return false
		}
		count := 0
		for _, s := range symbol {
			if char.StatusIsActive(s) {
				count++
			}
			char.DeleteStatus(s)
		}
		if count == 0 {
			return false
		}
		char.AddStatus(icdKey, 15*60, true)

		val[attributes.ATKP] = atkp * float64(count)

		// add atk buff
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag("range-gauge-atk-boost", duration),
			AffectedStat: attributes.ATKP,
			Amount: func() ([]float64, bool) {
				return val, true
			},
		})

		eleDMGP := baseEleDMGP * float64(count)
		for i := attributes.PyroP; i <= attributes.DendroP; i++ {
			m[i] = eleDMGP
		}

		// add elemental dmg buff
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag("range-gauge-eledmg-buff", duration),
			AffectedStat: attributes.NoStat,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})

		return false
	}, fmt.Sprintf("range-gauge-struggle-%v", char.Base.Key.String()))

	c.Events.Subscribe(event.OnSkill, func(args ...interface{}) bool {
		// check for active before deleting symbol
		if char.StatusIsActive(icdKey) {
			return false
		}
		if c.Player.Active() != char.Index {
			return false
		}
		count := 0
		for _, s := range symbol {
			if char.StatusIsActive(s) {
				count++
			}
			char.DeleteStatus(s)
		}
		if count == 0 {
			return false
		}
		char.AddStatus(icdKey, 15*60, true)
		val[attributes.ATKP] = atkp * float64(count)

		// add atk buff
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag("range-gauge-atk-boost", duration),
			AffectedStat: attributes.ATKP,
			Amount: func() ([]float64, bool) {
				return val, true
			},
		})

		eleDMGP := baseEleDMGP * float64(count)
		for i := attributes.PyroP; i <= attributes.DendroP; i++ {
			m[i] = eleDMGP
		}
		// add elemental dmg buff
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag("range-gauge-eledmg-buff", duration),
			AffectedStat: attributes.NoStat,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})

		return false
	}, fmt.Sprintf("dockhands-assistant-roused-%v", char.Base.Key.String()))
	return w, nil
}

func (w *Weapon) consumeEnergy(args ...interface{}) bool {
	em := 30 + 10*float64(w.refine)
	refund := 1.5 + 0.5*float64(w.refine)

	// check for active before deleting symbol
	if w.char.StatusIsActive(icdKey) {
		return false
	}
	if w.core.Player.Active() != w.char.Index {
		return false
	}

	count := 0
	for _, s := range symbol {
		if w.char.StatusIsActive(s) {
			count++
		}
		w.char.DeleteStatus(s)
	}
	if count == 0 {
		return false
	}

	w.char.AddStatus(icdKey, 15*60, true)
	w.buff[attributes.EM] = em * float64(count)

	// add atk buff
	w.char.AddStatMod(character.StatMod{
		Base:         modifier.NewBaseWithHitlag("range-gauge-atk-boost", duration),
		AffectedStat: attributes.ATKP,
		Amount: func() ([]float64, bool) {
			return w.buff, true
		},
	})

	eleDMGP := baseEleDMGP * float64(count)
	for i := attributes.PyroP; i <= attributes.DendroP; i++ {
		m[i] = eleDMGP
	}

	// add elemental dmg buff
	char.AddStatMod(character.StatMod{
		Base:         modifier.NewBaseWithHitlag("range-gauge-eledmg-buff", duration),
		AffectedStat: attributes.NoStat,
		Amount: func() ([]float64, bool) {
			return m, true
		},
	})

	return false
}

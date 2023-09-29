package powersaw

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

func init() {
	core.RegisterWeaponFunc(keys.PortablePowerSaw, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

// When the wielder is healed or heals others, they will gain a Stoic's Symbol that lasts 30s, up to a maximum of 3 Symbols.
//When using their Elemental Skill or Burst, all Symbols will be consumed and the Roused effect will be granted for 10s.
//For each Symbol consumed, gain 40/50/60/70/80 Elemental Mastery,
//and 2s after the effect occurs, 2/2.5/3/3.5/4 Energy per Symbol consumed will be restored for said character.
//The Roused effect can be triggered once every 15s,
//and Symbols can be gained even when the character is not on the field.

func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	const icdKey = "power-saw-roused-icd"
	symbol := []string{"stoic-symbol-0", "stoic-symbol-1", "stoic-symbol-2"}
	em := 30 + 10*float64(r)
	refund := 1.5 + 0.5*float64(r)
	duration := 10 * 60

	val := make([]float64, attributes.EndStatType)

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

		c.Log.NewEvent("powersaw proc'd", glog.LogWeaponEvent, char.Index).
			Write("index", idx)

		return false
	}, fmt.Sprintf("portable-power-saw-heal-%v", char.Base.Key.String()))

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
		val[attributes.EM] = em * float64(count)

		// add em buff
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag("portable-power-saw-em-boost", duration),
			AffectedStat: attributes.EM,
			Amount: func() ([]float64, bool) {
				return val, true
			},
		})

		// regen energy after 2 secs
		char.QueueCharTask(func() {
			char.AddEnergy("portable-power-saw-energy", refund*float64(count))
		}, 2*60)

		return false
	}, fmt.Sprintf("portable-power-saw-roused-%v", char.Base.Key.String()))

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
		val[attributes.EM] = em * float64(count)

		// add em buff
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag("portable-power-saw-em-boost", duration),
			AffectedStat: attributes.EM,
			Amount: func() ([]float64, bool) {
				return val, true
			},
		})

		// regen energy after 2 secs
		char.QueueCharTask(func() {
			char.AddEnergy("portable-power-saw-energy", refund*float64(count))
		}, 2*60)

		return false
	}, fmt.Sprintf("portable-power-saw-roused-%v", char.Base.Key.String()))
	return w, nil
}

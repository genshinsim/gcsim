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

const (
	symbolKey      = "portable-power-saw-symbol"
	symbolDuration = 30 * 60
	icdKey         = "portable-power-saw-icd"
	icdDuration    = 15 * 60
	buffKey        = "portable-power-saw-em"
	buffDuration   = 10 * 60
)

func init() {
	core.RegisterWeaponFunc(keys.PortablePowerSaw, NewWeapon)
}

type Weapon struct {
	stacks int
	Index  int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

// When the wielder is healed or heals others, they will gain a Stoic's Symbol that lasts 30s, up to a maximum of 3 Symbols.
// When using their Elemental Skill or Burst, all Symbols will be consumed and the Roused effect will be granted for 10s.
// For each Symbol consumed, gain 40/50/60/70/80 Elemental Mastery,
// and 2s after the effect occurs, 2/2.5/3/3.5/4 Energy per Symbol consumed will be restored for said character.
// The Roused effect can be triggered once every 15s,
// and Symbols can be gained even when the character is not on the field.
func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	// gain symbols
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

		if !char.StatusIsActive(symbolKey) {
			w.stacks = 0
		}
		if w.stacks < 3 {
			w.stacks++
		}
		c.Log.NewEvent("portable-power-saw adding stack", glog.LogWeaponEvent, char.Index).
			Write("stacks", w.stacks)
		char.AddStatus(symbolKey, symbolDuration, true)
		return false
	}, fmt.Sprintf("portable-power-saw-heal-%v", char.Base.Key.String()))

	// consume symbols
	em := 30 + 10*float64(r)
	refund := 1.5 + 0.5*float64(r)
	m := make([]float64, attributes.EndStatType)
	buffFunc := func(args ...interface{}) bool {
		// skip if no symbols (status not active implies symbols == 0)
		if !char.StatusIsActive(symbolKey) {
			return false
		}
		// skip if trigger on icd
		if char.StatusIsActive(icdKey) {
			return false
		}
		// check for active before deleting symbol
		if c.Player.Active() != char.Index {
			return false
		}
		// add icd
		char.AddStatus(icdKey, icdDuration, true)

		// consume symbols
		count := w.stacks
		char.DeleteStatus(symbolKey)
		w.stacks = 0

		// add em buff
		m[attributes.EM] = em * float64(count)
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag(buffKey, buffDuration),
			AffectedStat: attributes.EM,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})

		// regen energy after 2 secs
		char.QueueCharTask(func() {
			char.AddEnergy("portable-power-saw-energy", refund*float64(count))
		}, 2*60)

		return false
	}
	key := fmt.Sprintf("portable-power-saw-roused-%v", char.Base.Key.String())
	c.Events.Subscribe(event.OnBurst, buffFunc, key)
	c.Events.Subscribe(event.OnSkill, buffFunc, key)

	return w, nil
}

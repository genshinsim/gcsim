package forestregalia

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterWeaponFunc(keys.ForestRegalia, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

const (
	icdKey  = "forestregalia-icd"
	buffKey = "forest-sanctuary"
)

var procEvents = []event.Event{
	event.OnBurning,
	event.OnQuicken,
	event.OnAggravate,
	event.OnSpread,
	event.OnBloom,
	event.OnHyperbloom,
	event.OnBurgeon,
}

// After triggering Burning, Quicken, Aggravate, Spread, Bloom, Hyperbloom, or
// Burgeon, a Leaf of Consciousness will be created around the character for a
// maximum of 10s. When picked up, the Leaf will grant the character
// 60/75/90/105/120 Elemental Mastery for 12s. Only 1 Leaf can be generated
// this way every 20s. This effect can still be triggered if the character is
// not on the field. The Leaf of Consciousness' effect cannot stack.
func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	pickupDelay := p.Params["pickup_delay"]

	m := make([]float64, attributes.EndStatType)
	m[attributes.EM] = float64(45 + r*15)

	handleProc := func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.ActorIndex != char.Index {
			return false
		}
		if char.StatusIsActive(icdKey) {
			return false
		}
		char.AddStatus(icdKey, 1200, true)
		c.Log.NewEvent("forestregalia proc'd", glog.LogWeaponEvent, char.Index)
		if pickupDelay <= 0 {
			c.Log.NewEvent("forestregalia leaf ignored", glog.LogWeaponEvent, char.Index)
			return false
		}
		c.Tasks.Add(func() {
			active := c.Player.ActiveChar()
			active.AddStatMod(character.StatMod{
				Base:         modifier.NewBaseWithHitlag(buffKey, 720),
				AffectedStat: attributes.NoStat,
				Amount: func() ([]float64, bool) {
					return m, true
				},
			})
			c.Log.NewEvent(
				fmt.Sprintf("forestregalia leaf picked up by %v", active.Base.Key.String()),
				glog.LogWeaponEvent,
				char.Index,
			)
		}, pickupDelay)
		return false
	}
	for _, e := range procEvents {
		c.Events.Subscribe(e, handleProc, fmt.Sprintf("forestregalia-%v", char.Base.Key.String()))
	}
	return w, nil
}

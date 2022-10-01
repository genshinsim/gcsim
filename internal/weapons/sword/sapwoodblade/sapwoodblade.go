package sapwoodblade

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterWeaponFunc(keys.SapwoodBlade, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

const (
	icdKey  = "sapwoodblade-icd"
	buffKey = "sapwoodblade-buff"
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
// maximum of 10s. When picked up, the Leaf will grant the character 60
// Elemental Mastery for 12s. Only 1 Leaf can be generated this way every 20s.
// This effect can still be triggered if the character is not on the field. The
// Leaf of Consciousness' effect cannot stack.
func NewWeapon(c *core.Core, char *character.CharWrapper, p weapon.WeaponProfile) (weapon.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	handleProc := func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.ActorIndex != char.Index {
			return false
		}
		if char.StatusIsActive(icdKey) {
			return false
		}
		char.AddStatus(icdKey, 1200, false) // TODO: is this hitlag extendable?
		c.Log.NewEvent("sapwood blade proc'd", glog.LogWeaponEvent, char.Index).
			Write("seed expiry", c.F+600)
		c.Log.NewEvent("seed picked up", glog.LogWeaponEvent, char.Index)
		m := make([]float64, attributes.EndStatType)
		m[attributes.EM] = float64(45 + r*15)
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag(buffKey, 720),
			AffectedStat: attributes.NoStat,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})
		return false
	}
	for _, e := range procEvents {
		c.Events.Subscribe(e, handleProc, fmt.Sprintf("sapwoodblade-proc-%v", char.Base.Key.String()))
	}
	return w, nil
}

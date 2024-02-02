package mappa

import (
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
	core.RegisterWeaponFunc(keys.MappaMare, NewWeapon)
}

type Weapon struct {
	stacks int
	Index  int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	dmg := 0.06 + float64(r)*0.02

	const buffKey = "mappa-mare"
	buffDuration := 600 // 10s * 60

	addStack := func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.ActorIndex != char.Index {
			return false
		}
		if c.Player.Active() != char.Index {
			return false
		}

		if !char.StatusIsActive(buffKey) {
			w.stacks = 0
		}
		if w.stacks < 2 {
			w.stacks++
		}

		m := make([]float64, attributes.EndStatType)
		m[attributes.PyroP] = dmg * float64(w.stacks)
		m[attributes.HydroP] = dmg * float64(w.stacks)
		m[attributes.CryoP] = dmg * float64(w.stacks)
		m[attributes.ElectroP] = dmg * float64(w.stacks)
		m[attributes.AnemoP] = dmg * float64(w.stacks)
		m[attributes.GeoP] = dmg * float64(w.stacks)
		m[attributes.DendroP] = dmg * float64(w.stacks)

		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag(buffKey, buffDuration),
			AffectedStat: attributes.NoStat,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})

		c.Log.NewEvent("mappa-mare adding stack", glog.LogWeaponEvent, char.Index).
			Write("stacks", w.stacks)

		return false
	}

	for i := event.ReactionEventStartDelim + 1; i < event.OnShatter; i++ {
		c.Events.Subscribe(i, addStack, "mappa-mare-"+char.Base.Key.String())
	}

	return w, nil
}

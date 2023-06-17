package mappa

import (
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
	core.RegisterWeaponFunc(keys.MappaMare, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p weapon.WeaponProfile) (weapon.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	stacks := 0
	const stackKey = "mappa-mare-stacks"
	stackDuration := 600 // 10s * 60

	addStack := func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.ActorIndex != char.Index {
			return false
		}
		if c.Player.Active() != char.Index {
			return false
		}

		if !char.StatusIsActive(stackKey) {
			stacks = 1
			char.AddStatus(stackKey, stackDuration, true)
			c.Log.NewEvent("mappa proc'd", glog.LogWeaponEvent, char.Index).
				Write("stacks", stacks).
				Write("expiry (without hitlag)", c.F+stackDuration)
		} else if stacks < 2 {
			stacks++
			char.AddStatus(stackKey, stackDuration, true)
			c.Log.NewEvent("mappa proc'd", glog.LogWeaponEvent, char.Index).
				Write("stacks", stacks).
				Write("expiry (without hitlag)", c.F+stackDuration)
		}
		return false
	}

	for i := event.Event(event.ReactionEventStartDelim + 1); i < event.OnShatter; i++ {
		c.Events.Subscribe(i, addStack, "mappa"+char.Base.Key.String())
	}

	dmg := 0.06 + float64(r)*0.02
	m := make([]float64, attributes.EndStatType)
	char.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("mappa", -1),
		AffectedStat: attributes.NoStat,
		Amount: func() ([]float64, bool) {
			if !char.StatusIsActive(stackKey) {
				return nil, false
			}

			m[attributes.PyroP] = dmg * float64(stacks)
			m[attributes.HydroP] = dmg * float64(stacks)
			m[attributes.CryoP] = dmg * float64(stacks)
			m[attributes.ElectroP] = dmg * float64(stacks)
			m[attributes.AnemoP] = dmg * float64(stacks)
			m[attributes.GeoP] = dmg * float64(stacks)
			m[attributes.DendroP] = dmg * float64(stacks)
			return m, true
		},
	})

	return w, nil
}

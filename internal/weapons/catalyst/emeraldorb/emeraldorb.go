package emeraldorb

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/gadget"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterWeaponFunc(keys.EmeraldOrb, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

// Upon causing a Vaporize, Electro-Charged, Frozen, or a Hydro-infused Swirl reaction, increases ATK by 20/25/30/35/40% for 12s.
func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	m := make([]float64, attributes.EndStatType)
	m[attributes.ATKP] = 0.15 + float64(r)*0.05

	addBuff := func(args ...interface{}) bool {
		if _, ok := args[0].(*gadget.Gadget); ok {
			return false
		}

		atk := args[1].(*combat.AttackEvent)
		// don't proc if dmg not from weapon holder
		if atk.Info.ActorIndex != char.Index {
			return false
		}
		// don't proc if off-field
		if c.Player.Active() != char.Index {
			return false
		}

		// add buff
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag("emeraldorb", 720),
			AffectedStat: attributes.NoStat,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})

		return false
	}

	subKey := "emeraldorb-" + char.Base.Key.String()

	c.Events.Subscribe(event.OnVaporize, addBuff, subKey)
	c.Events.Subscribe(event.OnElectroCharged, addBuff, subKey)
	c.Events.Subscribe(event.OnFrozen, addBuff, subKey)
	c.Events.Subscribe(event.OnSwirlHydro, addBuff, subKey)

	return w, nil
}

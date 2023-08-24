package thrilling

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterWeaponFunc(keys.ThrillingTalesOfDragonSlayers, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	//When switching characters, the new character taking the field has their
	//ATK increased by 24% for 10s. This effect can only occur once every 20s.
	w := &Weapon{}
	r := p.Refine

	const icdKey = "ttds-icd"
	icd := 1200 // 20s * 60
	isActive := false
	key := fmt.Sprintf("ttds-%v", char.Base.Key.String())

	c.Events.Subscribe(event.OnInitialize, func(args ...interface{}) bool {
		isActive = c.Player.Active() == char.Index
		return true
	}, key)

	m := make([]float64, attributes.EndStatType)
	m[attributes.ATKP] = .18 + float64(r)*0.06

	c.Events.Subscribe(event.OnCharacterSwap, func(args ...interface{}) bool {
		if !isActive && c.Player.Active() == char.Index {
			isActive = true
			return false
		}

		if isActive && c.Player.Active() != char.Index {
			isActive = false
			if char.StatusIsActive(icdKey) {
				return false
			}
			char.AddStatus(icdKey, icd, true)
			active := c.Player.ActiveChar()
			active.AddStatMod(character.StatMod{
				Base:         modifier.NewBaseWithHitlag("ttds", 600),
				AffectedStat: attributes.NoStat,
				Amount: func() ([]float64, bool) {
					return m, true
				},
			})

			c.Log.NewEvent("ttds activated", glog.LogWeaponEvent, c.Player.Active()).
				Write("expiry (without hitlag)", c.F+600)
		}

		return false
	}, key)

	return w, nil
}

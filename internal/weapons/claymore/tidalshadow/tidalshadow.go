package tidalshadow

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterWeaponFunc(keys.TidalShadow, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

// After the wielder is healed, ATK will be increased by 24/30/36/42/48% for 8s.
// This can be triggered even when the character is not on the field.
func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	atk := 0.18 + float64(r)*0.06
	duration := 8 * 60
	val := make([]float64, attributes.EndStatType)
	val[attributes.ATKP] = atk
	c.Events.Subscribe(event.OnHeal, func(args ...interface{}) bool {
		index := args[1].(int)
		if index != char.Index {
			return false
		}
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag("tidal-shadow-atk-boost", duration),
			AffectedStat: attributes.ATKP,
			Amount: func() ([]float64, bool) {
				return val, true
			},
		})
		return false
	}, fmt.Sprintf("tidalshadow-%v", char.Base.Key.String()))
	return w, nil
}

package songofstillness

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterWeaponFunc(keys.SongOfStillness, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

// When the wielder is healed, ATK will be increased by 24/30/36/42/48% for 8s.
// This can be triggered even when the character is not on the field.
func NewWeapon(c *core.Core, char *character.CharWrapper, p weapon.WeaponProfile) (weapon.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	dmg := 0.12 + float64(r)*0.04
	duration := 480 //60 * 8s
	val := make([]float64, attributes.EndStatType)
	val[attributes.DmgP] = dmg
	c.Events.Subscribe(event.OnHeal, func(args ...interface{}) bool {
		healInfo := args[0].(*player.HealInfo)
		if healInfo.Target != -1 && healInfo.Target != char.Index {
			return false
		}
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBase("songofstillness-dmg-boost", duration),
			AffectedStat: attributes.NoStat,
			Amount: func() ([]float64, bool) {
				return val, true
			},
		})
		return false
	}, fmt.Sprintf("songofstillness-%v", char.Base.Key.String()))
	return w, nil
}

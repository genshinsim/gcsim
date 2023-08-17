package rightfulreward

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
)

func init() {
	core.RegisterWeaponFunc(keys.RightfulReward, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

// When the wielder is healed, restore 8/10/12/14/16 Energy.
// This effect can be triggered once every 10s, and can occur even when the character is not on the field.

func NewWeapon(c *core.Core, char *character.CharWrapper, p weapon.WeaponProfile) (weapon.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	refund := 6 + float64(r)*2
	icd := 600 //60 * 10s
	const icdKey = "rightfulreward-icd"

	c.Events.Subscribe(event.OnHeal, func(args ...interface{}) bool {
		healInfo := args[0].(*player.HealInfo)
		if healInfo.Target != -1 && healInfo.Target != char.Index {
			return false
		}
		if char.StatusIsActive(icdKey) {
			return false
		}
		char.AddStatus(icdKey, icd, true) //10s icd
		char.QueueCharTask(func() {
			char.AddEnergy("rightfulreward", refund)
		}, 0)

		return false
	}, fmt.Sprintf("rightfulreward-%v", char.Base.Key.String()))
	return w, nil
}

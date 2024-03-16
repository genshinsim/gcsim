package dialoguesofthedesertsages

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

func init() {
	core.RegisterWeaponFunc(keys.DialoguesOfTheDesertSages, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

const (
	energySrc = "dialoguesofthedesertsages"
	icdKey    = "dialoguesofthedesertsages-icd"
	icd       = 10 * 60
)

// When the wielder performs healing, restore 8/10/12/14/16 Energy.
// This effect can be triggered once every 10s and can occur even when the character is not on the field.
func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	energyRestore := 6 + float64(r)*2

	c.Events.Subscribe(event.OnHeal, func(args ...interface{}) bool {
		src := args[0].(*player.HealInfo)

		if src.Caller != char.Index {
			return false
		}
		if char.StatusIsActive(icdKey) {
			return false
		}
		char.AddStatus(icdKey, icd, true)
		char.AddEnergy(energySrc, energyRestore)

		return false
	}, fmt.Sprintf("dialoguesofthedesertsages-%v", char.Base.Key.String()))

	return w, nil
}

package cloudforged

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
	core.RegisterWeaponFunc(keys.Cloudforged, NewWeapon)
}

type Weapon struct {
	Index  int
	stacks int
}

const buffKey = "cloudforged-em"

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	var w Weapon
	refine := p.Refine

	m := make([]float64, attributes.EndStatType)
	c.Events.Subscribe(event.OnEnergyChange, func(args ...interface{}) bool {
		index := args[0].(*character.CharWrapper).Index
		amount := args[2].(float64)

		if char.Index != index || amount >= 0 {
			return false
		}
		if !char.StatModIsActive(buffKey) {
			w.stacks = 0
		}
		if w.stacks < 2 {
			w.stacks++
		}
		m[attributes.EM] = (30.0 + 10.0*float64(refine)) * float64(w.stacks)
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag(buffKey, 18*60),
			AffectedStat: attributes.EM,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})
		return false
	}, fmt.Sprintf("cloudforged-%v", char.Base.Key))

	return &w, nil
}

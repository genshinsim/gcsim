package absolution

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
	core.RegisterWeaponFunc(keys.Absolution, NewWeapon)
}

const (
	cdKey       = "absolution-crit-dmg"
	dmgBonusKey = "absolution-dmg-bonus"
)

type Weapon struct {
	Index  int
	stacks int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }
func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	var w Weapon
	refine := p.Refine

	perm := make([]float64, attributes.EndStatType)
	perm[attributes.CD] = 0.15 + 0.05*float64(refine)
	char.AddStatMod(character.StatMod{
		Base:         modifier.NewBase(cdKey, -1),
		AffectedStat: attributes.CD,
		Amount: func() ([]float64, bool) {
			return perm, true
		},
	})

	bonus := make([]float64, attributes.EndStatType)
	c.Events.Subscribe(event.OnHPDebt, func(args ...interface{}) bool {
		index := args[0].(int)
		amount := args[1].(float64)
		if char.Index != index || amount <= 0 {
			return false
		}
		if !char.StatModIsActive(dmgBonusKey) {
			w.stacks = 0
		}
		if w.stacks < 3 {
			w.stacks++
		}
		bonus[attributes.DmgP] = (0.12 + 0.04*float64(refine)) * float64(w.stacks)
		char.AddStatMod(character.StatMod{
			Base: modifier.NewBaseWithHitlag(dmgBonusKey, 6*60),
			Amount: func() ([]float64, bool) {
				return bonus, true
			},
		})

		return false
	}, fmt.Sprintf("absolution-%v", char.Base.Key))

	return &w, nil
}

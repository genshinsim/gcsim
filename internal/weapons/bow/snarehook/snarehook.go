package snarehook

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
	core.RegisterWeaponFunc(keys.SnareHook, NewWeapon)
}

type Weapon struct {
	stacks int
	Index  int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

// Upon causing an Elemental Reaction, increases Elemental Mastery by 60/75/90/105/120 for 12s.
// Moonsign: Ascendant Gleam: Elemental Mastery from this effect is further increased by 60/75/90/105/120.
// This effect can be triggered even if the equipping character is off-field.
func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	emBuff := 45 + 15*float64(r)
	m := make([]float64, attributes.EndStatType)

	add := func(args ...any) bool {
		atk := args[1].(*info.AttackEvent)

		if atk.Info.ActorIndex != char.Index() {
			return false
		}

		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag("snarehook-em", 12*60),
			AffectedStat: attributes.EM,
			Amount: func() ([]float64, bool) {
				m[attributes.EM] = emBuff * moonsignBonus(c)
				return m, true
			},
		})

		return false
	}

	// TODO: Does shatter count?
	for i := event.ReactionEventStartDelim + 1; i < event.ReactionEventEndDelim; i++ {
		c.Events.Subscribe(i, add, fmt.Sprintf("snarehook-%v", char.Base.Key.String()))
	}

	return w, nil
}

func moonsignBonus(c *core.Core) float64 {
	if c.Player.GetMoonsignLevel() >= 2 {
		return 2.0
	}
	return 1.0
}

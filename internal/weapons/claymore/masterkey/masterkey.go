package masterkey

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterWeaponFunc(keys.MasterKey, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

// Upon causing an Elemental Reaction, increases Elemental Mastery by 60/75/90/105/120 for 12s.
// Moonsign: Ascendant Gleam: Elemental Mastery from this effect is further increased by 60/75/90/105/120.
// This effect can be triggered even if the equipping character is off-field.
func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	const key = "masterkey-%v"
	emBuff := 45 + 15*float64(r)
	m := make([]float64, attributes.EndStatType)

	onReact := func(args ...any) {
		if _, ok := args[0].(*enemy.Enemy); !ok {
			return
		}

		atk := args[1].(*info.AttackEvent)
		if atk.Info.ActorIndex != char.Index() {
			return
		}

		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag(fmt.Sprintf(key, "em"), 12*60),
			AffectedStat: attributes.EM,
			Amount: func() []float64 {
				if c.Player.GetMoonsignLevel() >= 2 {
					m[attributes.EM] = emBuff * 2
				} else {
					m[attributes.EM] = emBuff
				}
				return m
			},
		})
	}
	for i := event.ReactionEventStartDelim + 1; i < event.ReactionEventEndDelim; i++ {
		c.Events.Subscribe(i, onReact, fmt.Sprintf(key, char.Base.Key.String()))
	}

	return w, nil
}

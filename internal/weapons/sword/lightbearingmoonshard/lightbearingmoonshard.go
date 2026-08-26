package lightbearingmoonshard

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	buffKey = "lightbearing-moonshard"
)

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	def := make([]float64, attributes.EndStatType)
	def[attributes.DEFP] = 0.15 + 0.05*float64(r)

	char.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("lightbearing-moonshard-def", -1),
		AffectedStat: attributes.DEFP,
		Amount: func() []float64 {
			return def
		},
	})

	bonus := 0.64 + 0.16*float64(r)

	moonshard := func(args ...any) {
		if c.Player.Active() != char.Index() {
			return
		}

		char.AddReactBonusMod(character.ReactBonusMod{
			Base: modifier.NewBaseWithHitlag(buffKey, 5*60),
			Amount: func(atk info.AttackInfo) float64 {
				switch atk.AttackTag {
				case attacks.AttackTagDirectLunarCrystallize,
					attacks.AttackTagReactionLunarCrystallize:
					return bonus
				default:
					return 0
				}
			},
		})
	}

	c.Events.Subscribe(
		event.OnSkill,
		moonshard,
		fmt.Sprintf("lightbearingmoonshard-%v", char.Base.Key.String()),
	)

	return w, nil
}

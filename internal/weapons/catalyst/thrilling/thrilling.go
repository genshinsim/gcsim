package thrilling

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	icdKey = "ttds-icd"
	icdDur = 20 * 60
	icdExt = 1 * 60
)

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	m := make([]float64, attributes.EndStatType)
	m[attributes.ATKP] = .18 + float64(r)*0.06

	c.Events.Subscribe(event.OnCharacterSwap, func(args ...any) {
		prev := args[0].(int)
		next := args[1].(int)
		if next == char.Index() && char.StatusDuration(icdKey) < icdExt {
			char.DeleteStatus(icdKey)
		}
		if prev != char.Index() {
			return
		}

		if char.StatusIsActive(icdKey) {
			return
		}
		char.AddStatus(icdKey, icdDur+icdExt, true)

		active := c.Player.ActiveChar()
		// When TTDS mod is active, don't reapply
		if active.StatModIsActive("ttds") {
			return
		}
		active.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag("ttds", 600),
			AffectedStat: attributes.ATKP,
			Amount: func() []float64 {
				return m
			},
		})
	}, fmt.Sprintf("ttds-swap-%v", char.Base.Key.String()))

	return w, nil
}

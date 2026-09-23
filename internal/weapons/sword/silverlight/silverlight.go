package silverlight

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const stackGainICDKey = "silver-light-stack-icd"

type Weapon struct {
	Index        int
	stackCounter int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

// Upon causing an Elemental Reaction, increases Max HP by 16%/20%/24%/28%/32% for 12s.
// Moonsign: Ascendant Gleam: Max HP from this effect is further increased by 16%/20%/24%/28%/32%
// This effect can be triggered even if the equipping character is off-field.
func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{
		stackCounter: 0,
	}
	r := p.Refine

	m := make([]float64, attributes.EndStatType)
	m[attributes.EM] = 39 + 13*float64(r)

	onSkill := func(args ...any) {
		if char.StatusIsActive(stackGainICDKey) {
			return
		}

		if char.Index() != c.Player.Active() {
			return
		}

		char.AddStatMod(character.StatMod{
			Base: modifier.NewBaseWithHitlag(fmt.Sprintf("silver-light-em-%v", w.stackCounter+1), 12*60),
			Amount: func() []float64 {
				return m
			},
		})

		char.AddStatus(stackGainICDKey, 0.2*60, true)

		w.stackCounter++
		w.stackCounter %= 2
	}

	c.Events.Subscribe(event.OnSkill, onSkill, "silver-light-on-skill-"+char.Base.Key.String())

	return w, nil
}

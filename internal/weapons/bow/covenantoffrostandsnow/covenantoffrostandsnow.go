package covenantoffrostandsnow

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	emBuff := make([]float64, attributes.EndStatType)
	emBuff[attributes.EM] = 90 + float64(r)*30

	onSkill := func(args ...any) {
		if c.Player.Active() != char.Index() {
			return
		}

		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag("covenant-of-frost-and-snow-em", 12*60),
			AffectedStat: attributes.EM,
			Amount: func() []float64 {
				return emBuff
			},
		})
	}

	c.Events.Subscribe(event.OnSkill, onSkill, "covenant-of-frost-and-snow-on-skill")

	return w, nil
}

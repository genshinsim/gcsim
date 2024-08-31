package fluteofezpitzal

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
	core.RegisterWeaponFunc(keys.FluteOfEzpitzal, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

// Using an Elemental Skill increases DEF by 16/20/24/28/32% for 15s.

func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	// TODO: check for ICD, refresh timing
	// const icdKey = "flute-of-ezpitzal-icd"
	def := 0.12 + float64(r)*0.04
	duration := 15 * 60

	m := make([]float64, attributes.EndStatType)
	m[attributes.DEFP] = def

	c.Events.Subscribe(event.OnSkill, func(args ...interface{}) bool {

		// TODO: check if wearer need to be on-field?
		if c.Player.Active() != char.Index {
			return false
		}

		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag("flute-of-ezpitzal-def-boost", duration),
			AffectedStat: attributes.DEFP,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})

		return false
	}, fmt.Sprintf("flute-of-ezpitzal-def%v", char.Base.Key.String()))

	return w, nil
}

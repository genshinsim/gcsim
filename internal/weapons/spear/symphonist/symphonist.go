package symphonist

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

const bufDur = 3 * 60
const buffKey = "symphonist-chanson-de-baies"

func init() {
	core.RegisterWeaponFunc(keys.SymphonistOfScents, NewWeapon)
}

type Weapon struct {
	Index int
	char  *character.CharWrapper
	c     *core.Core
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	// Gain 12% All Elemental DMG Bonus. Obtain Consummation for 20s after using
	// an Elemental Skill, causing ATK to increase by 3.2% per second. This ATK
	// increase has a maximum of 6 stacks. When the character equipped with this
	// weapon is not on the field, Consummation's ATK increase is doubled.

	// ATK is increased by 12%. When the equipping character is off-field, ATK
	// is increased by an additional 12%. When initiating healing, the equipping
	// character and healed character(s) will obtain the "Chanson de Baies"
	// effect, increasing their ATK by 32% for 3s. This effect can be triggered
	// even if the equipping character is off-field.
	w := &Weapon{
		char: char,
		c:    c,
	}
	r := p.Refine
	selfAtkP := 0.09 + float64(r)*0.03
	m := make([]float64, attributes.EndStatType)
	char.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("symphonist-atkp", -1),
		AffectedStat: attributes.ATKP,
		Amount: func() ([]float64, bool) {
			m[attributes.ATKP] = selfAtkP
			if c.Player.Active() != char.Index {
				m[attributes.ATKP] += selfAtkP
			}
			return m, true
		},
	})

	buffOnHeal := make([]float64, attributes.EndStatType)
	buffOnHeal[attributes.ATKP] = 0.24 + float64(r)*0.08

	c.Events.Subscribe(event.OnHeal, func(args ...interface{}) bool {
		source := args[0].(*info.HealInfo)
		index := args[1].(int)

		if source.Caller != char.Index {
			return false
		}

		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag(buffKey, bufDur),
			AffectedStat: attributes.ATKP,
			Amount: func() ([]float64, bool) {
				return buffOnHeal, true
			},
		})

		if index >= 0 {
			otherChar := c.Player.Chars()[index]
			otherChar.AddStatMod(character.StatMod{
				Base:         modifier.NewBaseWithHitlag(buffKey, bufDur),
				AffectedStat: attributes.ATKP,
				Amount: func() ([]float64, bool) {
					return buffOnHeal, true
				},
			})
			return false
		}

		for _, otherChar := range c.Player.Chars() {
			if otherChar.Index == char.Index {
				continue
			}
			otherChar.AddStatMod(character.StatMod{
				Base:         modifier.NewBaseWithHitlag(buffKey, bufDur),
				AffectedStat: attributes.ATKP,
				Amount: func() ([]float64, bool) {
					return buffOnHeal, true
				},
			})
		}
		return false
	}, fmt.Sprintf("symphonist-of-scents-%v", char.Base.Key.String()))

	return w, nil
}

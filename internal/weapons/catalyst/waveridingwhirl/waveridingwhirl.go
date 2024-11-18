package waveridingwhirl

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
	core.RegisterWeaponFunc(keys.WaveridingWhirl, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }

func (w *Weapon) Init() error {
	return nil
}

const (
	buffICD = "waveridingwhirl-icd"
	ICDDur  = 15 * 60
	buffDur = 10 * 60
	buffKey = "waveridingwhirl-buff"
)

func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{}
	r := float64(p.Refine)

	stacks := 0
	val := make([]float64, attributes.EndStatType)

	c.Events.Subscribe(event.OnInitialize, func(args ...interface{}) bool {
		for _, char := range c.Player.Chars() {
			if char.Base.Element == attributes.Hydro {
				stacks++
			}
		}
		val[attributes.HPP] = (0.15 + 0.05*r) + (0.09+0.03*r)*float64(min(2, stacks))
		return true
	}, fmt.Sprintf("waveridingwhirl-%v", char.Base.Key.String()))

	c.Events.Subscribe(event.OnSkill, func(args ...interface{}) bool {
		if c.Player.Active() != char.Index {
			return false
		}
		if char.StatusIsActive(buffICD) {
			return false
		}
		char.AddStatus(buffICD, ICDDur, true)

		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag(buffKey, buffDur),
			AffectedStat: attributes.HPP,
			Amount: func() ([]float64, bool) {
				return val, true
			},
		})
		return false
	}, fmt.Sprintf("waveridingwhirl-skill-%v", char.Base.Key.String()))

	return w, nil
}

package newbough

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

type Weapon struct {
	stacks int
	Index  int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

const (
	skillKey = "newbough-skill"
	buffKey  = "newbough-verdant"
	icdKey   = "newbough-icd"
)

// Radiance: Stellar Glimmer is not implemented because gcsim has no Stellar Glimmer reaction and current teams cannot activate it.
func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	c.Events.Subscribe(event.OnSkill, func(args ...any) {
		if c.Player.Active() == char.Index() {
			char.AddStatus(skillKey, 12*60, true)
		}
	}, fmt.Sprintf("newbough-skill-%v", char.Base.Key.String()))

	c.Events.Subscribe(event.OnEnemyDamage, func(args ...any) {
		atk, ok := args[1].(*info.AttackEvent)
		if !ok {
			return
		}
		if atk.Info.ActorIndex != char.Index() {
			return
		}
		if !char.StatusIsActive(skillKey) {
			return
		}
		if char.StatusIsActive(icdKey) {
			return
		}

		if !char.StatModIsActive(buffKey) {
			w.stacks = 0
		}
		if w.stacks < 3 {
			w.stacks++
		}

		val := make([]float64, attributes.EndStatType)
		val[attributes.ATKP] = (0.03 + 0.01*float64(r)) * float64(w.stacks)
		val[attributes.EM] = (15 + 5*float64(r)) * float64(w.stacks)
		char.AddStatMod(character.StatMod{
			Base: modifier.NewBaseWithHitlag(buffKey, 6*60),
			Amount: func() []float64 {
				return val
			},
		})

		char.AddStatus(icdKey, 60, true)
	}, fmt.Sprintf("newbough-verdant-%v", char.Base.Key.String()))

	return w, nil
}

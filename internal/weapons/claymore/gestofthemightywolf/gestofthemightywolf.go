package gestofthemightywolf

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
	Index  int
	char   *character.CharWrapper
	refine int
	stacks int
	mod    []float64
}

const (
	gestAttackSpeedKey = "gest-of-the-mighty-wolf-atkspd"
	gestStacksKey      = "gest-of-the-mighty-wolf-stacks"
)

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{
		char:   char,
		refine: p.Refine,
		mod:    make([]float64, attributes.EndStatType),
	}

	m := make([]float64, attributes.EndStatType)
	m[attributes.AtkSpd] = 0.1
	char.AddStatMod(character.StatMod{
		Base:         modifier.NewBase(gestAttackSpeedKey, -1),
		AffectedStat: attributes.AtkSpd,
		Amount: func() []float64 {
			return m
		},
	})

	c.Events.Subscribe(event.OnEnemyDamage, func(args ...any) {
		w.addStacks(1)
	}, fmt.Sprintf("gest-of-the-mighty-wolf-on-normal-attack-%v", char.Base.Key.String()))

	c.Events.Subscribe(event.OnChargeAttack, func(args ...any) {
		w.addStacks(2)
	}, fmt.Sprintf("gest-of-the-mighty-wolf-on-charge-%v", char.Base.Key.String()))

	c.Events.Subscribe(event.OnSkill, func(args ...any) {
		w.addStacks(2)
	}, fmt.Sprintf("gest-of-the-mighty-wolf-on-skill-%v", char.Base.Key.String()))

	return w, nil
}

func (w *Weapon) addStacks(amt int) {
	if !w.char.StatModIsActive(gestStacksKey) {
		w.stacks = 0
	}

	w.char.AddStatMod(character.StatMod{
		Base: modifier.NewBase(gestStacksKey, 4*60),
		Amount: func() []float64 {
			w.mod[attributes.DmgP] = (0.055 + 0.02*float64(w.refine)) * float64(w.stacks)

			if !w.char.IsHexerei {
				return w.mod
			}

			w.mod[attributes.CD] = (0.055 + 0.02*float64(w.refine)) * float64(w.stacks)

			return w.mod
		},
	})

	w.stacks = min(4, w.stacks+amt)
}

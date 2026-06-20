package gestofthemightywolf

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
	core.RegisterWeaponFunc(keys.GestOfTheMightyWolf, NewWeapon)
}

type Weapon struct {
	Index     int
	char      *character.CharWrapper
	refine    int
	stacks    int
	hexStacks int
}

const (
	gestAttackSpeedKey = "gest-of-the-mighty-wolf-atkspd"
	gestStacksKey      = "gest-of-the-mighty-wolf-stacks"
	gestHexStacksKey   = "gest-of-the-mighty-wolf-stacks-hexerei"
)

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{
		char:   char,
		refine: p.Refine,
	}
	r := p.Refine

	m := make([]float64, attributes.EndStatType)
	m[attributes.AtkSpd] = 0.1
	char.AddStatMod(character.StatMod{
		Base:         modifier.NewBase(gestAttackSpeedKey, -1),
		AffectedStat: attributes.AtkSpd,
		Amount: func() []float64 {
			return m
		},
	})

	m2 := make([]float64, attributes.EndStatType)
	char.AddStatMod(character.StatMod{
		Base:         modifier.NewBase(gestHexStacksKey, -1),
		AffectedStat: attributes.CD,
		Amount: func() []float64 {
			m2[attributes.CD] = (0.055 + 0.02*float64(r)) * float64(w.hexStacks)

			return m2
		},
	})

	m3 := make([]float64, attributes.EndStatType)

	c.Events.Subscribe(event.OnEnemyDamage, func(args ...any) {
		w.addStacks(1, m3)
		w.addHexStacks(1)
	}, fmt.Sprintf("gest-of-the-mighty-wolf-on-normal-attack-%v", char.Base.Key.String()))

	c.Events.Subscribe(event.OnChargeAttack, func(args ...any) {
		w.addStacks(2, m3)
		w.addHexStacks(2)
	}, fmt.Sprintf("gest-of-the-mighty-wolf-on-charge-%v", char.Base.Key.String()))

	c.Events.Subscribe(event.OnSkill, func(args ...any) {
		w.addStacks(2, m3)
		w.addHexStacks(2)
	}, fmt.Sprintf("gest-of-the-mighty-wolf-on-skill-%v", char.Base.Key.String()))

	return w, nil
}

func (w *Weapon) addStacks(amt int, m []float64) {
	if !w.char.StatModIsActive(gestStacksKey) {
		w.stacks = 0
	}

	w.char.AddStatMod(character.StatMod{
		Base:         modifier.NewBase(gestStacksKey, 4*60),
		AffectedStat: attributes.DmgP,
		Amount: func() []float64 {
			m[attributes.DmgP] = (0.055 + 0.02*float64(w.refine)) * float64(w.stacks)

			return m
		},
	})

	w.stacks = min(4, w.stacks+amt)
}

func (w *Weapon) addHexStacks(amt int) {
	addedAmt := min(4-w.hexStacks, amt)

	if addedAmt == 0 {
		return
	}

	w.hexStacks += addedAmt

	w.char.QueueCharTask(w.removeHexStacks(addedAmt), 4*60)
}

func (w *Weapon) removeHexStacks(amt int) func() {
	return func() {
		w.hexStacks -= amt
	}
}

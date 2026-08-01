package gestofthemightywolf

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

type Weapon struct {
	Index               int
	char                *character.CharWrapper
	core                *core.Core
	refine              int
	stacks              int
	mod                 []float64
	isHexereiSecretRite bool
}

const (
	gestAttackSpeedKey = "gest-of-the-mighty-wolf-atkspd"
	gestStacksKey      = "gest-of-the-mighty-wolf-stacks"
)

func (w *Weapon) SetIndex(idx int) { w.Index = idx }

func (w *Weapon) Init() error {
	w.isHexereiSecretRite = w.core.Player.GetHexereiCount() >= 2
	return nil
}

func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{
		char:   char,
		core:   c,
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
		atk := args[1].(*info.AttackEvent)
		if atk.Info.ActorIndex != char.Index() {
			return
		}
		if w.core.Player.Active() != char.Index() {
			return
		}
		if atk.Info.AttackTag != attacks.AttackTagNormal {
			return
		}
		w.addStacks(1)
		c.Log.NewEvent("gest adding stack- normal on hit", glog.LogWeaponEvent, char.Index()).
			Write("stacks", w.stacks)
	}, fmt.Sprintf("gest-of-the-mighty-wolf-on-normal-attack-%v", char.Base.Key.String()))

	c.Events.Subscribe(event.OnChargeAttack, func(args ...any) {
		w.addStacks(2)
		c.Log.NewEvent("gest adding stack- charge on start", glog.LogWeaponEvent, char.Index()).
			Write("stacks", w.stacks)
	}, fmt.Sprintf("gest-of-the-mighty-wolf-on-charge-%v", char.Base.Key.String()))

	c.Events.Subscribe(event.OnSkill, func(args ...any) {
		w.addStacks(2)
		c.Log.NewEvent("gest adding stack- skill on cast", glog.LogWeaponEvent, char.Index()).
			Write("stacks", w.stacks)
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

			if !w.isHexereiSecretRite {
				return w.mod
			}

			w.mod[attributes.CD] = (0.055 + 0.02*float64(w.refine)) * float64(w.stacks)

			return w.mod
		},
	})

	w.stacks = min(4, w.stacks+amt)
}

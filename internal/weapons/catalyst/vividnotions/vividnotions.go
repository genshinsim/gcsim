package vividnotions

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	plungeBuff     = "dawns-first-hue"
	skillBurstBuff = "twilights-splendor"
)

func init() {
	core.RegisterWeaponFunc(keys.VividNotions, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	m := make([]float64, attributes.EndStatType)
	m[attributes.ATKP] = 0.21 + float64(r)*0.07
	char.AddStatMod(character.StatMod{
		Base: modifier.NewBase("vividnotions-atk", -1),
		Amount: func() []float64 {
			return m
		},
	})

	mCD := make([]float64, attributes.EndStatType)
	plungeCD := 0.21 + float64(r)*0.07
	skillBurstCD := 0.3 + float64(r)*0.1

	char.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("vividnotions-cd", -1),
		Amount: func(atk *info.AttackEvent, t info.Target) []float64 {
			if atk.Info.AttackTag != attacks.AttackTagPlunge {
				return nil
			}

			mCD[attributes.CD] = 0
			if char.StatusIsActive(plungeBuff) {
				mCD[attributes.CD] += plungeCD
			}
			if char.StatusIsActive(skillBurstBuff) {
				mCD[attributes.CD] += skillBurstCD
			}
			return mCD
		},
	})

	c.Events.Subscribe(event.OnStateChange, func(args ...any) {
		next := args[1].(action.AnimationState)
		if next == action.PlungeAttackState {
			char.AddStatus(plungeBuff, 15*60, true)
		}
	}, fmt.Sprintf("vividnotions-plunge-%s", char.Base.Key.String()))

	c.Events.Subscribe(event.OnSkill, func(args ...any) {
		char.AddStatus(skillBurstBuff, 15*60, true)
	}, fmt.Sprintf("vividnotions-skill-%s", char.Base.Key.String()))

	c.Events.Subscribe(event.OnBurst, func(args ...any) {
		char.AddStatus(skillBurstBuff, 15*60, true)
	}, fmt.Sprintf("vividnotions-burst-%s", char.Base.Key.String()))

	c.Events.Subscribe(event.OnEnemyDamage, func(args ...any) {
		ae := args[1].(*info.AttackEvent)
		if ae.Info.ActorIndex != char.Index() {
			return
		}
		if ae.Info.AttackTag != attacks.AttackTagPlunge {
			return
		}
		if ae.Info.Durability == 0 {
			return
		}

		// TODO: hitlag affected?
		plungeF := char.StatusExpiry(plungeBuff)
		skillBurstF := char.StatusExpiry(skillBurstBuff)
		char.QueueCharTask(func() {
			if plungeF == char.StatusExpiry(plungeBuff) {
				char.DeleteStatus(plungeBuff)
			}
			if skillBurstF == char.StatusExpiry(skillBurstBuff) {
				char.DeleteStatus(skillBurstBuff)
			}
		}, 0.1*60)
	}, fmt.Sprintf("vividnotions-%s", char.Base.Key.String()))

	return w, nil
}

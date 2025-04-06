package vividnotions

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
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
		Amount: func() ([]float64, bool) {
			return m, true
		},
	})

	mCD := make([]float64, attributes.EndStatType)
	plungeCD := 0.21 + float64(r)*0.07
	skillBurstCD := 0.3 + float64(r)*0.1

	char.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("vividnotions-cd", -1),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			if atk.Info.AttackTag != attacks.AttackTagPlunge {
				return nil, false
			}

			mCD[attributes.CD] = 0
			if char.StatusIsActive(plungeBuff) {
				mCD[attributes.CD] += plungeCD
			}
			if char.StatusIsActive(skillBurstBuff) {
				mCD[attributes.CD] += skillBurstCD
			}
			return mCD, true
		},
	})

	c.Events.Subscribe(event.OnPlunge, func(args ...interface{}) bool {
		char.AddStatus(plungeBuff, 15*60, true)
		return false
	}, fmt.Sprintf("vividnotions-plunge-%s", char.Base.Key.String()))

	c.Events.Subscribe(event.OnSkill, func(args ...interface{}) bool {
		char.AddStatus(skillBurstBuff, 15*60, true)
		return false
	}, fmt.Sprintf("vividnotions-skill-%s", char.Base.Key.String()))

	c.Events.Subscribe(event.OnBurst, func(args ...interface{}) bool {
		char.AddStatus(skillBurstBuff, 15*60, true)
		return false
	}, fmt.Sprintf("vividnotions-burst-%s", char.Base.Key.String()))

	c.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		ae := args[1].(*combat.AttackEvent)
		if ae.Info.ActorIndex != char.Index {
			return false
		}
		if ae.Info.AttackTag != attacks.AttackTagPlunge {
			return false
		}
		if ae.Info.Durability == 0 {
			return false
		}

		// TODO: hitlag affected?
		char.QueueCharTask(func() {
			char.DeleteStatus(plungeBuff)
			char.DeleteStatus(skillBurstBuff)
		}, 0.05*60)

		return false
	}, fmt.Sprintf("vividnotions-%s", char.Base.Key.String()))

	return w, nil
}

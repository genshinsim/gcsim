package urakumisugiri

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

func init() {
	core.RegisterWeaponFunc(keys.UrakuMisugiri, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

const (
	buffKey = "urakumisugiri-increase-buff"
)

func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{}

	char.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("urakumisugiri-na", -1),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			if atk.Info.AttackTag != attacks.AttackTagNormal {
				return nil, false
			}
			value := make([]float64, attributes.EndStatType)
			value[attributes.DmgP] = 0.12 + 0.04*float64(p.Refine)
			if char.StatusIsActive(buffKey) {
				value[attributes.DmgP] = 0.14 + 0.08*float64(p.Refine)
			}
			return value, true
		},
	})

	char.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("urakumisugiri-skill", -1),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			if atk.Info.AttackTag != attacks.AttackTagElementalArt && atk.Info.AttackTag != attacks.AttackTagElementalArtHold {
				return nil, false
			}
			value := make([]float64, attributes.EndStatType)
			value[attributes.DmgP] = 0.18 + 0.06*float64(p.Refine)
			if char.StatusIsActive(buffKey) {
				value[attributes.DmgP] = 0.36 + 0.12*float64(p.Refine)
			}
			return value, true
		},
	})

	char.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("urakumisugiri-def", -1),
		AffectedStat: attributes.DEFP,
		Amount: func() ([]float64, bool) {
			value := make([]float64, attributes.EndStatType)
			value[attributes.DEFP] = 0.15 + 0.05*float64(p.Refine)
			return value, true
		},
	})

	c.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.Element != attributes.Geo {
			return false
		}
		char.AddStatus(buffKey, 900, true)
		return false
	}, fmt.Sprintf("urakumisugiri-%v", char.Base.Key.String()))

	return w, nil
}

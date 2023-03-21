package thundering

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterWeaponFunc(keys.ThunderingPulse, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p weapon.WeaponProfile) (weapon.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	m := make([]float64, attributes.EndStatType)
	m[attributes.ATKP] = 0.15 + float64(r)*0.05
	stack := 0.09 + float64(r)*0.03
	max := 0.3 + float64(r)*0.1

	const normalKey = "thundering-pulse-normal"
	normalDuration := 300 // 5s * 60
	const skillKey = "thundering-pulse-skill"
	skillDuration := 600 // 10s * 60

	key := fmt.Sprintf("thundering-pulse-%v", char.Base.Key.String())

	c.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		dmg := args[2].(float64)
		if atk.Info.ActorIndex != char.Index {
			return false
		}
		if atk.Info.AttackTag != attacks.AttackTagNormal {
			return false
		}
		if dmg == 0 {
			return false
		}
		char.AddStatus(normalKey, normalDuration, true)
		return false
	}, key)

	c.Events.Subscribe(event.OnSkill, func(args ...interface{}) bool {
		if c.Player.Active() != char.Index {
			return false
		}
		char.AddStatus(skillKey, skillDuration, true)
		return false
	}, key)

	char.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("thundering-pulse", -1),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			m[attributes.DmgP] = 0
			if atk.Info.AttackTag != attacks.AttackTagNormal {
				return m, true
			}
			count := 0
			if char.Energy < char.EnergyMax {
				count++
			}
			if char.StatusIsActive(normalKey) {
				count++
			}
			if char.StatusIsActive(skillKey) {
				count++
			}
			dmg := float64(count) * stack
			if count >= 3 {
				dmg = max
			}
			m[attributes.DmgP] = dmg
			return m, true
		},
	})

	return w, nil
}

package mistsplitter

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
	core.RegisterWeaponFunc(keys.MistsplitterReforged, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

const (
	normalBuffKey = "mistsplitter-normal"
	burstBuffKey  = "mistsplitter-burst"
)

// Gain a 12% Elemental DMG Bonus for all elements and receive the might of the
// Mistsplitter's Emblem. At stack levels 1/2/3, the Mistsplitter's Emblem
// provides a 8/16/28% Elemental DMG Bonus for the character's Elemental Type.
// The character will obtain 1 stack of Mistsplitter's Emblem in each of the
// following scenarios: Normal Attack deals Elemental DMG (stack lasts 5s),
// casting Elemental Burst (stack lasts 10s); Energy is less than 100% (stack
// disappears when Energy is full). Each stack's duration is calculated
// independently.
func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	//perm buff
	m := make([]float64, attributes.EndStatType)
	base := 0.09 + float64(r)*0.03
	for i := attributes.PyroP; i <= attributes.DendroP; i++ {
		m[i] = base
	}

	//stacking buff
	stack := 0.06 + float64(r)*0.02
	max := 0.03 + float64(r)*0.01
	bonus := attributes.EleToDmgP(char.Base.Element)

	//normal dealing dmg
	c.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.ActorIndex != char.Index {
			return false
		}
		if atk.Info.AttackTag != attacks.AttackTagNormal && atk.Info.AttackTag != attacks.AttackTagExtra {
			return false
		}
		if atk.Info.Element == attributes.Physical {
			return false
		}
		char.AddStatus(normalBuffKey, 300, true)
		return false
	}, fmt.Sprintf("mistsplitter-%v", char.Base.Key.String()))

	//using burst
	c.Events.Subscribe(event.OnBurst, func(args ...interface{}) bool {
		if c.Player.Active() != char.Index {
			return false
		}
		char.AddStatus(burstBuffKey, 600, true)
		return false

	}, fmt.Sprintf("mistsplitter-%v", char.Base.Key.String()))
	char.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("mistsplitter", -1),
		AffectedStat: attributes.NoStat,
		Amount: func() ([]float64, bool) {
			count := 0
			if char.Energy < char.EnergyMax {
				count++
			}
			if char.StatusIsActive(normalBuffKey) {
				count++
			}
			if char.StatusIsActive(burstBuffKey) {
				count++
			}
			dmg := float64(count) * stack
			if count >= 3 {
				dmg += max
			}
			m[bonus] = base + dmg
			return m, true
		},
	})

	return w, nil
}

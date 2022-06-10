package mistsplitter

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
)

func init() {
	core.RegisterWeaponFunc(keys.MistsplitterReforged, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

//Gain a 12% Elemental DMG Bonus for all elements and receive the might of the
//Mistsplitter's Emblem. At stack levels 1/2/3, the Mistsplitter's Emblem
//provides a 8/16/28% Elemental DMG Bonus for the character's Elemental Type.
//The character will obtain 1 stack of Mistsplitter's Emblem in each of the
//following scenarios: Normal Attack deals Elemental DMG (stack lasts 5s),
//casting Elemental Burst (stack lasts 10s); Energy is less than 100% (stack
//disappears when Energy is full). Each stack's duration is calculated
//independently.
func NewWeapon(c *core.Core, char *character.CharWrapper, p weapon.WeaponProfile) (weapon.Weapon, error) {
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
	normal := 0
	skill := 0

	//normal dealing dmg
	c.Events.Subscribe(event.OnDamage, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.ActorIndex != char.Index {
			return false
		}
		if atk.Info.AttackTag != combat.AttackTagNormal && atk.Info.AttackTag != combat.AttackTagExtra {
			return false
		}
		if atk.Info.Element == attributes.Physical {
			return false
		}
		normal = c.F + 300 // lasts 5 seconds
		return false
	}, fmt.Sprintf("mistsplitter-%v", char.Base.Name))

	//using burst
	c.Events.Subscribe(event.PreBurst, func(args ...interface{}) bool {
		if c.Player.Active() != char.Index {
			return false
		}
		skill = c.F + 600
		return false

	}, fmt.Sprintf("mistsplitter-%v", char.Base.Name))
	char.AddStatMod("mistsplitter", -1, attributes.NoStat, func() ([]float64, bool) {
		count := 0
		if char.Energy < char.EnergyMax {
			count++
		}
		if normal > c.F {
			count++
		}
		if skill > c.F {
			count++
		}
		dmg := float64(count) * stack
		if count >= 3 {
			dmg += max
		}
		m[bonus] = base + dmg
		return m, true
	})

	return w, nil
}

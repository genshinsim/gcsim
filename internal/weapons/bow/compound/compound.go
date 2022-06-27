package compound

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
	core.RegisterWeaponFunc(keys.CompoundBow, NewWeapon)
}

/*
* Normal Attack and Charged Attack hits increase ATK by 4/5/6/7/8% and Normal ATK SPD by
* 1.2/1.5/1.8/2.1/2.4% for 6s. Max 4 stacks. Can only occur once every 0.3s.
 */
type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p weapon.WeaponProfile) (weapon.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	m := make([]float64, attributes.EndStatType)

	incAtk := .03 + float64(r)*0.01
	incSpd := 0.009 + float64(r)*0.003

	stacks := 0
	maxStacks := 4
	stackExpiry := 0
	stackDuration := 360 // frames = 6s * 60 fps

	cd := 18 // frames = 0.3s * 60fps
	icd := 0

	c.Events.Subscribe(event.OnDamage, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)

		// Attack belongs to the equipped character
		if atk.Info.ActorIndex != char.Index {
			return false
		}

		// Active character has weapon equipped
		if c.Player.Active() != char.Index {
			return false
		}

		// Only apply on normal or charged attacks
		if (atk.Info.AttackTag != combat.AttackTagNormal) && (atk.Info.AttackTag != combat.AttackTagExtra) {
			return false
		}

		// Check if cd is up
		if icd > c.F {
			return false
		}

		// Reset stacks if they've expired
		if c.F > stackExpiry {
			stacks = 0
		}

		// Checks done, proc weapon passive
		// Increment stack count
		if stacks < maxStacks {
			stacks++
		}

		// trigger cd
		icd = c.F + cd
		stackExpiry = c.F + stackDuration

		//buff lasts 6 * 60 = 360 frames
		char.AddStatMod("compoundbow", stackDuration, attributes.NoStat, func() ([]float64, bool) {
			m[attributes.ATKP] = incAtk * float64(stacks)
			m[attributes.AtkSpd] = incSpd * float64(stacks)
			return m, true
		})

		return false
	}, fmt.Sprintf("compoundbow-%v", char.Base.Key.String()))

	return w, nil
}

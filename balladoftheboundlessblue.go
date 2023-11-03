package balladoftheboundlessblue

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
	core.RegisterWeaponFunc(keys.BalladOfTheBoundlessBlue, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

// Within 6s after Normal or Charged Attacks hit an opponent, Normal Attack DMG will be increased by 8% and 
// Charged Attack DMG will be increased by 6%. Max 3 stacks. This effect can be triggered once every 0.3s.
func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	m := make([]float64, attributes.EndStatType)

	incNADmg = .08 + float64(r-1)*.02
	incCADmg = .06 + float64(r-1)*0.015

	const buffKey = "boundless"

	stacks := 0
	maxStacks := 3
	const stackKey = "ballad-of-the-boundless-blue-stacks"
	stackDuration := 360 // frames = 6s * 60 fps
	const icdKey = "ballad-of-the-boundless-blue-icd"

	cd := int(0.3 * 60)

	c.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
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
		if (atk.Info.AttackTag != attacks.AttackTagNormal) && (atk.Info.AttackTag != attacks.AttackTagExtra) {
			return false
		}

		// check if cd is up
		if char.StatusIsActive(icdKey) {
			return false
		}

		// Reset stacks if they have expired
		if !char.StatusIsActive(stackKey) {
			stacks = 0
		}

		// Checks done
		// Increment stack count
		if stacks < maxStacks {
			stacks++
		}

		// trigger cd
		char.AddStatus(icdKey, cd, true)
		char.AddStatus(stackKey, stackDuration, true)

		if atk.Info.AttackTag == attacks.AttackTagNormal {
			char.AddAttackMod(character.AttackMod{
				Base: modifier.NewBaseWithHitlag("boundless-na", stackDuration),
				Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
					if atk.Info.AttackTag == attacks.AttackTagNormal {
						m[attributes.DmgP] = incNADmg * float64(stacks)
						return m, true
					}
					return nil, false
				},
			})
		}

		if atk.Info.AttackTag == attacks.AttackTagExtra {
			char.AddAttackMod(character.AttackMod{
				Base: modifier.NewBaseWithHitlag("boundless-ca", stackDuration),
				Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
					if atk.Info.AttackTag == attacks.AttackTagExtra {
						m[attributes.DmgP] = incCADmg * float64(stacks)
						return m, true
					}
					return nil, false
				},
			})
		}

		return false
	}, fmt.Sprintf("balladoftheboundlessblue-%v", char.Base.Key.String()))

	return w, nil
}



config.yml
package_name: balladoftheboundlessblue
genshin_id: 14511
key: balladoftheboundlessblue
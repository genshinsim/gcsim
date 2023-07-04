package ibispiercer

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
	core.RegisterWeaponFunc(keys.IbisPiercer, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

// The character's Elemental Mastery will increase by 40/50/60/70/80 within 6s after Charged Attacks hit opponents.
// Max 2 stacks. This effect can triggered once every 0.5s.
func NewWeapon(c *core.Core, char *character.CharWrapper, p weapon.WeaponProfile) (weapon.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	em := 30 + float64(r)*10

	m := make([]float64, attributes.EndStatType)

	stacks := 0
	maxStacks := 2
	const stackKey = "ibispiercer-stacks"
	stackDuration := 6 * 60
	const icdKey = "ibispiercer-icd"
	cd := int(0.5 * 60)

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

		// Only apply on charged attacks
		if atk.Info.AttackTag != attacks.AttackTagExtra {
			return false
		}

		// Check if cd is up
		if char.StatusIsActive(icdKey) {
			return false
		}

		// Reset stacks if they've expired
		if !char.StatusIsActive(stackKey) {
			stacks = 0
		}

		// Checks done, proc weapon passive
		// Increment stack count
		if stacks < maxStacks {
			stacks++
		}

		// trigger cd
		char.AddStatus(icdKey, cd, true)

		// add buff
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag(stackKey, stackDuration),
			AffectedStat: attributes.EM,
			Amount: func() ([]float64, bool) {
				m[attributes.EM] = em * float64(stacks)
				return m, true
			},
		})

		return false
	}, fmt.Sprintf("ibispiercer-%v", char.Base.Key.String()))

	return w, nil
}

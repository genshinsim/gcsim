package pines

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/weapons/common"
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
	core.RegisterWeaponFunc(keys.SongOfBrokenPines, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

// A part of the "Millennial Movement" that wanders amidst the winds.
// Increases ATK by 16%, and when Normal or Charged Attacks hit opponents,
// the character gains a Sigil of Whispers. This effect can be triggered once
// every 0.3s. When you possess 4 Sigils of Whispers, all of them will be
// consumed and all nearby party members will obtain the "Millennial
// Movement: Banner-Hymn" effect for 12s. "Millennial Movement: Banner-Hymn"
// increases Normal ATK SPD by 12% and increases ATK by 20%. Once this effect
// is triggered, you will not gain Sigils of Whispers for 20s. Of the many
// effects of the "Millennial Movement," buffs of the same type will not
// stack.
func NewWeapon(c *core.Core, char *character.CharWrapper, p weapon.WeaponProfile) (weapon.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	m := make([]float64, attributes.EndStatType)
	m[attributes.ATKP] = 0.12 + float64(r)*0.04
	char.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("pines-atk", -1),
		AffectedStat: attributes.NoStat,
		Amount: func() ([]float64, bool) {
			return m, true
		},
	})

	uniqueVal := make([]float64, attributes.EndStatType)
	uniqueVal[attributes.AtkSpd] = 0.09 + 0.03*float64(r)

	sharedVal := make([]float64, attributes.EndStatType)
	sharedVal[attributes.ATKP] = 0.15 + 0.05*float64(r)

	stacks := 0
	buffDuration := 12 * 60
	const icdKey = "songofbrokenpines-icd"
	icd := int(0.2 * 60)
	const cdKey = "songofbrokenpines-cooldown"
	cd := 20 * 60

	c.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.ActorIndex != char.Index {
			return false
		}
		if c.Player.Active() != char.Index {
			return false
		}
		if atk.Info.AttackTag != attacks.AttackTagNormal && atk.Info.AttackTag != attacks.AttackTagExtra {
			return false
		}
		if char.StatusIsActive(cdKey) {
			return false
		}
		if char.StatusIsActive(icdKey) {
			return false
		}

		char.AddStatus(icdKey, icd, true)
		stacks++
		if stacks == 4 {
			stacks = 0
			char.AddStatus(cdKey, cd, true)
			for _, char := range c.Player.Chars() {
				char.AddStatMod(character.StatMod{
					Base:         modifier.NewBaseWithHitlag("pines-proc", buffDuration),
					AffectedStat: attributes.AtkSpd,
					Amount: func() ([]float64, bool) {
						return uniqueVal, true
					},
				})
				char.AddStatMod(character.StatMod{
					Base:         modifier.NewBaseWithHitlag(common.MillennialKey, buffDuration),
					AffectedStat: attributes.ATKP,
					Amount: func() ([]float64, bool) {
						return sharedVal, true
					},
				})
			}
		}
		return false
	}, fmt.Sprintf("pines-%v", char.Base.Key.String()))

	return w, nil
}

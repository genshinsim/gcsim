package elegy

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
	core.RegisterWeaponFunc(keys.ElegyForTheEnd, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

// A part of the "Millennial Movement" that wanders amidst the winds.
// Increases Elemental Mastery by 60. When the Elemental Skills or Elemental Bursts
// of the character wielding this weapon hit opponents, that character gains a Sigil of Remembrance.
// This effect can be triggered once every 0.2s and can be triggered even if said character
// is not on the field. When you possess 4 Sigils of Remembrance, all of them will be consumed
// and all nearby party members will obtain the "Millennial Movement: Farewell Song" effect for 12s.
// "Millennial Movement: Farewell Song" increases Elemental Mastery by 100 and increases ATK by 20%.
// Once this effect is triggered, you will not gain Sigils of Remembrance for 20s.
// Of the many effects of the "Millennial Movement," buffs of the same type will not stack.
func NewWeapon(c *core.Core, char *character.CharWrapper, p weapon.WeaponProfile) (weapon.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	m := make([]float64, attributes.EndStatType)
	m[attributes.EM] = 45 + float64(r)*15
	char.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("elegy-em", -1),
		AffectedStat: attributes.NoStat,
		Amount: func() ([]float64, bool) {
			return m, true
		},
	})

	uniqueVal := make([]float64, attributes.EndStatType)
	uniqueVal[attributes.EM] = 75 + float64(r)*25

	sharedVal := make([]float64, attributes.EndStatType)
	sharedVal[attributes.ATKP] = .15 + float64(r)*0.05

	stacks := 0
	buffDuration := 12 * 60
	const icdKey = "elegy-sigil-icd"
	icd := int(0.2 * 60)
	const cdKey = "elegy-cd"
	cd := 20 * 60

	c.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.ActorIndex != char.Index {
			return false
		}
		switch atk.Info.AttackTag {
		case attacks.AttackTagElementalArt:
		case attacks.AttackTagElementalArtHold:
		case attacks.AttackTagElementalBurst:
		default:
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
					Base:         modifier.NewBaseWithHitlag("elegy-proc", buffDuration),
					AffectedStat: attributes.EM,
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
	}, fmt.Sprintf("elegy-%v", char.Base.Key.String()))

	return w, nil
}

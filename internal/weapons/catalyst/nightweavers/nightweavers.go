package nightweavers

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterWeaponFunc(keys.NightweaversLookingGlass, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

const (
	prayerKey  = "prayer-of-the-far-north"
	newMoonKey = "new-moon-verse"
)

// When the equipping character's Elemental Skill deals Hydro or Dendro DMG,
// they will gain Prayer of the Far North: Elemental Mastery is increased by
// 60/75/90/105/120 for 4.5s. When nearby party members trigger Lunar-Bloom reactions, the
// equipping character gains New Moon Verse: Elemental Mastery is increased
// by 60/75/90/105/120 for 10s. When both Prayer of the Far North and New Moon Verse are
// in effect, all nearby party members' Bloom DMG is increased by 120%/150%/180%/210%/240%, their
// Hyperbloom and Burgeon DMG is increased by 80%/100%/120%/140%/160%, and their Lunar-Bloom DMG
// is increased by 40%/50%/60%/70%/80%. This effect cannot stack. The aforementioned effects
// can be triggered even if the equipping character is off-field.
func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	m := make([]float64, attributes.EndStatType)
	m[attributes.EM] = 45 + float64(r)*15

	prayer := func(args ...any) {
		if _, ok := args[0].(*enemy.Enemy); !ok {
			return
		}

		atk := args[1].(*info.AttackEvent)
		if atk.Info.ActorIndex != char.Index() {
			return
		}
		switch atk.Info.AttackTag {
		case attacks.AttackTagElementalArt:
		case attacks.AttackTagElementalArtHold:
		default:
			return
		}

		switch atk.Info.Element {
		case attributes.Hydro:
		case attributes.Dendro:
		default:
			return
		}

		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag(prayerKey, 4.5*60),
			AffectedStat: attributes.EM,
			Amount: func() []float64 {
				return m
			},
		})
	}
	c.Events.Subscribe(event.OnEnemyDamage, prayer, fmt.Sprintf("prayer-of-the-far-north-%v", char.Base.Key.String()))

	newmoon := func(args ...any) {
		if _, ok := args[0].(*enemy.Enemy); !ok {
			return
		}
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag(newMoonKey, 10*60),
			AffectedStat: attributes.EM,
			Amount: func() []float64 {
				return m
			},
		})
	}
	c.Events.Subscribe(event.OnLunarBloom, newmoon, fmt.Sprintf("new-moon-verse-%v", char.Base.Key.String()))

	reactBuff := 0.3 + float64(r)*0.1
	// add reaction bonus when both of previous bonuses are active
	for _, otherChar := range c.Player.Chars() {
		otherChar.AddReactBonusMod(character.ReactBonusMod{
			Base: modifier.NewBase("nightweavers", -1),
			Amount: func(ai info.AttackInfo) float64 {
				if !char.StatusIsActive(prayerKey) || !char.StatusIsActive(newMoonKey) {
					return 0
				}

				switch ai.AttackTag {
				case attacks.AttackTagBloom, attacks.AttackTagBountifulCore:
					return reactBuff * 3
				case attacks.AttackTagHyperbloom, attacks.AttackTagBurgeon:
					return reactBuff * 2
				case attacks.AttackTagDirectLunarBloom:
					return reactBuff
				default:
					return 0
				}
			},
		})
	}

	return w, nil
}

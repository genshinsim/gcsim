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

func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	em := 45 + float64(r)*15
	m := make([]float64, attributes.EndStatType)
	m[attributes.EM] = em

	prayer := func(args ...any) {
		atk := args[1].(*info.AttackEvent)
		if atk.Info.ActorIndex != char.Index() {
			return
		}
		if atk.Info.AttackTag != attacks.AttackTagElementalArt && atk.Info.AttackTag != attacks.AttackTagElementalArtHold {
			return
		}
		if atk.Info.Element != attributes.Dendro && atk.Info.Element != attributes.Hydro {
			return
		}
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag(fmt.Sprintf("prayer-of-the-far-north-%v", char.Base.Key.String()), 10*60),
			Extra:        true,
			AffectedStat: attributes.EM,
			Amount: func() []float64 {
				return m
			},
		})
	}

	newmoon := func(args ...any) {
		atk := args[1].(*info.AttackEvent)
		if atk.Info.ActorIndex != char.Index() {
			return
		}
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag(fmt.Sprintf("new-moon-verse-%v", char.Base.Key.String()), 10*60),
			Extra:        true,
			AffectedStat: attributes.EM,
			Amount: func() []float64 {
				return m
			},
		})
	}

	// add reaction bonus when both of previous bonuses are active
	char.AddReactBonusMod(character.ReactBonusMod{
		Base: modifier.NewBase("nightweavers", -1),
		Amount: func(ai info.AttackInfo) float64 {
			if !char.StatusIsActive(fmt.Sprintf("prayer-of-the-far-north-%v", char.Base.Key.String())) ||
				!char.StatusIsActive(fmt.Sprintf("new-moon-verse-%v", char.Base.Key.String())) {
				return 0
			}

			switch ai.AttackTag {
			case attacks.AttackTagBloom, attacks.AttackTagBountifulCore:
				return 1.2
			case attacks.AttackTagHyperbloom, attacks.AttackTagBurgeon:
				return 0.8
			case attacks.AttackTagDirectLunarBloom:
				return 0.4
			default:
				return 0
			}
		},
	})

	c.Events.Subscribe(event.OnEnemyDamage, prayer, fmt.Sprintf("prayer-of-the-far-north-%v", char.Base.Key.String()))
	c.Events.Subscribe(event.OnLunarBloom, newmoon, fmt.Sprintf("new-moon-verse-%v", char.Base.Key.String()))

	return w, nil
}

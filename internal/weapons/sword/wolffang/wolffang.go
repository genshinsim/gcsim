package wolffang

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
	core.RegisterWeaponFunc(keys.WolfFang, NewWeapon)
}

type Weapon struct {
	Index  int
	refine int
	c      *core.Core
	char   *character.CharWrapper
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

// DMG dealt by Elemental Skill and Elemental Burst is increased by 16%. When an Elemental Skill hits an opponent,
// its CRIT Rate will be increased by 2%. When an Elemental Burst hits an opponent, its CRIT Rate will be increased by 2%.
// Both of these effects last 10s separately, have 4 max stacks, and can be triggered once every 0.1s.
func NewWeapon(c *core.Core, char *character.CharWrapper, p weapon.WeaponProfile) (weapon.Weapon, error) {
	w := &Weapon{
		refine: p.Refine,
		c:      c,
		char:   char,
	}

	mFirst := make([]float64, attributes.EndStatType)
	mFirst[attributes.DmgP] = 0.12 + 0.04*float64(p.Refine)
	char.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("wolf-fang", -1),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			switch atk.Info.AttackTag {
			case attacks.AttackTagElementalArt:
			case attacks.AttackTagElementalArtHold:
			case attacks.AttackTagElementalBurst:
			default:
				return nil, false
			}
			return mFirst, true
		},
	})

	w.addEvent("wolf-fang-skill", attacks.AttackTagElementalArt, attacks.AttackTagElementalArtHold)
	w.addEvent("wolf-fang-burst", attacks.AttackTagElementalBurst)

	return w, nil
}

func (w *Weapon) addEvent(name string, tags ...attacks.AttackTag) {
	stacks := 0
	cr := 0.015 + 0.005*float64(w.refine)
	m := make([]float64, attributes.EndStatType)
	icd := name + "-icd"

	w.c.Events.Subscribe(event.OnEnemyHit, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.ActorIndex != w.char.Index {
			return false
		}
		if w.c.Player.Active() != w.char.Index {
			return false
		}
		if !requiredTag(atk.Info.AttackTag, tags...) {
			return false
		}
		if w.char.StatusIsActive(icd) {
			return false
		}
		w.char.AddStatus(icd, 0.1*60, true)

		if !w.char.StatusIsActive(name) {
			stacks = 0
		}
		if stacks < 4 {
			stacks++
		}

		w.char.AddAttackMod(character.AttackMod{
			Base: modifier.NewBaseWithHitlag(name, 10*60),
			Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
				if !requiredTag(atk.Info.AttackTag, tags...) {
					return nil, false
				}
				m[attributes.CR] = cr * float64(stacks)
				return m, true
			},
		})
		return false
	}, fmt.Sprintf("%v-%v", name, w.char.Base.Key.String()))
}

func requiredTag(tag attacks.AttackTag, list ...attacks.AttackTag) bool {
	for _, value := range list {
		if tag == value {
			return true
		}
	}
	return false
}

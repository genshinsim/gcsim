package fruitfulhook

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
	core.RegisterWeaponFunc(keys.FruitfulHook, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

// Increase Plunging Attack CRIT Rate by 16%;
// After a Plunging Attack hits an opponent, Normal, Charged, and Plunging Attack DMG increased by 16% for 10s.
func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	// Increase Plunging Attack CRIT Rate
	mCR := make([]float64, attributes.EndStatType)
	mCR[attributes.CR] = 0.12 + 0.04*float64(r)
	char.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("fruitful-hook-cr", -1),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			if atk.Info.AttackTag == attacks.AttackTagPlunge {
				return mCR, true
			}
			return nil, false
		},
	})

	// After a Plunging Attack hits an opponent, Normal, Charged, and Plunging Attack DMG increased for 10s
	mDMG := make([]float64, attributes.EndStatType)
	mDMG[attributes.DmgP] = 0.12 + 0.04*float64(r)
	c.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.ActorIndex != char.Index {
			return false
		}
		if atk.Info.AttackTag != attacks.AttackTagPlunge {
			return false
		}
		char.AddAttackMod(character.AttackMod{
			Base: modifier.NewBaseWithHitlag("fruitful-hook-dmg%", 10*60),
			Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
				switch atk.Info.AttackTag {
				case attacks.AttackTagNormal:
				case attacks.AttackTagExtra:
				case attacks.AttackTagPlunge:
				default:
					return nil, false
				}
				return mDMG, true
			},
		})

		return false
	}, fmt.Sprintf("fruitful-hook-%v", char.Base.Key.String()))

	return w, nil
}

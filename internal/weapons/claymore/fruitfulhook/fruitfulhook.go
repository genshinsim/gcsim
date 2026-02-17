package fruitfulhook

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
		Amount: func(atk *info.AttackEvent, t info.Target) []float64 {
			if atk.Info.AttackTag == attacks.AttackTagPlunge {
				return mCR
			}
			return nil
		},
	})

	// After a Plunging Attack hits an opponent, Normal, Charged, and Plunging Attack DMG increased for 10s
	mDMG := make([]float64, attributes.EndStatType)
	mDMG[attributes.DmgP] = 0.12 + 0.04*float64(r)
	c.Events.Subscribe(event.OnEnemyDamage, func(args ...any) {
		atk := args[1].(*info.AttackEvent)
		if atk.Info.ActorIndex != char.Index() {
			return
		}
		if atk.Info.AttackTag != attacks.AttackTagPlunge {
			return
		}
		char.AddAttackMod(character.AttackMod{
			Base: modifier.NewBaseWithHitlag("fruitful-hook-dmg%", 10*60),
			Amount: func(atk *info.AttackEvent, t info.Target) []float64 {
				switch atk.Info.AttackTag {
				case attacks.AttackTagNormal:
				case attacks.AttackTagExtra:
				case attacks.AttackTagPlunge:
				default:
					return nil
				}
				return mDMG
			},
		})
	}, fmt.Sprintf("fruitful-hook-%v", char.Base.Key.String()))

	return w, nil
}

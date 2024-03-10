package cranesechoingcall

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
	core.RegisterWeaponFunc(keys.CranesEchoingCall, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

const (
	buffKey      = "crane-dmg%"
	buffDuration = 20 * 60
	energySrc    = "crane"
	energyIcdKey = "crane-energy-icd"
	energyIcd    = int(0.7 * 60)
)

// After the equipping character hits an opponent with a Plunging Attack,
// all nearby party members' Plunging Attacks will deal 28/41/54/67/80% increased DMG for 20s.
// When nearby party members hit opponents with Plunging Attacks,
// they will restore 2.5/2.75/3/3.25/3.5 Energy to the equipping character.
// Energy can be restored this way every 0.7s.
// This energy regain effect can be triggered even if the equipping character is not on the field.
func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	mDmg := make([]float64, attributes.EndStatType)
	mDmg[attributes.DmgP] = 0.15 + float64(r)*0.13

	energyRestore := 2.25 + float64(r)*0.25

	c.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)

		// can only trigger on plunge dmg
		if atk.Info.AttackTag != attacks.AttackTagPlunge {
			return false
		}

		// if dmg came from equipping char, then buff team plunge dmg
		if atk.Info.ActorIndex == char.Index {
			for _, char := range c.Player.Chars() {
				char.AddAttackMod(character.AttackMod{
					Base: modifier.NewBaseWithHitlag(buffKey, buffDuration),
					Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
						if atk.Info.AttackTag != attacks.AttackTagPlunge {
							return nil, false
						}
						return mDmg, true
					},
				})
			}
		}

		// restore energy regardless of who did plunge dmg
		if char.StatusIsActive(energyIcdKey) {
			return false
		}
		char.AddStatus(energyIcdKey, energyIcd, true)
		char.AddEnergy(energySrc, energyRestore)

		return false
	}, fmt.Sprintf("crane-onhit-%v", char.Base.Key.String()))

	return w, nil
}

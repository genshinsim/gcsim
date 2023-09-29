package eternalflow

import (
	"fmt"
	"math"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterWeaponFunc(keys.TomeOfTheEternalFlow, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

// HP is increased by 16%.
// When current HP increases or decreases, Charged Attack DMG will be increased by 14% for 4s.
// Max 3 stacks. This effect can be triggered once every 0.3s.
// When the character has 3 stacks or a third stack's duration refreshes, 8 Energy will be restored.
// This Energy restoration effect can be triggered once every 12s.

func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{}
	r := p.Refine
	stacks := 0
	const buffIcd = "eternalflow-ca-icd"
	const healKey = "eternalflow-heal"
	const drainKey = "eternalflow-drain"
	const energyIcd = "eternalflow-energy-icd"

	buffCA := make([]float64, attributes.EndStatType)

	hpp := 0.12 + float64(p.Refine)*0.04
	val := make([]float64, attributes.EndStatType)
	val[attributes.HPP] = hpp

	char.AddStatMod(character.StatMod{
		Base:         modifier.NewBaseWithHitlag("tome-of-the-eternal-flow-hpp", -1),
		AffectedStat: attributes.HPP,
		Amount: func() ([]float64, bool) {
			return val, true
		},
	})
	c.Events.Subscribe(event.OnPlayerHPDrain, func(args ...interface{}) bool {
		di := args[0].(player.DrainInfo)

		if c.Player.Active() != char.Index {
			return false
		}
		if di.ActorIndex != char.Index {
			return false
		}
		if di.Amount <= 0 {
			return false
		}
		if char.StatusIsActive(buffIcd) {
			return false
		}

		if !char.StatusIsActive(healKey) && !char.StatusIsActive(drainKey) {
			stacks = 0
		}
		stacks++
		if stacks > 3 {
			stacks = 3
		}

		char.AddStatus(buffIcd, 0.3*60, true)
		char.AddAttackMod(character.AttackMod{
			Base: modifier.NewBaseWithHitlag(drainKey, 4*60),
			Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
				buffCA[attributes.DmgP] = (0.105 + 0.035*float64(r)) * float64(stacks)
				switch atk.Info.AttackTag {
				case attacks.AttackTagExtra:
					return buffCA, true
				default:
					return nil, false
				}
			},
		})

		if stacks == 3 {
			if char.StatusIsActive(energyIcd) {
				return false
			}
			char.AddEnergy("eternal-flow-energy", 7+float64(r)*1)
			char.AddStatus(energyIcd, 12*60, true)
		}
		return false
	}, fmt.Sprintf("eternal-flow-ca-drain%v", char.Base.Key.String()))

	c.Events.Subscribe(event.OnHeal, func(args ...interface{}) bool {
		index := args[1].(int)
		amount := args[2].(float64)
		overheal := args[3].(float64)
		if c.Player.Active() != char.Index {
			return false
		}
		if index != char.Index {
			return false
		}
		if amount <= 0 {
			return false
		}
		if math.Abs(amount-overheal) <= 1e-9 {
			return false
		}
		if char.StatusIsActive(buffIcd) {
			return false
		}

		if !char.StatusIsActive(healKey) && !char.StatusIsActive(drainKey) {
			stacks = 0
		}
		stacks++
		if stacks > 3 {
			stacks = 3
		}

		char.AddStatus(buffIcd, 0.3*60, true)
		char.AddAttackMod(character.AttackMod{
			Base: modifier.NewBaseWithHitlag(drainKey, 4*60),
			Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
				buffCA[attributes.DmgP] = (0.105 + 0.035*float64(r)) * float64(stacks)
				switch atk.Info.AttackTag {
				case attacks.AttackTagExtra:
					return buffCA, true
				default:
					return nil, false
				}
			},
		})

		if stacks == 3 {
			if char.StatusIsActive(energyIcd) {
				return false
			}
			char.AddStatus(energyIcd, 12*60, true)
			char.AddEnergy("eternal-flow-energy", 7+float64(r)*1)
		}
		return false
	}, fmt.Sprintf("eternal-flow-ca-heal-%v", char.Base.Key.String()))
	return w, nil
}

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

const (
	buffKey   = "eternalflow-buff"
	buffIcd   = "eternalflow-icd"
	energyIcd = "eternalflow-energy-icd"
)

func init() {
	core.RegisterWeaponFunc(keys.TomeOfTheEternalFlow, NewWeapon)
}

type Weapon struct {
	stacks int
	core   *core.Core
	char   *character.CharWrapper
	refine int
	buffCA []float64
	Index  int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

// HP is increased by 16/20/24/28/32%.
// When current HP increases or decreases, Charged Attack DMG will be increased by 14/18/22/26/30% for 4s.
// Max 3 stacks. This effect can be triggered once every 0.3s.
// When the character has 3 stacks or a third stack's duration refreshes, 8/9/10/11/12 Energy will be restored.
// This Energy restoration effect can be triggered once every 12s.
func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{
		core:   c,
		char:   char,
		refine: p.Refine,
		buffCA: make([]float64, attributes.EndStatType),
	}

	hpp := 0.12 + float64(p.Refine)*0.04
	val := make([]float64, attributes.EndStatType)
	val[attributes.HPP] = hpp
	char.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("eternalflow-hpp", -1),
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

		w.onChangeHP()
		return false
	}, fmt.Sprintf("eternalflow-drain-%v", char.Base.Key.String()))

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
		// do not trigger if at max hp already
		if math.Abs(amount-overheal) <= 1e-9 {
			return false
		}

		w.onChangeHP()
		return false
	}, fmt.Sprintf("eternalflow-heal-%v", char.Base.Key.String()))
	return w, nil
}

func (w *Weapon) onChangeHP() {
	if w.char.StatusIsActive(buffIcd) {
		return
	}
	if !w.char.StatModIsActive(buffKey) {
		w.stacks = 0
	}
	if w.stacks < 3 {
		w.stacks++
	}

	w.char.AddStatus(buffIcd, 0.3*60, true)
	w.char.AddAttackMod(character.AttackMod{
		Base: modifier.NewBaseWithHitlag(buffKey, 4*60),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			w.buffCA[attributes.DmgP] = (0.10 + 0.04*float64(w.refine)) * float64(w.stacks)
			switch atk.Info.AttackTag {
			case attacks.AttackTagExtra:
				return w.buffCA, true
			default:
				return nil, false
			}
		},
	})

	if w.stacks == 3 && !w.char.StatusIsActive(energyIcd) {
		w.char.AddEnergy("eternalflow-energy", 7+float64(w.refine)*1)
		w.char.AddStatus(energyIcd, 12*60, true)
	}
}

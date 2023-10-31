package cashflow

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
	buffKey   = "cashflow-buff"
	buffIcd   = "cashflow-icd"
	atkSpdKey = "cashflow-atkspd"
)

func init() {
	core.RegisterWeaponFunc(keys.CashflowSupervision, NewWeapon)
}

type Weapon struct {
	stacks  int
	core    *core.Core
	char    *character.CharWrapper
	refine  int
	buffCA  []float64
	buffNA  []float64
	mAtkSpd []float64
	Index   int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

// ATK is increased by 16/20/24/28/32%. When current HP increases or decreases,
// Normal Attack DMG will be increased by 16/20/24/28/32%
// and Charged Attack DMG will be increased by 14/17.5/21/24.5/28% for 4s. Max 3 stacks.
// This effect can be triggered once every 0.3s.
// When the wielder has 3 stacks, ATK SPD will be increased by 8/10/12/14/16%.

func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{
		core:    c,
		char:    char,
		refine:  p.Refine,
		buffCA:  make([]float64, attributes.EndStatType),
		buffNA:  make([]float64, attributes.EndStatType),
		mAtkSpd: make([]float64, attributes.EndStatType),
	}

	atkp := 0.12 + float64(p.Refine)*0.04
	val := make([]float64, attributes.EndStatType)
	val[attributes.ATKP] = atkp

	char.AddStatMod(character.StatMod{
		Base:         modifier.NewBaseWithHitlag("cashflow-supervision-atkp", -1),
		AffectedStat: attributes.ATKP,
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

		w.onChangeHP()
		return false
	}, fmt.Sprintf("cashflow-na-ca-drain%v", char.Base.Key.String()))

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

		w.onChangeHP()
		return false
	}, fmt.Sprintf("cashflow-na-ca-heal-%v", char.Base.Key.String()))
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
			w.buffNA[attributes.DmgP] = (0.12 + 0.04*float64(w.refine)) * float64(w.stacks)
			w.buffCA[attributes.DmgP] = (0.105 + 0.035*float64(w.refine)) * float64(w.stacks)
			switch atk.Info.AttackTag {
			case attacks.AttackTagNormal:
				return w.buffNA, true
			case attacks.AttackTagExtra:
				return w.buffCA, true
			default:
				return nil, false
			}
		},
	})
	if w.stacks == 3 {
		w.char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag(atkSpdKey, 4*60),
			AffectedStat: attributes.NoStat,
			Amount: func() ([]float64, bool) {
				w.mAtkSpd[attributes.AtkSpd] = 0.06 + float64(w.refine)*0.02
				return w.mAtkSpd, true
			},
		})
	}
}

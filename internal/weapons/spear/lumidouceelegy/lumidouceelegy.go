package lumidouceelegy

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterWeaponFunc(keys.LumidouceElegy, NewWeapon)
}

const (
	atkBuffKey   = "lumidouceelegy-atk-buff"
	bonusBuffKey = "lumidouceelegy-bonus-buff"
	energyKey    = "lumidouceelegy-energy"
	energyICDKey = "lumidouceelegy-energy-icd"
)

type Weapon struct {
	Index  int
	refine int
	char   *character.CharWrapper
	stacks int
	buff   []float64
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := Weapon{
		refine: p.Refine,
		char:   char,
	}

	perm := make([]float64, attributes.EndStatType)
	perm[attributes.ATKP] = 0.11 + 0.04*float64(w.refine)
	char.AddStatMod(character.StatMod{
		Base:         modifier.NewBase(atkBuffKey, -1),
		AffectedStat: attributes.ATKP,
		Amount: func() ([]float64, bool) {
			return perm, true
		},
	})

	w.buff = make([]float64, attributes.EndStatType)

	c.Events.Subscribe(event.OnBurning, func(args ...interface{}) bool {
		_, ok := args[0].(*enemy.Enemy)
		atk := args[1].(*combat.AttackEvent)
		if !ok {
			return false
		}
		if atk.Info.ActorIndex != w.char.Index {
			return false
		}
		w.bonusCB()
		return false
	}, fmt.Sprintf("lumidouceelegy-on-burning-%v", char.Base.Key.String()))

	c.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		t, ok := args[0].(*enemy.Enemy)
		atk := args[1].(*combat.AttackEvent)
		if !ok {
			return false
		}
		if !t.IsBurning() {
			return false
		}
		if atk.Info.Element != attributes.Dendro {
			return false
		}
		if atk.Info.ActorIndex != w.char.Index {
			return false
		}
		w.bonusCB()
		return false
	}, fmt.Sprintf("lumidouceelegy-on-damage-%v", char.Base.Key.String()))

	return &w, nil
}

func (w *Weapon) bonusCB() {
	if !w.char.StatModIsActive(bonusBuffKey) {
		w.stacks = 0
	}
	if w.stacks < 2 {
		w.stacks++
	}

	if w.stacks == 2 && !w.char.StatusIsActive(energyICDKey) {
		w.char.AddStatus(energyICDKey, 12*60, true)
		w.char.AddEnergy(energyKey, float64(w.refine)+11.0)
	}

	w.char.AddAttackMod(character.AttackMod{
		Base: modifier.NewBaseWithHitlag(bonusBuffKey, 8*60),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			w.buff[attributes.DmgP] = (0.05*float64(w.refine) + 0.13) * float64(w.stacks)
			return w.buff, true
		},
	})
}

package nocturnes

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	ICDKey    = "nocturnes-curtain-call-icd"
	buffName  = "nocturnes-curtain-call-buff"
	hpBuffKey = buffName + "-hp"
	cdBuffKey = buffName + "-cd"
)

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	m := make([]float64, attributes.EndStatType)
	m[attributes.HPP] = 0.08 + float64(r)*0.02
	char.AddStatMod(character.StatMod{
		Base: modifier.NewBase("nocturnes-curtain-call-hp", -1),
		Amount: func() []float64 {
			return m
		},
	})

	hpBuff := make([]float64, attributes.EndStatType)
	hpBuff[attributes.HPP] = 0.12 + float64(r)*0.02

	critBuff := make([]float64, attributes.EndStatType)
	critBuff[attributes.CD] = 0.4 + float64(r)*0.2

	energy := 13.0 + float64(r)

	onDmgF := func(args ...any) {
		atk := args[1].(*info.AttackEvent)
		if atk.Info.ActorIndex != char.Index() {
			return
		}
		if !atk.Info.AttackTag.IsLunar() {
			return
		}
		nocturneBuff(char, energy, hpBuff, critBuff)
	}

	onReactF := func(args ...any) {
		atk := args[1].(*info.AttackEvent)
		if atk.Info.ActorIndex != char.Index() {
			return
		}

		nocturneBuff(char, energy, hpBuff, critBuff)
	}

	onLunarReactionAttackF := func(args ...any) {
		atk := args[1].(*info.AttackEvent)
		if atk.Info.ActorIndex != char.Index() {
			return
		}

		if !atk.Info.AttackTag.IsLunar() {
			return
		}

		if !char.StatusIsActive(cdBuffKey) {
			return
		}

		atk.Snapshot.Stats[attributes.CD] += critBuff[attributes.CD]
	}

	key := fmt.Sprintf("%v-%v", buffName, char.Base.Key.String())
	c.Events.Subscribe(event.OnSpecialReactionAttack, onLunarReactionAttackF, key)
	c.Events.Subscribe(event.OnEnemyDamage, onDmgF, key)
	c.Events.Subscribe(event.OnLunarCharged, onReactF, key)
	c.Events.Subscribe(event.OnLunarBloom, onReactF, key)
	c.Events.Subscribe(event.OnLunarCrystallize, onReactF, key)

	return w, nil
}

func nocturneBuff(char *character.CharWrapper, energy float64, hpBuff, critBuff []float64) {
	if !char.StatusIsActive(ICDKey) {
		char.AddEnergy("nocturnes-curtain-call", energy)
		char.AddStatus(ICDKey, 18*60, true)
	}

	char.AddStatMod(character.StatMod{
		Base: modifier.NewBaseWithHitlag(hpBuffKey, 12*60),
		Amount: func() []float64 {
			return hpBuff
		},
	})

	char.AddAttackMod(character.AttackMod{
		Base: modifier.NewBaseWithHitlag(cdBuffKey, 12*60),
		Amount: func(atk *info.AttackEvent, t info.Target) []float64 {
			if !atk.Info.AttackTag.IsLunar() {
				return nil
			}
			return critBuff
		},
	})
}

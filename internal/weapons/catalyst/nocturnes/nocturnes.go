package nocturnes

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
	core.RegisterWeaponFunc(keys.NocturnesCurtainCall, NewWeapon)
}

const (
	ICDKey  = "nocturnes-curtain-call-icd"
	buffKey = "nocturnes-curtain-call-buff"
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

	onHitF := func(args ...any) {
		atk := args[1].(*info.AttackEvent)
		if atk.Info.ActorIndex != char.Index() {
			return
		}
		if !attacks.AttackTagIsLunar(atk.Info.AttackTag) {
			return
		}
		w.nocturneBuff(char, r, hpBuff, critBuff)
	}

	onReactF := func(args ...any) {
		atk := args[1].(*info.AttackEvent)
		if atk.Info.ActorIndex != char.Index() {
			return
		}

		w.nocturneBuff(char, r, hpBuff, critBuff)
	}

	onLunarChargedF := func(args ...any) {
		atk := args[1].(*info.AttackEvent)
		if atk.Info.ActorIndex != char.Index() {
			return
		}

		if char.StatModIsActive(buffKey) {
			atk.Snapshot.Stats[attributes.CD] += 0.4 + float64(r)*0.2
		}
	}

	c.Events.Subscribe(event.OnLunarChargedReactionAttack, onLunarChargedF, buffKey)
	// c.Events.Subscribe(event.OnLunarCrystallizeReactionAttack, onLunarChargedF, buffKey)
	c.Events.Subscribe(event.OnEnemyDamage, onHitF, buffKey)
	c.Events.Subscribe(event.OnLunarCharged, onReactF, buffKey)
	c.Events.Subscribe(event.OnLunarBloom, onReactF, buffKey)
	// c.Events.Subscribe(event.OnLunarCrystallize, onReactF, buffKey)

	return w, nil
}

func (w *Weapon) nocturneBuff(char *character.CharWrapper, r int, hpBuff []float64, critBuff []float64) {
	if !char.StatusIsActive(ICDKey) {
		char.AddEnergy("nocturnes-curtain-call", 13+float64(r))
		char.AddStatus(ICDKey, 18*60, true)
	}

	char.AddStatMod(character.StatMod{
		Base: modifier.NewBaseWithHitlag(fmt.Sprintf("%v-hp", buffKey), 12*60),
		Amount: func() []float64 {
			return hpBuff
		},
	})

	char.AddAttackMod(character.AttackMod{
		Base: modifier.NewBaseWithHitlag(fmt.Sprintf("%v-cd", buffKey), 12*60),
		Amount: func(atk *info.AttackEvent, t info.Target) []float64 {
			if attacks.AttackTagIsLunar(atk.Info.AttackTag) {
				return nil
			}
			return critBuff
		},
	})
}

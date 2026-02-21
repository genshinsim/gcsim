package nocturnes

import (
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
	core.RegisterWeaponFunc(keys.NocturnesCourtainCall, NewWeapon)
}

const ICDKey = "nocturnes-courtain-call-icd"

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	m := make([]float64, attributes.EndStatType)
	m[attributes.HPP] = 0.075 + float64(r)*0.025
	char.AddStatMod(character.StatMod{
		Base: modifier.NewBase("nocturnes-courtain-call-hp", -1),
		Amount: func() []float64 {
			return m
		},
	})

	onHitF := func(args ...any) {
		atk := args[1].(*info.AttackEvent)
		if atk.Info.ActorIndex != char.Index() {
			return
		}
		if atk.Info.AttackTag < attacks.LunarReactionStartDelim || atk.Info.AttackTag > attacks.DirectLunarReactionEndDelim {
			return
		}
		w.nocturneBuff(c, char, r)
	}

	onReactF := func(args ...any) {
		atk := args[1].(*info.AttackEvent)
		if atk.Info.ActorIndex != char.Index() {
			return
		}

		w.nocturneBuff(c, char, r)
	}

	c.Events.Subscribe(event.OnEnemyHit, onHitF, "nocturnes-courtain-call-buff")
	c.Events.Subscribe(event.OnLunarCharged, onReactF, "nocturnes-courtain-call-buff")
	c.Events.Subscribe(event.OnLunarBloom, onReactF, "nocturnes-courtain-call-buff")
	// c.Events.Subscribe(event.OnLunarCrystallize, onReactF, "nocturnes-courtain-call-buff")

	return w, nil
}

func (w *Weapon) nocturneBuff(c *core.Core, char *character.CharWrapper, r int) {
	if !char.StatusIsActive(ICDKey) {
		char.AddEnergy("nocturnes-courtain-call", 15)
		char.AddStatus(ICDKey, 18*60, true)
	}

	m := make([]float64, attributes.EndStatType)
	m[attributes.HPP] = 0.09 + float64(r)*0.03
	char.AddStatMod(character.StatMod{
		Base:         modifier.NewBaseWithHitlag("nocturnes-courtain-call-hp", 12*60),
		AffectedStat: attributes.HPP,
		Amount: func() []float64 {
			return m
		},
	})

	c.Events.Subscribe(event.OnEnemyHit, func(args ...any) {
		atk := args[1].(*info.AttackEvent)

		if atk.Info.AttackTag < attacks.LunarReactionStartDelim || atk.Info.AttackTag > attacks.DirectLunarReactionEndDelim {
			return
		}
		atk.Snapshot.Stats[attributes.CD] += 0.6
	}, "nocturnes-courtain-call-buff")
}

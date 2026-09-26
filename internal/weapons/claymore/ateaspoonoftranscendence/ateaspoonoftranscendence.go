package ateaspoonoftranscendence

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	stackGainICDKey = "a-teaspoon-of-transcendence-stack-icd"
	stellarBuffKey  = "a-teaspoon-of-transcendence-stellar"
)

type Weapon struct {
	Index  int
	stacks int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }

func (w *Weapon) Init() error {
	return nil
}

func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{
		stacks: 0,
	}
	r := p.Refine

	atkBuff := make([]float64, attributes.EndStatType)
	atkBuff[attributes.ATKP] = 0.21 + 0.07*float64(r)

	stellarBuff := 0.12 + 0.04*float64(r)

	char.AddStatMod(character.StatMod{
		Base: modifier.NewBase("a-teaspoon-of-transcendence-atk", -1),
		Amount: func() []float64 {
			return atkBuff
		},
	})

	onChargeHit := func(args ...any) {
		atk, ok := args[1].(*info.AttackEvent)
		if !ok {
			return
		}

		if atk.Info.ActorIndex != char.Index() {
			return
		}

		if atk.Info.AttackTag != attacks.AttackTagExtra {
			return
		}

		if char.StatusIsActive(stackGainICDKey) {
			return
		}

		if !char.StatusIsActive(stellarBuffKey) {
			w.stacks = 0
		}

		w.stacks = min(w.stacks+1, 3)

		char.AddStatus(stackGainICDKey, 0.2*60, true)

		char.AddReactBonusMod(character.ReactBonusMod{
			Base: modifier.NewBaseWithHitlag(stellarBuffKey, 5*60),
			Amount: func(ai info.AttackInfo) float64 {
				if !ai.AttackTag.IsStellar() {
					return 0
				}

				return stellarBuff * float64(w.stacks)
			},
		})
	}

	c.Events.Subscribe(event.OnEnemyHit, onChargeHit, "a-teaspoon-of-transcendence-on-charge-hit-"+char.Base.Key.String())

	return w, nil
}

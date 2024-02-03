package common

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/model"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

type Wavebreaker struct {
	Index int
	data  *model.WeaponData
}

func (w *Wavebreaker) SetIndex(idx int)        { w.Index = idx }
func (w *Wavebreaker) Init() error             { return nil }
func (w *Wavebreaker) Data() *model.WeaponData { return w.data }

func NewWavebreaker(data *model.WeaponData) *Wavebreaker {
	return &Wavebreaker{data: data}
}

func (w *Wavebreaker) NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	r := p.Refine

	per := 0.09 + 0.03*float64(r)
	max := 0.3 + 0.1*float64(r)

	var amt float64

	c.Events.Subscribe(event.OnInitialize, func(args ...interface{}) bool {
		var energy float64

		for _, x := range c.Player.Chars() {
			energy += x.EnergyMax
		}

		amt = energy * per / 100
		if amt > max {
			amt = max
		}
		c.Log.NewEvent("wavebreaker dmg calc", glog.LogWeaponEvent, char.Index).
			Write("total", energy).
			Write("per", per).
			Write("max", max).
			Write("amt", amt)
		m := make([]float64, attributes.EndStatType)
		m[attributes.DmgP] = amt
		char.AddAttackMod(character.AttackMod{
			Base: modifier.NewBase("wavebreaker", -1),
			Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
				if atk.Info.AttackTag == attacks.AttackTagElementalBurst {
					return m, true
				}
				return nil, false
			},
		})
		return true
	}, fmt.Sprintf("wavebreaker-%v", char.Base.Key.String()))

	return w, nil
}

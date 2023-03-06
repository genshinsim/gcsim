package common

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

type Wavebreaker struct {
	Index int
}

func (b *Wavebreaker) SetIndex(idx int) { b.Index = idx }
func (b *Wavebreaker) Init() error      { return nil }

func NewWavebreaker(c *core.Core, char *character.CharWrapper, p weapon.WeaponProfile) (weapon.Weapon, error) {
	r := p.Refine
	b := &Wavebreaker{}

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

	return b, nil
}

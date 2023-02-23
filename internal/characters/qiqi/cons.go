package qiqi

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func (c *char) c1(a combat.AttackCB) {
	e, ok := a.Target.(*enemy.Enemy)
	if !ok {
		return
	}

	if !e.StatusIsActive(talismanKey) {
		return
	}

	c.AddEnergy("qiqi-c1", 2)
	c.Core.Log.NewEvent("Qiqi C1 Activation - Adding 2 energy", glog.LogCharacterEvent, c.Index).
		Write("target", a.Target.Key())
}

// Qiqi's Normal and Charge Attack DMG against opponents affected by Cryo is increased by 15%.
func (c *char) c2() {
	m := make([]float64, attributes.EndStatType)
	m[attributes.DmgP] = .15
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("qiqi-c2", -1),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			if atk.Info.AttackTag != attacks.AttackTagNormal && atk.Info.AttackTag != attacks.AttackTagExtra {
				return nil, false
			}

			e, ok := t.(*enemy.Enemy)
			if !ok {
				return nil, false
			}
			if !e.AuraContains(attributes.Cryo, attributes.Frozen) {
				return nil, false
			}

			return m, true
		},
	})
}

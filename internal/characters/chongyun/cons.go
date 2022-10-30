package chongyun

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func (c *char) c4() {
	const icdKey = "chongyun-c4-icd"
	c.Core.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		t, ok := args[0].(core.Reactable)
		if !ok {
			return false
		}
		if atk.Info.ActorIndex != c.Index {
			return false
		}
		if c.Core.Player.Active() != c.Index {
			return false
		}
		if c.StatusIsActive(icdKey) {
			return false
		}
		if !t.AuraContains(attributes.Cryo) {
			return false
		}

		c.AddEnergy("chongyun-c4", 2)

		c.Core.Log.NewEvent("chongyun c4 recovering 2 energy", glog.LogCharacterEvent, c.Index).
			Write("final energy", c.Energy)
		c.AddStatus(icdKey, 120, true)

		return false
	}, "chongyun-c4")
}

func (c *char) c6() {
	if c.Core.Combat.DamageMode {
		m := make([]float64, attributes.EndStatType)
		m[attributes.DmgP] = 0.15
		c.AddAttackMod(character.AttackMod{
			Base: modifier.NewBase("chongyun-c6", -1),
			Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
				if atk.Info.AttackTag != combat.AttackTagElementalBurst {
					return nil, false
				}
				x, ok := t.(*enemy.Enemy)
				if !ok {
					return nil, false
				}
				if x.HP()/x.MaxHP() < c.HPCurrent/c.MaxHP() {
					return m, true
				}
				return nil, false
			},
		})
	}
}

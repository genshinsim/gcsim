package chongyun

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const c4ICDKey = "chongyun-c4-icd"

// Chongyun regenerates 1 Energy every time he hits an opponent affected by Cryo.
// This effect can only occur once every 2s.
func (c *char) makeC4Callback() combat.AttackCBFunc {
	if c.Base.Cons < 4 {
		return nil
	}
	return func(a combat.AttackCB) {
		e, ok := a.Target.(*enemy.Enemy)
		if !ok {
			return
		}
		if c.Core.Player.Active() != c.Index {
			return
		}
		if !e.AuraContains(attributes.Cryo) {
			return
		}
		if c.StatusIsActive(c4ICDKey) {
			return
		}
		c.AddStatus(c4ICDKey, 2*60, true)
		c.AddEnergy("chongyun-c4", 2)
		c.Core.Log.NewEvent("chongyun c4 recovering 2 energy", glog.LogCharacterEvent, c.Index).
			Write("final energy", c.Energy)
	}
}

func (c *char) c6() {
	if c.Core.Combat.DamageMode {
		m := make([]float64, attributes.EndStatType)
		m[attributes.DmgP] = 0.15
		c.AddAttackMod(character.AttackMod{
			Base: modifier.NewBase("chongyun-c6", -1),
			Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
				if atk.Info.AttackTag != attacks.AttackTagElementalBurst {
					return nil, false
				}
				x, ok := t.(*enemy.Enemy)
				if !ok {
					return nil, false
				}
				if x.HP()/x.MaxHP() < c.CurrentHPRatio() {
					return m, true
				}
				return nil, false
			},
		})
	}
}

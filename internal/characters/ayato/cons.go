package ayato

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func (c *char) c1() {
	m := make([]float64, attributes.EndStatType)
	m[attributes.DmgP] = 0.4
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("ayato-c1", -1),
		Amount: func(a *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			if a.Info.AttackTag != combat.AttackTagNormal || t.HP()/t.MaxHP() > 0.5 {
				return nil, false
			}
			return m, true
		},
	})
}

func (c *char) c2() {
	m := make([]float64, attributes.EndStatType)
	m[attributes.HPP] = 0.5
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("ayato-c2", -1),
		AffectedStat: attributes.HPP,
		Amount: func() ([]float64, bool) {
			if c.stacks >= 3 {
				return m, true
			}
			return nil, false
		},
	})
}

func (c *char) c6() {
	c.Core.Events.Subscribe(event.OnDamage, func(args ...interface{}) bool {
		if c.Core.Player.Active() != c.Index {
			return false
		}
		if !c.c6ready {
			return false
		}
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.AttackTag != combat.AttackTagNormal {
			return false
		}
		ai := combat.AttackInfo{
			Abil:       fmt.Sprintf("Normal %v", c.NormalCounter),
			ActorIndex: c.Index,
			AttackTag:  combat.AttackTagNormal,
			ICDTag:     combat.ICDTagNormalAttack,
			ICDGroup:   combat.ICDGroupDefault,
			Element:    attributes.Hydro,
			Durability: 25,
			Mult:       4.5,
		}
		c.Core.QueueAttack(ai, combat.NewDefCircHit(2, false, combat.TargettableEnemy), 20, 20)
		c.Core.QueueAttack(ai, combat.NewDefCircHit(2, false, combat.TargettableEnemy), 22, 22)

		c.Core.Log.NewEvent("ayato c6 proc'd", glog.LogCharacterEvent, c.Index)
		c.c6ready = false
		return false
	}, "ayato-c6")
}

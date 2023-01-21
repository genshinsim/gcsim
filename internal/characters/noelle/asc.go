package noelle

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/core/player/shield"
)

const a1IcdKey = "noelle-a1-icd"

func (c *char) a1() {
	c.Core.Events.Subscribe(event.OnPlayerHPDrain, func(args ...interface{}) bool {
		di := args[0].(player.DrainInfo)
		if di.Amount <= 0 {
			return false
		}
		if c.StatusIsActive(a1IcdKey) {
			return false
		}
		active := c.Core.Player.ActiveChar()
		if active.HPCurrent/active.MaxHP() >= 0.3 {
			return false
		}
		c.AddStatus(a1IcdKey, 3600, false)
		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "A1 Shield",
			AttackTag:  combat.AttackTagNone,
		}
		snap := c.Snapshot(&ai)

		//add shield
		x := snap.BaseDef*(1+snap.Stats[attributes.DEFP]) + snap.Stats[attributes.DEF]
		c.Core.Player.Shields.Add(&shield.Tmpl{
			Src:        c.Core.F,
			ShieldType: shield.ShieldNoelleA1,
			Name:       "Noelle A1",
			HP:         4 * x,
			Ele:        attributes.Cryo,
			Expires:    c.Core.F + 1200, //20 sec
		})
		return false
	}, "noelle-a1")
}

func (c *char) a4() combat.AttackCBFunc {
	done := false
	return func(a combat.AttackCB) {
		if a.Target.Type() != combat.TargettableEnemy {
			return
		}
		if done {
			return
		}
		done = true
		c.a4Counter++
		if c.a4Counter == 4 {
			c.a4Counter = 0
			if c.Cooldown(action.ActionSkill) > 0 {
				c.ReduceActionCooldown(action.ActionSkill, 60)
			}
		}
	}
}

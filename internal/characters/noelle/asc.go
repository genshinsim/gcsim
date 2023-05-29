package noelle

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/core/player/shield"
	"github.com/genshinsim/gcsim/pkg/core/targets"
)

const a1IcdKey = "noelle-a1-icd"

// When Noelle is in the party but not on the field,
// this ability triggers automatically when your active character's HP falls below 30%:
// Creates a shield for your active character that lasts for 20s and absorbs DMG equal to 400% of Noelle's DEF.
// The shield has a 150% DMG Absorption effectiveness against all Elemental and Physical DMG.
// This effect can only occur once every 60s.
func (c *char) a1() {
	if c.Base.Ascension < 1 {
		return
	}
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
			AttackTag:  attacks.AttackTagNone,
		}
		snap := c.Snapshot(&ai)

		//add shield
		x := snap.BaseDef*(1+snap.Stats[attributes.DEFP]) + snap.Stats[attributes.DEF]
		c.Core.Player.Shields.Add(&shield.Tmpl{
			ActorIndex: c.Index,
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

// Noelle will decrease the CD of Breastplate by 1s for every 4 Normal or Charged Attack hits she scores on opponents.
// One hit may be counted every 0.1s.
func (c *char) makeA4CB() combat.AttackCBFunc {
	if c.Base.Ascension < 4 {
		return nil
	}
	done := false
	return func(a combat.AttackCB) {
		if a.Target.Type() != targets.TargettableEnemy {
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

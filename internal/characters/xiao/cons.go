package xiao

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

// Implements Xiao C2:
// When in the party and not on the field, Xiao's Energy Recharge is increased by 25%
func (c *char) c2() {
	m := make([]float64, attributes.EndStatType)
	m[attributes.ER] = 0.25
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("xiao-c2", -1),
		AffectedStat: attributes.ER,
		Amount: func() ([]float64, bool) {
			if c.Core.Player.Active() != c.Index {
				return m, true
			}
			return nil, false
		},
	})
}

// Implements Xiao C4:
// When Xiao's HP falls below 50%, he gains a 100% DEF Bonus.
func (c *char) c4() {
	//TODO: in game this is actually a check every 0.3s. if hp is < 50% then buff is active until
	//the next time check takes places
	m := make([]float64, attributes.EndStatType)
	m[attributes.DEFP] = 1
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("xiao-c4", -1),
		AffectedStat: attributes.DEFP,
		Amount: func() ([]float64, bool) {
			if c.HPCurrent/c.MaxHP() <= 0.5 {
				return m, true
			}
			return nil, false
		},
	})
}

const c6BuffKey = "xiao-c6"

// Implements Xiao C6:
// While under the effect of Bane of All Evil, hitting at least 2 opponents with Xiao's Plunge Attack will immediately grant him 1 charge of Lemniscatic Wind Cycling, and for the next 1s, he may use Lemniscatic Wind Cycling while ignoring its CD.
// Adds an OnDamage event checker - if we record two or more instances of plunge damage, then activate C6
func (c *char) c6() {
	c.Core.Events.Subscribe(event.OnDamage, func(evt event.EventPayload) bool {
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.ActorIndex != c.Index {
			return false
		}
		if !((atk.Info.Abil == "High Plunge") || (atk.Info.Abil == "Low Plunge")) {
			return false
		}
		if !c.StatusIsActive(burstBuffKey) {
			return false
		}
		// Stops after reaching 2 hits on a single plunge.
		// Plunge frames are greater than duration of C6 so this will always refresh properly.
		if c.StatusIsActive(c6BuffKey) {
			//do nothing if buff still on
			return false
		}
		if c.c6Src != atk.SourceFrame {
			c.c6Src = atk.SourceFrame
			c.c6Count = 0
			return false
		}

		c.c6Count++

		// Prevents activation more than once in a single plunge attack
		if c.c6Count == 2 {
			c.ResetActionCooldown(action.ActionSkill)

			c.AddStatus(c6BuffKey, 60, true)
			c.Core.Log.NewEvent("Xiao C6 activated", glog.LogCharacterEvent, c.Index).
				Write("new E charges", c.Tags["eCharge"]).
				Write("expiry", c.Core.F+60)

			c.c6Count = 0
			return false
		}
		return false
	}, "xiao-c6")
}

package mona

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

// A1:
// After she has used Illusory Torrent for 2s, if there are any opponents nearby,
// Mona will automatically create a Phantom.
// A Phantom created in this manner lasts for 2s, and its explosion DMG is equal to 50% of Mirror Reflection of Doom.
func (c *char) a1() func() {
	return func() {
		// do nothing if not Mona
		if c.Core.Player.Active() != c.Index {
			return
		}
		// do nothing if we aren't dashing anymore
		if c.Core.Player.CurrentState() != action.DashState {
			return
		}
		c.Core.Log.NewEvent("mona-a1 phantom added", glog.LogCharacterEvent, c.Index).
			Write("expiry:", c.Core.F+120)
		// queue up phantom explosion
		c.Core.Tasks.Add(func() {
			aiExplode := combat.AttackInfo{
				ActorIndex: c.Index,
				Abil:       "Mirror Reflection of Doom (A1 Explode)",
				AttackTag:  combat.AttackTagElementalArt,
				ICDTag:     combat.ICDTagNone,
				ICDGroup:   combat.ICDGroupDefault,
				StrikeType: combat.StrikeTypeDefault,
				Element:    attributes.Hydro,
				Durability: 25,
				Mult:       0.5 * skill[c.TalentLvlSkill()],
			}
			c.Core.QueueAttack(aiExplode, combat.NewCircleHit(c.Core.Combat.Player(), 5), 0, 0)
		}, 120)
		// queue up next A1 check because Mona's still dashing
		// different Phantoms coexist and don't overwrite each other
		c.Core.Tasks.Add(c.a1(), 120) // check again in 2s
	}
}

// A4:
// Increases Mona's Hydro DMG Bonus by a degree equivalent to 20% of her Energy Recharge rate.
func (c *char) a4() {
	m := make([]float64, attributes.EndStatType)
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("mona-a4", -1),
		Amount: func(atk *combat.AttackEvent, _ combat.Target) ([]float64, bool) {
			m[attributes.HydroP] = .2 * (1 + atk.Snapshot.Stats[attributes.ER])
			return m, true
		},
	})
}

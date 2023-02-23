package mona

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

// After she has used Illusory Torrent for 2s, if there are any opponents nearby,
// Mona will automatically create a Phantom.
// A Phantom created in this manner lasts for 2s, and its explosion DMG is equal to 50% of Mirror Reflection of Doom.
//
// - checks for ascension level in dash.go to avoid queuing this up only to fail the ascension level check
func (c *char) a1() {
	// do nothing if not Mona
	if c.Core.Player.Active() != c.Index {
		return
	}
	// do nothing if we aren't dashing anymore
	if c.Core.Player.CurrentState() != action.DashState {
		return
	}
	enemies := c.Core.Combat.EnemiesWithinArea(combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 15), nil)
	if enemies != nil {
		c.Core.Log.NewEvent("mona-a1 phantom added", glog.LogCharacterEvent, c.Index).
			Write("expiry:", c.Core.F+120)
		// queue up phantom explosion
		phantomPos := c.Core.Combat.Player()
		c.Core.Tasks.Add(func() {
			aiExplode := combat.AttackInfo{
				ActorIndex: c.Index,
				Abil:       "Mirror Reflection of Doom (A1 Explode)",
				AttackTag:  attacks.AttackTagElementalArt,
				ICDTag:     combat.ICDTagNone,
				ICDGroup:   combat.ICDGroupDefault,
				StrikeType: attacks.StrikeTypeDefault,
				Element:    attributes.Hydro,
				Durability: 25,
				Mult:       0.5 * skill[c.TalentLvlSkill()],
			}
			c.Core.QueueAttack(aiExplode, combat.NewCircleHitOnTarget(phantomPos, nil, 5), 0, 0)
		}, 120)
	}
	// queue up next A1 check because Mona's still dashing
	// different Phantoms coexist and don't overwrite each other
	c.Core.Tasks.Add(c.a1, 120) // check again in 2s
}

// Increases Mona's Hydro DMG Bonus by a degree equivalent to 20% of her Energy Recharge rate.
func (c *char) a4() {
	if c.Base.Ascension < 4 {
		return
	}
	m := make([]float64, attributes.EndStatType)
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("mona-a4", -1),
		Amount: func(atk *combat.AttackEvent, _ combat.Target) ([]float64, bool) {
			m[attributes.HydroP] = .2 * (1 + atk.Snapshot.Stats[attributes.ER])
			return m, true
		},
	})
}

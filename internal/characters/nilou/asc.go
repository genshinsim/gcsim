package nilou

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
	"github.com/genshinsim/gcsim/pkg/reactable"
)

const (
	a1Status = "nilou-a1"
	a4Mod    = "nilou-a4"
)

// When all characters in the party are all Dendro or Hydro, and there are at least one Dendro character and one Hydro character:
// The completion of the third dance step of Nilou’s Dance of Haftkarsvar will grant all nearby characters the Golden Chalice’s Bounty
// for 30s upon its completion.
// Characters under the effect of Golden Chalice’s Bounty will increase the Elemental Mastery of all nearby characters by 100 for 10s
// whenever they are hit by Dendro attacks. Also, triggering the Bloom reaction will create Bountiful Cores instead of Dendro Cores.
// Such Cores will burst very quickly after being created, and they have larger AoEs.
// Bountiful Cores cannot trigger Hyperbloom or Burgeon, and they share an upper numerical limit with Dendro Cores. Bountiful Core DMG
// is considered DMG dealt by Dendro Cores produced by Bloom.
// Should the party not meet the conditions for this Passive Talent, any existing Golden Chalice’s Bounty effects will be canceled.
func (c *char) a1() {
	if c.Base.Ascension < 1 || !c.onlyBloomTeam {
		return
	}

	for _, this := range c.Core.Player.Chars() {
		this.AddStatus(a1Status, 30*60, true)
	}
	c.a4()

	// Bountiful Cores
	c.Core.Events.Subscribe(event.OnDendroCore, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		char := c.Core.Player.ByIndex(atk.Info.ActorIndex)
		if !char.StatusIsActive(a1Status) {
			return false
		}
		g, ok := args[0].(*reactable.DendroCore)
		if !ok {
			return false
		}

		b := newBountifulCore(c.Core, g.Gadget.Pos(), atk)
		b.Gadget.SetKey(g.Gadget.Key())
		c.Core.Combat.ReplaceGadget(g.Key(), b)
		//prevent blowing up
		g.Gadget.OnExpiry = nil
		g.Gadget.OnKill = nil

		return false
	}, "nilou-a1-cores")

	c.Core.Events.Subscribe(event.OnPlayerHit, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.Element != attributes.Dendro {
			return false
		}
		char := c.Core.Player.ActiveChar()
		if !char.StatusIsActive(a1Status) {
			return false
		}

		m := make([]float64, attributes.EndStatType)
		m[attributes.EM] = 100
		for _, this := range c.Core.Player.Chars() {
			this.AddStatMod(character.StatMod{
				Base:         modifier.NewBaseWithHitlag("nilou-a1-em", 10*60),
				AffectedStat: attributes.EM,
				Amount: func() ([]float64, bool) {
					return m, true
				},
			})
		}

		return false
	}, "nilou-a1")
}

// Every 1,000 points of Nilou’s Max HP above 30,000 will cause the DMG dealt by Bountiful Cores created by characters affected
// by Golden Chalice’s Bounty to increase by 9%.
// The maximum increase in Bountiful Core DMG that can be achieved this way is 400%.
func (c *char) a4() {
	if c.Base.Ascension < 4 {
		return
	}
	for _, this := range c.Core.Player.Chars() {
		this.AddReactBonusMod(character.ReactBonusMod{
			Base: modifier.NewBaseWithHitlag(a4Mod, 30*60),
			Amount: func(ai combat.AttackInfo) (float64, bool) {
				if ai.AttackTag != combat.AttackTagBloom {
					return 0, false
				}

				// check is bountiful core?
				var t combat.Gadget
				for _, v := range c.Core.Combat.Gadgets() {
					if v != nil && v.Key() == ai.DamageSrc {
						t = v
					}
				}
				if _, ok := t.(*BountifulCore); !ok {
					return 0, false
				}

				c.Core.Combat.Log.NewEvent("adding nilou a4 bonus", glog.LogCharacterEvent, c.Index).Write("bonus", c.a4Bonus)
				return c.a4Bonus, false
			},
		})
	}
	c.a4Src = c.Core.F
	c.QueueCharTask(c.updateA4Bonus(c.a4Src), 0.5*60)
}

func (c *char) updateA4Bonus(src int) func() {
	return func() {
		if c.a4Src != src {
			return
		}
		if !c.ReactBonusModIsActive(a4Mod) {
			return
		}

		c.a4Bonus = (c.MaxHP() - 30000) * 0.001 * 0.09
		if c.a4Bonus < 0 {
			c.a4Bonus = 0
		} else if c.a4Bonus > 4 {
			c.a4Bonus = 4
		}

		c.QueueCharTask(c.updateA4Bonus(src), 0.5*60)
	}
}

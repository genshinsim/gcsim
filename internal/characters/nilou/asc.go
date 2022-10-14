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
	a1Status = "golden-chalice" // TODO: or "nilou-a1"?
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
	chars := c.Core.Player.Chars()
	//count number of ele first
	count := make(map[attributes.Element]int)
	for _, this := range chars {
		count[this.Base.Element]++
	}
	if count[attributes.Hydro] == 0 || count[attributes.Dendro] == 0 || count[attributes.Hydro]+count[attributes.Dendro] != len(chars) {
		return
	}

	for _, this := range chars {
		this.AddStatus(a1Status, 30*60, true)

		// Every 1,000 points of Nilou’s Max HP above 30,000 will cause the DMG dealt by Bountiful Cores created by characters affected
		// by Golden Chalice’s Bounty to increase by 9%.
		// The maximum increase in Bountiful Core DMG that can be achieved this way is 400%.
		this.AddReactBonusMod(character.ReactBonusMod{
			Base: modifier.NewBaseWithHitlag("nilou-a4", 30*60),
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

				a4Bonus := float64(int(c.MaxHP()-30000)/1000) * 0.09
				if a4Bonus < 0 {
					a4Bonus = 0
				} else if a4Bonus > 4 {
					a4Bonus = 4
				}
				c.Core.Combat.Log.NewEvent("adding a4 bonus", glog.LogCharacterEvent, c.Index).Write("bonus", a4Bonus)
				return a4Bonus, false
			},
		})
	}

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

		x, y := g.Gadget.Pos()
		b := newBountifulCore(c.Core, x, y, atk)
		b.Gadget.SetKey(g.Gadget.Key())
		c.Core.Combat.ReplaceGadget(g.Key(), b)
		//prevent blowing up
		g.Gadget.OnExpiry = nil
		g.Gadget.OnKill = nil

		return false
	}, "nilou-a1-cores")

	c.Core.Events.Subscribe(event.OnDamage, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		t, ok := args[0].(combat.Target)
		if !ok {
			return false
		}
		if t.Type() != combat.TargettablePlayer || atk.Info.Element != attributes.Dendro {
			return false
		}
		char := c.Core.Player.ByIndex(t.Index())
		if !char.StatusIsActive(a1Status) {
			return false
		}

		m := make([]float64, attributes.EndStatType)
		m[attributes.EM] = 100
		for _, this := range chars {
			this.AddStatMod(character.StatMod{
				Base:         modifier.NewBaseWithHitlag("nilou-a1", 10*60),
				AffectedStat: attributes.EM,
				Amount: func() ([]float64, bool) {
					return m, true
				},
			})
		}

		return false
	}, "nilou-a1")
}

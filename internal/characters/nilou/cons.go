package nilou

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

// Dance of Haftkarsvar will be enhanced as follows:
// · Luminous Illusion DMG is increased by 65%.
// · The Tranquility Aura’s duration is extended by 6s.
func (c *char) c1() {
	m := make([]float64, attributes.EndStatType)
	m[attributes.DmgP] = 0.65

	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("nilou-c1", -1),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			if atk.Info.Abil != "Luminous Illusion" {
				return nil, false
			}
			return m, true
		},
	})
}

// After characters affected by the Golden Chalice’s Bounty deal Hydro DMG to opponents, that opponent’s Hydro RES will be decreased by 35% for 10s.
// After a triggered Bloom reaction deals DMG to opponents, their Dendro RES will be decreased by 35% for 10s.
// You need to have unlocked the “Court of Dancing Petals” Talent.
func (c *char) c2() {
	c.Core.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		dmg := args[2].(float64)
		t, ok := args[0].(*enemy.Enemy)
		if !ok {
			return false
		}
		if dmg == 0 {
			return false
		}

		char := c.Core.Player.ByIndex(atk.Info.ActorIndex)
		if !char.StatusIsActive(a1Status) {
			return false
		}

		if atk.Info.Element == attributes.Hydro {
			t.AddResistMod(combat.ResistMod{
				Base:  modifier.NewBaseWithHitlag("nilou-c2-hydro", 10*60),
				Ele:   attributes.Hydro,
				Value: -0.35,
			})
		} else if atk.Info.AttackTag == attacks.AttackTagBloom {
			t.AddResistMod(combat.ResistMod{
				Base:  modifier.NewBaseWithHitlag("nilou-c2-dendro", 10*60),
				Ele:   attributes.Dendro,
				Value: -0.35,
			})
		}

		return false
	}, "nilou-c2")
}

// After the third dance step of Dance of Haftkarsvar‘s Pirouette hits opponents, Nilou will gain 15 Elemental Energy,
// and DMG from her Dance of Abzendegi: Distant Dreams, Listening Spring will be increased by 50% for 8s.
func (c *char) c4() {
	c.AddEnergy("nilou-c4", 15)

	m := make([]float64, attributes.EndStatType)
	m[attributes.DmgP] = 0.5
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBaseWithHitlag("nilou-c4", 8*60),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			if atk.Info.AttackTag != attacks.AttackTagElementalBurst {
				return nil, false
			}
			return m, true
		},
	})
}

func (c *char) c4cb() combat.AttackCBFunc {
	if c.Base.Cons < 4 {
		return nil
	}
	if c.Tag(skillStep) != 2 || !c.StatusIsActive(pirouetteStatus) {
		return nil
	}

	done := false
	return func(a combat.AttackCB) {
		if a.Target.Type() != combat.TargettableEnemy {
			return
		}
		if done {
			return
		}
		c.c4()
		done = true
	}
}

// For every 1,000 points of Max HP, Nilou’s CRIT Rate and CRIT DMG will increase by 0.6% and 1.2% respectively.
// The maximum increase in CRIT Rate and CRIT DMG is 30% and 60% respectively.
func (c *char) c6() {
	// cr and cd separately to avoid stack overflow due to NoStat attribute
	mCR := make([]float64, attributes.EndStatType)
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("nilou-c6-cr", -1),
		AffectedStat: attributes.CR,
		Amount: func() ([]float64, bool) {
			cr := c.MaxHP() * 0.001 * 0.006
			if cr > 0.3 {
				cr = 0.3
			}
			mCR[attributes.CR] = cr
			return mCR, true
		},
	})

	mCD := make([]float64, attributes.EndStatType)
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("nilou-c6-cd", -1),
		AffectedStat: attributes.CD,
		Amount: func() ([]float64, bool) {
			cd := c.MaxHP() * 0.001 * 0.012
			if cd > 0.6 {
				cd = 0.6
			}
			mCD[attributes.CD] = cd
			return mCD, true
		},
	})
}

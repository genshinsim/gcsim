package razor

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

// Picking up an Elemental Orb or Particle increases Razor's DMG by 10% for 8s.
func (c *char) c1() {
	if c.Base.Cons < 1 {
		return
	}

	val := make([]float64, attributes.EndStatType)
	val[attributes.DmgP] = 0.1

	c.Core.Events.Subscribe(event.OnParticleReceived, func(args ...interface{}) bool {
		c.AddStatMod(character.StatMod{Base: modifier.NewBase("razor-c1", 8*60), AffectedStat: attributes.DmgP, Amount: func() ([]float64, bool) {
			return val, true
		}})
		return false
	}, "razor-c1")
}

// Increases CRIT Rate against opponents with less than 30% HP by 10%.
func (c *char) c2() {
	if c.Base.Cons < 2 {
		return
	}

	m := make([]float64, attributes.EndStatType)
	m[attributes.CR] = 0.1

	c.AddAttackMod(character.AttackMod{Base: modifier.NewBase("razor-c2", -1), Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
		if t.HP()/t.MaxHP() < 0.3 {
			return m, true
		}
		return nil, false
	}})
}

// When casting Claw and Thunder (Press), opponents hit will have their DEF decreased by 15% for 7s.
func (c *char) c4cb(a combat.AttackCB) {
	if c.Base.Cons < 4 {
		return
	}

	e, ok := a.Target.(*enemy.Enemy)
	if !ok {
		return
	}
	e.AddDefMod(enemy.DefMod{Base: modifier.NewBase("razor-c4", 7*60), Value: -0.15})
}

// Every 10s, Razor's sword charges up, causing the next Normal Attack to release lightning that deals 100% of Razor's ATK as Electro DMG.
// When Razor is not using Lightning Fang, a lightning strike on an opponent will grant Razor an Electro Sigil for Claw and Thunder.
func (c *char) c6() {
	if c.Base.Cons < 6 {
		return
	}

	dur := 0
	c.Core.Events.Subscribe(event.OnDamage, func(args ...interface{}) bool {
		if c.Core.Player.Active() != c.Index {
			return false
		}
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.AttackTag != combat.AttackTagNormal {
			return false
		}
		if dur > c.Core.F {
			return false
		}

		dur = c.Core.F + 10*60
		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Lupus Fulguris",
			AttackTag:  combat.AttackTagNormal, // or combat.AttackTagNone?
			ICDTag:     combat.ICDTagNormalAttack,
			ICDGroup:   combat.ICDGroupDefault,
			Element:    attributes.Electro,
			Durability: 25,
			Mult:       1,
		}
		c.Core.QueueAttack(
			ai,
			combat.NewDefCircHit(0.5, false, combat.TargettableEnemy),
			1,
			1,
		)

		if c.Core.Status.Duration("razorburst") == 0 {
			c.addSigil()
		}

		return false
	}, "razor-c6")
}

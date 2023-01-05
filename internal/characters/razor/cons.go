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
	c.c1bonus = make([]float64, attributes.EndStatType)
	c.c1bonus[attributes.DmgP] = 0.1

	c.Core.Events.Subscribe(event.OnParticleReceived, func(_ ...interface{}) bool {
		// ignore if character not on field
		if c.Core.Player.Active() != c.Index {
			return false
		}
		c.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag("razor-c1", 8*60),
			AffectedStat: attributes.DmgP,
			Amount: func() ([]float64, bool) {
				return c.c1bonus, true
			},
		})
		return false
	}, "razor-c1")
}

// Increases CRIT Rate against opponents with less than 30% HP by 10%.
func (c *char) c2() {
	if c.Core.Combat.DamageMode {
		c.c2bonus = make([]float64, attributes.EndStatType)
		c.c2bonus[attributes.CR] = 0.1

		c.AddAttackMod(character.AttackMod{
			Base: modifier.NewBase("razor-c2", -1),
			Amount: func(_ *combat.AttackEvent, t combat.Target) ([]float64, bool) {
				x, ok := t.(*enemy.Enemy)
				if !ok {
					return nil, false
				}
				if x.HP()/x.MaxHP() < 0.3 {
					return c.c2bonus, true
				}
				return nil, false
			},
		})
	}
}

// When casting Claw and Thunder (Press), opponents hit will have their DEF decreased by 15% for 7s.
func (c *char) c4cb(a combat.AttackCB) {
	e, ok := a.Target.(*enemy.Enemy)
	if !ok {
		return
	}
	e.AddDefMod(enemy.DefMod{
		Base:  modifier.NewBaseWithHitlag("razor-c4", 7*60),
		Value: -0.15,
	})
}

const c6ICDKey = "razor-c6-icd"

// Every 10s, Razor's sword charges up, causing the next Normal Attack to release lightning that deals 100% of Razor's ATK as Electro DMG.
// When Razor is not using Lightning Fang, a lightning strike on an opponent will grant Razor an Electro Sigil for Claw and Thunder.
func (c *char) c6cb(a combat.AttackCB) {
	if a.Target.Type() != combat.TargettableEnemy {
		return
	}
	//effect can only happen every 10s
	if c.StatusIsActive(c6ICDKey) {
		return
	}

	c.AddStatus(c6ICDKey, 600, true)

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Lupus Fulguris",
		AttackTag:  combat.AttackTagNone, // TODO: it has another tag?
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeDefault,
		Element:    attributes.Electro,
		Durability: 25,
		Mult:       1,
	}

	sigilcb := func(a combat.AttackCB) {
		//add sigil only outside burst
		if c.StatusIsActive(burstBuffKey) {
			return
		}
		c.addSigil(false)(a)
	}

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHit(c.Core.Combat.Player(), a.Target, combat.Point{Y: 0.7}, 1.5),
		1,
		1,
		sigilcb,
	)

}

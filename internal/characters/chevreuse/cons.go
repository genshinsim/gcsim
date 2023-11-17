package chevreuse

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	buffStack1Key = "chev-c6-1"
	buffStack2Key = "chev-c6-2"
	buffStack3Key = "chev-c6-3"
)

func (c *char) c1() {

	if c.Base.Cons < 1 {
		return
	}

	if !c.onlyPyroElectro {
		return
	}

	c.Core.Events.Subscribe(event.OnOverload, c.OnC1Overload, "chev-c1")

}

func (c *char) c6TeamHeal() func() {

	return func() {
		if c.Base.Cons < 6 {
			return
		}

		for _, char := range c.Core.Player.Chars() {
			c.c6(char)
		}

		c.Core.Player.Heal(player.HealInfo{
			Caller:  c.Index,
			Target:  -1,
			Message: "In Pursuit of Ending Evil (C6)",
			Src:     0.1 * c.MaxHP(),
			Bonus:   c.Stat(attributes.Heal),
		})
	}
}

func (c *char) c6(char *character.CharWrapper) {
	if c.Base.Cons < 6 {
		return
	}

	m := make([]float64, attributes.EndStatType)
	m[attributes.PyroP] = 0.20
	m[attributes.ElectroP] = 0.20

	buffDuration := 8 * 60

	// 3 stackable buffs with independant timer. Temp hack.
	if !char.StatModIsActive(buffStack1Key) {
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag(buffStack1Key, buffDuration),
			AffectedStat: attributes.NoStat,
			Extra:        true,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})
		return
	}

	if !char.StatModIsActive(buffStack2Key) {
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag(buffStack2Key, buffDuration),
			AffectedStat: attributes.NoStat,
			Extra:        true,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})
		return
	}

	if !char.StatModIsActive(buffStack3Key) {
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag(buffStack3Key, buffDuration),
			AffectedStat: attributes.NoStat,
			Extra:        true,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})
		return
	}

	char.AddStatMod(character.StatMod{
		Base:         modifier.NewBaseWithHitlag(buffStack1Key, buffDuration),
		AffectedStat: attributes.NoStat,
		Extra:        true,
		Amount: func() ([]float64, bool) {
			return m, true
		},
	})

}

func (c *char) C2() combat.AttackCBFunc {

	return func(a combat.AttackCB) {
		if c.Base.Cons < 2 {
			return
		}

		if c.Core.F > c.c2Icd {
			c.c2Icd = c.Core.F + (60 * 10) // 10s icd
			ai := combat.AttackInfo{
				ActorIndex: c.Index,
				Abil:       "Sniper Induced Explosion (C2)",
				AttackTag:  attacks.AttackTagElementalArt,
				ICDTag:     attacks.ICDTagElementalArt,
				ICDGroup:   attacks.ICDGroupDefault, // TBC
				StrikeType: attacks.StrikeTypeDefault,
				Element:    attributes.Pyro,
				Durability: 25,
				Mult:       1.2,
			}

			// c2 1st hit
			c.Core.QueueAttack(
				ai,
				combat.NewCircleHit(c.Core.Combat.Player(), c.Core.Combat.PrimaryTarget(), nil, 1),
				skillHoldHitmark,
				skillHoldHitmark+skillHoldTravel,
			)

			// c2 2nd hit
			c.Core.QueueAttack(
				ai,
				combat.NewCircleHit(c.Core.Combat.Player(), c.Core.Combat.PrimaryTarget(), nil, 1),
				skillHoldHitmark,
				skillHoldHitmark+skillHoldTravel,
			)
		}
	}
}

func (c *char) OnC1Overload(args ...interface{}) bool {

	if c.Core.F > c.c1Icd {
		c.c1Icd = c.Core.F + (60 * 10)
		atk := args[1].(*combat.AttackEvent)
		atkCharIndex := atk.Info.ActorIndex

		// chev can't trigger her own c1
		if atk.Info.ActorIndex != c.Index {
			olTriggerChar := c.Core.Player.Chars()[atkCharIndex]
			olTriggerChar.AddEnergy("chev-c1", 6)
		}
	}
	return false
}

func (c *char) ApplyBuff() {

}

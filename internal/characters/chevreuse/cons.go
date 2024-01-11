package chevreuse

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	c4StatusKey = "chev-c4"
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

		if c.Core.F > c.c6Icd {
			c.c6Icd = c.Core.F + (60 * 15) // 15s. will only activate once regardless of c4
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
}

func (c *char) c4() {
	if c.Base.Cons < 4 {
		return
	}
	c.AddStatus(c4StatusKey, 6*60, false)
	c.c4ShotsLeft = 2
}

func (c *char) c6(char *character.CharWrapper) {
	if c.Base.Cons < 6 {
		return
	}

	m := make([]float64, attributes.EndStatType)
	m[attributes.PyroP] = 0.20
	m[attributes.ElectroP] = 0.20

	addStackMod := func(idx int, duration int) {
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag(fmt.Sprintf("chev-c6-%v-stack", idx+1), duration),
			AffectedStat: attributes.NoStat,
			Extra:        true,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})
	}

	buffDuration := 8 * 60
	addStackMod(c.c6Stack, buffDuration)
	c.c6Stack = (c.c6Stack + 1) % 3

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
				ICDGroup:   attacks.ICDGroupDefault,
				StrikeType: attacks.StrikeTypeBlunt,
				Element:    attributes.Pyro,
				Durability: 25,
				Mult:       1.2,
			}

			skillPos := c.Core.Combat.PrimaryTarget().Pos()
			// c2 1st hit
			c.Core.QueueAttack(
				ai,
				combat.NewCircleHitOnTarget(skillPos, nil, 4.5),
				skillHoldHitmark,
				skillHoldHitmark+skillHoldTravel,
			)

			// c2 2nd hit
			c.Core.QueueAttack(
				ai,
				combat.NewCircleHitOnTarget(skillPos, nil, 4.5),
				skillHoldHitmark,
				skillHoldHitmark+skillHoldTravel,
			)
		}
	}
}

func (c *char) OnC1Overload(args ...interface{}) bool {

	if c.Core.F > c.c1Icd {
		atk := args[1].(*combat.AttackEvent)
		atkCharIndex := atk.Info.ActorIndex

		// chev can't trigger her own c1
		if atk.Info.ActorIndex != c.Index && atk.Info.ActorIndex == c.Core.Player.Active() {
			c.c1Icd = c.Core.F + (60 * 10)
			olTriggerChar := c.Core.Player.Chars()[atkCharIndex]
			olTriggerChar.AddEnergy("chev-c1", 6)
		}
	}
	return false
}

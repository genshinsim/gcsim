package dahlia

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/gmod"
)

const (
	onFrozenReaction          = "dahlia-a1-onfrozen"
	benisonStackGenerationIcd = "dahlia-a1-benison-stack-icd"
	attackSpeedKey            = "dahlia-a4-atk-speed"
)

// A1
// When characters affected by the Favonian Favor effect of the Elemental Burst Radiant Psalter trigger Frozen
// reactions on opponents, they will grant Dahlia 2 Benison stacks. This effect can trigger once every 8s.
func (c *char) a1() {
	if c.Base.Ascension < 1 {
		return
	}

	c.Core.Events.Subscribe(event.OnFrozen, func(args ...any) bool {
		if !c.StatusIsActive(burstFavonianFavor) {
			return false
		}

		_, ok := args[0].(*enemy.Enemy)
		if !ok {
			return false
		}

		ae := args[1].(*info.AttackEvent)
		if !c.StatusIsActive(benisonStackGenerationIcd) {
			c.AddStatus(benisonStackGenerationIcd, 8*60, true)
			c.addBenisonStack(2, ae.Info.ActorIndex)
		}

		return false
	}, onFrozenReaction)
}

// A4
// When your active party member is affected by the Favonian Favor effect of the Elemental Burst their ATK SPD will
// increase based on Dahlia's Max HP: Every 1,000 Max HP will cause an increase of 0.5%, up to a maximum of 20%.
func (c *char) a4() {
	if c.Base.Ascension < 4 {
		return
	}

	c.attackSpeedBuff = make([]float64, attributes.EndStatType)

	// Swaps update Attack Speed buff and provide a new 0.5s timer for the on-fielder
	c.Core.Events.Subscribe(event.OnCharacterSwap, func(args ...any) bool {
		if !c.StatusIsActive(burstFavonianFavor) {
			return false
		}

		prev, next := args[0].(int), args[1].(int)

		// Remove Attack Speed buff from swapped out character and give it to swapped in character
		for _, char := range c.Core.Player.Chars() {
			if char.Index() == prev && char.StatusIsActive(attackSpeedKey) {
				char.DeleteStatMod(attackSpeedKey)
			}
			if char.Index() == next && !char.StatusIsActive(attackSpeedKey) {
				c.addAttackSpeedbuff(char)
			}
		}

		c.updateSpeedBuff(next)() // ID of swapped in character

		return false
	}, attackSpeedKey)
}

func (c *char) addAttackSpeedbuff(char *character.CharWrapper) {
	char.AddStatMod(character.StatMod{
		Base:         gmod.NewBase(attackSpeedKey, -1),
		AffectedStat: attributes.AtkSpd,
		Amount: func() ([]float64, bool) {
			// No Attack Speed buff if Favonian Favor from Dahlia's Burst is not active
			if !c.StatusIsActive(burstFavonianFavor) {
				return nil, false
			}
			// No Attack Speed for Charged Attacks
			if c.Core.Player.CurrentState() != action.NormalAttackState {
				return nil, false
			}
			return c.attackSpeedBuff, true
		},
	})
}

func (c *char) updateSpeedBuff(src int) func() {
	return func() {
		// If no longer on-field, stop updating the Attack Speed buff based on that character's hitlag
		if c.Core.Player.Active() != src {
			return
		}

		// Calcuate Attack Speed buff (max 20%)
		// If C6, add another 10% Attack Speed
		burstAttackSpeed := min(c.MaxHP()*0.001*0.005, 0.2) + c.c6AtkSpd()
		c.attackSpeedBuff[attributes.AtkSpd] = burstAttackSpeed

		// Update Attack Speed buff after 0.5s (affected by on-fielder hitlag)
		c.Core.Player.ByIndex(src).QueueCharTask(c.updateSpeedBuff(src), 0.5*60)
	}
}

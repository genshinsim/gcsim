package dahlia

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
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

	c.Core.Events.Subscribe(event.OnFrozen, c.a1Hook, onFrozenReaction)
}

func (c *char) a1Hook(args ...any) bool {
	char := c.Core.Player.ActiveChar()
	if !char.StatusIsActive(burstFavonianFavor) {
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
}

// A4
// When your active party member is affected by the Favonian Favor effect of the Elemental Burst their ATK SPD will
// increase based on Dahlia's Max HP: Every 1,000 Max HP will cause an increase of 0.5%, up to a maximum of 20%.
func (c *char) a4() {
	if c.Base.Ascension < 4 {
		return
	}

	c.Core.Events.Subscribe(event.OnCharacterSwap, c.a4Hook, attackSpeedKey)
}

func (c *char) a4Hook(args ...any) bool {
	prev, next := args[0].(int), args[1].(int)

	if !c.StatusIsActive(burstFavonianFavor) {
		return false
	}

	// Remove buff from swapped out character and give it to swapped in character
	for _, char := range c.Core.Player.Chars() {
		if char.Index() == prev && char.StatusIsActive(attackSpeedKey) {
			char.DeleteStatus(attackSpeedKey)
		}
		if char.Index() == next && !char.StatusIsActive(attackSpeedKey) {
			c.addAttackSpeedbuff(char)
		}
	}

	return false
}

func (c *char) addAttackSpeedbuff(char *character.CharWrapper) {
	attackSpeedInc := 0.005 * (c.MaxHP() / 1000)
	if attackSpeedInc > 0.2 {
		attackSpeedInc = 0.2
	}
	// If C6, add another 10% Attack Speed
	if c.Base.Cons >= 6 {
		attackSpeedInc += 0.1
	}
	attackSpeedBuff := make([]float64, attributes.EndStatType)
	attackSpeedBuff[attributes.AtkSpd] = attackSpeedInc

	char.AddStatMod(character.StatMod{
		Base:         modifier.NewBaseWithHitlag(attackSpeedKey, favonianFavorExpiry-c.Core.F),
		AffectedStat: attributes.AtkSpd,
		Amount: func() ([]float64, bool) {
			return attackSpeedBuff, true
		},
	})
}

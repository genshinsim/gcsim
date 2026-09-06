package odette

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const a1Key = "odette-a1"

// When Odette summons her Solo Dance Double, she also obtains 4 stacks of Marvelous Splendor.
// Marvelous Splendor
// After unlocking the Ascension Talent "Spring Rite of the Chosen One," Odette will obtain a
// special effect when she summons her Solo Dance Double.
// Every stack of Marvelous Splendor active increases the character's Stellar Glimmer DMG by 15%.
// This lasts until her Dance Double exits the field or when she summons her Solo Dance Double again.
// When Odette is off-field, she loses 1 stack of Marvelous Splendor every 1 second, while other
// nearby party members gain the corresponding number of Marvelous Splendor stacks at the same time.
func (c *char) a1Init() {
	if c.Base.Ascension < 1 {
		return
	}

	c.AddReactBonusMod(character.ReactBonusMod{
		Base: modifier.NewBase(a1Key+"-buff", -1),
		Amount: func(ai info.AttackInfo) float64 {
			switch ai.AttackTag {
			case attacks.AttackTagDirectStellarConduct,
				attacks.AttackTagDirectStellarSwirl,
				attacks.AttackTagReactionStellarSwirl:
			default:
				return 0
			}
			if !c.StatusIsActive(danceDoubleKey) {
				return 0
			}
			return float64(c.a1StacksSelf) * 0.15
		},
	})

	for _, char := range c.Core.Player.Chars() {
		if char.Index() == c.Index() {
			continue
		}
		char.AddReactBonusMod(character.ReactBonusMod{
			Base: modifier.NewBase(a1Key+"-buff", -1),
			Amount: func(ai info.AttackInfo) float64 {
				switch ai.AttackTag {
				case
					attacks.AttackTagDirectStellarConduct,
					attacks.AttackTagDirectStellarSwirl,
					attacks.AttackTagReactionStellarSwirl:
				default:
					return 0
				}
				if !c.StatusIsActive(danceDoubleKey) {
					return 0
				}
				return float64(c.a1StacksOthers) * 0.15
			},
		})
	}

	c.Core.Events.Subscribe(event.OnCharacterSwap, func(args ...any) {
		prev := args[0].(int)
		next := args[1].(int)
		if prev == c.Index() {
			src := c.Core.F
			c.a1Src = src
			c.Core.Tasks.Add(func() { c.a1Ticker(src) }, 60)
		} else if next == c.Index() {
			// cancel the a1Ticker
			c.a1Src = -1
		}
	}, a1Key)
}

func (c *char) a1Ticker(src int) {
	// don't need to check asc because it's only called by a1Init()
	if c.a1Src != src {
		return
	}

	if c.a1StacksSelf == 0 {
		return
	}

	if c.a1StacksOthers == c.a1MaxStacks() {
		return
	}

	// TODO: This check isn't needed because we cancel the task when we swap back to Odette
	if c.Core.Player.Active() == c.Index() {
		return
	}

	if !c.StatusIsActive(danceDoubleKey) {
		return
	}

	stacks := min(c.c1a1Remove(), c.a1StacksSelf)
	c.a1StacksSelf -= stacks * c.c6a1ReduceMod()
	c.a1StacksOthers = min(c.a1StacksOthers+stacks, c.a1MaxStacks())

	c.Core.Tasks.Add(func() { c.a1Ticker(src) }, 60)
}

func (c *char) a1OnDanceSummon() {
	if c.Base.Ascension < 1 {
		return
	}

	c.a1StacksSelf = c.a1MaxStacks()
	c.a1StacksOthers = 0
}

func (c *char) a1MaxStacks() int {
	return 4 + c.c1a1Stacks()
}

// For every 100 ATK Odette has over 1,000, her Stellar Glimmer DMG is additionally increased by
// 1.5% of the original DMG. She can deal up to 30% more additional DMG in this way.
// Does not affect reaction SSW
func (c *char) a4StellarGlimmerMult() float64 {
	if c.Base.Ascension < 1 {
		return 1
	}

	scaling := max(c.TotalAtk()-1000, 0)
	buff := min(scaling/100*0.015, 0.3)

	return 1.0 + buff
}

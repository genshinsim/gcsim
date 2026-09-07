package shenhe

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

// A4:
// After Shenhe uses Spring Spirit Summoning, she will grant all nearby party members the following effects:
//
// - Press: Elemental Skill and Elemental Burst DMG increased by 15% for 10s.
func (c *char) a4PressBuff() {
	if c.Base.Ascension < 4 {
		return
	}
	for _, other := range c.Core.Player.Chars() {
		other.AddAttackMod(character.AttackMod{
			Base: modifier.NewBaseWithHitlag("shenhe-a4-press", 10*60),
			Amount: func(a *info.AttackEvent, _ info.Target) []float64 {
				switch a.Info.AttackTag {
				case attacks.AttackTagElementalArt:
				case attacks.AttackTagElementalArtHold:
				case attacks.AttackTagElementalBurst:
				default:
					return nil
				}
				return c.skillBuff
			},
		})
	}
}

// A4:
// After Shenhe uses Spring Spirit Summoning, she will grant all nearby party members the following effects:
//
// - Hold: Normal, Charged, and Plunging Attack DMG increased by 15% for 15s.
func (c *char) a4HoldBuff() {
	if c.Base.Ascension < 4 {
		return
	}
	for _, other := range c.Core.Player.Chars() {
		other.AddAttackMod(character.AttackMod{
			Base: modifier.NewBaseWithHitlag("shenhe-a4-hold", 15*60),
			Amount: func(a *info.AttackEvent, _ info.Target) []float64 {
				switch a.Info.AttackTag {
				case attacks.AttackTagNormal:
				case attacks.AttackTagExtra:
				case attacks.AttackTagPlunge:
				default:
					return nil
				}
				return c.skillBuff
			},
		})
	}
}

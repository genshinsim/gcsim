package shenhe

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const c4BuffKey = "shenhe-c4"

func (c *char) c2(active *character.CharWrapper, dur int) {
	active.AddAttackMod(character.AttackMod{
		Base: modifier.NewBaseWithHitlag("shenhe-c2", dur),
		Amount: func(ae *combat.AttackEvent, _ combat.Target) ([]float64, bool) {
			if ae.Info.Element != attributes.Cryo {
				return nil, false
			}
			return c.c2buff, true
		},
	})
}

// When characters under the effect of Icy Quill applied by Shenhe trigger its DMG Bonus effects, Shenhe will gain a Skyfrost Mantra stack:
//
// - When Shenhe uses Spring Spirit Summoning, she will consume all stacks of Skyfrost Mantra, increasing the DMG of that Spring Spirit Summoning by 5% for each stack consumed.
//
// - Max 50 stacks. Stacks last for 60s.
func (c *char) c4Init() {
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("shenhe-c4-dmg", -1),
		Amount: func(atk *combat.AttackEvent, _ combat.Target) ([]float64, bool) {
			if atk.Info.AttackTag != combat.AttackTagElementalArt && atk.Info.AttackTag != combat.AttackTagElementalArtHold {
				return nil, false
			}
			if !c.StatusIsActive(c4BuffKey) {
				c.c4count = 0
				return nil, false
			}
			c.c4bonus[attributes.DmgP] = 0.05 * float64(c.c4count)
			c.c4count = 0
			return c.c4bonus, true
		},
	})
}

// C4 stacks are gained after the damage has been dealt and not before
// https://library.keqingmains.com/evidence/characters/cryo/shenhe?q=shenhe#c4-insight
func (c *char) c4CB(a combat.AttackCB) {
	//reset stacks to zero if all expired
	if !c.StatusIsActive(c4BuffKey) {
		c.c4count = 0
	}
	if c.c4count < 50 {
		c.c4count++
		c.Core.Log.NewEvent("shenhe-c4 stack gained", glog.LogCharacterEvent, c.Index).
			Write("stacks", c.c4count)
	}
	c.AddStatus(c4BuffKey, 3600, true) // 60 s
}

func (c *char) makeC4ResetCB() combat.AttackCBFunc {
	if c.Base.Cons < 4 {
		return nil
	}
	return func(a combat.AttackCB) {
		if c.Core.Player.Active() != c.Index {
			return
		}
		c.DeleteStatus(c4BuffKey)
	}
}

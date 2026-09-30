package shenhe

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	c2BuffKey = "shenhe-c2"
	c4BuffKey = "shenhe-c4"
)

func (c *char) c2(active *character.CharWrapper, dur int) {
	active.AddAttackMod(character.AttackMod{
		Base: modifier.NewBaseWithHitlag(c2BuffKey, dur),
		Amount: func(ae *info.AttackEvent, _ info.Target) []float64 {
			if ae.Info.Element != attributes.Cryo {
				return nil
			}
			return c.c2buff
		},
	})
}

func (c *char) c2Init() {
	if c.Base.Cons < 2 {
		return
	}

	c.c2buff = make([]float64, attributes.EndStatType)
	c.c2buff[attributes.CD] = 0.15

	c.Core.Events.Subscribe(event.OnSpecialReactionAttack, func(args ...any) {
		atk, ok := args[1].(*info.AttackEvent)
		if !ok {
			return
		}

		if !atk.Info.AttackTag.IsStellarReact() {
			return
		}

		char := c.Core.Player.Chars()[atk.Info.ActorIndex]

		if !char.StatModIsActive(c2BuffKey) {
			return
		}

		atk.Snapshot.Stats[attributes.CD] += c.c2buff[attributes.CD]
	}, "shenhe-c2-on-stellar")
}

// When characters under the effect of Icy Quill applied by Shenhe trigger its DMG Bonus effects, Shenhe will gain a Skyfrost Mantra stack:
//
// - When Shenhe uses Spring Spirit Summoning, she will consume all stacks of Skyfrost Mantra, increasing the DMG of that Spring Spirit Summoning by 5% for each stack consumed.
//
// - Max 50 stacks. Stacks last for 60s.
func (c *char) c4() float64 {
	if c.Base.Cons < 4 {
		return 0
	}
	if !c.StatusIsActive(c4BuffKey) {
		c.c4count = 0
		return 0
	}
	dmgBonus := 0.05 * float64(c.c4count)
	c.Core.Log.NewEvent("shenhe-c4 adding dmg bonus", glog.LogCharacterEvent, c.Index()).
		Write("stacks", c.c4count).
		Write("dmg_bonus", dmgBonus)
	c.c4count = 0
	c.DeleteStatus(c4BuffKey)
	return dmgBonus
}

// C4 stacks are gained after the damage has been dealt and not before
// https://library.keqingmains.com/evidence/characters/cryo/shenhe?q=shenhe#c4-insight
func (c *char) c4CB(a info.AttackCB) {
	// reset stacks to zero if all expired
	if !c.StatusIsActive(c4BuffKey) {
		c.c4count = 0
	}
	if c.c4count < 50 {
		c.c4count++
		c.Core.Log.NewEvent("shenhe-c4 stack gained", glog.LogCharacterEvent, c.Index()).
			Write("stacks", c.c4count)
	}
	c.AddStatus(c4BuffKey, 3600, true) // 60 s
}

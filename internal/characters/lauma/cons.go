package lauma

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	c1Key     = "lauma-threads-of-life"
	c1HitMark = 5
)

func (c *char) c1() {
	// on lb proc heal
	c.Core.Events.Subscribe(event.OnLunarBloom, func(args ...any) bool {
		if !c.StatusIsActive(c1Key) {
			return false
		}
		_, ok := args[0].(*enemy.Enemy)
		if !ok {
			return false
		}
		healAmt := 5.0 * c.Stat(attributes.EM)

		// heal active character
		c.Core.Tasks.Add(func() {
			c.Core.Player.Heal(info.HealInfo{
				Type:    info.HealTypeAbsolute,
				Message: "Lauma C1 (Heal)",
				Src:     healAmt,
			})
		}, c1HitMark)

		c.Core.Tasks.Add(func() {
			c.c1()
		}, 19*6)

		return true
	}, "lauma-c1")
}

func (c *char) c2() {
	// if !c.ascendantGleam {
	// 	return
	// }

	bonus := 0.4

	for _, x := range c.Core.Player.Chars() {
		this := x
		// special case if applying 4 Noblesse to holder to fix this mess:
		// https://library.keqingmains.com/evidence/general-mechanics/bugs#noblesse-oblige-4pc-bonus-not-applying-to-some-bursts
		// https://docs.google.com/spreadsheets/d/1jhIP3C6B16nL1unX9DL_-LhSNaOy_wwhdr29pzikpcg/edit?usp=sharing
		// TODO: Does the char snapshot 4 Noblesse if 4 Noblesse is already up and they're refreshing the duration? (rn they would snapshot it)

		this.AddReactBonusMod(character.ReactBonusMod{
			Base: modifier.NewBase("lauma-c2-lunarbloom-buff", -1),
			Amount: func(atk info.AttackInfo) (float64, bool) {
				if atk.AttackTag != attacks.AttackTagDirectLunarBloom {
					return 0, false
				}
				return bonus, false
			},
		})
	}
}

func (c *char) c4Refund(a info.AttackCB) {
	if c.Base.Cons < 4 {
		return
	}
	if c.StatusIsActive(laumaC4RefundKey) {
		return
	}
	c.AddEnergy("lauma-c4-refund", 5)
	c.AddStatus(laumaC4RefundKey, 5*60, true)
}

func (c *char) c6Elevation() {
	if !c.ascendantGleam {
		return
	}

	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...any) bool {
		atk := args[1].(*info.AttackEvent)

		switch atk.Info.AttackTag {
		case attacks.AttackTagDirectLunarBloom:
		default:
			return false
		}

		bonus := 0.25

		if c.Core.Flags.LogDebug {
			c.Core.Log.NewEvent("lauma-c6-elevation", glog.LogCharacterEvent, c.Index()).Write("bonus", bonus)
		}

		atk.Info.Elevation += bonus
		return false
	}, lunarbloomBonusKey)
}

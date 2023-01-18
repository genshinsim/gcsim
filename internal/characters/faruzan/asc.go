package faruzan

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

const (
	a4Key    = "faruzan-a4"
	a4ICDKey = "faruzan-a4-icd"
)

// When characters affected by The Wind's Secret Ways' Prayerful Wind's Gift
// deal Anemo DMG using Normal, Charged, Plunging Attacks, Elemental Skills, or
// Elemental Bursts to opponents, they will gain the Hurricane Guard effect:
// This DMG will be increased based on 32% of Faruzan's Base ATK. 1 instance of
// Hurricane Guard can occur once every 0.8s. This DMG Bonus will be cleared
// after Prayerful Wind's Benefit expires or after the effect is triggered
// once.
func (c *char) a4() {
	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.Element != attributes.Anemo {
			return false
		}

		switch atk.Info.AttackTag {
		case combat.AttackTagNormal,
			combat.AttackTagExtra,
			combat.AttackTagPlunge,
			combat.AttackTagElementalArt,
			combat.AttackTagElementalArtHold,
			combat.AttackTagElementalBurst:
		default:
			return false
		}

		active := c.Core.Player.ByIndex(atk.Info.ActorIndex)
		if active.StatusIsActive(burstBuffKey) && !c.StatusIsActive(a4ICDKey) {
			amt := 0.32 * (c.Base.Atk + c.Weapon.Atk)
			if c.Core.Flags.LogDebug {
				c.Core.Log.NewEvent("faruzan a4 proc dmg add", glog.LogPreDamageMod, atk.Info.ActorIndex).
					Write("before", atk.Info.FlatDmg).
					Write("addition", amt)
			}
			atk.Info.FlatDmg += amt
			c.AddStatus(a4ICDKey, 48, false)
		}

		return false
	}, "faruzan-a4-hook")
}

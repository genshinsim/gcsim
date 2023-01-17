package travelergeo

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/construct"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

// C1:
// Party members within the radius of Wake of Earth have their CRIT Rate increased by 10%
//
//	and have increased resistance against interruption.
func (c *char) c1(ticks int) func() {
	return func() {
		// different Q fields can co-exist at C6 if you do the following:
		// - cast another Q after Q cooldown is up (after 15s) but before Q field expires (before 20s)
		// so it's ok that they're both queueing up c1 ticks

		// if Q construct isn't up, then don't apply buff / queue another tick
		if c.Core.Constructs.CountByType(construct.GeoConstructTravellerBurst) == 0 {
			return
		}

		// this makes sure that every Q field only ticks 15/20 times at <C6/C6
		if ticks > c.c1TickCount {
			return
		}

		c.Core.Log.NewEvent("geo-traveler field ticking", glog.LogCharacterEvent, -1).
			Write("tick_number", ticks)

		// apply C1 buff to active char for 2s
		if combat.TargetIsWithinArea(c.Core.Combat.Player().Pos(), c.burstArea) {
			m := make([]float64, attributes.EndStatType)
			m[attributes.CR] = .1

			active := c.Core.Player.ActiveChar()
			active.AddStatMod(character.StatMod{
				Base:         modifier.NewBaseWithHitlag("geo-traveler-c1", 120), // 2s
				AffectedStat: attributes.CR,
				Amount: func() ([]float64, bool) {
					return m, true
				},
			})
		}

		// check again in 1s
		ticks += 1
		c.Core.Tasks.Add(c.c1(ticks), 60)
	}
}

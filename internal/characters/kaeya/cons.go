package kaeya

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/shield"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

// C1:
// The CRIT Rate of Kaeya's Normal and Charge Attacks against opponents affected by Cryo is increased by 15%.
func (c *char) c1() {
	m := make([]float64, attributes.EndStatType)
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("kaeya-c1", -1),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			e, ok := t.(*enemy.Enemy)
			if !ok {
				return nil, false
			}
			if atk.Info.AttackTag != attacks.AttackTagNormal && atk.Info.AttackTag != attacks.AttackTagExtra {
				return nil, false
			}
			if !e.AuraContains(attributes.Cryo, attributes.Frozen) {
				return nil, false
			}
			m[attributes.CR] = 0.15
			return m, true
		},
	})
}

// C2:
// Every time Glacial Waltz defeats an opponent during its duration, its duration is increased by 2.5s, up to a maximum of 15s.
func (c *char) c2() {
	c.Core.Events.Subscribe(event.OnTargetDied, func(args ...interface{}) bool {
		_, ok := args[0].(*enemy.Enemy)
		// ignore if not an enemy
		if !ok {
			return false
		}
		// ignore if burst isn't up
		if c.Core.Status.Duration(burstKey) == 0 {
			return false
		}
		// ignore if extension limit reached
		if c.c2ProcCount > 2 {
			return false
		}
		// burst duration steps
		// 8s
		// 10.5s (+2.5s from previous)
		// 13s (+2.5s from previous)
		// 15s (+2.0s from previous because extension is capped to 15s)
		extension := 150
		if c.c2ProcCount == 2 {
			extension = 120
		}
		c.Core.Status.Extend(burstKey, extension)
		c.c2ProcCount++
		c.Core.Log.NewEvent("kaeya-c2 proc'd", glog.LogCharacterEvent, c.Index).
			Write("c2ProcCount", c.c2ProcCount).
			Write("extension", extension)
		return false
	}, "kaeya-c2")
}

// C4:
// Triggers automatically when Kaeya's HP falls below 20%:
// Creates a shield that absorbs damage equal to 30% of Kaeya's Max HP. Lasts for 20s.
// This shield absorbs Cryo DMG with 250% efficiency.
// Can only occur once every 60s.
func (c *char) c4() {
	c.Core.Events.Subscribe(event.OnPlayerHPDrain, func(args ...interface{}) bool {
		di := args[0].(player.DrainInfo)
		if di.Amount <= 0 {
			return false
		}
		if c.Core.F < c.c4icd && c.c4icd != 0 {
			return false
		}
		maxhp := c.MaxHP()
		if c.CurrentHPRatio() < 0.2 {
			c.c4icd = c.Core.F + 3600
			c.Core.Player.Shields.Add(&shield.Tmpl{
				ActorIndex: c.Index,
				Src:        c.Core.F,
				ShieldType: shield.ShieldKaeyaC4,
				Name:       "Kaeya C4",
				HP:         .3 * maxhp,
				Ele:        attributes.Cryo,
				Expires:    c.Core.F + 1200,
			})
		}
		return false
	}, "kaeya-c4")
}

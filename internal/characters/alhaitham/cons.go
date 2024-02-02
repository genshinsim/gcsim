package alhaitham

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	c1IcdKey    = "alhaitham-c1-icd"
	c2MaxStacks = 4
)

// When a Projection Attack hits an opponent, Universality: An Elaboration on Form's CD is decreased by 1.2s.
// This effect can be triggered once every 1s.
func (c *char) c1(a combat.AttackCB) {
	// ignore if c1 on icd
	if c.StatusIsActive(c1IcdKey) {
		return
	}
	c.ReduceActionCooldown(action.ActionSkill, 72) // reduced by 1.2s
	c.AddStatus(c1IcdKey, 60, true)                // 1s icd affected by hitlag
}

// When Alhaitham generates a Chisel-Light Mirror, his Elemental Mastery will be increased by 50 for 8 seconds,
// max 4 stacks.
// Each stack's duration is counted independently.
// This effect can be triggered even when the maximum number of Chisel-Light Mirrors has been reached.
func (c *char) c2(generated int) {
	m := make([]float64, attributes.EndStatType)
	m[attributes.EM] = 50
	for i := 0; i < generated; i++ {
		c.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag(c2ModName(c.c2Counter+1), 480), // 8s
			AffectedStat: attributes.EM,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})
		c.c2Counter = (c.c2Counter + 1) % c2MaxStacks // stacks are independent from each other, this will cycle them
	}
}

func c2ModName(num int) string {
	return fmt.Sprintf("alhaitham-c2-%v-stack", num)
}

// When Particular Field: Fetters of Phenomena is unleashed, the following effects will become active
// based on the number of Chisel-Light Mirrors consumed and created this time around:
// ·Each Mirror consumed will increase the Elemental Mastery of all other nearby party members by 30 for 15s.
// ·Each Mirror generated will grant Alhaitham a 10% Dendro DMG Bonus for 15s.
// The pre-existing duration of the aforementioned effects will be cleared if you use Particular Field: Fetters of Phenomena again while they are in effect
func (c *char) c4Loss(consumed int) {
	if consumed <= 0 {
		return
	}
	m := make([]float64, attributes.EndStatType)
	m[attributes.EM] = 30.0 * float64(consumed)
	for i, char := range c.Core.Player.Chars() {
		// skip Alhaitham
		if i == c.Index {
			continue
		}
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag("alhaitham-c4-loss", 900),
			AffectedStat: attributes.EM,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})
	}
}

func (c *char) c4Gain(generated int) {
	if generated <= 0 {
		return
	}
	m := make([]float64, attributes.EndStatType)
	m[attributes.DendroP] = 0.1 * float64(generated)
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBaseWithHitlag("alhaitham-c4-gain", 900),
		AffectedStat: attributes.DendroP,
		Amount: func() ([]float64, bool) {
			return m, true
		},
	})
}

// Alhaitham gains the following effects:
// · 2 seconds after Particular Field: Fetters of Phenomena is unleashed,
// he will generate 3 Chisel-Light Mirrors regardless of the number of mirrors consumed.
//
// · If Alhaitham generates Chisel-Light Mirrors when their numbers have already maxed out,
// his CRIT Rate and CRIT DMG will increase by 10% and 70% respectively for 6s.
// If this effect is triggered again during its initial duration, the duration remaining will be increased by 6s.
const c6key = "alhaitham-c6"

func (c *char) c6(generated int) {
	m := make([]float64, attributes.EndStatType)
	m[attributes.CR] = 0.1
	m[attributes.CD] = 0.7
	for i := 0; i < generated; i++ {
		if c.StatModIsActive(c6key) {
			c.ExtendStatus(c6key, 360)
			c.Core.Log.NewEvent("c6 buff extended", glog.LogCharacterEvent, c.Index).Write("c6 expiry on", c.StatusExpiry(c6key))
		} else {
			c.AddStatMod(character.StatMod{
				Base:         modifier.NewBaseWithHitlag((c6key), 360), // 6s
				AffectedStat: attributes.CR,
				Amount: func() ([]float64, bool) {
					return m, true
				},
			})
		}
	}
}

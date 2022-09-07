package cyno

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

// After using Sacred Rite: Wolf's Swiftness, Cyno's Normal Attack SPD will be increased by 20% for 10s.
// If the Judication effect of his Passive Talent Featherfall Judgment is triggered during Secret Rite: Chasmic Soulfarer,
// the duration of this increase will be refreshed.
//
// You need to unlock the Passive Talent "Featherfall Judgment."
func (c *char) c1() {
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBaseWithHitlag("cyno-c1", 300), //5s
		AffectedStat: attributes.AtkSpd,
		Amount: func() ([]float64, bool) {
			return c.c1buff, true
		},
	})
	c.AddStatus(c1key, 600, true)

}

// When Cyno's Normal Attacks hit opponents, his Normal Attack CRIT Rate and CRIT DMG will be increased by 3% and 6% respectively for 4s.
// This effect can be triggered once every 0.1s.
// Max 5 stacks. Each stack's duration is counted independently.
func (c *char) c2() {
	const c2Icd = "cyno-c2-icd"
	c.Core.Events.Subscribe(event.OnAttackWillLand, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.ActorIndex != c.Index {
			return false
		}
		if c.StatusIsActive(c2Icd) {
			return false
		}
		if atk.Info.AttackTag != combat.AttackTagNormal {
			return false
		}
		m := make([]float64, attributes.EndStatType)
		m[attributes.CR] = 0.03
		m[attributes.CD] = 0.06

		c.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag(fmt.Sprintf("cyno-c2-%v-stack", c.c2counter+1), 240), //4s
			AffectedStat: attributes.CR,
			Amount: func() ([]float64, bool) {

				return m, true
			},
		})
		c.AddStatus(c2Icd, 6, true)         //0.1s icd
		c.c2counter = (c.c2counter + 1) % 5 //stacks are independent from each other, this will cycle them
		return false
	}, "cyno-c2")

}

// When Cyno is in the Pactsworn Pathclearer state triggered by Sacred Rite: Wolf's Swiftness,
// after he triggers Electro-Charged, Superconduct, Overloaded, Quicken, Aggravate, Hyperbloom, an Electro Swirl
// or an Electro Crystallization reaction, he will restore 3 Elemental Energy for all nearby party members (except himself.)
// This effect can occur 5 times within one use of Sacred Rite: Wolf's Swiftness.
func (c *char) c4() {

}

// After using Sacred Rite: Wolf's Swiftness or triggering the Judication effect of the Passive Talent "Featherfall Judgment,"
// Cyno will gain 4 stacks of the "Day of the Jackal" effect. When he hits opponents with Normal Attacks,
// he will consume 1 stack of "Day of the Jackal" to fire off one Duststalker Bolt.
// "Day of the Jackal" lasts for 8s. Max 8 stacks. It will be canceled once Pactsworn Pathclearer ends.
// A maximum of 1 Duststalker Bolt can be unleashed this way every 0.4s.
// You must first unlock the Passive Talent "Featherfall Judgment."
func (c *char) c6() {

}

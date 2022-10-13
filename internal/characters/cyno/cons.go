package cyno

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const c1Key = "cyno-c1"
const c6Key = "cyno-c6"

// After using Sacred Rite: Wolf's Swiftness, Cyno's Normal Attack SPD will be increased by 20% for 10s.
// If the Judication effect of his Passive Talent Featherfall Judgment is triggered during Secret Rite: Chasmic Soulfarer,
// the duration of this increase will be refreshed.
//
// You need to unlock the Passive Talent "Featherfall Judgment."
func (c *char) c1() {
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBaseWithHitlag(c1Key, 600), // 10s
		AffectedStat: attributes.AtkSpd,
		Amount: func() ([]float64, bool) {
			m := make([]float64, attributes.EndStatType)
			m[attributes.AtkSpd] = 0.2
			return m, true
		},
	})
}

// When Cyno's Normal Attacks hit opponents, his Electro DMG Bonus will
// increase by 10% for 4s. This effect can be triggered once every 0.1s. Max 5
// stacks.
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
		m[attributes.ElectroP] = 0.1

		c.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag(fmt.Sprintf("cyno-c2-%v-stack", c.c2counter+1), 240), // 4s
			AffectedStat: attributes.CR,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})
		c.AddStatus(c2Icd, 6, true)         // 0.1s icd
		c.c2counter = (c.c2counter + 1) % 5 // stacks are independent from each other, this will cycle them
		return false
	}, "cyno-c2")
}

// When Cyno is in the Pactsworn Pathclearer state triggered by Sacred Rite: Wolf's Swiftness,
// after he triggers Electro-Charged, Superconduct, Overloaded, Quicken, Aggravate, Hyperbloom, an Electro Swirl
// or an Electro Crystallization reaction, he will restore 3 Elemental Energy for all nearby party members (except himself.)
// This effect can occur 5 times within one use of Sacred Rite: Wolf's Swiftness.
func (c *char) c4() {
	restore := func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.ActorIndex != c.Index {
			return false
		}
		if c.c4counter > 4 { // counting from 0 to 4, 5 instances max
			return false
		}
		c.c4counter++
		for _, this := range c.Core.Player.Chars() {
			// not for cyno
			if this.Index != c.Index {
				this.AddEnergy("cyno-c4", 3)
			}
		}

		return false
	}
	c.Core.Events.Subscribe(event.OnOverload, restore, "cyno-c4")
	c.Core.Events.Subscribe(event.OnElectroCharged, restore, "cyno-c4")
	c.Core.Events.Subscribe(event.OnSuperconduct, restore, "cyno-c4")
	c.Core.Events.Subscribe(event.OnQuicken, restore, "cyno-c4")
	c.Core.Events.Subscribe(event.OnAggravate, restore, "cyno-c4")
	c.Core.Events.Subscribe(event.OnHyperbloom, restore, "cyno-c4")
	c.Core.Events.Subscribe(event.OnSwirlElectro, restore, "cyno-c4")
}

// After using Sacred Rite: Wolf's Swiftness or triggering the Judication effect of the Passive Talent "Featherfall Judgment,"
// Cyno will gain 4 stacks of the "Day of the Jackal" effect. When he hits opponents with Normal Attacks,
// he will consume 1 stack of "Day of the Jackal" to fire off one Duststalker Bolt.
// "Day of the Jackal" lasts for 8s. Max 8 stacks. It will be canceled once Pactsworn Pathclearer ends.
// A maximum of 1 Duststalker Bolt can be unleashed this way every 0.4s.
// You must first unlock the Passive Talent "Featherfall Judgment."
func (c *char) c6() {
	c.Core.Events.Subscribe(event.OnAttackWillLand, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.ActorIndex != c.Index {
			return false
		}
		if c.c6stacks == 0 {
			return false
		}
		if atk.Info.AttackTag != combat.AttackTagNormal {
			return false
		}
		if !c.StatusIsActive(c6Key) {
			return false
		}
		// Queue the attack
		ai := combat.AttackInfo{ // TODO: idk about the ICD and attack on this one being the same as the normal dust bolt
			ActorIndex:   c.Index,
			Abil:         "Raiment: Just Scales (C6)",
			AttackTag:    combat.AttackTagElementalArt,
			ICDTag:       combat.ICDTagNone,
			ICDGroup:     combat.ICDGroupDefault,
			Element:      attributes.Electro,
			Durability:   25,
			IsDeployable: true,
			Mult:         1.0,
			FlatDmg:      c.Stat(attributes.EM) * 2.5, // this is the A4
		}

		c.Core.QueueAttack(
			ai,
			combat.NewCircleHit(c.Core.Combat.Player(), 1, false, combat.TargettableEnemy),
			0,
			0,
		)

		c.c6stacks--
		return false
	}, "cyno-c6")
}

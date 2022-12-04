package cyno

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/gadget"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	c1Key    = "cyno-c1"
	c6Key    = "cyno-c6"
	c6ICDKey = "cyno-c6-icd"
)

// After using Sacred Rite: Wolf's Swiftness, Cyno's Normal Attack SPD will be increased by 20% for 10s.
// If the Judication effect of his Passive Talent Featherfall Judgment is triggered during Secret Rite: Chasmic Soulfarer,
// the duration of this increase will be refreshed.
//
// You need to unlock the Passive Talent "Featherfall Judgment."
func (c *char) c1() {
	m := make([]float64, attributes.EndStatType)
	m[attributes.AtkSpd] = 0.2
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBaseWithHitlag(c1Key, 600), // 10s
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			if atk.Info.AttackTag != combat.AttackTagNormal {
				return nil, false
			}
			return m, true
		},
	})
}

// When Cyno's Normal Attacks hit opponents, his Electro DMG Bonus will
// increase by 10% for 4s. This effect can be triggered once every 0.1s. Max 5
// stacks.
func (c *char) c2() {
	const c2Key = "cyno-c2"
	const c2Icd = "cyno-c2-icd"
	stacks := 0
	m := make([]float64, attributes.EndStatType)
	c.Core.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
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

		if !c.StatModIsActive(c2Key) {
			stacks = 0
		}
		stacks++
		if stacks > 5 {
			stacks = 5
		}

		c.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag(c2Key, 240), // 4s
			AffectedStat: attributes.ElectroP,
			Amount: func() ([]float64, bool) {
				m[attributes.ElectroP] = 0.1 * float64(stacks)
				return m, true
			},
		})
		c.AddStatus(c2Icd, 6, true) // 0.1s icd
		return false
	}, "cyno-c2")
}

// When Cyno is in the Pactsworn Pathclearer state triggered by Sacred Rite:
// Wolf's Swiftness, after he triggers Electro-Charged, Superconduct,
// Overloaded, Quicken, Aggravate, Hyperbloom, or an Electro Swirl reaction, he
// will restore 3 Elemental Energy for all nearby party members (except
// himself.)
// This effect can occur 5 times within one use of Sacred Rite: Wolf’s Swiftness.
func (c *char) c4() {
	restore := func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.ActorIndex != c.Index {
			return false
		}
		if c.c4Counter > 4 { // counting from 0 to 4, 5 instances max
			return false
		}
		c.c4Counter++
		for _, this := range c.Core.Player.Chars() {
			// not for cyno
			if this.Index != c.Index {
				this.AddEnergy("cyno-c4", 3)
			}
		}

		return false
	}

	restoreNoGadget := func(args ...interface{}) bool {
		if _, ok := args[0].(*gadget.Gadget); ok {
			return false
		}
		return restore(args...)
	}
	c.Core.Events.Subscribe(event.OnOverload, restoreNoGadget, "cyno-c4")
	c.Core.Events.Subscribe(event.OnElectroCharged, restoreNoGadget, "cyno-c4")
	c.Core.Events.Subscribe(event.OnSuperconduct, restoreNoGadget, "cyno-c4")
	c.Core.Events.Subscribe(event.OnQuicken, restoreNoGadget, "cyno-c4")
	c.Core.Events.Subscribe(event.OnAggravate, restoreNoGadget, "cyno-c4")
	c.Core.Events.Subscribe(event.OnHyperbloom, restore, "cyno-c4")
	c.Core.Events.Subscribe(event.OnSwirlElectro, restoreNoGadget, "cyno-c4")
}

// After using Sacred Rite: Wolf's Swiftness or triggering the Judication effect of the Passive Talent "Featherfall Judgment,"
// Cyno will gain 4 stacks of the "Day of the Jackal" effect. When he hits opponents with Normal Attacks,
// he will consume 1 stack of "Day of the Jackal" to fire off one Duststalker Bolt.
// "Day of the Jackal" lasts for 8s. Max 8 stacks. It will be canceled once Pactsworn Pathclearer ends.
// A maximum of 1 Duststalker Bolt can be unleashed this way every 0.4s.
// You must first unlock the Passive Talent "Featherfall Judgment."
func (c *char) c6() {
	c.Core.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.ActorIndex != c.Index {
			return false
		}
		if c.c6Stacks == 0 {
			return false
		}
		if atk.Info.AttackTag != combat.AttackTagNormal {
			return false
		}
		if !c.StatusIsActive(c6Key) {
			return false
		}
		if c.StatusIsActive(c6ICDKey) {
			return false
		}

		c.AddStatus(c6ICDKey, 24, true)

		// technically should use ICDGroupCynoC6, but it's just reskinned standard ICD
		ai := combat.AttackInfo{
			ActorIndex:   c.Index,
			Abil:         "Raiment: Just Scales (C6)",
			AttackTag:    combat.AttackTagElementalArt,
			ICDTag:       combat.ICDTagElementalArt,
			ICDGroup:     combat.ICDGroupDefault,
			StrikeType:   combat.StrikeTypeSlash,
			Element:      attributes.Electro,
			Durability:   25,
			IsDeployable: true,
			Mult:         1.0,
			FlatDmg:      c.Stat(attributes.EM) * 2.5, // this is the A4
		}

		c.Core.QueueAttack(
			ai,
			combat.NewCircleHit(
				c.Core.Combat.Player(),
				c.Core.Combat.PrimaryTarget(),
				nil,
				0.3,
			),
			0,
			0,
		)

		c.c6Stacks--
		return false
	}, "cyno-c6")
}

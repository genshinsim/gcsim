package cyno

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
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
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBaseWithHitlag(c1Key, 600), // 10s
		AffectedStat: attributes.AtkSpd,
		Amount: func() ([]float64, bool) {
			if c.Core.Player.CurrentState() != action.NormalAttackState {
				return nil, false
			}
			return m, true
		},
	})
}

const c2Key = "cyno-c2"
const c2ICD = "cyno-c2-icd"

// When Cyno's Normal Attacks hit opponents, his Electro DMG Bonus will
// increase by 10% for 4s. This effect can be triggered once every 0.1s. Max 5
// stacks.
func (c *char) makeC2CB() combat.AttackCBFunc {
	if c.Base.Cons < 2 {
		return nil
	}
	return func(a combat.AttackCB) {
		if a.Target.Type() != combat.TargettableEnemy {
			return
		}
		if c.StatusIsActive(c2ICD) {
			return
		}
		c.AddStatus(c2ICD, 0.1*60, true)

		if !c.StatModIsActive(c2Key) {
			c.c2Stacks = 0
		}
		c.c2Stacks++
		if c.c2Stacks > 5 {
			c.c2Stacks = 5
		}

		m := make([]float64, attributes.EndStatType)
		c.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag(c2Key, 4*60),
			AffectedStat: attributes.ElectroP,
			Amount: func() ([]float64, bool) {
				m[attributes.ElectroP] = 0.1 * float64(c.c2Stacks)
				return m, true
			},
		})
	}
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
// Cyno will gain 4 stacks of the "Day of the Jackal" effect.
func (c *char) c6Init() {
	if c.Base.Cons < 6 {
		return
	}
	c.AddStatus(c6Key, 8*60, true)
	c.c6Stacks += 4
	if c.c6Stacks > 8 {
		c.c6Stacks = 8
	}
}

// When he hits opponents with Normal Attacks, he will consume 1 stack of "Day of the Jackal" to fire off one Duststalker Bolt.
// "Day of the Jackal" lasts for 8s. Max 8 stacks. It will be canceled once Pactsworn Pathclearer ends.
// A maximum of 1 Duststalker Bolt can be unleashed this way every 0.4s.
// You must first unlock the Passive Talent "Featherfall Judgment."
func (c *char) makeC6CB() combat.AttackCBFunc {
	if c.Base.Cons < 6 || c.c6Stacks == 0 || !c.StatusIsActive(c6Key) {
		return nil
	}
	return func(a combat.AttackCB) {
		if a.Target.Type() != combat.TargettableEnemy {
			return
		}
		if c.c6Stacks == 0 {
			return
		}
		if !c.StatusIsActive(c6Key) {
			return
		}
		if c.StatusIsActive(c6ICDKey) {
			return
		}
		c.AddStatus(c6ICDKey, 0.4*60, true)
		c.c6Stacks--

		// technically should use ICDGroupCynoC6, but it's just reskinned standard ICD
		ai := combat.AttackInfo{
			ActorIndex:   c.Index,
			Abil:         "Raiment: Just Scales (C6)",
			AttackTag:    attacks.AttackTagElementalArt,
			ICDTag:       attacks.ICDTagElementalArt,
			ICDGroup:     combat.ICDGroupDefault,
			StrikeType:   attacks.StrikeTypeSlash,
			Element:      attributes.Electro,
			Durability:   25,
			IsDeployable: true,
			Mult:         1.0,
			FlatDmg:      c.a4Bolt(),
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
	}
}
